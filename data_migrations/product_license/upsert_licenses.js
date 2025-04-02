require("dotenv").config();
const { Pool } = require("pg");
const fs = require("fs");
const path = require("path");
const fetch = (...args) =>
  import("node-fetch").then(({ default: fetch }) => fetch(...args));

// Database configuration
const pool = new Pool({
  user: process.env.DB_USER || "postgres",
  host: process.env.DB_HOST || "localhost",
  database: process.env.DB_NAME || "beetle",
  password: process.env.DB_PASSWORD || "postgres",
  port: process.env.DB_PORT || 5432,
});

pool.connect((err, client, release) => {
  if (err) {
    console.error("Error connecting to the database:", err);
    return;
  }
  console.log("Successfully connected to database");
});

const BATCH_SIZE = 1000;
const DOWNLOAD_URL =
  "https://health-products.canada.ca/api/natural-licences/productlicence/?lang=en&type=json";

// Function to convert flag values to proper booleans
function convertFlagToBoolean(value) {
  if (value === null || value === undefined) return null;
  if (typeof value === "boolean") return value;
  if (typeof value === "string") {
    value = value.toLowerCase();
    if (value === "true" || value === "1" || value === "yes") return true;
    if (value === "false" || value === "0" || value === "no") return false;
  }
  if (typeof value === "number") {
    return value === 1;
  }
  return false; // Default to false for any other value
}

const upsertQuery = `
  INSERT INTO product_license (
    lnhpd_id,
    license_number,
    license_date,
    revised_date,
    time_receipt,
    date_start,
    product_name_id,
    product_name,
    dosage_form,
    company_id,
    company_name,
    company_name_id,
    sub_submission_type_code,
    sub_submission_type_desc,
    flag_primary_name,
    flag_product_status,
    flag_attested_monograph
  )
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
  ON CONFLICT (lnhpd_id) DO UPDATE SET
    license_number = EXCLUDED.license_number,
    license_date = EXCLUDED.license_date,
    revised_date = EXCLUDED.revised_date,
    time_receipt = EXCLUDED.time_receipt,
    date_start = EXCLUDED.date_start,
    product_name_id = EXCLUDED.product_name_id,
    product_name = EXCLUDED.product_name,
    dosage_form = EXCLUDED.dosage_form,
    company_id = EXCLUDED.company_id,
    company_name = EXCLUDED.company_name,
    company_name_id = EXCLUDED.company_name_id,
    sub_submission_type_code = EXCLUDED.sub_submission_type_code,
    sub_submission_type_desc = EXCLUDED.sub_submission_type_desc,
    flag_primary_name = EXCLUDED.flag_primary_name,
    flag_product_status = EXCLUDED.flag_product_status,
    flag_attested_monograph = EXCLUDED.flag_attested_monograph;
`;

async function downloadFile(url, outputPath) {
  console.log("Starting download from Health Canada API...");

  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(
      `Failed to download file: ${response.status} ${response.statusText}`
    );
  }

  const contentLength = response.headers.get("content-length");
  const total = parseInt(contentLength, 10);
  let downloaded = 0;

  const fileStream = fs.createWriteStream(outputPath);

  return new Promise((resolve, reject) => {
    response.body.on("data", (chunk) => {
      downloaded += chunk.length;
      const progress = ((downloaded / total) * 100).toFixed(2);
      process.stdout.write(
        `\rDownloading: ${progress}% (${(downloaded / 1024 / 1024).toFixed(
          2
        )} MB / ${(total / 1024 / 1024).toFixed(2)} MB)`
      );
      fileStream.write(chunk);
    });

    response.body.on("end", () => {
      console.log("\nDownload complete!");
      fileStream.end();
      resolve();
    });

    response.body.on("error", (err) => {
      fileStream.end();
      reject(err);
    });
  });
}

async function upsertLicenses(jsonFilePath) {
  const client = await pool.connect();

  try {
    // Get file stats to calculate total size
    const stats = await fs.promises.stat(jsonFilePath);
    const totalSize = stats.size;
    let processedSize = 0;

    console.log(
      `Starting to read file of size ${(totalSize / 1024 / 1024).toFixed(
        2
      )} MB...`
    );

    // Create a promise that resolves when the file is fully read
    const fileReadPromise = new Promise((resolve, reject) => {
      let jsonData = "";
      const readStream = fs.createReadStream(jsonFilePath);

      readStream.on("data", (chunk) => {
        processedSize += chunk.length;
        const progress = ((processedSize / totalSize) * 100).toFixed(2);
        process.stdout.write(`\rReading file: ${progress}%`);
        jsonData += chunk;
      });

      readStream.on("end", () => {
        console.log("\nFile read complete!");
        resolve(jsonData);
      });

      readStream.on("error", (error) => {
        reject(error);
      });
    });

    // Wait for file to be read
    const jsonData = await fileReadPromise;
    const licenses = JSON.parse(jsonData);

    console.log(`Starting to upsert ${licenses.length} licenses...`);

    await client.query("BEGIN");

    // Process licenses in batches
    for (let i = 0; i < licenses.length; i += BATCH_SIZE) {
      const batch = licenses.slice(i, i + BATCH_SIZE);
      console.log(
        `Processing batch ${Math.floor(i / BATCH_SIZE) + 1} of ${Math.ceil(
          licenses.length / BATCH_SIZE
        )}`
      );

      // Create a batch of promises for concurrent processing
      const batchPromises = batch.map((license) => {
        const values = [
          license.lnhpd_id,
          license.license_number,
          license.license_date,
          license.revised_date,
          license.time_receipt,
          license.date_start,
          license.product_name_id,
          license.product_name,
          license.dosage_form,
          license.company_id,
          license.company_name,
          license.company_name_id,
          license.sub_submission_type_code,
          license.sub_submission_type_desc,
          convertFlagToBoolean(license.flag_primary_name),
          convertFlagToBoolean(license.flag_product_status),
          convertFlagToBoolean(license.flag_attested_monograph),
        ];
        return client.query(upsertQuery, values);
      });

      await Promise.all(batchPromises);
    }

    await client.query("COMMIT");
    console.log("Successfully upserted all licenses");
  } catch (error) {
    await client.query("ROLLBACK");
    console.error("Error upserting licenses:", error);
    throw error;
  } finally {
    client.release();
  }
}

const jsonFilePath = path.join(__dirname, "productlicense.json");

async function main() {
  try {
    await downloadFile(DOWNLOAD_URL, jsonFilePath);
    await upsertLicenses(jsonFilePath);

    console.log("Script completed successfully");
    process.exit(0);
  } catch (error) {
    console.error("Script failed:", error);
    process.exit(1);
  }
}

main();

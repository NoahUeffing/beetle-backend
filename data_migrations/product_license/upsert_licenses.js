// Script to add product licenses and associated entities to the database from the Health Canada

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

const BATCH_SIZE = 10000;
const DOWNLOAD_URL =
  "https://health-products.canada.ca/api/natural-licences/productlicence/?lang=en&type=json";

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
  return false;
}

const upsertCompanyQuery = `
  INSERT INTO companies (company_id, company_name, company_name_id)
  SELECT * FROM UNNEST ($1::integer[], $2::text[], $3::integer[])
  ON CONFLICT (company_id) DO UPDATE SET
    company_name = EXCLUDED.company_name,
    company_name_id = EXCLUDED.company_name_id
  RETURNING id::text, company_id;
`;

const upsertDosageFormQuery = `
  INSERT INTO dosage_forms (name)
  SELECT UNNEST($1::text[])
  ON CONFLICT (name) DO UPDATE SET
    name = EXCLUDED.name
  RETURNING id::text, name;
`;

const upsertSubmissionTypeQuery = `
  INSERT INTO submission_types (code, description)
  SELECT * FROM UNNEST($1::integer[], $2::text[])
  ON CONFLICT (code) DO UPDATE SET
    description = EXCLUDED.description
  RETURNING id::text, code;
`;

const upsertProductLicenseQuery = `
  INSERT INTO product_licenses (
    lnhpd_id,
    license_number,
    license_date,
    revised_date,
    time_receipt,
    date_start,
    product_name_id,
    product_name,
    dosage_form_id,
    company_id,
    submission_type_id,
    flag_primary_name,
    flag_product_status,
    flag_attested_monograph
  )
  SELECT * FROM UNNEST(
    $1::integer[],
    $2::integer[],
    $3::timestamp[],
    $4::timestamp[],
    $5::timestamp[],
    $6::timestamp[],
    $7::integer[],
    $8::text[],
    $9::uuid[],
    $10::uuid[],
    $11::uuid[],
    $12::boolean[],
    $13::boolean[],
    $14::boolean[]
  )
  ON CONFLICT (lnhpd_id) DO UPDATE SET
    license_number = EXCLUDED.license_number,
    license_date = EXCLUDED.license_date,
    revised_date = EXCLUDED.revised_date,
    time_receipt = EXCLUDED.time_receipt,
    date_start = EXCLUDED.date_start,
    product_name_id = EXCLUDED.product_name_id,
    product_name = EXCLUDED.product_name,
    dosage_form_id = EXCLUDED.dosage_form_id,
    company_id = EXCLUDED.company_id,
    submission_type_id = EXCLUDED.submission_type_id,
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
    const jsonData = await fs.promises.readFile(jsonFilePath, "utf8");
    const licenses = JSON.parse(jsonData);

    console.log(`Starting to upsert ${licenses.length} licenses...`);

    await client.query("BEGIN");

    for (let i = 0; i < licenses.length; i += BATCH_SIZE) {
      const batch = licenses.slice(i, i + BATCH_SIZE);
      console.log(
        `Processing batch ${Math.floor(i / BATCH_SIZE) + 1} of ${Math.ceil(
          licenses.length / BATCH_SIZE
        )}`
      );

      // Batch upsert companies - deduplicate by company_id
      const companyMap = new Map();
      batch.forEach((l) => {
        if (!companyMap.has(l.company_id)) {
          companyMap.set(l.company_id, {
            id: l.company_id,
            name: l.company_name,
            nameId: parseInt(l.company_name_id) || null,
          });
        }
      });
      const companies = Array.from(companyMap.values());
      const companyIds = companies.map((c) => c.id);
      const companyNames = companies.map((c) => c.name);
      const companyNameIds = companies.map((c) => c.nameId);
      const companyResult = await client.query(upsertCompanyQuery, [
        companyIds,
        companyNames,
        companyNameIds,
      ]);
      const companyIdMap = new Map(
        companyResult.rows.map((row) => [row.company_id, row.id])
      );

      // Batch upsert dosage forms - already using Set for deduplication
      const uniqueDosageForms = [
        ...new Set(batch.map((l) => l.dosage_form).filter(Boolean)),
      ];
      let dosageFormIdMap = new Map();
      if (uniqueDosageForms.length > 0) {
        const dosageFormResult = await client.query(upsertDosageFormQuery, [
          uniqueDosageForms,
        ]);
        dosageFormIdMap = new Map(
          dosageFormResult.rows.map((row) => [row.name, row.id])
        );
      }

      // Batch upsert submission types - already using Set for deduplication
      const submissionTypeCodes = [];
      const submissionTypeDescs = [];
      const uniqueSubmissionTypes = new Set();
      batch.forEach((l) => {
        if (
          l.sub_submission_type_code &&
          !uniqueSubmissionTypes.has(l.sub_submission_type_code)
        ) {
          const code = parseInt(l.sub_submission_type_code);
          if (!isNaN(code)) {
            uniqueSubmissionTypes.add(code);
            submissionTypeCodes.push(code);
            submissionTypeDescs.push(l.sub_submission_type_desc);
          }
        }
      });
      let submissionTypeIdMap = new Map();
      if (submissionTypeCodes.length > 0) {
        const submissionTypeResult = await client.query(
          upsertSubmissionTypeQuery,
          [submissionTypeCodes, submissionTypeDescs]
        );
        submissionTypeIdMap = new Map(
          submissionTypeResult.rows.map((row) => [row.code, row.id])
        );
      }

      // Batch upsert product licenses - deduplicate by lnhpd_id
      const licenseMap = new Map();
      batch.forEach((license) => {
        if (!licenseMap.has(license.lnhpd_id)) {
          licenseMap.set(license.lnhpd_id, license);
        }
      });
      const uniqueLicenses = Array.from(licenseMap.values());

      const lnhpdIds = [];
      const licenseNumbers = [];
      const licenseDates = [];
      const revisedDates = [];
      const timeReceipts = [];
      const dateStarts = [];
      const productNameIds = [];
      const productNames = [];
      const dosageFormIds = [];
      const companyIds2 = [];
      const submissionTypeIds = [];
      const flagPrimaryNames = [];
      const flagProductStatuses = [];
      const flagAttestedMonographs = [];

      uniqueLicenses.forEach((license) => {
        lnhpdIds.push(license.lnhpd_id);
        licenseNumbers.push(parseInt(license.license_number) || null);
        licenseDates.push(license.license_date);
        revisedDates.push(license.revised_date);
        timeReceipts.push(license.time_receipt);
        dateStarts.push(license.date_start);
        productNameIds.push(parseInt(license.product_name_id) || null);
        productNames.push(license.product_name);
        dosageFormIds.push(
          license.dosage_form ? dosageFormIdMap.get(license.dosage_form) : null
        );
        companyIds2.push(companyIdMap.get(license.company_id));
        submissionTypeIds.push(
          license.sub_submission_type_code
            ? submissionTypeIdMap.get(
                parseInt(license.sub_submission_type_code)
              )
            : null
        );
        flagPrimaryNames.push(convertFlagToBoolean(license.flag_primary_name));
        flagProductStatuses.push(
          convertFlagToBoolean(license.flag_product_status)
        );
        flagAttestedMonographs.push(
          convertFlagToBoolean(license.flag_attested_monograph)
        );
      });

      await client.query(upsertProductLicenseQuery, [
        lnhpdIds,
        licenseNumbers,
        licenseDates,
        revisedDates,
        timeReceipts,
        dateStarts,
        productNameIds,
        productNames,
        dosageFormIds,
        companyIds2,
        submissionTypeIds,
        flagPrimaryNames,
        flagProductStatuses,
        flagAttestedMonographs,
      ]);
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

# This script checks that new migration files are named correctly
newMigrations=$(git diff --name-only origin/main | grep '^assets/db/migrations/' | sed 's/[^0-9]*//g' | sort -V)
startAtLine=1 # 1 
if [ ${#newMigrations} == 0 ]; then
  exit 0
fi
# Check length of timestamp
for file in $newMigrations; do
  if [ ${#file} != 14 ]; then
    echo "Migration $file timestamp is not correct length"
    exit 1
  fi
  startAtLine=$((startAtLine+1))
done
newMigrationsArr=($newMigrations)
# Check that the new migration timestamps are greater than existing migrations
for file in $(ls -1 -r ../../../assets/db/migrations/*); do
  number=$(echo $file | sed 's/[^0-9]*//g')
  if [ $number -ge ${newMigrationsArr[0]} ] && [[ ! " ${newMigrationsArr[*]} " =~ " ${number} " ]]; then
    echo "Migration $number is newer than ${newMigrationsArr[0]}"
    exit 1
  fi
done
echo "All goose checks passed"

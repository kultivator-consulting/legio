#!/usr/bin/env bash

CSV_PATH=database/migration/data

export PGPASSWORD=secret

migrate -path database/migration/ -database "postgres://root:$PGPASSWORD@localhost:5432/cortexdb?sslmode=disable" up

rm -rf $CSV_PATH/*.tmp

csv_files=$(ls $CSV_PATH|sort -Vt - -k2,2)

for file in $csv_files
do
  echo "Copying $file to postgres"
  grep -v -e '^[[:space:]]*$' $CSV_PATH/$file | sed '/^#/d' > $CSV_PATH/$file.tmp
  docker cp $CSV_PATH/$file.tmp postgres:/var/lib/postgresql/$file
  docker exec -it postgres psql -U cortexdbuser -d cortexdb -c "SET DateStyle TO ISO, DMY;" -c "\copy $(echo $file | cut -d'-' -f1) FROM '/var/lib/postgresql/$file' WITH DELIMITER '|' NULL AS 'NULL' CSV;"
done

rm -rf $CSV_PATH/*.tmp

#!/bin/bash

migrationLocation=$1
migrationName=$2

if [[ -z "${migrationName}" ]]; then 
  echo "Error: migrationName is an empty string"; 
  exit 1;
fi

migrate create -ext sql -dir ${migrationLocation} ${migrationName}
createMigrationCodeRes=$?
if [[ $createMigrationCodeRes -ne 0 ]]; then
  echo "Error: create migration process failed with code: ${createMigrationCodeRes}"
  exit 1;
fi
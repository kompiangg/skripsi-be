#!/bin/bash

source script/validate_migration.sh

username=$1
password=$2
dbHost=$3
databaseName=$4
state=$5
qty=$6
migrationLocation=$7

isValid=1

if [[ -z "${username}" ]]; then 
  echo "Error: username is an empty string"; 
  isValid=0; 
fi

if [[ -z "${password}" ]]; then 
  echo "Error: password is an empty string"; 
  isValid=0; 
fi

if [[ -z "${dbHost}" ]]; then 
  echo "Error: dbHost is an empty string"; 
  isValid=0; 
fi

if [[ -z "${databaseName}" ]]; then 
  echo "Error: databaseName is an empty string"; 
  isValid=0; 
fi

if [[ $isValid -eq 0 ]]; then 
  exit 1; 
fi

migrate -source file://${migrationLocation} -database postgres://${username}:${password}@${dbHost}/${databaseName}?sslmode=disable ${state} ${qty}
migrateCodeRes=$?
if [[ $migrateCodeRes -ne 0 ]]; then
  echo "Error: migration ${state} process failed with code: ${migrateCodeRes}"
  exit 1;
fi
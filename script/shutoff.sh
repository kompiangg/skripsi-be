#!/bin/bash

for i in $(virsh list --state-running | awk 'NR > 2 {print $1}')
do
  virsh shutdown $i

  if [[ $? -ne 0 ]]; then
    echo "Error: VM $i could not be shutdown"
    exit 1
  fi
done
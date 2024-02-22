#!/bin/bash

for i in $(virsh list --state-shutoff | awk 'NR > 2 {print $2}')
do
  echo "Turning on VM $i"

  virsh start $i

  if [[ $? -ne 0 ]]; then
    echo "Error: VM $i could not be started"
    exit 1
  fi

  echo "VM $i has been started"
done
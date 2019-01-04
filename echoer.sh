#! /usr/bin/env bash

echo "This script echoes input back to you"

while read -r input; do
  echo "You entered $input"
done

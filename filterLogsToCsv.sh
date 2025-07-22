#!/bin/bash

input="$1"
echo "Input filename: $input"
prefix=$(echo "$input" | sed -rnE 's/^(.*).log/\1/p')
if [[ -z "$prefix" ]]; then
    echo "Error: Unable to extract prefix from input filename."
    exit 1
fi
outputTime="$prefix-time.csv"
outputData="$prefix-data.csv"

echo "Output 1: $outputTime"
echo "Output 2: $outputData"

if [[ -z "$input" ]];  then
    echo "Usage: $0 <input_filename>"
    exit 1
fi

touch "$outputTime"
echo "Date,EndTime,Name,Runtime" > "$outputTime"
touch "$outputData"
echo "Date,EndTime,Name,Alloc,TotalAlloc,Sys" > "$outputData"

while IFS= read -r line; do
  timeData=$(echo "$line" | sed -rnE 's/^([0-9]{4}\/[0-9]{2}\/[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{6}) (.*)-time: ([0-9]+\w*)/\1,\2,\3,\4/p')
  dataData=$(echo "$line" | sed -rnE 's/^([0-9]{4}\/[0-9]{2}\/[0-9]{2}) ([0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{6}) (.*) - Alloc: ([0-9]+), TotalAlloc: ([0-9]+), Sys: ([0-9]+)/\1,\2,\3,\4,\5,\6/p')
  if [[ -n "$timeData" ]]; then
    echo "$timeData" >> "$outputTime"
  elif [[ -n "$dataData" ]]; then
    echo "$dataData" >> "$outputData"
  else
    echo "No match found for line: $line"
  fi

done < "$input"

echo $filename
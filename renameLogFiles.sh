#!/bin/bash

input="$1"
echo "Input directory: $input"

if [[ -z "$input" ]];  then
    echo "Usage: $0 <input_directory>"
    exit 1
fi

echo "Renamed $input"/ds-aggregation-server*.log "$input/he-aggregation-service.log"
mv "$input"/ds-aggregation-server*.log "$input/he-aggregation-service.log"

for file in "$input"/ds-c*-python.log; do
    if [[ -f "$file" ]]; then
        # keep the directory portion and append filename using the client number
        prefix=$(dirname "$file")
        # grab the number after "ds-c" from the basename
        num=$(basename "$file" | sed -rnE 's/^ds-c([0-9]+)-.*-python\.log/\1/p')
        new_name="$prefix/he-client-$num-python.log"
        mv "$file" "$new_name"
        echo "Renamed $file to $new_name"
    else
        echo "No files matching the pattern found in $input."
    fi
done

for file in "$input"/ds-c*-go.log; do
    if [[ -f "$file" ]]; then
        # keep the directory portion and append filename using the client number
        prefix=$(dirname "$file")
        # grab the number after "ds-c" from the basename
        num=$(basename "$file" | sed -rnE 's/^ds-c([0-9]+)-.*-go\.log/\1/p')
        new_name="$prefix/he-client-$num-go.log"
        mv "$file" "$new_name"
        echo "Renamed $file to $new_name"
    else
        echo "No files matching the pattern found in $input."
    fi
done



echo $prefix

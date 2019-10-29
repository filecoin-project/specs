#!/bin/bash

system_file=$1
input_file=$2
output_file=$3

mkdir -p $output_file && rmdir $output_file # Dirty hack to ensure output_file directory exists.

echo $output_file

orient/bin/orient solve --system=$system_file --in=$input_file | jq > $output_file

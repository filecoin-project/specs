#!/bin/bash

system_file=$1
input_file=$2
output_file=$3

if [[ -z "${ORIENT_DCALC}" ]]; then
    orient_bin=orient/bin/orient
else
    orient_bin=orient/bin/dcalc
    system_file=/orientd/$system_file
    input_file=/orientd/$input_file
fi;

mkdir -p $output_file && rmdir $output_file # Dirty hack to ensure output_file directory exists.

echo $output_file

$orient_bin solve --system=$system_file --in=$input_file | jq > $output_file

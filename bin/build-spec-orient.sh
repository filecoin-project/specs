output_dir=build/orient
input_dir=src/orient

orient/bin/orient solve --system=$input_dir/filecoin.orient --in=$input_dir/params.json | jq > $output_dir/solved-parameters.json

orient/bin/orient report --system=$input_dir/filecoin.orient --in=$input_dir/params.json  > $output_dir/filecoin-report.html

orient/bin/orient dump --system=$input_dir/filecoin.orient | jq > $output_dir/filecoin.json

orient/bin/orient graph --system=$input_dir/filecoin.orient --in=$input_dir/params.json > $output_dir/filecoin.dot

cat $output_dir/filecoin.dot | dot -Tsvg > $output_dir/filecoin.svg


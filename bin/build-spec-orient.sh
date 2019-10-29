output_dir=build/orient
input_dir=src/orient
diagram_dir=src/diagrams/orient

mkdir -p $output_dir
mkdir -p $diagram_dir

orient/bin/orient solve --system=$input_dir/filecoin.orient --in=$input_dir/snark-table.json | jq > $output_dir/snark-table.json
orient/bin/orient solve --system=$input_dir/filecoin.orient --in=$input_dir/filecoin.json | jq > $output_dir/solved-parameters.json
orient/bin/orient solve --system=$input_dir/filecoin.orient --in=$input_dir/multi-params.json | jq > $output_dir/multi-solved-parameters.json

orient/bin/orient report --system=$input_dir/filecoin.orient --in=$input_dir/filecoin.json  > $output_dir/filecoin-report.html

orient/bin/orient dump --system=$input_dir/filecoin.orient | jq > $output_dir/filecoin.json

orient/bin/orient graph --system=$input_dir/filecoin.orient --in=$input_dir/filecoin.json > $diagram_dir/filecoin.dot

output_dir=build/orient
input_dir=src/orient

orient/bin/orient solve --system=$input_dir/filecoin.orient --in=$input_dir/params.json | jq > $output_dir/solved-parameters.json

orient/bin/orient dump --system=$input_dir/filecoin.orient | jq > $output_dir/filecoin.json

orient/bin/orient graph --system=$input_dir/filecoin.orient $input_dir/params.json | dot -Tsvg > $output_dir/filecoin.svg


output_dir=build/orient
input_dir=src/orient
diagram_dir=src/diagrams/orient

# preprocessing step
./$(dirname $0)/preprocessing/preprocessing.sh $(ls input_dir/*json)

if [[ -z "${ORIENT_DCALC}" ]]; then
    orient_bin=orient/bin/orient
else
    orient_bin=orient/bin/dcalc
    input_dir=/orientd/$input_dir
    diagram_dir=/orientd/$diagram_dir
fi;

mkdir -p $output_dir > /dev/null 2>&1 
mkdir -p $diagram_dir > /dev/null 2>&1

$orient_bin solve --system=$input_dir/filecoin.orient --in=$input_dir/snark-table.json | jq > $output_dir/snark-table.json
$orient_bin solve --system=$input_dir/filecoin.orient --in=$input_dir/filecoin.json | jq > build/orient/solved-parameters.json
$orient_bin solve --system=$input_dir/filecoin.orient --in=$input_dir/multi-params.json | jq > $output_dir/multi-solved-parameters.json
$orient_bin solve --system=$input_dir/fast-porep.orient --in=$input_dir/fast-porep.json | jq > $output_dir/fast-porep.json
$orient_bin report --system=$input_dir/filecoin.orient --in=$input_dir/filecoin.json  > $output_dir/filecoin-report.html
$orient_bin dump --system=$input_dir/filecoin.orient | jq > $output_dir/filecoin.json
$orient_bin graph --system=$input_dir/filecoin.orient --in=$input_dir/filecoin.json > $diagram_dir/filecoin.dot

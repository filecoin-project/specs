#!/usr/bin/sh

scripts_dir=$(dirname $0)
script_name=$(basename $0)
tmp=$(mktemp)

for jinput in "$@"
do
    echo "[+] Preprocessing JSON file $(basename $jinput)"
    for script in $(ls $scripts_dir)
    do
        if [ $script == $script_name ]; then
            # dont execute this script
            continue
        fi
        fullpath=$scripts_dir/$script
        # because the following issue, we need to make a temp file
        # https://stackoverflow.com/questions/6696842/how-can-i-use-a-file-in-a-command-and-redirect-output-to-the-same-file-without-t
        ./$fullpath $jinput > $tmp
        exit_code=$?
        test $exit_code -eq 0 || (echo "Error with script $(basename $script).  Exit" && exit 1)
        # Format via jq by default
        cat $tmp | jq > $jinput
        echo -e "\t-> $(basename $script) OK"
    done
done
echo "[+] All preprocessing done"
rm $tmp

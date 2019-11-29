#!/usr/bin/env python3
#
# filter.py is a wrapper around jq to facilitate the search of 
# fields inside a JSON. It takes a JSON input file and a list 
# of field to print as arguments. It returns a JSON containing 
# only the fields relevant.
# cat json | <filter.py> field1 .... fieldn
# It can also be called with --filter (or -f) with a filename that 
# contains a list of fields to output

import sys
import subprocess
import argparse

parser = argparse.ArgumentParser()
parser.add_argument("fields",nargs="*")
parser.add_argument("-f","--filter",help="filter file, one field per line")
parsed = parser.parse_args()
if parsed.filter is None:
    args = sys.argv[1:]
else:
    with open(parsed.filter,"r") as f:
        args = f.read().strip().split("\n")

fields = " , ".join(map(lambda x: "%s: .%s?" % (x,x),args))
cmd = 'jq " .. | { ' + fields + '}"'
subprocess.call(cmd,shell=True,stdout=sys.stdout,stdin=sys.stdin)

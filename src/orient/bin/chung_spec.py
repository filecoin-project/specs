#!/usr/bin/env python3

import sys
import json
import argparse
import warnings
import math

warnings.filterwarnings("error")

def frange(start, stop=None, step=None):
    #Use float number in range() function
    # if stop and step argument is null set start=0.0 and step = 1.0
    if stop == None:
        stop = start + 0.0
        start = 0.0
    if step == None:
        step = 1.0
    while True:
        if step > 0 and start >= stop:
            break
        elif step < 0 and start <= stop:
            break
        yield float(start)
        start = start + step

def bin_entropy(x):
    return -x * math.log2(x) - (1 - x) * math.log2(1 - x)

def chung_formula(d,alpha,beta):
    h = bin_entropy
    try:
        return h(alpha) + h(beta) + d * (beta * h(alpha/beta) - h(alpha))
    except:
        return 0

def find_max_beta(d,alpha):
    betas = frange(0.01,1,0.01)
    max_beta = 0
    for b in betas:
        r = chung_formula(d,alpha,b)
        if r < 0 and b > max_beta:
            max_beta = b
    return max_beta

## find_optimal returns the lowest degree and highest beta given such an alpha
## Priority is given to lowest degree: 
## If (d1,b1) and (d2,b2) with d1 < d2, but target_beta < b1 < b2, 
## then (d1,b1) is chosen
def find_optimal(alpha, target_beta):
    degree = 1
    while True:
        beta = find_max_beta(degree,alpha) 
        if beta > target_beta:
            return (degree,beta)
        degree += 1

def parse():
    parser = argparse.ArgumentParser()
    parser.add_argument('json', nargs='?', type=argparse.FileType('r'), 
            default=sys.stdin)
    parser.add_argument("-d","--degree",help="field name in JSON representing degree")
    parser.add_argument("-a","--alpha",help="field name in JSON representing alpha")
    parser.add_argument("-b","--beta",help="field name to output for beta")
    args = parser.parse_args() 
    alpha = args.alpha if args.alpha is not None else "chung_alpha"
    beta = args.beta if args.beta is not None else "chung_beta"
    degree = args.degree if args.degree is not None else "expander_parents"
    # sys.stderr.write("parsing json %s" %args.json)
    try:
        json_input = json.load(args.json)
    except Exception as e:
        sys.exit("json reading error",e)
    finally:
        args.json.close()

    return [json_input, alpha,beta,degree]

# From https://hackersandslackers.com/extract-data-from-complex-json-python/
def extract_value(obj, key):
    arr = []

    def extract(obj, arr, key):
        if isinstance(obj, dict):
            for k, v in obj.items():
                if isinstance(v, (dict, list)):
                    extract(v, arr, key)
                elif k == key:
                    arr.append(v)
        elif isinstance(obj, list):
            for item in obj:
                extract(item, arr, key)
        return arr

    results = extract(obj, arr, key)
    if len(results) > 1:
        sys.exit("Multiple %s values inside JSON file. Exit." % key)
    elif len(results) == 0:
        return None
    return results[0]

def main():
    jinput, alphaT,betaT, degreeT = parse()
    alpha = extract_value(jinput,alphaT)
    degree = extract_value(jinput,degreeT)

    if alpha is None:
        sys.exit(1)

    v = {}
    if degree is None:
        # find min degree and max beta for this given alpha
        degree,beta = find_optimal(alpha,0.80)
        v["expander_parents"] = degree
    else:
        # find beta for this given alpha + degree
        beta = find_max_beta(degree,alpha)
    
    v["chung_beta"] = round(beta,5)
    json.dump(v, sys.stdout)

main()

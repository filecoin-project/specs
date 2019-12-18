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


def find_optimal_degree(alpha, target_beta,max_degree=1000):
    degree = 1
    while degree < max_degree:
        r = chung_formula(degree,alpha,target_beta)
        if r < 0:
            return degree
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

def inject_value(input_json,search_key,inject_key,inject_value):
    def extract(obj):
        if isinstance(obj, dict):
            for k, v in obj.items():
                if isinstance(v, (dict, list)):
                    extract(v)
                elif k == search_key:
                    # add the key/value here
                    obj[inject_key] = inject_value
                    return
        elif isinstance(obj, list):
            for item in obj:
                extract(item)
    extract(input_json)

def main():
    jinput, alphaT,betaT, degreeT = parse()
    alpha = extract_value(jinput,alphaT)
    beta = extract_value(jinput,betaT)
    degree = extract_value(jinput,degreeT)

    # find the optimal degree (lowest)
    if alpha is not None and beta is not None:
        degree = find_optimal_degree(alpha,beta)
        # print("{\"chung_degree\": %s}" % degree)
        inject_value(jinput,alphaT,degreeT,degree)
    elif alpha is not None and degree is not None:   
        beta = find_max_beta(degree,alpha)
        rounded = round(beta,5)
        inject_value(jinput,alphaT,betaT,rounded)
        # print("{\"chung_beta\": %s}" % rounded)
    # default behavior: return same thing if nothing to be done

    json.dump(jinput, sys.stdout)
    
main()

#!/bin/bash

source "$(dirname $0)/lib.sh"

# test require
require emacs emacs
require go go

# test version compare
compare_versions 10.10 10.11
compare_versions 1.2.3 1.3.2
compare_versions 1 999.999.999

# test require version
require_version "$(emacs -version)" emacs 26.3 "recommended install from you package manager"
require_version "$(go version)" go 1.10 "recommended install from https://golang.org/dl/"
require_version "$(go version)" go 1.13 "recommended install from https://golang.org/dl/"

#!/bin/bash

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/

export TMPDIR=~/tmp/
mkdir ~/tmp

# 测试
cd ${REPO_ROOT}/ipam
/usr/local/go/bin/go test .
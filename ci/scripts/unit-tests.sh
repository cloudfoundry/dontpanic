#!/bin/bash
set -e

args=$@
[ "$args" == "" ] && args="-r"
ginkgo -mod vendor -p -nodes 5 -race -skipPackage integration $args

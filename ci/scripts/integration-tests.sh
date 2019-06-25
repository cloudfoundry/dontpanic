#!/usr/bin/env bash
set -eo pipefail

args=$@
[ "$args" == "" ] && args="-r integration"
ginkgo -mod vendor -p -nodes 5 -race $args

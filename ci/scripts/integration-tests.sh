#!/usr/bin/env bash
set -eo pipefail

# source "$( dirname "$BASH_SOURCE" )/test/utils.sh"

# trap unmount_storage EXIT

# mount_storage

# make
# make prefix=/usr/bin install

# chmod +s /usr/bin/newuidmap
# chmod +s /usr/bin/newgidmap

# umask 077

args=$@
[ "$args" == "" ] && args="-r integration"
ginkgo -mod vendor -p -nodes 5 -race $args

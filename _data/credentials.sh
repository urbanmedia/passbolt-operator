#!/bin/bash

set -ex

scriptDir=$(dirname -- "$(readlink -f -- "$BASH_SOURCE")")

export PASSBOLT_URL="http://localhost:8088"
export PASSBOLT_PASSWORD="TestTest123!"
export PASSBOLT_GPG=$(cat $scriptDir/passbolt_gpg.key)
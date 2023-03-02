#!/bin/bash

set -e

# color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

function getPassboltSecret() {
    local apiVersion=$1
    kubectl get \
        passboltsecrets.${apiVersion}.passbolt.tagesspiegel.de \
        passboltsecret-sample-${apiVersion} \
        -o json
}

function isSyncStatusSuccess() {
    local jsonRsp=${1}
    local syncStatus=$(echo ${jsonRsp} | jq -r '.status.syncStatus')
    if [ "$syncStatus" != "Success" ]; then
        echo "${RED}Sync status is not Success: $syncStatus${NC}"
        exit 1
    fi
}

function isSecretExist() {
    local jsonRsp=${1}
    local secretName=$(echo ${jsonRsp} | jq -r '.metadata.name')
    echo -e "${BLUE}Checking if secret ${YELLOW}${secretName}${BLUE} exists${NC}"
    kubectl get secret ${secretName}
}

#################
# Execute tests #
#################

# apiVersions represents a list of API versions to test
# The list is ordered by priority, the first version is tested first
# Example: apiVersions="v1alpha1 v1alpha2 v1beta1 ..."
apiVersions="v1alpha1"

for apiVersion in $apiVersions; do
    echo -e "${BLUE}Testing API version: ${YELLOW}${apiVersion}${NC}"
    jsonRsp=$(getPassboltSecret $apiVersion)
    isSyncStatusSuccess "${jsonRsp}"
    isSecretExist "${jsonRsp}"
    echo -e "${GREEN}Tests passed for version ${YELLOW}${apiVersion}${NC}"
done

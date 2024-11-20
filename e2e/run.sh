#!/bin/bash -e

if [ "$DEBUG" == "true" ]; then
    set -x
fi

# load lib.sh from the same directory
source "$(dirname $0)/lib.sh"

############################
########## Test 1 ##########
############################
api_version="v1alpha2"
echo -e "${color_magenta}1: Testing API version: ${api_version}${color_reset}"
secret_name="${api_version}-simple"
createPassboltSecretV1alpha2 ${secret_name}

api_version="v1alpha2"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "4" ${payload_length}

api_version="v1alpha3"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.passboltSecrets | length")
compareLength "4" ${payload_length}

api_version="v1"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.passboltSecrets | length")
compareLength "4" ${payload_length}

echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} exists and has the right .data length${NC}"
payload_length=$(getKubernetesSecret ${secret_name} | jq -r ".data | length")
compareLength "4" ${payload_length}

############################
########## Test 2 ##########
############################
api_version="v1alpha3"
echo -e "${color_magenta}2: Testing API version: ${api_version}${color_reset}"
secret_name="${api_version}-simple"
createPassboltSecretV1alpha3 ${secret_name}

api_version="v1alpha2"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "6" ${payload_length}

api_version="v1alpha3"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r "(.spec.passboltSecrets | length) + (.spec.plainTextFields | length)")
compareLength "6" ${payload_length}

api_version="v1"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r "(.spec.passboltSecrets | length) + (.spec.plainTextFields | length)")
compareLength "6" ${payload_length}

echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} exists and has the right .data length${NC}"
payload_length=$(getKubernetesSecret ${secret_name} | jq -r ".data | length")
compareLength "6" ${payload_length}

############################
########## Test 3 ##########
############################
api_version="v1"
echo -e "${color_magenta}3: Testing API version: ${api_version}${color_reset}"
secret_name="${api_version}-simple"
createPassboltSecretV1 ${secret_name}

api_version="v1alpha2"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "6" ${payload_length}

api_version="v1alpha3"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r "(.spec.passboltSecrets | length) + (.spec.plainTextFields | length)")
compareLength "6" ${payload_length}

api_version="v1"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r "(.spec.passboltSecrets | length) + (.spec.plainTextFields | length)")
compareLength "6" ${payload_length}

echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} exists and has the right .data length${NC}"
payload_length=$(getKubernetesSecret ${secret_name} | jq -r ".data | length")
compareLength "6" ${payload_length}

############################
########## Test 4 ##########
############################
api_version="v1"
echo -e "${color_magenta}4: Testing API version: ${api_version}${color_reset}"
secret_name="${api_version}-backoff-check"
createPassboltSecretV1WithSecretNotFound ${secret_name}

echo -e $(getPassboltSecret ${secret_name} ${api_version})

sync_status=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".status.lastSync")

# check if status is not Success
if [ "${sync_status}" != "Success" ]; then
    echo -e "${color_red}Expected status to be not ${color_yellow}Success${color_red} but got ${color_yellow}${sync_status}${color_reset}"
    exit 1
fi

payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r "(.spec.passboltSecrets | length) + (.spec.plainTextFields | length)")
compareLength "0" ${payload_length}

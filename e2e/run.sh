#!/bin/bash -e

if [ "$DEBUG" == "true" ]; then
    set -x
fi

# load lib.sh from the same directory
source "$(dirname $0)/lib.sh"

############################
########## Test 3 ##########
############################
api_version="v1alpha1"
echo -e "${color_magenta}1: Testing API version: ${api_version}${color_reset}"
secret_name="${api_version}-simple"
createPassboltSecretV1alpha1 ${secret_name}

api_version="v1alpha1"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "3" ${payload_length}

api_version="v1alpha2"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "3" ${payload_length}

api_version="v1alpha3"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.passboltSecrets | length")
compareLength "3" ${payload_length}

############################
########## Test 2 ##########
############################
api_version="v1alpha2"
echo -e "${color_magenta}2: Testing API version: ${api_version}${color_reset}"
secret_name="${api_version}-simple"
createPassboltSecretV1alpha2 ${secret_name}

api_version="v1alpha1"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "3" ${payload_length}

api_version="v1alpha2"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "4" ${payload_length}

api_version="v1alpha3"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.passboltSecrets | length")
compareLength "4" ${payload_length}

############################
########## Test 3 ##########
############################
api_version="v1alpha3"
echo -e "${color_magenta}3: Testing API version: ${api_version}${color_reset}"
secret_name="${api_version}-simple"
createPassboltSecretV1alpha3 ${secret_name}

api_version="v1alpha1"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "3" ${payload_length}

api_version="v1alpha2"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r ".spec.secrets | length")
compareLength "6" ${payload_length}

api_version="v1alpha3"
echo -e "${color_blue}Checking if Kubernetes secret ${color_yellow}${secret_name}${color_blue} in version ${color_yellow}${api_version}${color_blue} exists${NC}"
payload_length=$(getPassboltSecret ${secret_name} ${api_version} | jq -r "(.spec.passboltSecrets | length) + (.spec.plainTextFields | length)")
compareLength "6" ${payload_length}

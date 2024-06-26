#!/bin/bash

color_reset="\e[0m"
color_red="\e[31m"
color_green="\e[32m"
color_yellow="\e[33m"
color_blue="\e[34m"
color_magenta="\e[35m"

# createPassboltSecret <secret-payload>
# Creates a secret from the given payload (complete manifest)
function createPassboltSecret() {
    cat <<EOF | kubectl apply -f - > /dev/null
$1
EOF
}

# getPassboltSecret <name> <api-version>
# Returns the secret in json format
function getPassboltSecret() {
    local api_version=passboltsecrets.${2}.passbolt.tagesspiegel.de
    kubectl get ${api_version} $1 -o json
}

# getKubernetesSecret <secret-name>
# Returns the Kubernetes secret in json format
function getKubernetesSecret() {
    local secretName=$1
    kubectl get secret ${secretName} -o json
}

# kubernetesSecretHasData <json-payload>
function kubernetesSecretHasData() {
    local jsonPayload=$1
    local dataLength=$(echo ${jsonPayload} | jq -r ".data | length")
    if [ "${dataLength}" -eq "0" ]; then
        echo -e "${color_red}Kubernetes secret has no data${color_reset}"
        exit 1
    fi
}

# compareLength <expected> <actual>
function compareLength() {
    local expected=$1
    local actual=$2
    if [ "${expected}" -ne "${actual}" ]; then
        echo -e "${color_red}Expected ${color_yellow}${expected}${color_red} but got ${color_yellow}${actual}${color_blue}${color_reset}"
        exit 1
    fi
    echo -e "${color_green}OK${color_reset}"
}

# createPassboltSecretV1alpha2 <name>
function createPassboltSecretV1alpha2() {
    createPassboltSecret "$(cat <<EOF
apiVersion: passbolt.tagesspiegel.de/v1alpha2
kind: PassboltSecret
metadata:
  name: ${1}
spec:
  leaveOnDelete: false
  secretType: Opaque
  secrets:
    - kubernetesSecretKey: s3_access_key
      passboltSecret:
        name: APP_EXAMPLE
        field: username
    - kubernetesSecretKey: s3_secret_key
      passboltSecret:
        name: APP_EXAMPLE
        field: password
    - kubernetesSecretKey: s3_endpoint
      passboltSecret:
        name: APP_EXAMPLE
        field: uri
    - kubernetesSecretKey: dsn
      passboltSecret:
        name: APP_EXAMPLE
        value: postgres://{{.Username}}@{{.URI}}/passbolt?sslmode=disable&password={{.Password}}&connect_timeout=10
EOF
)"
    sleep 5
}

# createPassboltSecretV1alpha3 <name>
function createPassboltSecretV1alpha3() {
    createPassboltSecret "$(cat <<EOF
apiVersion: passbolt.tagesspiegel.de/v1alpha3
kind: PassboltSecret
metadata:
  name: ${1}
spec:
  leaveOnDelete: false
  secretType: Opaque
  passboltSecrets:
    s3_access_key:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      field: username
    s3_secret_key:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      field: password
    s3_endpoint:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      field: uri
    dsn:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      value: postgres://{{.Username}}@{{.URI}}/passbolt?sslmode=disable&password={{.Password}}&connect_timeout=10
  plainTextFields:
    key: value
    foo: bar
EOF
)"
    sleep 5
}

# createPassboltSecretV1 <name>
function createPassboltSecretV1() {
    createPassboltSecret "$(cat <<EOF
apiVersion: passbolt.tagesspiegel.de/v1
kind: PassboltSecret
metadata:
  name: ${1}
spec:
  leaveOnDelete: false
  secretType: Opaque
  passboltSecrets:
    s3_access_key:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      field: username
    s3_secret_key:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      field: password
    s3_endpoint:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      field: uri
    dsn:
      id: 184734ea-8be3-4f5a-ba6c-5f4b3c0603e8
      value: postgres://{{.Username}}@{{.URI}}/passbolt?sslmode=disable&password={{.Password}}&connect_timeout=10
  plainTextFields:
    key: value
    foo: bar
EOF
)"
    sleep 5
}

# createPassboltSecretV1 <name>
function createPassboltSecretV1WithSecretNotFound() {
    createPassboltSecret "$(cat <<EOF
apiVersion: passbolt.tagesspiegel.de/v1
kind: PassboltSecret
metadata:
  name: ${1}
spec:
  leaveOnDelete: false
  secretType: Opaque
  passboltSecrets:
    secret:
      id: 00000000-0000-0000-0000-000000000000
      field: username
  plainTextFields:
    key: value
    foo: bar
EOF
)"
    sleep 5
}

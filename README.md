# Kubernetes Passbolt Operator

[![tests](https://github.com/urbanmedia/passbolt-operator/actions/workflows/tests.yaml/badge.svg?branch=main)](https://github.com/urbanmedia/passbolt-operator/actions/workflows/tests.yaml)

This repository contains the Kubernetes Operator for Passbolt. Passbolt is an open source password manager for teams. It is a self-hosted solution that allows you to manage your passwords securely and share them with your team. The Passbolt Operator is a Kubernetes Operator that allows you to synchronize your Passbolt credentials with Kubernetes Secrets. It is based on [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) framework.

## Getting Started

### Synchronize credentials with Kubernetes Secrets

The Passbolt Operator allows you to synchronize your Passbolt credentials with Kubernetes Secrets. To do so, you need to create a `PassboltSecret` resource. The `PassboltSecret` resource is a Kubernetes Custom Resource Definition (CRD) that allows you to define the Passbolt credentials that you want to synchronize with Kubernetes Secrets. The Passbolt Operator will then synchronize the Passbolt credentials with Kubernetes Secrets.

```yaml
apiVersion: passbolt.tagesspiegel.de/v1
kind: PassboltSecret
metadata:
  name: passboltsecret-sample
spec:
  leaveOnDelete: false
  secretType: Opaque
  passboltSecrets:
    s3_access_key:
      name: EXAMPLE_APP
      field: username
    dsn:
      name: EXAMPLE_APP
      value: postgres://{{ .Username }}@{{ .URI }}/mydb?sslmode=disable&password={{ .Password }}
  plainTextFields:
    key: value
    foo: bar
```

The `PassboltSecret` resource contains the following fields:

| Field | Type | Default | Required | Condition | Description |
| ----- | ---- | ------- | -------- | --------- | ----------- |
| `leaveOnDelete` | `bool` | `false` | false | - | A boolean that indicates if the Passbolt Operator should leave the Kubernetes Secret on deletion of the `PassboltSecret` resource. |
| `secretType` | `string` | `Opaque` | false | - | The type of the Kubernetes Secret. Can be either `Opaque` or `kubernetes.io/dockerconfigjson`. If the `secretType` is `kubernetes.io/dockerconfigjson`, the `passboltSecretName` field is required. If the `secretType` is `Opaque`, the `secrets` field is required. |
| `passboltSecretID` | `string` | - | false | `secretType` is `kubernetes.io/dockerconfigjson` | The ID of the Passbolt credential that contains the Docker configuration (URI, Username, Password). |
| `passboltSecrets` | `map[string]PassboltSecrets` | - | false | `secretType` is `Opaque` | A mapping of Passbolt credentials that you want to synchronize with Kubernetes Secrets. The key represents the name of the key in the Kubernetes secret to be added. |
| `passboltSecrets[*].id` | `string` | - | true | - | The ID of the Passbolt credential that you want to synchronize with Kubernetes Secrets. |
| `passboltSecrets[*].field` | `string` | - | false | - | The field of the Passbolt credential that you want to synchronize with Kubernetes Secrets. Can be one of: `username`, `password`, `uri` |
| `passboltSecrets[*].value` | `string` | - | false | - | A Go template value of the Passbolt credential that you want to synchronize with Kubernetes Secrets. Supported variables are: `Username`, `Password`, `URI`. The `secrets[*].passboltSecret.value` field is mutually exclusive with the `passboltSecrets[*].field` field. |
| `plainTextFields` | `map[string]string` | - | false | - | Assignment of plain text fields that you want to synchronize with Kubernetes Secrets. The key represents the name of the key in the Kubernetes secret to be added and the corresponding value. It is not recommended to store "secret" values such as passwords in it. |

The Passbolt Operator will then synchronize the Passbolt credentials with Kubernetes Secrets. The Passbolt Operator will create a Kubernetes Secret with the name `passbolt-secret-name` in the namespace `default`. The resulting Kubernetes Secret is defined as follows:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: passbolt-secret-sample
  namespace: default
type: Opaque
data:
  username: example
  pg_dsn: postgres://example@db.example.com:5432/mydb?sslmode=disable&password=example
```

Under the hood, the Passbolt Operator does the following during the reconciliation loop:

1. Retrieve the Passbolt CRD from Kubernetes.
2. Retrieve the `secrets[*].passboltSecret.name` credentials from Passbolt.
3. Create a Kubernetes secret with the name `passbolt-secret-name` in the namespace `default` with the `secrets[*].kubernetesSecretKey` key and the `secrets[*].passboltSecret.name` value.

If an error occurs during the reconciliation loop, the Passbolt Operator will update the `.status.syncStatus` field to `Error` and adds the error message to the `.status.syncErrors` field of the `PassboltSecret` resource. If the reconciliation loop is successful, the Passbolt Operator will update the `.status.syncStatus` field of the `PassboltSecret` resource with the message `Success`.

### Installation

For both installation methods, you need to create a Kubernetes Secret with the Passbolt credentials. To do so, you need to run the following command:

```bash
kubectl create secret generic controller-passbolt-secret \
  --from-file=gpg_key='/path/to/my/gpg.key' \
  --from-literal=password='<my-user-password>' \
  --from-literal=url='<my-passbolt-url>' \
  --namespace system
```

#### Prerequisites

In order to install the Passbolt Operator, you need to have the following prerequisites installed on your K8s cluster:

- [cert-manager](https://cert-manager.io/docs/installation/kubernetes/) >= v1.7

#### Kustomize

**ATTENTION**: This installation method is not recommended for production environments. If you want to install the Passbolt Operator in a production environment, please refer to the [Helm](#helm) installation method.

To install the Passbolt Operator with Kustomize, for example in a local Kind cluster, we expect that you have a working Kubernetes cluster with the `controller-passbolt-secret` secret created. To install the Passbolt Operator, you need to run the following commands:

1. Load the controller-manager image into the Kind cluster. See [Kind Load](#in-cluster-testing-kind)

2.

```bash
make deploy
```

#### Helm

For the installation instructions regarding Helm, please refer to the [Helm-Chart](https://github.com/urbanmedia/passbolt-operator-helm) repository.

### Configuration

The Passbolt Operator can be configured with the following environement variables:

- `PASSBOLT_URL`: The URL of the Passbolt instance.
- `PASSBOLT_GPG`: The GPG key to identify the user.
- `PASSBOLT_PASSWORD`: The password of the Passbolt user.

## Development

### Prerequisites

- [Go](https://golang.org/dl/) >= v1.19
- [Docker](https://docs.docker.com/get-docker/) >= 20.10
- [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) >= v3.7
- [Kubectl](https://kubernetes.io/docs/tasks/tools/) >= v1.25
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) >= v0.17
- mysql-client >= 15.1 (`mysql --version` => `mysql  Ver 15.1 Distrib 10.6.11-MariaDB`)
  - [Prometheus Operator CRDs installed in the cluster](https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-operator-crds) >= 5.1
  - [cert-manager](https://cert-manager.io/docs/installation/helm/) >= v1.12

### Setup the development environment

To setup the development environment, you need to run the following commands:

```bash
docker-compose up -d
```

Restore the database

```bash
mysql \
  --host=127.0.0.1 \
  --port=13306 \
  --database=passbolt \
  --user=passbolt \
  --password=P4ssb0lt < _data/passbolt_db.sql
```

### Create another API (Version)

Kubebuilder allows you to bootstrap a new API Version. To do so, you need to run the following command:

```bash
kubebuilder create api --group passbolt --version v1 --kind PassboltSecret
```

### Create validation and defaulting webhooks

Kubebuilder allows you to bootstrap a new webhook. To do so, you need to run the following command:

```bash
kubebuilder create webhook --group passbolt --version v1 --kind PassboltSecret --defaulting --programmatic-validation
```

### Start the Operator

Since the Operator requires a running instance of Passbolt, we will use the [Passbolt Docker image](https://hub.docker.com/r/passbolt/passbolt) to start a Passbolt instance. To start the Passbolt instance, you need to run the following command:

```bash
docker-compose up -d
```

When the Passbolt instance is up and running, the second step would be to expose the Passbolt instance credentials to the Operator. To do so, you need to run the following command:

```bash
./_data/credentials.sh
```

The last step would be to start the Operator. To do so, you need to run the following command:

```bash
make run
```

### Create a User in Passbolt

To create a user in Passbolt, you need to run the following command:

```bash
docker exec -ti passbolt /usr/share/php/passbolt/bin/cake \
  passbolt register_user \
  -u user.example@mydomain.local \
  -f user \
  -l example \
  -r admin
```

We already created a user with the above command. You can retrieve the password and the GPG key from the _data/credentials.sh file.

### Test the Operator

In order to run the end-to-end and unit tests, we need to start passbolt locally. To do so, you need to run the following command:

```bash
docker-compose up -d
```

When the Passbolt instance is up and running, the second step would be to execute the end-to-end and unit tests. To do so, you need to run the following command:

```bash
make test
```

Note:

When we bootstraped passbolt, we created two secrets:

| ID | Name |
|----|------|
| `184734ea-8be3-4f5a-ba6c-5f4b3c0603e8` | `APP_EXAMPLE` |
| `cec328ec-cb1f-48f6-be1e-1ca35fc3c62d` | `APP2_EXAMPLE` |

### In Cluster Testing (Kind)

In order to run the end-to-end tests in a Kubernetes cluster, we need to load the Passbolt Operator image into the Kind cluster. To do so, you need to run the following command:

```bash
kind load docker-image <image-name> --name <cluster-name> --nodes <node-names-separated-by-a-comma>
```

When the Passbolt Operator image is loaded into the Kind cluster, the second step would be to deploy the Passbolt Operator in the Kind cluster. To do so, you need to run the following command:

```bash
IMG="my-img-name:latest" make deploy
```

### Continuous Integration (CI)

During the continuous integration, we automatically run the end-to-end and unit tests. To do so, we use the [GitHub Actions](https://github.com/features/actions) platform. The GitHub Actions configuration is defined in the [.github/workflows](.github/workflows) directory.

During the end-to-end and unit tests, we restore the the mysql Passbolt database from the `_data/passbolt_db.sql` file. The '_data/passbolt_db.sql' file was generated with the following command:

```bash
mysqldump \
  --host=127.0.0.1 \
  --port=13306 \
  --databases passbolt \
  --user=passbolt \
  --password=P4ssb0lt > _data/passbolt_db.sql
```

## Contributing

### Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [leon.steinhaeuser@tagesspiegel.de].

### Contributing Guide

If you want to contribute to this project, please read the [Contributing Guide](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).

## Security

If you discover a security vulnerability within this project, please use issue form [Report a security vulnerability](https://github.com/urbanmedia/passbolt-operator/security/advisories/new). All security vulnerabilities will be promptly addressed. Please see [Security Policy](SECURITY.md) for more information.

## Support

If you need help with Passbolt, please contact [Passbolt Support](https://www.passbolt.com/support).

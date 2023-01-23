# Kubernetes Passbolt Operator

[![tests](https://github.com/leonsteinhaeuser/passbolt-operator/actions/workflows/tests.yaml/badge.svg?branch=main)](https://github.com/leonsteinhaeuser/passbolt-operator/actions/workflows/tests.yaml)

This repository contains the Kubernetes Operator for Passbolt. Passbolt is an open source password manager for teams. It is a self-hosted solution that allows you to manage your passwords securely and share them with your team. The Passbolt Operator is a Kubernetes Operator that allows you to synchronize your Passbolt credentials with Kubernetes Secrets. It is based on [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) framework.

## Getting Started

### Synchronize credentials with Kubernetes Secrets

The Passbolt Operator allows you to synchronize your Passbolt credentials with Kubernetes Secrets. To do so, you need to create a `PassboltSecret` resource. The `PassboltSecret` resource is a Kubernetes Custom Resource Definition (CRD) that allows you to define the Passbolt credentials that you want to synchronize with Kubernetes Secrets. The Passbolt Operator will then synchronize the Passbolt credentials with Kubernetes Secrets.

```yaml
apiVersion: passbolt.tagesspiegel.de/v1alpha1
kind: PassboltSecret
metadata:
  name: passbolt-secret-sample
  namespace: default
spec:
  secrets:
    - kubernetesSecretKey: "" # The key of the Kubernetes Secret that you want to synchronize with the Passbolt credential.
      passboltSecret:
        name: "" # The name of the Passbolt credential that you want to synchronize with Kubernetes Secrets.
        field: username # The field of the Passbolt credential that you want to synchronize with Kubernetes Secrets.
```

The `PassboltSecret` resource contains the following fields:

- `secrets`: A list of Passbolt credentials that you want to synchronize with Kubernetes Secrets. Each Passbolt credential is defined by the following fields:
  - `passboltSecret.name`: The name of the Passbolt credential that you want to synchronize with Kubernetes Secrets.
  - `passboltSecret.field`: The field of the Passbolt credential that you want to synchronize with Kubernetes Secrets.
  - `kubernetesSecretKey`: The key of the Kubernetes Secret that you want to synchronize with the Passbolt credential.

The Passbolt Operator will then synchronize the Passbolt credentials with Kubernetes Secrets. The Passbolt Operator will create a Kubernetes Secret with the name `passbolt-secret-name` in the namespace `default`. The resulting Kubernetes Secret is defined as follows:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: passbolt-secret-sample
  namespace: default
type: Opaque
data:
  KUBERNETES_SECRET_KEY: PASSBOLT_SECRET_NAME
```

Under the hood, the Passbolt Operator does the following during the reconciliation loop:

1. Retrieve the Passbolt CRD from Kubernetes.
2. Retrieve the `secrets[*].name` credentials from Passbolt.
3. Create a Kubernetes secret with the name `passbolt-secret-name` in the namespace `default` with the `secrets[*].kubernetesSecretKey` key and the `secrets[*].name` value.

If an error occurs during the reconciliation loop, the Passbolt Operator will retry to reconcile the `PassboltSecret` after 30 seconds and increments the `status.reconcileErrorCount` field. If the `status.reconcileErrorCount` field is greater than 5, the Passbolt Operator will stop reconciling the `PassboltSecret` for a period of 30 minutes (except configuration changes).

### Installation

For both installation methods, you need to create a Kubernetes Secret with the Passbolt credentials. To do so, you need to run the following command:

```bash
kubectl create secret generic controller-passbolt-secret \
  --from-files=gpg_key='/path/to/my/gpg.key' \
  --from-literal=password='<my-user-password>' \
  --from-literal=url='<my-passbolt-url>' \
  --namespace system
```

#### Kustomize

// TODO

#### Helm

// TODO

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
kubebuilder create api --group passbolt --version v1alpha1 --kind PassboltSecret
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

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to [leonsteinhaeuser@gmail.com].

### Contributing Guide

If you want to contribute to this project, please read the [Contributing Guide](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).

## Security

If you discover a security vulnerability within this project, please send an e-mail to [leonsteinhaeuser@gmail.com] instead of using the issue tracker. All security vulnerabilities will be promptly addressed. Please see [Security Policy](SECURITY.md) for more information.

## Support

If you need help with Passbolt, please contact [Passbolt Support](https://www.passbolt.com/support).

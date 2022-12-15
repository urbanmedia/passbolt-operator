# Kubernetes Passbolt Operator

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
    - name: PASSBOLT_SECRET_NAME # The name of the Passbolt credential that you want to synchronize with Kubernetes Secrets.
      kubernetesSecretKey: KUBERNETES_SECRET_KEY # The key of the Kubernetes Secret that you want to synchronize with the Passbolt credential.
```

The `PassboltSecret` resource contains the following fields:

- `secrets`: A list of Passbolt credentials that you want to synchronize with Kubernetes Secrets. Each Passbolt credential is defined by the following fields:
  - `name`: The name of the Passbolt credential that you want to synchronize with Kubernetes Secrets.
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
2. Create a Passbolt client with only `read` permissions to retrieve the `secrets[*].name` credentials.
3. Retrieve the `secrets[*].name` credentials from Passbolt.
4. Create a Kubernetes secret with the name `passbolt-secret-name` in the namespace `default` with the `secrets[*].kubernetesSecretKey` key and the `secrets[*].name` value.
5. Delete the Passbolt client.

If an error occurs during the reconciliation loop, the Passbolt Operator will retry to reconcile the `PassboltSecret` after 30 seconds and increments the `status.reconcileErrorCount` field. If the `status.reconcileErrorCount` field is greater than 5, the Passbolt Operator will stop reconciling the `PassboltSecret` for a period of 30 minutes (except configuration changes).

### Installation

TODO

### Configuration

The Passbolt Operator can be configured with the following environement variables:

- `PASSBOLT_URL`: The URL of the Passbolt instance.
- `PASSBOLT_USERNAME`: The username of the Passbolt user.
- `PASSBOLT_PASSWORD`: The password of the Passbolt user.
- `PASSBOLT_KEY`: The key of the Passbolt user.

## Development

### Prerequisites

- [Go](https://golang.org/dl/) >= v1.19
- [Docker](https://docs.docker.com/get-docker/) >= 20.10
- [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) >= v3.7
- [Kubectl](https://kubernetes.io/docs/tasks/tools/) >= v1.25
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) >= v0.17

### Create another API (Version)

Kubebuilder allows you to bootstrap a new API Version. To do so, you need to run the following command:

```bash
kubebuilder create api --group passbolt --version v1alpha1 --kind PassboltSecret
```

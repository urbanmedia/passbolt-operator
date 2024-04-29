package util

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	passboltv1alpha3 "github.com/urbanmedia/passbolt-operator/api/v1alpha3"
	"github.com/urbanmedia/passbolt-operator/pkg/passbolt"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

// UpdateSecret updates the kubernetes secret with the data from passbolt
// The thrown error is of type SyncError
func UpdateSecret(ctx context.Context, clnt *passbolt.Client, scheme *runtime.Scheme, pbscrt *passboltv1alpha3.PassboltSecret, secret *corev1.Secret) func() error {
	fmt.Println(pbscrt.Spec.SecretType)
	secret.Data = make(map[string][]byte)
	return func() error {
		switch pbscrt.Spec.SecretType {
		case corev1.SecretTypeDockerConfigJson:
			// get secret from passbolt
			secretData, err := clnt.GetSecret(ctx, *pbscrt.Spec.PassboltSecretID)
			if err != nil {
				return passboltv1alpha3.SyncError{
					Message:          err.Error(),
					PassboltSecretID: *pbscrt.Spec.PassboltSecretID,
					Time:             v1.Now(),
				}
			}
			dockerConfigJson, err := getSecretDockerConfigJson(secretData)
			if err != nil {
				return passboltv1alpha3.SyncError{
					Message:          err.Error(),
					PassboltSecretID: *pbscrt.Spec.PassboltSecretID,
					Time:             v1.Now(),
				}
			}
			secret.Data = dockerConfigJson
		case corev1.SecretTypeOpaque:
			for key, value := range pbscrt.Spec.PlainTextFields {
				secret.Data[key] = []byte(value)
			}

			// iterate over all secrets and get secret from passbolt
			for secretKeyName, pbSecret := range pbscrt.Spec.PassboltSecrets {
				secretData, err := clnt.GetSecret(ctx, pbSecret.ID)
				if err != nil {
					return passboltv1alpha3.SyncError{
						Message:          err.Error(),
						PassboltSecretID: pbSecret.ID,
						SecretKey:        secretKeyName,
						Time:             v1.Now(),
					}
				}

				// TODO: normally this should not be necessary because we already assigned an empty map
				// TODO: for some reason, it is necessary to check if the map is nil at this point
				if secret.Data == nil {
					secret.Data = make(map[string][]byte)
				}

				switch {
				// check if field is set
				// if field is set, get field value from passbolt secret and set it as kubernetes secret value
				case pbSecret.Field != "":
					secret.Data[secretKeyName] = []byte(secretData.FieldValue(pbSecret.Field))
					continue
				// check if value is set
				// if value is set, parse value as template and set it as kubernetes secret value
				case pbSecret.Value != nil:
					bts, err := getSecretTemplateValueData(*pbSecret.Value, secretData)
					if err != nil {
						return passboltv1alpha3.SyncError{
							Message:          err.Error(),
							PassboltSecretID: pbSecret.ID,
							SecretKey:        secretKeyName,
							Time:             v1.Now(),
						}
					}
					secret.Data[secretKeyName] = bts
					continue
					// neither field nor value is set
				default:
					return passboltv1alpha3.SyncError{
						Message:          "either field or value must be set",
						PassboltSecretID: pbSecret.ID,
						SecretKey:        secretKeyName,
						Time:             v1.Now(),
					}
				}
			}
		// secret type is not supported
		default:
			return passboltv1alpha3.SyncError{
				Message: fmt.Sprintf("secret type %s is not supported", pbscrt.Spec.SecretType),
				Time:    v1.Now(),
			}
		}
		// set owner reference if LeaveOnDelete was set to false
		if !pbscrt.Spec.LeaveOnDelete {
			// set owner reference
			err := ctrl.SetControllerReference(pbscrt, secret, scheme)
			if err != nil {
				return passboltv1alpha3.SyncError{
					Message: err.Error(),
					Time:    v1.Now(),
				}
			}
		}
		return nil
	}
}

func getSecretDockerConfigJson(secret *passbolt.PassboltSecretDefinition) (map[string][]byte, error) {
	// create docker auth config
	dockerAuthConfig := map[string]any{
		"auths": map[string]any{
			secret.URI: map[string]string{
				"auth": base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", secret.Username, secret.Password))),
			},
		},
	}
	bts, err := json.Marshal(dockerAuthConfig)
	if err != nil {
		return nil, err
	}
	return map[string][]byte{
		corev1.DockerConfigJsonKey: bts,
	}, nil
}

func getSecretTemplateValueData(templateStr string, secret *passbolt.PassboltSecretDefinition) ([]byte, error) {
	tmpl, err := template.New("value").Funcs(sprig.FuncMap()).Parse(templateStr)
	if err != nil {
		return nil, err
	}
	target := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(target, *secret)
	if err != nil {
		return nil, err
	}
	return target.Bytes(), nil
}

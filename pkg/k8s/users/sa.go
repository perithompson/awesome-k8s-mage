package users

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/perithompson/awesome-k8s-mage/pkg/k8s/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func CreateServiceAccount(name, namespace string) error {
	return sh.RunV(constants.KubectlCmd, "create", "serviceaccount", name, "-n", namespace)
}

func DeleteServiceAccount(name, namespace string) error {
	return sh.RunV(constants.KubectlCmd, "delete", "serviceaccount", name, "-n", namespace)
}

func CreateAccountTokenSecret(name, namespace string) (string, error) {
	selector := fmt.Sprintf("kubernetes.io/service-account.name=%s", name)
	secretName, err := sh.Output(constants.KubectlCmd, "get", "secret", "-l", selector, "-n", namespace, constants.OutName)
	if err != nil {
		return "", err
	}
	if secretName != "" {
		return strings.Split(secretName, "/")[1], nil
	}
	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-token-", name),
			Namespace:    namespace,
			Annotations: map[string]string{
				"kubernetes.io/service-account.name": name,
			},
			Labels: map[string]string{
				"kubernetes.io/service-account.name": name,
			},
		},
		Immutable:  new(bool),
		Data:       map[string][]byte{},
		StringData: map[string]string{},
		Type:       "kubernetes.io/service-account-token",
	}
	tempFile, err := os.CreateTemp(os.TempDir(), constants.TempPrefix)
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())
	y, err := yaml.Marshal(secret)
	if err != nil {
		return "", err
	}
	if err = os.WriteFile(tempFile.Name(), y, fs.ModeAppend); err != nil {
		return "", err
	}
	if err = sh.RunV(constants.KubectlCmd, "create", "-f", tempFile.Name()); err != nil {
		return "", err
	}
	secretName, err = sh.Output(constants.KubectlCmd, "get", "secret", "-l", selector, "-n", namespace, constants.OutName)
	if err != nil {
		return "", err
	}
	return strings.Split(secretName, "/")[1], err
}

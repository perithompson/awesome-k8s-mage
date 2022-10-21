package k8s

import (
	"github.com/magefile/mage/sh"
	k8sConstants "github.com/perithompson/awesome-k8s-mage/pkg/k8s/constants"
)

func GetKubectlClientVersion() (string, error) {
	version, err := sh.Output(k8sConstants.KubectlCmd, "version", k8sConstants.OutJson)
	if err != nil {
		return "", err
	}
	return version, nil
}

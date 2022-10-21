package users

import (
	"github.com/magefile/mage/sh"
	"github.com/perithompson/awesome-k8s-mage/pkg/k8s/constants"
)

func CreateServiceAccount(name, namespace string) error {
	return sh.RunV(constants.KubectlCmd, "create", "serviceaccount", name, "-n", namespace)
}

func DeleteServiceAccount(name, namespace string) error {
	return sh.RunV(constants.KubectlCmd, "delete", "serviceaccount", name, "-n", namespace)
}

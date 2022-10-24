package users

import (
	"os"

	"github.com/magefile/mage/sh"
	"github.com/perithompson/awesome-k8s-mage/pkg/k8s/constants"
)

func AuthCanI(kubeconfig, verb, resource string) error {
	f, err := os.CreateTemp(os.TempDir(), constants.TempPrefix)
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	f.WriteString(kubeconfig)
	return sh.RunV(constants.KubectlCmd, "auth", "can-i", constants.KubeconfigArg(f.Name()), resource, verb)
}

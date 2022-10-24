package users

import (
	"fmt"
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
	out, err := sh.Output(constants.KubectlCmd, "auth", "can-i", constants.KubeconfigArg(f.Name()), resource, verb)
	fmt.Println(out)
	if err != nil {
		if out == "no" {
			return nil
		}
		return err
	}
	return nil
}

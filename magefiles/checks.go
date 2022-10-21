//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	k8s "github.com/perithompson/awesome-k8s-mage/pkg/k8s/misc"
)

type Check mg.Namespace

// K8sVersions prints current k8s versions
func (Check) K8sVersions() error {
	version, err := k8s.GetKubectlClientVersion()
	if err != nil {
		return err
	}
	fmt.Println(version)
	return nil
}

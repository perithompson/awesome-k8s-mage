//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	k8sRoles "github.com/perithompson/awesome-k8s-mage/pkg/k8s/roles"
	k8sUser "github.com/perithompson/awesome-k8s-mage/pkg/k8s/users"
	"github.com/perithompson/awesome-k8s-mage/pkg/prompt"
)

type User mg.Namespace

// Create creates a user
func (User) Create(username, namespace string) error {
	ok := prompt.YesNo(fmt.Sprintf("Are you sure you want to create %s/%s", namespace, username), true)
	if ok {
		return k8sUser.CreateServiceAccount(username, namespace)
	} else {
		return nil
	}
}

// Delete deletes a user
func (User) Delete(username, namespace string) error {
	ok := prompt.YesNo(fmt.Sprintf("Are you sure you want to delete %s/%s", namespace, username), false)
	if ok {
		return k8sUser.DeleteServiceAccount(username, namespace)
	} else {
		return nil
	}
}

// CanThey runs can-i as the users kubeconfig
func (User) CanThey(username, namespace, verb, resource string) error {
	kubeconfig, err := k8sUser.Kubeconfig(username, namespace)
	if err != nil {
		return err
	}
	return k8sUser.AuthCanI(kubeconfig, resource, verb)
}

// GetRoles Gets a list of roles where the user is a member
func (User) GetRoles(username, namespace string) error {
	return k8sRoles.GetRoleBindings(username, namespace)
}

// Kubeconfig gets the user kubeconfig
func (User) Kubeconfig(username, namespace string) error {
	kubeconfig, err := k8sUser.Kubeconfig(username, namespace)
	if err != nil {
		return err
	}
	fmt.Println(kubeconfig)
	return nil
}

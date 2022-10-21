package roles

import (
	"encoding/json"
	"fmt"

	"github.com/magefile/mage/sh"
	"github.com/perithompson/awesome-k8s-mage/pkg/k8s/constants"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/yaml"
)

func CreateRoleBinding(username, namespace, role string) error {
	saArg := fmt.Sprintf("--serviceaccount=%s:%s", namespace, username)
	crArg := fmt.Sprintf("--clusterrole=%s", role)
	crbName := fmt.Sprintf("%s:%s", username, role)
	return sh.RunV(constants.KubectlCmd, "create", "clusterrolebinding", crbName, saArg, crArg)
}

func DeleteRoleBinding(username, namespace, role string) error {
	crbName := fmt.Sprintf("%s:%s:%s", namespace, username, role)
	return sh.RunV(constants.KubectlCmd, "delete", "clusterrolebinding", crbName)
}

func GetRoleBindings(username, namespace string) error {
	clusterRoles := &rbacv1.ClusterRoleBindingList{}
	roles := &rbacv1.RoleBindingList{}
	clusterroleJson, err := sh.Output(constants.KubectlCmd, "get", "clusterrolebinding", constants.OutJson)
	if err != nil {
		return err
	}
	filteredClusterRoles := rbacv1.ClusterRoleBindingList{}
	filteredRoles := rbacv1.RoleBindingList{}
	if err = json.Unmarshal([]byte(clusterroleJson), clusterRoles); err != nil {
		return err
	}
	for _, cr := range clusterRoles.Items {
		for _, s := range cr.Subjects {
			if s.Kind == "ServiceAccount" && s.Name == username && s.Namespace == namespace {
				filteredClusterRoles.Items = append(filteredClusterRoles.Items, cr)
				continue
			}
		}
	}

	roleJson, err := sh.Output(constants.KubectlCmd, "get", "rolebinding", constants.OutJson)
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(roleJson), roles); err != nil {
		return err
	}
	for _, r := range roles.Items {
		for _, s := range r.Subjects {
			if s.Kind == "ServiceAccount" && s.Name == username && s.Namespace == namespace {
				filteredRoles.Items = append(filteredRoles.Items, r)
				continue
			}
		}
	}
	ycrs, err := yaml.Marshal(filteredClusterRoles)
	if err != nil {
		return err
	}
	fmt.Println("Resolving Cluster Roles")
	fmt.Println("-----------------------")
	fmt.Println(string(ycrs))
	yrs, err := yaml.Marshal(filteredRoles)
	if err != nil {
		return err
	}
	fmt.Println("Resolving Roles")
	fmt.Println("---------------")
	fmt.Println(string(yrs))
	return nil
}
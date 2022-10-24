package users

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"

	"github.com/magefile/mage/sh"
	"github.com/perithompson/awesome-k8s-mage/pkg/k8s/constants"
)

func GetCurrentContextString() (string, error) {
	return sh.Output(constants.KubectlCmd, "config", "current-context", "get")
}

func GetCurrentKubeconfig() (string, error) {
	return sh.Output(constants.KubectlCmd, "config", "config", "view", constants.OutJson)
}

func GetClusterFromContext(context string) (string, error) {
	query := fmt.Sprintf("{.contexts[?(@.name=='%s')].context}", context)
	c := &clientcmdapi.Context{}
	jsonContext, err := sh.Output(constants.KubectlCmd, "config", "view", constants.OutJsonPath(query))
	if err != nil {
		return "", err
	}
	if err = json.Unmarshal([]byte(jsonContext), c); err != nil {
		return "", err
	}
	return c.Cluster, nil

}

func GetCurrentCluster() (string, error) {
	currentContext, err := GetCurrentContextString()
	if err != nil {
		return "", err
	}
	fmt.Println(currentContext)
	cluster, err := GetClusterFromContext(currentContext)
	if err != nil {
		return "", err
	}
	fmt.Println(cluster)
	query := fmt.Sprintf("{.clusters[?(@.name=='%s')]}", cluster)
	return sh.Output(constants.KubectlCmd, "config", "view", constants.OutJsonPath(query), "--allow-missing-template-keys=true", "--raw")
}

type saSecret struct {
	Name string `json:"name"`
}

func GetSATokenSecret(name, namespace string) (string, error) {
	query := fmt.Sprintf("{.secrets[?(@.name contains %s-token-)]}", name)
	tokenName, err := sh.Output(constants.KubectlCmd, "get", "serviceaccount", name, "-n", namespace, constants.OutJsonPath(query))
	if err != nil {
		return "", err
	}
	if tokenName == "" {
		tokenName, err = CreateAccountTokenSecret(name, namespace)
		if err != nil {
			return "", err
		}
		return tokenName, nil
	}
	secretName := &saSecret{}
	if err := json.Unmarshal([]byte(tokenName), secretName); err != nil {
		return "", err
	}

	return secretName.Name, nil
}

func GetToken(user, namespace string) (string, error) {
	secret, err := GetSATokenSecret(user, namespace)
	if err != nil {
		return "", err
	}

	token, err := sh.Output(constants.KubectlCmd, "get", "secret", secret, "-n", namespace, constants.OutJsonPath("{.data.token}"))
	if err != nil {
		return "", err
	}
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// Get gets a user Kubeconfig in the current cluster
func Kubeconfig(user, namespace string) (string, error) {
	cluster, err := GetCurrentCluster()
	if err != nil {
		return "", fmt.Errorf("could not get current cluster name: %v", err)
	}

	token, err := GetToken(user, namespace)
	if err != nil {
		return "", fmt.Errorf("could not get sa token: %v", err)
	}
	clusterObj := &clientcmdapi.NamedCluster{}
	fmt.Println(cluster)
	if err = json.Unmarshal([]byte(cluster), clusterObj); err != nil {
		return "", err
	}
	contextName := fmt.Sprintf("%s@%s", user, clusterObj.Name)
	config := clientcmdapi.Config{
		Kind:        "Config",
		APIVersion:  "v1",
		Preferences: clientcmdapi.Preferences{},
		Clusters: []clientcmdapi.NamedCluster{
			*clusterObj,
		},
		AuthInfos: []clientcmdapi.NamedAuthInfo{
			{
				Name: user,
				AuthInfo: clientcmdapi.AuthInfo{
					Token: token,
				},
			},
		},
		Contexts: []clientcmdapi.NamedContext{
			{
				Name: contextName,
				Context: clientcmdapi.Context{
					Cluster:    clusterObj.Name,
					AuthInfo:   user,
					Namespace:  namespace,
					Extensions: []clientcmdapi.NamedExtension{},
				},
			},
		},
		CurrentContext: contextName,
		Extensions:     []clientcmdapi.NamedExtension{},
	}

	y, err := yaml.Marshal(config)
	return string(y), err
}

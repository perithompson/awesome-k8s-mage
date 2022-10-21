package users

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api/v1"

	"github.com/magefile/mage/sh"
	"github.com/perithompson/awesome-k8s-mage/pkg/k8s/constants"
)

func GetCurrentContext() (string, error) {
	return sh.Output(constants.KubectlCmd, "config", "current-context", "get")
}

func GetCluster(name string) (string, error) {
	query := fmt.Sprintf("{.clusters[?(@.name=='%s')]}", name)
	return sh.Output(constants.KubectlCmd, "config", "view", constants.OutJsonPath(query), "--allow-missing-template-keys=true", "--raw")
}

type saSecret struct {
	Name string `json:"name"`
}

func GetSATokenSecret(name string) (string, error) {
	query := fmt.Sprintf("{.secrets[?(@.name contains %s-token-)]}", name)
	tokenName, err := sh.Output(constants.KubectlCmd, "get", "serviceaccount", name, constants.OutJsonPath(query))
	if err != nil {
		return "", err
	}
	secretName := &saSecret{}
	if err := json.Unmarshal([]byte(tokenName), secretName); err != nil {
		return "", err
	}

	return secretName.Name, nil
}

func GetToken(user string) (string, error) {
	secret, err := GetSATokenSecret(user)
	if err != nil {
		return "", err
	}
	token, err := sh.Output(constants.KubectlCmd, "get", "secret", secret, constants.OutJsonPath("{.data.token}"))
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
	currentContext, err := GetCurrentContext()
	if err != nil {
		return "", fmt.Errorf("could not get current context: %v", err)
	}
	currentCluster := strings.Split(currentContext, "@")[1]

	cluster, err := GetCluster(currentCluster)
	if err != nil {
		return "", fmt.Errorf("could not get current cluster name from current context %s: %v", currentContext, err)
	}

	token, err := GetToken(user)
	if err != nil {
		return "", fmt.Errorf("could not get sa token: %v", err)
	}

	clusterObj := &clientcmdapi.NamedCluster{}

	if err = json.Unmarshal([]byte(cluster), clusterObj); err != nil {
		return "", err
	}

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
				Name: fmt.Sprintf("%s@%s", user, currentCluster),
				Context: clientcmdapi.Context{
					Cluster:    clusterObj.Name,
					AuthInfo:   user,
					Namespace:  namespace,
					Extensions: []clientcmdapi.NamedExtension{},
				},
			},
		},
		CurrentContext: fmt.Sprintf("%s@%s", user, currentCluster),
		Extensions:     []clientcmdapi.NamedExtension{},
	}

	y, err := yaml.Marshal(config)
	return string(y), err
}

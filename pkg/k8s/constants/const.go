package constants

import (
	"fmt"
)

// commands
const KubectlCmd = "kubectl"

// outputs
const OutJson = "-ojson"
const OutYaml = "-oyaml"
const JsonPath = "-ojsonpath"
const OutName = "-oname"

const TempPrefix = "k8smage-"

func OutJsonPath(query string) string {
	return fmt.Sprintf("%s=%s", JsonPath, query)
}

func KubeconfigArg(path string) string {
	return fmt.Sprintf("--kubeconfig=%s", path)
}

package k8s

import "strings"

func IsDeployment(kind string) bool {
	lowerKind := strings.ToLower(strings.Trim(kind, " "))
	return lowerKind == "deployment" || lowerKind == "deployments" || lowerKind == "deploy"
}

func IsStatefulSet(kind string) bool {
	lowerKind := strings.ToLower(strings.Trim(kind, " "))
	return lowerKind == "statefulset" || lowerKind == "statefulsets" || lowerKind == "sts"
}

package utils

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	rayv1 "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
)

func GetCurrentNSOrDefault() string {
	ns, err := GetCurrentNS()
	if err != nil {
		return "default"
	}
	return ns
}

// IsClusterReady tells whether the cluster status in 'Ready' condition.
func IsRayServiceReady(rayServiceStatus *rayv1.RayServiceStatuses) bool {
	return rayServiceStatus != nil && rayServiceStatus.ServiceStatus == rayv1.Running
}

// EqualRayService
func EqualRayService(a, b *rayv1.RayService) bool {
	return a.Name == b.Name && a.Namespace == b.Namespace
}

// GetCurrentNS fetch namespace the current pod running in. reference to client-go (config *inClusterClientConfig) Namespace() (string, bool, error).
func GetCurrentNS() (string, error) {
	if ns := os.Getenv("POD_NAMESPACE"); ns != "" {
		return ns, nil
	}

	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns, nil
		}
	}
	return "", fmt.Errorf("can not get namespace where pods running in")
}

type TemplateName string

func RenderTemplate(templateName TemplateName, pattern string, data interface{}) (string, error) {
	t, err := template.New(string(templateName)).Parse(pattern)
	if err != nil {
		return "", fmt.Errorf("parse template %s but %w", templateName, err)
	}
	buffer := bytes.Buffer{}
	err = t.Execute(&buffer, data)
	if err != nil {
		return "", fmt.Errorf("execute template %s but %w", templateName, err)
	}
	result := buffer.Bytes()
	return string(result), nil
}

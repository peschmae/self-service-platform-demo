package k8s

import (
	"context"
	"fmt"
	"self-service-platform/internal/forms"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(name string, labels []forms.Label) error {
	// create a namespace
	client, err := getKubeClient()

	if err != nil {
		return err
	}

	// get existing namespaces

	namespaces, err := getNamespaces(client)
	if err != nil {
		return err
	}

	// check if namespace already exists
	for _, namespace := range namespaces {
		if namespace.Name == name {
			return fmt.Errorf("namespace %s already exists", name)
		}
	}

	// create namespace

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	ns.ObjectMeta.Labels = make(map[string]string)
	for _, label := range labels {
		ns.ObjectMeta.Labels[label.Key] = label.Value
	}

	_, err = client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("error creating namespace %s: %v", name, err)
	}

	return nil

}

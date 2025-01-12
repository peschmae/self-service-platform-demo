package k8s

import (
	"context"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateNamespace(name string, labels []string) error {
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
		l := strings.Split(label, "=")
		if len(l) != 2 {
			continue
		}
		ns.ObjectMeta.Labels[l[0]] = l[1]
	}

	_, err = client.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})

	if err != nil {
		return fmt.Errorf("error creating namespace %s: %v", name, err)
	}

	return nil

}

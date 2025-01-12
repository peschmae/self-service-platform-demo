package k8s

import (
	"context"

	"github.com/labstack/gommon/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

func getKubeClient() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()

	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", "/Users/peschmae/.kube/config")

		if err != nil {
			log.Error("Couldn't load configuration to connect to cluster!")
			return nil, err
		}
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Couldn't create k8s client")
		return nil, err
	}

	return clientset, nil

}

func getNamespaces(clientset *kubernetes.Clientset) ([]corev1.Namespace, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Error("Couldn't get namespaces")
		return nil, err
	}

	return namespaces.Items, nil
}

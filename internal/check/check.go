package check

import (
	_ "embed"
	"fmt"
	"self-service-platform/internal/k8s"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//go:embed script.sh
var checkScript string

func DeployCheckScript(namespace string, checkEndpoints []string) error {
	err := k8s.CreateConfigMap(namespace, "check-script", map[string]string{"check.sh": checkScript})
	if err != nil {
		return err
	}

	urls := ""
	for _, url := range checkEndpoints {
		urls += url + "\n"
	}

	err = k8s.CreateConfigMap(namespace, "url-list", map[string]string{"url-list": urls})
	if err != nil {
		return err
	}

	deploy := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "check-script",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "check-script"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "check-script"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    "check-script",
							Image:   "quay.io/curl/curl:8.11.1",
							Command: []string{"sh", "/opt/check.sh"},
							SecurityContext: &v1.SecurityContext{
								RunAsUser: int64Ptr(0),
							},
							Env: []v1.EnvVar{
								{
									Name:  "INTERVAL",
									Value: "10",
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "check-script",
									MountPath: "/opt/check.sh",
									SubPath:   "check.sh",
								},
								{
									Name:      "url-list",
									MountPath: "/opt/url-list",
									SubPath:   "url-list",
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: "check-script",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "check-script",
									},
									DefaultMode: int32Ptr(0755),
								},
							},
						},
						{
							Name: "url-list",
							VolumeSource: v1.VolumeSource{
								ConfigMap: &v1.ConfigMapVolumeSource{
									LocalObjectReference: v1.LocalObjectReference{
										Name: "url-list",
									},
									DefaultMode: int32Ptr(0420),
								},
							},
						},
					},
				},
			},
		},
	}

	err = k8s.CreateDeployment(deploy)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil

}

func int32Ptr(i int32) *int32 { return &i }

func int64Ptr(i int64) *int64 { return &i }

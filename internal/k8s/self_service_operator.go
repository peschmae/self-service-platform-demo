package k8s

import (
	"encoding/json"
	"self-service-platform/internal/forms"
	"strconv"
	"strings"

	k8smpetermannchv1beta1 "github.com/peschmae/self-service-operator-demo/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateSelfServiceNamespace(nsForm forms.NamespaceForm) error {

	operatorNamespace := &k8smpetermannchv1beta1.SelfServiceNamespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "k8s.mpetermann.ch/v1beta1",
			Kind:       "SelfServiceNamespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   nsForm.Name,
			Labels: map[string]string{},
		},
		Spec: k8smpetermannchv1beta1.SelfServiceNamespaceSpec{
			EgressConfigurations: []k8smpetermannchv1beta1.EgressConfigurationSpec{},
			AdditionalLabels:     map[string]string{},
			NetworkChecks:        []k8smpetermannchv1beta1.NetworkCheckConfigurationSpec{},
		},
	}

	for _, label := range nsForm.Labels {
		l := strings.Split(label, "=")
		if len(l) != 2 {
			continue
		}
		operatorNamespace.Spec.AdditionalLabels[l[0]] = l[1]
	}

	for _, egress := range nsForm.Egress {

		e := strings.Split(egress, ":")

		cidr := e[0]
		port, err := strconv.Atoi(e[1])
		if err != nil {
			return err
		}

		operatorNamespace.Spec.EgressConfigurations = append(operatorNamespace.Spec.EgressConfigurations, k8smpetermannchv1beta1.EgressConfigurationSpec{
			Cidr:     cidr,
			Port:     int32(port),
			Protocol: "TCP",
		})
	}

	if nsForm.Checks {
		operatorNamespace.Spec.NetworkChecksEnabled = true

		for _, check := range nsForm.CheckEndpoints {
			operatorNamespace.Spec.NetworkChecks = append(operatorNamespace.Spec.NetworkChecks, k8smpetermannchv1beta1.NetworkCheckConfigurationSpec{
				Url: check,
			})
		}
	}

	resourceString, err := json.Marshal(operatorNamespace)
	if err != nil {
		return err
	}

	err = ApplyUnstructured(nsForm.Name, string(resourceString))

	if err != nil {
		return err
	}

	return nil
}

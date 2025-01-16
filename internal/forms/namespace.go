package forms

import (
	"fmt"
	"strconv"
	"strings"

	k8smpetermannchv1beta1 "github.com/peschmae/self-service-operator-demo/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceForm struct {
	Name           string   `form:"name" validate:"required"`
	Environment    string   `form:"environment" validate:"required"`
	Labels         []Label  `form:"labels[]"`
	Egress         []string `form:"egress[]"`
	Checks         bool     `form:"enableChecks"`
	CheckEndpoints []string `form:"checks[]"`
}

type Label struct {
	Key   string `form:"key" validate:"required"`
	Value string `form:"value" validate:"required"`
}

func (nsForm *NamespaceForm) UnmarshalParam(param string) error {
	fmt.Println(param)
	return nil
}

func (nsForm *NamespaceForm) MapToSelfServiceNamespace() (*k8smpetermannchv1beta1.SelfServiceNamespace, error) {
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
		operatorNamespace.Spec.AdditionalLabels[label.Key] = label.Value
	}

	for _, egress := range nsForm.Egress {

		e := strings.Split(egress, ":")

		cidr := e[0]
		port, err := strconv.Atoi(e[1])
		if err != nil {
			return nil, err
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

	return operatorNamespace, nil
}

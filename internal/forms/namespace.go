package forms

import (
	k8smpetermannchv1beta1 "github.com/peschmae/self-service-operator-demo/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceForm struct {
	Name           string   `form:"name" json:"name" validate:"required"`
	Environment    string   `form:"environment" json:"environment" validate:"required"`
	Labels         []Label  `form:"labels[]" json:"labels[]"`
	Egress         []Egress `form:"egress[]" json:"egress[]"`
	Checks         bool     `form:"enableChecks" json:"enableChecks"`
	CheckEndpoints []string `form:"checks[]" json:"checks[]"`
}

type Label struct {
	Key   string `form:"key" json:"key" validate:"required"`
	Value string `form:"value" json:"value" validate:"required"`
}

type Egress struct {
	Cidr string `form:"cidr" json:"cidr" validate:"required"`
	Port int32  `form:"port" json:"port" validate:"required"`
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

		operatorNamespace.Spec.EgressConfigurations = append(operatorNamespace.Spec.EgressConfigurations, k8smpetermannchv1beta1.EgressConfigurationSpec{
			Cidr:     egress.Cidr,
			Port:     egress.Port,
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

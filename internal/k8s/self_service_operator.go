package k8s

import (
	"encoding/json"
	"self-service-platform/internal/forms"
)

func CreateSelfServiceNamespace(nsForm forms.NamespaceForm) error {

	operatorNamespace, err := nsForm.MapToSelfServiceNamespace()
	if err != nil {
		return err
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

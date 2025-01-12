package forms

type NamespaceForm struct {
	Name           string   `form:"name" validate:"required"`
	Environment    string   `form:"environment" validate:"required"`
	Labels         []string `form:"labels[]"`
	Egress         []string `form:"egress[]"`
	Checks         bool     `form:"enableChecks"`
	CheckEndpoints []string `form:"checks[]"`
}

package forms

type NamespaceForm struct {
	Name   string   `form:"name" validate:"required"`
	Labels []string `form:"labels[]"`
	Egress []string `form:"egress[]"`
}

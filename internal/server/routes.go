package server

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"self-service-platform/internal/check"
	"self-service-platform/internal/forms"
	"self-service-platform/internal/k8s"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRegistry struct {
	templates map[string]*template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		err := errors.New("Template not found -> " + name)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base.html", data)
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Validator = &CustomValidator{validator: validator.New()}

	e.Static("/static", "assets")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Instantiate a template registry with an array of template set
	// Ref: https://gist.github.com/rand99/808e6e9702c00ce64803d94abff65678
	templates := make(map[string]*template.Template)
	templates["index.html"] = template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	templates["confirmation.html"] = template.Must(template.ParseFiles("templates/confirmation.html", "templates/base.html"))
	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.GET("/", s.IndexHandler)
	e.POST("/create", s.FormHandler)

	return e
}

func (s *Server) IndexHandler(c echo.Context) error {

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

func (s *Server) FormHandler(c echo.Context) error {
	nsForm := new(forms.NamespaceForm)
	if err := c.Bind(nsForm); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	if err := c.Validate(nsForm); err != nil {
		return err
	}

	nsForm.Labels = append(nsForm.Labels, "k8s.mpetermann.ch/environment="+nsForm.Environment)

	err := k8s.CreateNamespace(nsForm.Name, nsForm.Labels)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	err = k8s.CreateDefaultNetpols(nsForm.Name)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	err = k8s.CreateEgressNetpol(nsForm.Name, nsForm.Egress)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	if nsForm.Checks {
		err = check.DeployCheckScript(nsForm.Name, nsForm.CheckEndpoints)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	return c.Render(http.StatusOK, "confirmation.html", map[string]interface{}{"Namespace": nsForm.Name, "Checks": nsForm.Checks})
}

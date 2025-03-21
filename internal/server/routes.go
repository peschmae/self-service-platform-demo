package server

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"self-service-platform/internal/check"
	"self-service-platform/internal/forms"
	"self-service-platform/internal/git"
	"self-service-platform/internal/k8s"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/itchyny/json2yaml"
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
	templates["namespace-form.html"] = template.Must(template.ParseFiles("templates/namespace-form.html", "templates/base.html"))
	templates["confirmation.html"] = template.Must(template.ParseFiles("templates/confirmation.html", "templates/base.html"))
	e.Renderer = &TemplateRegistry{
		templates: templates,
	}

	e.GET("/", s.IndexHandler)
	e.GET("/create", s.NamespaceForm)
	e.POST("/create", s.RegularFormHandler)

	e.GET("/create-operator", s.NamespaceForm)
	e.POST("/create-operator", s.OperatorFormHandler)

	e.GET("/create-operator-gitops", s.NamespaceForm)
	e.POST("/create-operator-gitops", s.GitopsFormHandler)

	return e
}

func (s *Server) IndexHandler(c echo.Context) error {

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}

func (s *Server) getFormContext() map[string]any {

	ctx := make(map[string]interface{})
	ctx["Environments"] = map[string]string{"dev": "Development", "test": "Test", "staging": "Staging", "prod": "Production"}
	ctx["RequiredLabels"] = map[string]string{"k8wms.swisscom.com/customer": "Customer/Project name", "k8wms.swisscom.com/psp": "PSP"}

	return ctx
}

func (s *Server) NamespaceForm(c echo.Context) error {

	ctx := s.getFormContext()
	ctx["FormAction"] = c.Path()

	return c.Render(http.StatusOK, "namespace-form.html", ctx)
}

func (s *Server) mapForm(c echo.Context) (*forms.NamespaceForm, error) {
	nsForm := new(forms.NamespaceForm)
	if err := c.Bind(nsForm); err != nil {
		return nil, err
	}
	if err := c.Validate(nsForm); err != nil {
		return nil, err
	}

	nsForm.Labels = append(nsForm.Labels, "k8s.mpetermann.ch/environment="+nsForm.Environment)

	return nsForm, nil
}

func (s *Server) RegularFormHandler(c echo.Context) error {
	nsForm, err := s.mapForm(c)

	err = k8s.CreateNamespace(nsForm.Name, nsForm.Labels)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	defaultResourcePath := os.Getenv("DEFAULT_RESOURCES")

	if defaultResourcePath != "" {
		var defaultResources []string

		filepath.WalkDir(defaultResourcePath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				c.Logger().Info("Skipping directory, nested directory not supported: ", d.Name())
				return nil
			}
			if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				defaultResources = append(defaultResources, string(content))
			}
			return nil
		})

		for _, res := range defaultResources {
			err = k8s.ApplyUnstructured(nsForm.Name, res)
			if err != nil {
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}
		}
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

func (s *Server) OperatorFormHandler(c echo.Context) error {
	nsForm, err := s.mapForm(c)

	err = k8s.CreateSelfServiceNamespace(*nsForm)
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.Render(http.StatusOK, "confirmation.html", map[string]interface{}{"Namespace": nsForm.Name, "Checks": nsForm.Checks})
}

func (s *Server) GitopsFormHandler(c echo.Context) error {
	nsForm, err := s.mapForm(c)

	operatorNamespace, err := nsForm.MapToSelfServiceNamespace()
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	namspaceJson, err := json.Marshal(operatorNamespace)
	if err != nil {
		return err
	}

	var output strings.Builder
	if err := json2yaml.Convert(&output, strings.NewReader(string(namspaceJson))); err != nil {
		return err
	}

	gitOps := git.GitOps{
		RepoPath: os.Getenv("GITOPS_REPO"),
		RepoURL:  os.Getenv("GITOPS_REPO_URL"),
	}

	if !filepath.IsAbs(gitOps.RepoPath) {
		c.Logger().Error("RepoPath must be an absolute path")
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	if _, err = os.Stat(gitOps.RepoPath); os.IsNotExist(err) {
		err = gitOps.Clone()
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	} else {
		gitOps.Pull()
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	}

	os.WriteFile(filepath.Join(gitOps.RepoPath, nsForm.Name+".yaml"), []byte(output.String()), 0644)

	gitOps.Commit("Add " + nsForm.Name + " selfServiceNamespace")
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	err = gitOps.Push()
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.Render(http.StatusOK, "confirmation.html", map[string]interface{}{"Namespace": nsForm.Name, "Checks": nsForm.Checks})
}

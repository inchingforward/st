package main

import (
	"io"
	"net/http"
	"path"

	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"

	"github.com/labstack/echo"
)

type Renderer struct {
	TemplateDir   string
	Reload        bool
	TemplateCache map[string]*pongo2.Template
}

func renderTemplate(c echo.Context, templateName string) error {
	return c.Render(http.StatusOK, templateName, pongo2.Context{})
}

// GetTemplate returns a template, loading it every time if reload is true.
func (r *Renderer) GetTemplate(name string, reload bool) *pongo2.Template {
	filename := path.Join(r.TemplateDir, name)

	if r.Reload {
		return pongo2.Must(pongo2.FromFile(filename))
	}

	return pongo2.Must(pongo2.FromCache(filename))
}

// Render renders a pongo2 template.
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	template := r.GetTemplate(name, debug)

	pctx := data.(pongo2.Context)

	pctx["csrf"] = c.Get("csrf")
	pctx["debug"] = debug

	return template.ExecuteWriter(pctx, w)
}

func setRenderer(e *echo.Echo, debug bool) {
	e.Renderer = &Renderer{TemplateDir: "templates", Reload: debug, TemplateCache: make(map[string]*pongo2.Template)}
}

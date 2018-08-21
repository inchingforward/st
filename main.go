package main

import (
    "net/http"

    "flag"
    "fmt"
    "io"
    "path"

    "github.com/flosch/pongo2"
    _ "github.com/flosch/pongo2-addons"

    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
)

var (
    debug = false
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

func renderHome(c echo.Context) error {
    return renderTemplate(c, "home.html")
}

func renderAbout(c echo.Context) error {
    return renderTemplate(c, "about.html")
}

func renderCreateStory(c echo.Context) error {
    // FIXME: Initial page to create a story...allows user to set title, visibility, etc.
    return renderTemplate(c, "story_create.html")
}

func createStory(c echo.Context) error {
    // FIXME: Create Story record, redirect to edit page using story uuid. 
    return renderTemplate(c, "story_create.html")
}

func renderEditStory(c echo.Context) error {
    // FIXME: Look up story by uuid, return story edit template with story details.
    // FIXME: 404 if the uuid is not found. 
    return renderTemplate(c, "story_edit.html")
}

func renderPublishStory(c echo.Context) error {
    // FIXME: Summarize story information in a template, allowing the author to change details.
    return renderTemplate(c, "story_publish.html")
}

func publishStory(c echo.Context) error {
    // FIXME: Finalize the story by setting the publish date.
    return renderTemplate(c, "story_publish.html")
}

func renderStoryList(c echo.Context) error {
    // FIXME: Pull up list of publicly-visible stories.
    return renderTemplate(c, "story_list.html")
}

func renderStory(c echo.Context) error {
    // FIXME: Pull up story by uuid.
    return renderTemplate(c, "story.html")
}

func main() {
    flag.BoolVar(&debug, "debug", false, "true to enable debug")
    flag.Parse()

    fmt.Printf("debug: %v\n", debug)

    e := echo.New()
    e.Pre(middleware.RemoveTrailingSlash())

    e.Renderer = &Renderer{TemplateDir: "templates", Reload: debug, TemplateCache: make(map[string]*pongo2.Template)}

    e.GET("/", renderHome)
    e.GET("/about", renderAbout)
    e.GET("/stories/create", renderCreateStory)
    e.POST("/stories/:uuid/create", createStory)
    e.GET("/stories/:uuid", renderStory)
    e.GET("/stories/:uuid/edit", renderEditStory)
    e.GET("/stories/:uuid/publish", renderPublishStory)
    e.POST("/stories/:uuid/publish", publishStory)
    e.GET("/stories", renderStoryList)

    e.Logger.Fatal(e.Start(":8009"))
}

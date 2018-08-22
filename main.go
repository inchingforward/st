package main

import (
	"flag"
	"fmt"

	"github.com/flosch/pongo2"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	debug = false
)

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	fmt.Printf("debug: %v\n", debug)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
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

	e.Logger.Fatal(e.Start(":8011"))
}

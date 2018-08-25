package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/flosch/pongo2"

	"github.com/jmoiron/sqlx"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	debug = false
)

func init() {
	x, err := sqlx.Connect("postgres", "user=storytellers dbname=storytellers sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	x.SetMaxOpenConns(2)
	x.SetMaxIdleConns(2)

	fmt.Println("connected to db")
	db = x
}

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	fmt.Printf("debug: %v\n", debug)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Renderer = &Renderer{TemplateDir: "templates", Reload: debug, TemplateCache: make(map[string]*pongo2.Template)}

	e.GET("/", getHome)
	e.GET("/about", getAbout)
	e.GET("/stories/create", getCreateStory)
	e.POST("/stories/create", createStory)
	e.GET("/stories/:uuid", getStory)
	e.GET("/stories/:uuid/edit", getEditStory)
	e.GET("/stories/:uuid/publish", getPublishStory)
	e.POST("/stories/publish", publishStory)
	e.GET("/stories", getStoryList)

	e.Logger.Fatal(e.Start(":8011"))
}

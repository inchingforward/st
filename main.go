package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/flosch/pongo2"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	debug = false
)

func init() {
	initDB()
}

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	fmt.Printf("debug: %v\n", debug)

	rand.Seed(time.Now().UnixNano())

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Renderer = &Renderer{TemplateDir: "templates", Reload: debug, TemplateCache: make(map[string]*pongo2.Template)}

	addHandlers(e)

	e.Logger.Fatal(e.Start(":8011"))
}

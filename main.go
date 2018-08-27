package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

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

	setRenderer(e)
	addHandlers(e)

	e.Logger.Fatal(e.Start(":8011"))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

var (
	debug = false
)

func init() {
	initDB()
}

func main() {
	var sessionSecret string

	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.StringVar(&sessionSecret, "secret", "", "session secret")
	flag.Parse()

	fmt.Printf("debug: %v\n", debug)

	if sessionSecret == "" {
		log.Fatal("usage: st secret [-debug]")
	}

	fmt.Println(sessionSecret)

	rand.Seed(time.Now().UnixNano())

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))
	e.Static("/static", "static")

	setRenderer(e, debug)
	addHandlers(e)

	e.Logger.Fatal(e.Start(":8011"))
}

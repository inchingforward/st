package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"

	"gopkg.in/olahol/melody.v1"
)

func addHandlers(e *echo.Echo) {
	m := melody.New()

	e.GET("/", getHome)
	e.GET("/about", getAbout)
	e.GET("/stories/create", getCreateStory)
	e.POST("/stories/create", createStory)
	e.GET("/stories/:uuid/join", getJoinStory)
	e.GET("/stories/join", getJoinStory)
	e.POST("/stories/join", joinStory)
	e.GET("/stories/:uuid", getStory)
	e.GET("/stories/:uuid/edit", getEditStory)
	e.GET("/stories/:uuid/publish", getPublishStory)
	e.POST("/stories/publish", publishStory)
	e.GET("/stories", getStoryList)
	e.GET("/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response().Writer, c.Request())
		return nil
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})
}

func getHome(c echo.Context) error {
	return renderTemplate(c, "home.html")
}

func getAbout(c echo.Context) error {
	return renderTemplate(c, "about.html")
}

func getCreateStory(c echo.Context) error {
	return renderTemplate(c, "story_create.html")
}

func createStory(c echo.Context) error {
	story := new(Story)
	if err := c.Bind(story); err != nil {
		return c.Render(http.StatusBadRequest, "story_create.html", pongo2.Context{
			"Error": "Invalid fields",
		})
	}

	if story.Title == "" {
		return c.Render(http.StatusBadRequest, "story_create.html", pongo2.Context{
			"Error": "Title is required",
		})
	}

	author := c.FormValue("author")
	if author == "" {
		return c.Render(http.StatusBadRequest, "story_create.html", pongo2.Context{
			"Error": "Your name is required",
		})
	}

	sess, _ := session.Get("session", c)
	sess.Values["Author"] = author
	sess.Values["Creator"] = true
	sess.Save(c.Request(), c.Response())

	// At this point we have a valid story.
	story.StartedAt = time.Now()
	story.UUID = generateStoryUUID()
	story.Private = false

	err := insertStory(story)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "story_create.html", pongo2.Context{
			"Error": err.Error(),
		})
	}

	editURL := fmt.Sprintf("/stories/%s/edit", story.UUID)

	return c.Redirect(http.StatusSeeOther, editURL)
}

func getJoinStory(c echo.Context) error {
	uuid := c.Param("uuid")
	return c.Render(http.StatusOK, "story_join.html", pongo2.Context{
		"UUID": uuid,
	})
}

func joinStory(c echo.Context) error {
	uuid := c.FormValue("uuid")
	authorName := c.FormValue("author")

	if uuid == "" {
		return c.Render(http.StatusBadRequest, "story_join.html", pongo2.Context{
			"ErrorTitle": "Missing Story Code",
			"Error":      "Please enter a valid Story Code",
			"Author":     authorName,
		})
	}

	if authorName == "" {
		return c.Render(http.StatusBadRequest, "story_join.html", pongo2.Context{
			"ErrorTitle": "Missing Name",
			"Error":      "Please enter your name",
			"UUID":       uuid,
		})
	}

	sess, _ := session.Get("session", c)
	sess.Values["Author"] = authorName
	sess.Values["Creator"] = false
	sess.Save(c.Request(), c.Response())

	editURL := fmt.Sprintf("/stories/%s/edit", uuid)

	return c.Redirect(http.StatusSeeOther, editURL)
}

func getEditStory(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return c.Render(http.StatusBadRequest, "story_edit.html", pongo2.Context{
			"ErrorTitle": "Invalid Story ID",
			"Error":      "Please use a valid Story ID.",
		})
	}

	story, err := selectEditableStory(uuid)
	if err != nil {
		return c.Render(http.StatusBadRequest, "story_edit.html", pongo2.Context{
			"ErrorTitle": "Not Found",
			"Error":      "The Story ID was not found.",
		})
	}

	sess, _ := session.Get("session", c)

	return c.Render(http.StatusOK, "story_edit.html", pongo2.Context{
		"Story":   story,
		"Session": sess,
	})
}

func getPublishStory(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return c.Render(http.StatusBadRequest, "story_edit.html", pongo2.Context{
			"ErrorTitle": "Invalid Story ID",
			"Error":      "Please use a valid Story ID.",
		})
	}

	story, err := selectEditableStory(uuid)
	if err != nil {
		return c.Render(http.StatusBadRequest, "story_edit.html", pongo2.Context{
			"ErrorTitle": "Not Found",
			"Error":      "The Story ID was not found.",
		})
	}

	return c.Render(http.StatusOK, "story_publish.html", pongo2.Context{
		"Story": story,
	})
}

func publishStory(c echo.Context) error {
	uuid := c.FormValue("uuid")
	if uuid == "" {
		return c.Render(http.StatusBadRequest, "story_publish.html", pongo2.Context{
			"Error": "Invalid fields",
		})
	}

	story, err := selectEditableStory(uuid)
	if err != nil {
		return c.Render(http.StatusNotFound, "error.html", pongo2.Context{
			"ErrorTitle": "Invalid Story",
			"Error":      err.Error(),
		})
	}

	updatedTitle := c.FormValue("title")
	if updatedTitle == "" {
		return c.Render(http.StatusBadRequest, "story_publish.html", pongo2.Context{
			"Error": "Invalid fields",
		})
	}

	// FIXME: combine all related authors.

	updatedPrivate := c.FormValue("private") == "on"

	story.Title = updatedTitle
	story.Private = updatedPrivate

	updatePublishStory(&story)

	storyURL := fmt.Sprintf("/stories/%s", uuid)

	return c.Redirect(http.StatusSeeOther, storyURL)
}

func getStoryList(c echo.Context) error {
	stories, err := selectPublishedStories()
	if err != nil {
		return c.Render(http.StatusInternalServerError, "story_list.html", pongo2.Context{
			"ErrorTitle": "Invalid fields",
			"Error":      err.Error(),
		})
	}

	return c.Render(http.StatusOK, "story_list.html", pongo2.Context{
		"Stories": stories,
	})
}

func getStory(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		return c.Render(http.StatusNotFound, "error.html", pongo2.Context{
			"ErrorTitle": "Invalid Story",
			"Error":      "Missing story id",
		})
	}

	story, err := selectPublishedStory(uuid)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", pongo2.Context{
			"ErrorTitle": "Error",
			"Error":      err.Error(),
		})
	}

	parts, err := selectPublishedStoryParts(story.ID)
	if err != nil {
		return c.Render(http.StatusInternalServerError, "error.html", pongo2.Context{
			"ErrorTitle": "Error",
			"Error":      err.Error(),
		})
	}

	return c.Render(http.StatusOK, "story.html", pongo2.Context{
		"Story": story,
		"Parts": parts,
	})
}

func generateStoryUUID() string {
	length := 6
	bytes := make([]byte, length)
	lowerA := 97
	lowerZ := 122

	for i := 0; i < length; i++ {
		bytes[i] = byte(lowerA + rand.Intn(lowerZ-lowerA))
	}

	// FIXME: Add db check.

	return string(bytes)
}

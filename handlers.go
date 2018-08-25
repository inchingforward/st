package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
)

func getHome(c echo.Context) error {
	return renderTemplate(c, "home.html")
}

func getAbout(c echo.Context) error {
	return renderTemplate(c, "about.html")
}

func getCreateStory(c echo.Context) error {
	// FIXME: Initial page to create a story...allows user to set title, visibility, etc.
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

	if story.Authors == "" {
		return c.Render(http.StatusBadRequest, "story_create.html", pongo2.Context{
			"Error": "Authors is required",
		})
	}

	// At this point we have a valid story.
	story.StartedAt = time.Now()

	uuid, _ := uuid.NewV4()
	story.UUID = uuid.String()

	err := insertStory(story)

	if err != nil {
		return c.Render(http.StatusInternalServerError, "story_create.html", pongo2.Context{
			"Error": err.Error(),
		})
	}

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

	return c.Render(http.StatusOK, "story_edit.html", pongo2.Context{
		"Story": story,
	})
}

func getPublishStory(c echo.Context) error {
	// FIXME: Summarize story information in a template, allowing the author to change details.
	// FIXME: Don't allow an already published story to be accessed?
	return renderTemplate(c, "story_publish.html")
}

func publishStory(c echo.Context) error {
	// FIXME: Don't let an already published story be published.
	// FIXME: Redirect to the actual story uuid.
	return c.Redirect(http.StatusSeeOther, "/stories/foo")
}

func getStoryList(c echo.Context) error {
	// FIXME: Pull up list of publicly-visible stories.
	return renderTemplate(c, "story_list.html")
}

func getStory(c echo.Context) error {
	// FIXME: Pull up story by uuid regardless of visibility.
	// FIXME: Don't pull up the story if it hasn't been published.
	return renderTemplate(c, "story.html")
}

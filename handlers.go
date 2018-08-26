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

func getJoinStory(c echo.Context) error {
	return renderTemplate(c, "story_join.html")
}

func joinStory(c echo.Context) error {
	uuid := c.FormValue("uuid")
	authorName := c.FormValue("author_name")

	if uuid == "" {
		return c.Render(http.StatusBadRequest, "story_join.html", pongo2.Context{
			"ErrorTitle": "Missing Story Code",
			"Error":      "Please enter a valid Story Code",
			"AuthorName": authorName,
		})
	}

	if authorName == "" {
		return c.Render(http.StatusBadRequest, "story_join.html", pongo2.Context{
			"ErrorTitle": "Missing Name",
			"Error":      "Please enter your name",
			"UUID":       uuid,
		})
	}

	// FIXME:  add author_name to session

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

	updatedAuthors := c.FormValue("authors")
	if updatedAuthors == "" {
		return c.Render(http.StatusBadRequest, "story_publish.html", pongo2.Context{
			"Error": "Authors is required",
		})
	}

	updatedPrivate := c.FormValue("private") == "on"

	story.Title = updatedTitle
	story.Authors = updatedAuthors
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

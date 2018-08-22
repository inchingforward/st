package main

import (
	"net/http"

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
	// FIXME: Create Story record, redirect to edit page using story uuid.
	return c.Redirect(http.StatusSeeOther, "/stories/foo/edit")
}

func getEditStory(c echo.Context) error {
	// FIXME: Look up story by uuid, return story edit template with story details.
	// FIXME: Redirect to the view story page if the story has already been published.
	// FIXME: 404 if the uuid is not found.
	return renderTemplate(c, "story_edit.html")
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

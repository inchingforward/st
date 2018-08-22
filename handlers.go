package main

import "github.com/labstack/echo"

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

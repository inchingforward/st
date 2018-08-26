package main

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	db *sqlx.DB
)

type Story struct {
	ID          uint64      `db:"id" form:"story_id"`
	Title       string      `db:"title" form:"title"`
	UUID        string      `db:"uuid" form:"uuid"`
	Authors     string      `db:"authors" form:"authors"`
	Private     bool        `db:"private" form:"private"`
	StartedAt   time.Time   `db:"started_at"`
	Published   bool        `db:"published"`
	PublishedAt pq.NullTime `db:"published_at"`
}

type StoryPart struct {
	ID        uint64    `db:"id"`
	StoryID   uint64    `db:"story_id" json:"story_id"`
	PartNum   int       `db:"part_num" json:"part_num"`
	PartText  string    `db:"part_text" json:"part_text"`
	WrittenBy string    `db:"written_by" json:"written_by"`
	WrittenAt time.Time `db:"written_at"`
}

func insertStory(story *Story) error {
	_, err := db.Exec("insert into story values (default, $1, $2, $3, $4, $5, false, null)", story.Title, story.UUID, story.Authors, story.Private, story.StartedAt)

	return err
}

func selectEditableStory(uuid string) (Story, error) {
	var story Story

	err := db.Get(&story, "select * from story where uuid = $1 and published = false", uuid)

	return story, err
}

func updatePublishStory(story *Story) error {
	_, err := db.Exec(`
		update story 
		set    title = $1, 
				authors = $2, 
				private = $3, 
				published = true, 
				published_at = now() 
		where   id = $4`, story.Title, story.Authors, story.Private, story.ID)

	return err
}

func selectPublishedStories() ([]Story, error) {
	stories := []Story{}

	err := db.Select(&stories, "select * from story where published = true order by published_at")

	return stories, err
}

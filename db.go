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
	ID          uint64       `db:"id"`
	Title       string       `db:"title"`
	UUID        string       `db:"uuid"`
	Authors     string       `db:"authors"`
	Private     bool         `db:"private"`
	StartedAt   time.Time    `db:"started_at"`
	Published   bool         `db:"published"`
	PublishedAt pq.NullTime `db:"published_at"`
}

type StoryPart struct {
	ID        uint64    `db:"id"`
	StoryID   uint64    `db:"story_id"`
	PartNum   int       `db:"part_num"`
	PartText  string    `db:"part_text"`
	WrittenBy string    `db:"written_by"`
	WrittenAt time.Time `db:"written_at"`
}

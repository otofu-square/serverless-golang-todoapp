package main

import (
	"time"

	"github.com/satori/go.uuid"
)

type Todo struct {
	ID        string `dynamo:"ID"`
	Title     string
	Completed bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewTodo(title string, completed bool) *Todo {
	id := uuid.NewV4().String()
	now := time.Now()
	return &Todo{
		ID:        id,
		Title:     title,
		Completed: completed,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

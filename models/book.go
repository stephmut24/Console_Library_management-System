package models

import "time"

type Book struct {
	ID int
	Title string
	Author string
	Status string // "Available" or "Borrowed"

	// Reservation fields
	ReservedBy int
	ReservedUntil time.Time
}

func NewBook(id int, title, author string) *Book {
	return &Book{
		ID: id,
		Title: title,
		Author: author,
		Status: "Available",
	}
}
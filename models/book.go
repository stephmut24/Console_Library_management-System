package models

type Book struct {
	ID int 
	Title string
	Author string
	Status string // "Available" or "Borrowed"
}

func NewBook(id int, title, author string) *Book {
	return &Book{
		ID: id,
		Title: title,
		Author: author,
		Status: "Available",
	}
}
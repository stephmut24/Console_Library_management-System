package models

type Member struct {
	ID int
	Name string
	BorrowedBooks []Book
}

func NewMember(id int, name string) *Member {
	return &Member{
		ID: id,
		Name: name,
		BorrowedBooks: make([]Book,0),
		}
	}
	
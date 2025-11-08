package services

import( 
	"errors"
	
	"library_management/models"
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int) error
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
	AddMember(member models.Member)
	FindMember(memberID int) (*models.Member, error)
	FIndBook(bookID int) (*models.Book, error)
}

type Library struct {
	books map[int]models.Book
	members map[int]models.Member
	nextBookID int
	nextMemberID int
}

//NewLibrary
func NewLibrary() *Library {
	return &Library{
		books: make(map[int]models.Book),
		members: make(map[int]models.Member),
		nextBookID: 1,
		nextMemberID: 1,
	}
}

// AddBook 
func (l *Library) AddBook(book models.Book) {
	book.ID = l.nextBookID
	l.books[l.nextBookID] = book
	l.nextBookID++
}

//RemoveBook
func (l *Library) RemoveBook(bookID int) error {
	book, exists := l.books[bookID]
	if !exists {
		return errors.New("livre non trouve")
	}
	if book.Status == "Borrowed" {
		return errors.New("impossible de supprier un livre emprunte")
	}

	delete(l.books, bookID)
	return nil
}

//BorrowBook
func (l *Library) BorrowBook(bookID int, memberID int) error{
	//check if the book exists

	book, exists := l.books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	// check if the memeber exists

	member, exists := l.members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	// check if the book is available
	if book.Status != "Available" {
		return errors.New("book borrowed")
	}
	// Update the book status
	book.Status = "Borrowed"
	l.books[bookID] = book

	//Add the book to the member's borrowed books
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] =member
	
	return nil
}

//ReturnBook function 
func (l *Library) ReturnBook(bookID int, memberID int) error {
	//check if the book exists
	book, exists := l.books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	// check if the member exists
	member, exists := l.members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	//Check if the book is borrowed
	if book.Status != "Borrowed" {
		return errors.New("this book is not borrowed")
	}

	//Remove the book from the member's borrowed books
	for i, borrowedBook := range member.BorrowedBooks {
		if borrowedBook.ID == bookID {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			break
		}
	}


	// Update book's status
	book.Status = "Available"
	l.books[bookID] = book
	l.members[memberID] = member

	return nil
}

//ListAvailableBooks
func (l *Library) ListAvailableBooks() []models.Book {
	var availableBooks []models.Book
	for _, book := range l.books {
		if book.Status == "Available" {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks

}

//ListBorrowedBooks
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	member, exists := l.members[memberID]
	if !exists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

//Addmember
func (l *Library) AddMember(member models.Member) {
	member.ID = l.nextMemberID
	l.members[l.nextMemberID] = member
	l.nextMemberID++
}

//Find member
func (l *Library) FindMember(memberID int) (*models.Member, error) {
	member, exists := l.members[memberID]
	if !exists {
		return nil, errors.New("member not found")
	}
	return &member, nil
}

//FindBook 
func (l *Library) FIndBook(bookID int) (*models.Book, error) {
	book, exists := l.books[bookID]
	if !exists {
		return nil, errors.New("book not found")
	}
	return &book, nil
}
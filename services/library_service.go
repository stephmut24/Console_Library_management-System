package services

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"library_management/concurrency"
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
	ReserveBook(bookID int, memberID int) error
}

type Library struct {
	books             map[int]models.Book
	members           map[int]models.Member
	nextBookID        int
	nextMemberID      int
	mu                sync.Mutex
	reservationCenter *concurrency.ReservationCenter
}

// NewLibrary
func NewLibrary() *Library {
	l := &Library{
		books:        make(map[int]models.Book),
		members:      make(map[int]models.Member),
		nextBookID:   1,
		nextMemberID: 1,
	}

	// create reservation center and workers which will call l.handleReservation
	rc := concurrency.NewReservationCenter(5, l.handleReservation)
	rc.Start()
	l.reservationCenter = rc

	// seed randomness used in simulated async borrowing
	rand.Seed(time.Now().UnixNano())

	return l
}

// AddBook
func (l *Library) AddBook(book models.Book) {
	l.mu.Lock()
	defer l.mu.Unlock()

	book.ID = l.nextBookID
	l.books[l.nextBookID] = book
	l.nextBookID++
}

// RemoveBook
func (l *Library) RemoveBook(bookID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.books[bookID]
	if !exists {
		return errors.New("book not found")
	}
	if book.Status == "Borrowed" {
		return errors.New("impossible to delete a Borrowed book")
	}

	delete(l.books, bookID)
	return nil
}

// BorrowBook
func (l *Library) BorrowBook(bookID int, memberID int) error {
	//check if the book exists

	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.books[bookID]
	if !exists {
		return errors.New("book not found")
	}

	// check if the memeber exists

	member, exists := l.members[memberID]
	if !exists {
		return errors.New("member not found")
	}

	// check if the book is available, or reserved by this member
	if book.Status != "Available" {
		return errors.New("book borrowed")
	}

	// if the book is reserved, make sure the requester is the reserver
	if book.ReservedBy != 0 && book.ReservedBy != memberID {
		return errors.New("book reserved by another member")
	}
	// Update the book status and clear reservation
	book.Status = "Borrowed"
	book.ReservedBy = 0
	book.ReservedUntil = time.Time{}
	l.books[bookID] = book

	//Add the book to the member's borrowed books
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.members[memberID] = member

	return nil
}

// ReserveBook enqueues a reservation request to the reservation worker
func (l *Library) ReserveBook(bookID int, memberID int) error {
	if l.reservationCenter == nil {
		return errors.New("reservation center not initialized")
	}

	resp := make(chan error)
	req := concurrency.ReservationRequest{BookID: bookID, MemberID: memberID, RespChan: resp}
	l.reservationCenter.Enqueue(req)

	// wait for worker to process reservation and return result
	err := <-resp
	return err
}

// handleReservation is the actual processing logic executed by worker goroutines
func (l *Library) handleReservation(req concurrency.ReservationRequest) {
	// lock to safely read/update maps
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.books[req.BookID]
	if !exists {
		req.RespChan <- errors.New("book not found")
		return
	}

	_, memberExists := l.members[req.MemberID]
	if !memberExists {
		req.RespChan <- errors.New("member not found")
		return
	}

	// If the book is borrowed or already reserved, reject
	if book.Status != "Available" {
		req.RespChan <- errors.New("book is already borrowed")
		return
	}
	if book.ReservedBy != 0 {
		req.RespChan <- errors.New("book already reserved")
		return
	}

	// Reserve it
	book.ReservedBy = req.MemberID
	book.ReservedUntil = time.Now().Add(5 * time.Second)
	l.books[req.BookID] = book

	// notify success
	req.RespChan <- nil

	// start goroutine that will auto-cancel reservation after 5s unless borrowed
	go func(bookID int, memberID int, expiry time.Time) {
		// Simulate potential asynchronous borrowing by the member (random delay)
		// We'll sometimes attempt to auto-borrow to demonstrate async processing.
		simulatedDelay := time.Duration(rand.Intn(8000)) * time.Millisecond // 0-8s

		// If simulated delay is small, attempt to auto-borrow within that time
		time.Sleep(simulatedDelay)

		l.mu.Lock()
		curBook := l.books[bookID]
		// If the book is no longer reserved by this member or has been borrowed, nothing to do
		if curBook.ReservedBy != memberID || curBook.Status != "Available" {
			l.mu.Unlock()
			return
		}

		// If borrow happened before expiry, nothing to do. If we hit expiry before simulated borrow, auto-cancel.
		if time.Now().After(expiry) {
			// cancel the reservation
			curBook.ReservedBy = 0
			curBook.ReservedUntil = time.Time{}
			l.books[bookID] = curBook
			l.mu.Unlock()
			fmt.Printf("Reservation for book %d by member %d expired and was cancelled\n", bookID, memberID)
			return
		}

		// Simulate automatic borrow (this demonstrates asynchronous processing)
		curBook.Status = "Borrowed"
		curBook.ReservedBy = 0
		curBook.ReservedUntil = time.Time{}
		l.books[bookID] = curBook

		// add book to member
		member := l.members[memberID]
		member.BorrowedBooks = append(member.BorrowedBooks, curBook)
		l.members[memberID] = member
		l.mu.Unlock()

		fmt.Printf("Auto-borrow: book %d successfully borrowed by member %d\n", bookID, memberID)
	}(req.BookID, req.MemberID, book.ReservedUntil)
}

// ReturnBook function
func (l *Library) ReturnBook(bookID int, memberID int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

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

// ListAvailableBooks
func (l *Library) ListAvailableBooks() []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	var availableBooks []models.Book
	for _, book := range l.books {
		if book.Status == "Available" && book.ReservedBy == 0 {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks

}

// ListBorrowedBooks
func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, exists := l.members[memberID]
	if !exists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}

// Addmember
func (l *Library) AddMember(member models.Member) {
	l.mu.Lock()
	defer l.mu.Unlock()

	member.ID = l.nextMemberID
	l.members[l.nextMemberID] = member
	l.nextMemberID++
}

// Find member
func (l *Library) FindMember(memberID int) (*models.Member, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	member, exists := l.members[memberID]
	if !exists {
		return nil, errors.New("member not found")
	}
	return &member, nil
}

// FindBook
func (l *Library) FIndBook(bookID int) (*models.Book, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	book, exists := l.books[bookID]
	if !exists {
		return nil, errors.New("book not found")
	}
	return &book, nil
}

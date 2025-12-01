package controllers

import (
	"bufio"
	"fmt"
	"library_management/models"
	"library_management/services"
	"os"
	"strconv"
	"strings"
)

type LibraryController struct {
	library services.LibraryManager
	scanner *bufio.Scanner
}

func NewLibraryController(library services.LibraryManager) *LibraryController {
	return &LibraryController{
		library: library,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// run the console interface
func (lc *LibraryController) reserveBook() {
	fmt.Println("\n---- Reserve a book ----")
	bookIDStr := lc.getInput("Book ID: ")
	memberIDStr := lc.getInput("Member ID: ")

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		fmt.Println("Error: Invalid Book ID")
		return
	}

	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		fmt.Println("Error: Invalid Member ID")
		return
	}

	err = lc.library.ReserveBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println("Book reserved successfully!")
}

func (lc *LibraryController) Run() {
	for {
		lc.showMenu()
		choice := lc.getInput("Choose an option: ")

		switch choice {
		case "1":
			lc.AddBook()
		case "2":
			lc.removeBook()
		case "3":
			lc.addMember()
		case "4":
			lc.borrowBook()
		case "5":
			lc.returnBook()
		case "6":
			lc.listAvailableBooks()
		case "7":
			lc.listBorrowedBooks()
		case "8":
			lc.listAllBooks()
		case "9":
			lc.listAllMember()
		case "10":
			lc.reserveBook()
		case "0":
			fmt.Println("Good bay!")
			return
		default:
			fmt.Println("Invalid option")
		}

		fmt.Println("\nPress Enter to continue...")
		lc.scanner.Scan()
	}
}

func (lc *LibraryController) showMenu() {
	fmt.Println("\n===Library Management System===")
	fmt.Println("1. Add a book")
	fmt.Println("2. Remove a book")
	fmt.Println("3. Add a member")
	fmt.Println("4. Borrow a book")
	fmt.Println("5. Return a book")
	fmt.Println("6. List available books")
	fmt.Println("7. List books borrowed by a member")
	fmt.Println("8. List all books")
	fmt.Println("9. List all members")
	fmt.Println("10. Reserve a book")
	fmt.Println("0. Exit")

}

func (lc *LibraryController) getInput(prompt string) string {
	fmt.Println(prompt)
	lc.scanner.Scan()
	return strings.TrimSpace(lc.scanner.Text())
}

func (lc *LibraryController) AddBook() {
	fmt.Println("\n---- Add a book ----")
	title := lc.getInput("Title: ")
	author := lc.getInput("Author: ")

	if title == "" || author == "" {
		fmt.Println("Error: Title and author are required")
		return
	}

	book := models.Book{
		Title:  title,
		Author: author,
		Status: "Available",
	}

	lc.library.AddBook(book)
	fmt.Println("Book added successfully!")
}

func (lc *LibraryController) removeBook() {
	fmt.Println("\n---- Remove a book ----")
	bookIDStr := lc.getInput("Book ID to delete: ")
	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		fmt.Println("Error: Invalid ID")
		return
	}

	err = lc.library.RemoveBook(bookID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println("Book removed successfully!")
}

func (lc *LibraryController) addMember() {
	fmt.Println("\n---- Add a member ----")
	name := lc.getInput("Full-Name: ")

	if name == "" {
		fmt.Println("error: Name is required")
		return
	}
	member := models.Member{
		Name: name,
	}
	lc.library.AddMember(member)
	fmt.Println("Members adds successfully!")
}

func (lc *LibraryController) borrowBook() {
	fmt.Println("\n---- Borrow a book ----")
	bookIDStr := lc.getInput("book ID: ")
	memberIDStr := lc.getInput("Member ID: ")

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		fmt.Println("Error: Invalid Book ID")
		return
	}
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		fmt.Println("Error: Invalid Member ID")
		return
	}
	err = lc.library.BorrowBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println("Book borrowed successfully!")
}

func (lc *LibraryController) returnBook() {
	fmt.Println("\n---- Return a book ----")
	bookIDStr := lc.getInput("Book ID: ")
	memberIDstr := lc.getInput("Member ID: ")

	bookID, err := strconv.Atoi(bookIDStr)
	if err != nil {
		fmt.Println("Error: Invalid Book ID")
		return
	}

	memberID, err := strconv.Atoi(memberIDstr)
	if err != nil {
		fmt.Println("Error: Invalid Member ID")
		return
	}

	err = lc.library.ReturnBook(bookID, memberID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println("Book return successfully")
}

func (lc *LibraryController) listAvailableBooks() {
	fmt.Println("\n---- Available books ----")
	books := lc.library.ListAvailableBooks()

	if len(books) == 0 {
		fmt.Println("No books available")
		return
	}

	for _, book := range books {
		fmt.Printf("ID:%d | Title: %s | Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func (lc *LibraryController) listBorrowedBooks() {
	fmt.Println("\n---- Books borrowed by a member ----")
	memberIDStr := lc.getInput("Member ID: ")
	memberID, err := strconv.Atoi(memberIDStr)
	if err != nil {
		fmt.Println("Error: Invalid member ID")
		return
	}
	books := lc.library.ListBorrowedBooks(memberID)

	if len(books) == 0 {
		fmt.Println("No borrowed book")
		return
	}

	for _, book := range books {
		fmt.Printf("ID:%d | Title: %s | Author: %s\n", book.ID, book.Title, book.Author)
	}
}

func (lc *LibraryController) listAllBooks() {
	fmt.Println("\n--- Tous les livres ---")

	availableBooks := lc.library.ListAvailableBooks()

	if len(availableBooks) == 0 {
		fmt.Println("No books in the library")
		return
	}

	for _, book := range availableBooks {
		fmt.Printf("ID: %d | Titre: %s | Auteur: %s | Statut: %s\n",
			book.ID, book.Title, book.Author, book.Status)
	}
}

func (lc *LibraryController) listAllMember() {
	fmt.Println("\n---- All members ----")
}

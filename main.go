package main

import (
	"fmt"
	"time"

	"library_management/controllers"
	"library_management/models"
	"library_management/services"
)

func main() {
	// Create a library
	library := services.NewLibrary()

	// seed some data
	library.AddBook(models.Book{Title: "Concurrency in Go", Author: "K. Author", Status: "Available"})
	library.AddBook(models.Book{Title: "The Go Programming Language", Author: "A. Donovan", Status: "Available"})

	library.AddMember(models.Member{Name: "Alice"})
	library.AddMember(models.Member{Name: "Bob"})
	library.AddMember(models.Member{Name: "Carol"})

	fmt.Println("Starting reservation simulation for book ID 1 by multiple members concurrently...")

	// Simulate multiple concurrent reservation attempts for the same book
	go func() {
		if err := library.ReserveBook(1, 1); err != nil {
			fmt.Printf("Member 1 reservation error: %v\n", err)
		} else {
			fmt.Println("Member 1 reserved book 1")
		}
	}()

	go func() {
		if err := library.ReserveBook(1, 2); err != nil {
			fmt.Printf("Member 2 reservation error: %v\n", err)
		} else {
			fmt.Println("Member 2 reserved book 1")
		}
	}()

	go func() {
		if err := library.ReserveBook(1, 3); err != nil {
			fmt.Printf("Member 3 reservation error: %v\n", err)
		} else {
			fmt.Println("Member 3 reserved book 1")
		}
	}()

	// let the reservation workers and timers run
	time.Sleep(8 * time.Second)

	// Show status of books and borrowed lists
	fmt.Println("\nFinal state: available books:")
	for _, b := range library.ListAvailableBooks() {
		fmt.Printf("ID:%d | Title:%s | Author:%s | Status:%s\n", b.ID, b.Title, b.Author, b.Status)
	}

	fmt.Println("\nBorrowed books by members:")
	for i := 1; i <= 3; i++ {
		borrowed := library.ListBorrowedBooks(i)
		for _, b := range borrowed {
			fmt.Printf("Member %d borrowed -> ID:%d | %s\n", i, b.ID, b.Title)
		}
	}

	// create a controller (interactive) - kept for backward compatibility
	controller := controllers.NewLibraryController(library)
	// controller.Run() // keep interactive run commented in demo mode
	_ = controller
}

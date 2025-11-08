# Library Management System Documentation

## Project Structure

library_management/
├── main.go
├── controllers/
│   └── library_controller.go
├── models/
│   └── book.go
│   └── member.go
├── services/
│   └── library_service.go
├── docs/
│   └── documentation.md
└── go.mod
## Components Description

1. **Models**

    -Book: Represents a book in the library with ID, Title, Author, and Status fields

    -Member: Represents a library member with ID, Name, and BorrowedBooks fields

2. **Services**

    -LibraryManager: Interface defining the contract for library operations

    -Library: Concrete implementation of the LibraryManager interface

3. **Controllers**

    -LibraryController: Handles console interaction and user input

4. **Features**
    1. Book Management
    Add a new book

    Remove an existing book

    List all available books

    2. Member Management
    Add a new member

    List all members

    3. Borrowing System
    Borrow a book

    Return a book

## List books borrowed by a specific member

**Usage**
    - Run the application with go run main.go and follow the menu instructions.

**Go Concepts Used**
        - Structs: For defining Book and Member data structures

        - Interfaces: LibraryManager interface for abstraction

        - Maps: For storing books and members with ID-based access

        - Slices: For managing collections of borrowed books

        - Error Handling: Comprehensive error checking and reporting

        - Packages: Modular code organization with separate concerns

        - Methods: Functions associated with struct types

**Key Methods**

 *Library Service Methods*

    - AddBook(book Book): Adds a new book to the library

    - RemoveBook(bookID int) error: Removes a book by ID

    - BorrowBook(bookID int, memberID int) error: Handles book borrowing

    - ReturnBook(bookID int, memberID int) error: Handles book returns

    - ListAvailableBooks() []Book: Returns all available books

    - ListBorrowedBooks(memberID int) []Book: Returns books borrowed by a member

 *Error Handling*

The system handles various error scenarios:

Book not found

Member not found

Book already borrowed

Attempt to remove a borrowed book

Invalid input formats

# Console Interface
*The application provides a user-friendly console menu with the following options:*

    - Add a book

    - Remove a book

    - Add a member

    - Borrow a book

    - Return a book

    - List available books

    - List borrowed books by member

    - List all books

    - List all members

    - Exit



## This project demonstrates practical application of Go's core features while building a functional library management system.
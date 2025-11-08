package main

import (
	"library_management/controllers"
	"library_management/services"
)

func main() {
	// Create a library
	library := services.NewLibrary()

	// Create a controller
	controller := controllers.NewLibraryController(library)

	// Start the app
	controller.Run()
}

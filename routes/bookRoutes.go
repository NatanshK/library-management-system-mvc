package routes

import (
	"library-management-system-mvc/controllers"
	"library-management-system-mvc/middleware"
	"net/http"
)

func RegisterBookRoutes(mux *http.ServeMux) {
	// --- CLIENT ROUTES ---
	mux.HandleFunc("/api/books", middleware.RequireAuth(controllers.GetCatalog))

	// --- ADMIN ROUTES ---
	mux.HandleFunc("/api/books/add", middleware.RequireAdmin(controllers.AddBook))
	mux.HandleFunc("/api/books/update", middleware.RequireAdmin(controllers.UpdateBook))
	mux.HandleFunc("/api/books/delete", middleware.RequireAdmin(controllers.DeleteBook))
}

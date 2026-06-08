package routes

import (
	"library-management-system-mvc/controllers"
	"library-management-system-mvc/middleware"
	"net/http"
)

func RegisterTxRoutes(mux *http.ServeMux) {
	// --- CLIENT ROUTES ---
	mux.HandleFunc("/api/transactions/request", middleware.RequireAuth(controllers.RequestCheckout))
	mux.HandleFunc("/api/transactions/history", middleware.RequireAuth(controllers.GetUserHistory))
	mux.HandleFunc("/api/transactions/return", middleware.RequireAuth(controllers.RequestCheckin))

	// --- ADMIN ROUTES ---
	mux.HandleFunc("/api/transactions/queue", middleware.RequireAdmin(controllers.GetAdminQueue))
	mux.HandleFunc("/api/transactions/approve", middleware.RequireAdmin(controllers.ApproveCheckout))
	mux.HandleFunc("/api/transactions/checkin/approve", middleware.RequireAdmin(controllers.ApproveCheckin))
	mux.HandleFunc("/api/transactions/reject/checkout", middleware.RequireAdmin(controllers.RejectCheckout))
	mux.HandleFunc("/api/transactions/reject/checkin", middleware.RequireAdmin(controllers.RejectCheckin))
}

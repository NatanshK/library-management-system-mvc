package routes

import (
	"library-management-system-mvc/controllers"
	"library-management-system-mvc/middleware"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/register", controllers.Register)
	mux.HandleFunc("/api/login", controllers.Login)
	mux.HandleFunc("/api/users/promote/request", middleware.RequireAuth(controllers.RequestPromotion))
	mux.HandleFunc("/api/users/promote/approve", middleware.RequireAdmin(controllers.ApprovePromotion))
	mux.HandleFunc("/api/users/promote/queue", middleware.RequireAdmin(controllers.GetPendingPromotions))
}

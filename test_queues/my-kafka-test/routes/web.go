package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"my-kafka-test/app/controllers"
)

// RegisterWebRoutes 注册 Web 路由
func RegisterWebRoutes(router *mux.Router) {
	// 首页路由
	router.HandleFunc("/", controllers.NewHomeController().Index).Methods("GET")
	router.HandleFunc("/health", controllers.NewHomeController().Health).Methods("GET")
	
	// API 路由
	api := router.PathPrefix("/api").Subrouter()
	
	// 用户路由
	userController := controllers.NewUserController()
	api.HandleFunc("/users", userController.Index).Methods("GET")
	api.HandleFunc("/users", userController.Store).Methods("POST")
	api.HandleFunc("/users/{id}", userController.Show).Methods("GET")
	
	// 静态文件
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
}
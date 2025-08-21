package router

import (
	"went-framework/app/controllers"

	"github.com/gorilla/mux"
)

func setupUserRoutes(api *mux.Router) {

	// User routes
	api.HandleFunc("/users", controllers.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", controllers.GetUser).Methods("GET")
	api.HandleFunc("/users", controllers.CreateUser).Methods("POST")
	api.HandleFunc("/users/{id}", controllers.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", controllers.DeleteUser).Methods("DELETE")

}

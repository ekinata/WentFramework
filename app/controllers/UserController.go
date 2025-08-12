package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"went-framework/app/database"
	"went-framework/app/models"

	"github.com/gorilla/mux"
)

// Response structure for JSON responses
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// GetAllUsers handles GET /api/users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Connect to database if not already connected
	if database.DB == nil {
		database.Connect()
	}

	// Get users from database
	users, err := models.GetAllUsers(database.DB)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Failed to retrieve users: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    users,
	}

	json.NewEncoder(w).Encode(response)
}

// GetUser handles GET /api/users/{id}
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Connect to database if not already connected
	if database.DB == nil {
		database.Connect()
	}

	// Get user from database
	user, err := models.GetUserByID(database.DB, uint(id))
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "User retrieved successfully",
		Data:    user,
	}

	json.NewEncoder(w).Encode(response)
}

// CreateUser handles POST /api/users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid JSON data",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate required fields
	if userData.Name == "" || userData.Email == "" {
		response := Response{
			Status:  "error",
			Message: "Name and email are required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Connect to database if not already connected
	if database.DB == nil {
		database.Connect()
	}

	// Create new user
	user := models.User{
		Name:  userData.Name,
		Email: userData.Email,
	}

	if err := user.Create(database.DB); err != nil {
		response := Response{
			Status:  "error",
			Message: "Failed to create user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    user,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
} // UpdateUser handles PUT /api/users/{id}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var userData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid JSON data",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Connect to database if not already connected
	if database.DB == nil {
		database.Connect()
	}

	// Get existing user
	user, err := models.GetUserByID(database.DB, uint(id))
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update user fields
	if userData.Name != "" {
		user.Name = userData.Name
	}
	if userData.Email != "" {
		user.Email = userData.Email
	}

	// Save updated user
	if err := user.Update(database.DB); err != nil {
		response := Response{
			Status:  "error",
			Message: "Failed to update user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    user,
	}

	json.NewEncoder(w).Encode(response)
}

// DeleteUser handles DELETE /api/users/{id}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid user ID",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Connect to database if not already connected
	if database.DB == nil {
		database.Connect()
	}

	// Get existing user to verify it exists
	user, err := models.GetUserByID(database.DB, uint(id))
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "User not found",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete user
	if err := user.Delete(database.DB); err != nil {
		response := Response{
			Status:  "error",
			Message: "Failed to delete user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: fmt.Sprintf("User %d deleted successfully", id),
	}

	json.NewEncoder(w).Encode(response)
}

package controllers

import (
	"fmt"
	"net/http"

    "went-framework/database"
	"went-framework/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	if db != nil {
		return nil
	}

	database.Connect()
	db = database.DB

	return nil
}

func GetAll{{.ModelName}}s(w http.ResponseWriter, r *http.Request) {
	Connect()
	var models []models.{{.ModelName}}
	err := db.Find(&models).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(models)
}

func Get{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    Connect()
    var model models.{{.ModelName}}
    id := mux.Vars(r)["id"]
    err := db.First(&model, id).Error
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(model)
}

func Create{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    Connect()
    var model models.{{.ModelName}}
    if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := db.Create(&model).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(model)
}

func Update{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    Connect()
    var model models.{{.ModelName}}
    id := mux.Vars(r)["id"]
    if err := db.First(&model, id).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    if err := db.Save(&model).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(model)
}

func Delete{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    Connect()
    var model models.{{.ModelName}}
    id := mux.Vars(r)["id"]
    if err := db.First(&model, id).Error; err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    if err := db.Delete(&model).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

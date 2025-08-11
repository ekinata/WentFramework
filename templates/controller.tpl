package controllers

import (
	"fmt"
	"net/http"
)

func GetAll{{.ModelName}}s(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "All {{.ModelName}}s")
}

func Get{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Get {{.ModelName}}")
}

func Create{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Create {{.ModelName}}")
}

func Update{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Update {{.ModelName}}")
}

func Delete{{.ModelName}}(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Delete {{.ModelName}}")
}

package handler

import (
	"encoding/json"
	"fmt"
	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("web", "templates", "layout.html")
	navbar := filepath.Join("web", "templates", "navbar.html")
	register := filepath.Join("web", "templates", "register.html")
	search := filepath.Join("web", "templates", "search.html")

	email := r.FormValue("email")
	password := r.FormValue("password")
	username := r.FormValue("username")

	user := models.User{
		Email:    email,
		Password: password,
		Username: username,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "JSON verisi oluşturulurken hata oluştu", http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles(layout, navbar, register, search)

	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	apiURL := apiUrlManager.GetUsersApiURL()
	response, errR := apiManager.Post(apiURL, jsonData)

	if errR != nil {
		http.Error(w, "Kayıt işlemi sırasında hata oluştu", http.StatusInternalServerError)
		return
	} else if response.StatusCode == http.StatusCreated {
		fmt.Println("Register : ", response.StatusCode)
		data := struct {
			Email    string
			Password string
			Username string
			Success  bool
		}{
			Email:    email,
			Password: "",
			Username: username,
			Success:  true,
		}

		_ = tmpl.Execute(w, data)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = tmpl.Execute(w, nil)

}

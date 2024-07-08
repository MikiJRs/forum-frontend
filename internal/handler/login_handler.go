package handler

import (
	"encoding/json"
	"fmt"
	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("web", "templates", "layout.html")
	navbar := filepath.Join("web", "templates", "navbar.html")
	login := filepath.Join("web", "templates", "login.html")
	search := filepath.Join("web", "templates", "search.html")
	fmt.Println("Login Page Handler")
	tmpl, err := template.ParseFiles(layout, navbar, login, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = tmpl.Execute(w, nil)
}
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login Form Handler")
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println(email, password)

	apiManager := manager.NewAPIManager()
	apiUrls := manager.NewAPIUrls()

	apiURL := apiUrls.GetLoginApiURL()

	user := models.User{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, errR := apiManager.Post(apiURL, jsonData)
	if errR != nil {
		fmt.Println(errR)
		http.Error(w, errR.Error(), http.StatusInternalServerError)
		return
	}

	if response.StatusCode == http.StatusOK {
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var loginResponse models.LoginResponse
		err = json.Unmarshal(body, &loginResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userToken := loginResponse.Token
		fmt.Println(userToken)
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    userToken,
			HttpOnly: true,
		})
		fmt.Println("Login successful")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Hata durumunda yanıtı ve durum kodunu yazdırın
	fmt.Println("API yanıtı:", response.Status)
	body, err := io.ReadAll(response.Body)
	if err == nil {
		fmt.Println("API yanıt gövdesi:", string(body))
	}

	wrong := struct {
		Success bool
		Message string
	}{
		Success: false,
		Message: "Mail veya şifre yanlış",
	}
	fmt.Println(wrong)

	layout := filepath.Join("web", "templates", "layout.html")
	navbar := filepath.Join("web", "templates", "navbar.html")
	login := filepath.Join("web", "templates", "login.html")
	search := filepath.Join("web", "templates", "search.html")

	tmpl, err := template.ParseFiles(layout, navbar, login, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = tmpl.Execute(w, wrong)
}

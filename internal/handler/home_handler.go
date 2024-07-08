package handler

import (
	"encoding/json"
	"fmt"
	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
	"golang-forum-frontend/internal/repository"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("web", "templates", "layout.html")
	index := filepath.Join("web", "templates", "index.html")
	navbar := filepath.Join("web", "templates", "navbar.html")
	search := filepath.Join("web", "templates", "search.html")
	post := filepath.Join("web", "templates", "posts.html")

	tmpl, err := template.ParseFiles(layout, index, navbar, search, post)
	if err != nil {
		fmt.Println("Error parsing templates:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var username string
	var posts []models.Post

	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	token, _ := r.Cookie("token")

	if token != nil {
		apiManager.SetUserToken(token.Value)

		apiURL := apiUrlManager.GetUserApiURL()

		userResponse, err := apiManager.Get(apiURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if userResponse.StatusCode == http.StatusOK {
			defer userResponse.Body.Close()
			body, err := io.ReadAll(userResponse.Body)
			fmt.Println("User Response Body:", string(body)) // Log the response body
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var user models.User
			fmt.Println("User Response Body:", string(body)) // Log the response body
			err = json.Unmarshal(body, &user)
			fmt.Println("User : ", user)

			if err != nil {
				fmt.Println("Error unmarshalling user data")
				http.Error(w, "Error unmarshalling user data", http.StatusInternalServerError)
				return
			}
			username = user.Username

		}

	}

	fmt.Println("4444")

	posts, _ = repository.GetPosts()
	data := struct {
		Username string
		Posts    []models.Post
	}{
		Username: username,
		Posts:    posts,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

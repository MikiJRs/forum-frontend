package handler

import (
	"encoding/json"
	"fmt"
	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
	"golang-forum-frontend/internal/repository"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	layout := filepath.Join("web", "templates", "layout.html")
	navbar := filepath.Join("web", "templates", "navbar.html")
	search := filepath.Join("web", "templates", "search.html")
	postHtml := filepath.Join("web", "templates", "post.html")

	cookie, err := r.Cookie("token")
	apiManager := manager.NewAPIManager()
	var token string
	if err != nil {
		guestToken := apiManager.GetGuestToken()
		cookie = &http.Cookie{
			Name:  "token",
			Value: guestToken,
		}
	}

	token = cookie.Value

	vars := mux.Vars(r)
	postID := vars["id"]
	var post models.Post
	var comments []models.Comment

	post, errPost := repository.GetPost(postID)
	comments, errComment := repository.GetPostComments(postID, token)

	if errComment != nil {
		http.Error(w, errComment.Error(), http.StatusInternalServerError)
		return
	}

	if errPost != nil {
		http.Error(w, errPost.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(layout, navbar, search, postHtml)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username, err := repository.GetUsername(r)

	if err != nil {
		fmt.Println("Error 22:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Username:", username)
	data := struct {
		Post     models.Post
		Comments []models.Comment
		Username string
	}{
		Post:     post,
		Comments: comments,
		Username: username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetPostCommentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Post Comments Handler")
	vars := mux.Vars(r)
	postID := vars["id"]

	// Create a channel to receive the comments
	commentsChan := make(chan []models.Comment)

	cookie, err := r.Cookie("token")

	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := cookie.Value

	// Start a goroutine to fetch the comments
	go func() {
		comments, err := repository.GetPostComments(postID, token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		commentsChan <- comments
	}()

	// Wait for the comments to be fetched and then send them as a response
	comments := <-commentsChan
	fmt.Println("Comments:", comments)
	for i, comment := range comments {
		fmt.Println("Comment:", comment)
		username, err := repository.GetUserNameFromId(comment.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		comments[i].Author = username
		fmt.Println("Username:", username)
	}

	jsonData, err := json.Marshal(comments)
	if err != nil {
		http.Error(w, "JSON verisi oluşturulurken hata oluştu", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

}

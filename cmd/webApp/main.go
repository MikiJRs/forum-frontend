package main

import (
	"golang-forum-frontend/internal/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// File server for serving static files
	fileServer := http.FileServer(http.Dir("./web/static/"))
	router.PathPrefix("/web/static/").Handler(http.StripPrefix("/web/static", fileServer))

	// Route for the home handler
	router.HandleFunc("/", handler.HomeHandler)
	router.HandleFunc("/login.html", handler.LoginPageHandler).Methods("GET")
	router.HandleFunc("/login.html", handler.LoginFormHandler).Methods("POST")
	router.HandleFunc("/register.html", handler.RegisterHandler)
	router.HandleFunc("/register", handler.RegisterHandler)
	router.HandleFunc("/logout", handler.LogoutHandler)
	router.HandleFunc("/post/{id}", handler.PostHandler).Methods("GET")
	router.HandleFunc("/add-comment", handler.AddCommentHandler).Methods("POST")
	router.HandleFunc("/get-post-comments", handler.GetPostCommentsHandler).Methods("GET")
	router.HandleFunc("/add-post", handler.AddPostHandler).Methods("GET")
	router.HandleFunc("/add-post", handler.AddPostSubmitHandler).Methods("POST")
	//router.HandleFunc("/search", handler.SearchHandler).Methods("GET")
	router.HandleFunc("/add-comment", handler.AddCommentHandler).Methods("POST")

	// Starting the server
	log.Println("Starting server on :8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

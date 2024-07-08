package handler

import (
	"fmt"
	"golang-forum-frontend/internal/models"
	"golang-forum-frontend/internal/repository"
	"net/http"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		postID := r.FormValue("postID")
		content := r.FormValue("content")
		userId := repository.GetUserId(token.Value)

		if postID == "" {
			http.Error(w, "Post ID eksik", http.StatusBadRequest)
			return
		}

		comment := models.Comment{
			Content: content,
			UserID:  userId,
		}
		fmt.Println("Post ID:", postID)
		fmt.Println("Comment:", comment)
		fmt.Println("User ID:", userId)
		err = repository.AddComment(postID, comment, userId)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/post/%s", postID), http.StatusSeeOther)
	}
}

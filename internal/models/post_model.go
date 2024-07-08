package models

type Post struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Categories    []string `json:"categories"`
	UserID        int      `json:"user_id"`
	Likes         int      `json:"likes"`
	Dislikes      int      `json:"dislikes"`
	CommentLength int      `json:"comment_length"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	Author        string   `json:"author"` // Added Author field
}

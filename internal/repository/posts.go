package repository

import (
	"encoding/json"
	"fmt"
	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
	"io"
	"net/http"
)

func GetPosts() ([]models.Post, error) {
	var posts []models.Post

	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	postApiUrl := apiUrlManager.GetPostsApiURL()
	fmt.Println("Post API URL:", postApiUrl) // Log the post API URL

	postResponse, errPost := apiManager.Get(postApiUrl)
	if errPost != nil {

		return nil, errPost
	}
	if postResponse.StatusCode == http.StatusOK {
		defer postResponse.Body.Close()
		body, err := io.ReadAll(postResponse.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &posts)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(posts); i++ {
			posts[i].Author, _ = GetUserNameFromId(posts[i].UserID)
		}

	}

	return posts, nil
}

func GetPost(postId string) (models.Post, error) {
	var post models.Post

	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	postApiUrl := apiUrlManager.GetPostsApiURL() + "/" + postId

	postResponse, errPost := apiManager.Get(postApiUrl)
	if errPost != nil {

		return post, errPost
	}
	if postResponse.StatusCode == http.StatusOK {
		defer postResponse.Body.Close()
		body, err := io.ReadAll(postResponse.Body)

		postAuthor, errUser := GetUserNameFromId(post.UserID)
		if errUser != nil {
			return post, errUser
		}

		if err != nil {
			return post, err
		}

		err = json.Unmarshal(body, &post)
		postAuthor, _ = GetUserNameFromId(post.UserID)

		newPost := models.Post{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			UpdatedAt: post.UpdatedAt,
			CreatedAt: post.CreatedAt,
			Author:    postAuthor,
		}
		return newPost, err

	}

	return post, nil
}

func AddPost(post models.Post) error {
	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	postApiUrl := apiUrlManager.GetPostsApiURL()
	fmt.Println("POST : ", post)
	body, err := json.Marshal(post)
	if err != nil {
		return err
	}
	fmt.Println("Body :", string(body))
	postResponse, errPost := apiManager.Post(postApiUrl, body)

	if errPost != nil {

		return errPost
	}
	if postResponse.StatusCode == http.StatusOK {
		defer postResponse.Body.Close()
		body, err := io.ReadAll(postResponse.Body)
		if err != nil {
			return err
		}

		fmt.Println("API Response Body:", string(body)) // Log the response body
	}
	fmt.Println("API Response Status:", postResponse.StatusCode) // Log the response status code
	return nil
}

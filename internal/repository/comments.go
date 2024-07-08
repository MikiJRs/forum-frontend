package repository

import (
	"encoding/json"
	"fmt"
	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
	"io"
	"net/http"
)

func GetComments(token string) ([]models.Comment, error) {
	var comments []models.Comment

	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	apiManager.SetUserToken(token)
	commentsApiUrl := apiUrlManager.GetCommentsApiURL()
	fmt.Println("Comments API URL:", commentsApiUrl) // Log the comments API URL
	commentsResponse, errComment := apiManager.Get(commentsApiUrl)

	fmt.Println("GetComments Status:", commentsResponse.StatusCode)

	if errComment != nil {
		return nil, errComment
	}

	if commentsResponse.StatusCode == http.StatusOK {
		defer commentsResponse.Body.Close()
		body, err := io.ReadAll(commentsResponse.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &comments)
		if err != nil {
			return nil, err
		}

	}

	return comments, nil
}

func GetPostComments(commentId string, token string) ([]models.Comment, error) {
	var comments []models.Comment
	var commentResponse *http.Response

	apiManager := manager.NewAPIManager()
	apiManager.SetUserToken(token)
	apiUrlManager := manager.NewAPIUrls()
	commentApiUrl := apiUrlManager.GetPostsApiURL() + "/" + commentId + "/comments"
	commentResponse, errComment := apiManager.Get(commentApiUrl)

	if errComment != nil {
		return nil, errComment
	}

	if commentResponse.StatusCode == http.StatusOK {
		defer commentResponse.Body.Close()
		body, err := io.ReadAll(commentResponse.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &comments)
		if err != nil {
			return nil, err
		}

		for i, comment := range comments {
			fmt.Println("Comment:", comment.UserID)
			username, _ := GetUserNameFromId(comment.UserID)

			comments[i].Author = username
			fmt.Println("Username:", username)
		}

		return comments, nil
	}

	return comments, nil
}

func AddComment(postId string, comment models.Comment, userId int) error {
	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	commentApiUrl := apiUrlManager.GetPostsApiURL() + "/" + postId + "/comments"
	comment.UserID = userId
	fmt.Println("Comment:", comment)
	fmt.Println("Comment API URL:", commentApiUrl)
	fmt.Println("userId :", userId)
	body, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	commentResponse, errComment := apiManager.Post(commentApiUrl, body)
	if errComment != nil {
		return errComment
	}

	if commentResponse.StatusCode == http.StatusOK {
		defer commentResponse.Body.Close()
		body, err := io.ReadAll(commentResponse.Body)
		if err != nil {
			return err
		}

		var comment models.Comment
		err = json.Unmarshal(body, &comment)
		if err != nil {
			return err
		}

	}

	return nil
}

package repository

import (
	"encoding/json"
	"fmt"
	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
	"io"
	"net/http"
)

func GetUsername(r *http.Request) (string, error) {
	fmt.Println("Get Username")
	var username string
	// Get token from cookie

	cookie, err := r.Cookie("token")

	if err != nil {
		username = "Guest"
		return username, nil
	}

	token := cookie.Value

	apiManager := manager.NewAPIManager()
	apiManager.SetUserToken(token)
	apiUrlManager := manager.NewAPIUrls()
	userApiUrl := apiUrlManager.GetUserApiURL()

	userResponse, err := apiManager.Get(userApiUrl)

	fmt.Println("Token :", token)
	fmt.Println("User API URL 21:", userApiUrl)

	if err != nil {
		fmt.Println("Error:", err)
		return username, err
	}

	fmt.Println("User Response:", userResponse)

	if userResponse.StatusCode == http.StatusOK {
		defer userResponse.Body.Close()
		body, err := io.ReadAll(userResponse.Body)
		fmt.Println("Body:", string(body))

		if err != nil {
			return username, err
		}

		var user models.User

		err = json.Unmarshal(body, &user)
		if err != nil {
			return username, err
		}

		fmt.Println("User: test", user)

		username = user.Username
	}
	return username, nil

}

func GetUserNameFromId(userId int) (string, error) {
	var username string

	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	userApiUrl := apiUrlManager.GetUsersApiURL() + "/" + fmt.Sprint(userId)

	userResponse, err := apiManager.Get(userApiUrl)

	if err != nil {
		fmt.Println("Error:", err)
		return username, err
	}

	if userResponse.StatusCode == http.StatusOK {
		defer userResponse.Body.Close()
		body, err := io.ReadAll(userResponse.Body)
		fmt.Println("Body:", string(body))

		if err != nil {
			return username, err
		}

		var user models.User

		err = json.Unmarshal(body, &user)
		if err != nil {
			return username, err
		}

		username = user.Username
	}
	return username, nil

}

func GetUserId(token string) int {
	apiManager := manager.NewAPIManager()
	apiUrlManager := manager.NewAPIUrls()
	apiManager.SetUserToken(token)

	userApiUrl := apiUrlManager.GetUserApiURL()

	userResponse, err := apiManager.Get(userApiUrl)
	fmt.Println("-----------------")
	fmt.Println("User API URL:", userApiUrl)
	fmt.Println("User Response:", userResponse)
	fmt.Println("Token :", token)
	fmt.Println("-----------------")

	if err != nil {
		return 0
	}

	var user models.User
	err = json.NewDecoder(userResponse.Body).Decode(&user)
	if err != nil {
		return 0
	}
	return user.ID

}

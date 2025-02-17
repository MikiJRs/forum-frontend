package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang-forum-frontend/internal/manager"
	"golang-forum-frontend/internal/models"
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

// google github facebook

var (
	googleClientID       string
	googleClientSecret   string
	githubClientID       string
	githubClientSecret   string
	facebookClientID     string
	facebookClientSecret string
)

func LoadEnv() {
	file, err := os.Open("/home/ertugrulmert/MikiJR/forum-frontend/.env")
	if err != nil {
		log.Fatalf("Error opening .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := parts[1]
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	googleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	githubClientID = os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
	facebookClientID = os.Getenv("FACEBOOK_CLIENT_ID")
	facebookClientSecret = os.Getenv("FACEBOOK_CLIENT_SECRET")
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := "https://accounts.google.com/o/oauth2/auth?client_id=" + googleClientID +
		"&redirect_uri=http://localhost:8081/callback/google" +
		"&scope=https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email&prompt=select_account" +
		"&response_type=code" +
		"&state=random-string"
	http.Redirect(w, r, url, http.StatusFound)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "random-string" {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in response", http.StatusBadRequest)
		return
	}

	token, err := ExchangeGoogleCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := GetGoogleUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	email, emailOk := userInfo["email"].(string)
	name, nameOk := userInfo["name"].(string)

	if !emailOk || !nameOk || email == "" {
		log.Printf("Google user info is missing required fields: %+v", userInfo)
		http.Error(w, "Failed to get valid user info", http.StatusInternalServerError)
		return
	}

	apiManager := manager.NewAPIManager()
	apiUrls := manager.NewAPIUrls()

	apiURL := apiUrls.GetThirdPartyLoginApiURL()

	user := models.User{
		Email:    email,
		Username: name,
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
			Path:     "/",
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

func ExchangeGoogleCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", googleClientID)
	data.Set("client_secret", googleClientSecret)
	data.Set("redirect_uri", "http://localhost:8081/callback/google")
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("no access token in response")
}

func GetGoogleUserInfo(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)
	return userInfo, nil
}

func HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := "https://github.com/login/oauth/authorize?client_id=" + githubClientID +
		"&redirect_uri=http://localhost:8081/callback/github" +
		"&scope=read:user user:email&prompt=select_account" +
		"&state=random-string"
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "random-string" {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in response", http.StatusBadRequest)
		return
	}

	token, err := ExchangeGitHubCodeForToken(code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := GetGitHubUserInfo(token)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	email, emailOk := userInfo["email"].(string)
	if !emailOk || email == "" {
		// Email doğrudan alınamadıysa, ek endpoint'ten email bilgilerini al
		emails, err := GetGitHubUserEmails(token)
		if err != nil || len(emails) == 0 {
			log.Printf("Failed to get valid user info: %v", err)
			http.Error(w, "Failed to get valid user info", http.StatusInternalServerError)
			return
		}
		email = emails[0]
	}

	username, usernameOk := userInfo["login"].(string)
	if !usernameOk || username == "" {
		http.Error(w, "Failed to get valid user info", http.StatusInternalServerError)
		return
	}

	apiManager := manager.NewAPIManager()
	apiUrls := manager.NewAPIUrls()

	apiURL := apiUrls.GetThirdPartyLoginApiURL()
	user := models.User{
		Email:    email,
		Username: username,
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
			Path:     "/",
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

func ExchangeGitHubCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", githubClientID)
	data.Set("client_secret", githubClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:8080/callback/github")

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		return "", err
	}
	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("no access token in response")
}

func GetGitHubUserInfo(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)
	return userInfo, nil
}

func GetGitHubUserEmails(token string) ([]string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var emails []struct {
		Email      string `json:"email"`
		Primary    bool   `json:"primary"`
		Verified   bool   `json:"verified"`
		Visibility string `json:"visibility"`
	}
	err = json.NewDecoder(resp.Body).Decode(&emails)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, e := range emails {
		if e.Verified {
			result = append(result, e.Email)
		}
	}
	return result, nil
}

//FACEBOOK

func HandleFacebookLogin(w http.ResponseWriter, r *http.Request) {
	url := "https://www.facebook.com/v12.0/dialog/oauth?client_id=" + facebookClientID +
		"&redirect_uri=http://localhost:8081/callback/facebook" +
		"&scope=email&state=random-string"
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleFacebookCallback(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state != "random-string" {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in response", http.StatusBadRequest)
		return
	}

	token, err := ExchangeFacebookCodeForToken(code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userInfo, err := GetFacebookUserInfo(token)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	email, emailOk := userInfo["email"].(string)
	name, nameOk := userInfo["name"].(string)

	if !emailOk || !nameOk || email == "" {
		log.Printf("Facebook user info is missing required fields: %+v", userInfo)
		http.Error(w, "Failed to get valid user info", http.StatusInternalServerError)
		return
	}

	apiManager := manager.NewAPIManager()
	apiUrls := manager.NewAPIUrls()

	apiURL := apiUrls.GetThirdPartyLoginApiURL()
	user := models.User{
		Email:    email,
		Username: name,
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
			Path:     "/",
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

func ExchangeFacebookCodeForToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", facebookClientID)
	data.Set("client_secret", facebookClientSecret)
	data.Set("redirect_uri", "http://localhost:8081/callback/facebook")
	data.Set("code", code)

	resp, err := http.PostForm("https://graph.facebook.com/v12.0/oauth/access_token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if token, ok := result["access_token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("no access token in response")
}

func GetFacebookUserInfo(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", "https://graph.facebook.com/me?fields=email,name", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userInfo)
	return userInfo, nil
}

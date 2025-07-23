package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func registerUser(username, password string) error {
	url := "http://localhost:8083/api/v5/authentication/built_in_database/users"

	data := &RegisterUserRequest{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	if resp.StatusCode == 201 {
		return nil
	} else {
		return fmt.Errorf("failed to register user")
	}
}

func main() {
	username := "zqy"
	password := "zqy123456"

	err := registerUser(username, password)
	if err != nil {
		fmt.Printf("Error registering user: %v\n", err)
	} else {
		fmt.Println("User registered successfully.")
	}
}

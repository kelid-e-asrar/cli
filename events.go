package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

// SyncPasswords Updates local Passwords
func SyncPasswords(userID string, deviceID string) error {

	marshaled, err := json.Marshal(map[string]string{
		"userID":   userID,
		"deviceID": deviceID,
	})

	if err != nil {
		log.Panic("error in marshaling json")
		return err
	}
	timeout := time.Duration(30 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	payload := bytes.NewBuffer(marshaled)

	request, err := http.NewRequest("GET", "localhost:8080/events/?user_id="+userID, payload)
	if err != nil {
		log.Panic("error in creating request to sync messages")
		return err
	}
	request.Header.Set("Content-type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal("error in sending request to server")
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		return nil
	}

	return errors.New("error in syncing password with server")
}

    

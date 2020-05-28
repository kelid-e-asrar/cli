package main

import "github.com/google/uuid"

// CreateUUID for userID and deviceID
func CreateUUID() string {

	id := uuid.New()
	return id.String()
}

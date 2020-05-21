package main

import "os"

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var ServerAddres = GetEnv("SERVER_ADDR", "localhost")
var ServerPort = GetEnv("SERVER_PORT", "8080")

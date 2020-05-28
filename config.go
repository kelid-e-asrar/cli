package main

import "os"

// ServerURL is the url of the server
var ServerURL = os.Getenv("SERVER_URL", "localhost")

// ServerPort is the server port
var ServerPort = os.Getenv("SERVER_PORT", "8080")

package main

import (
	"log"
	"os"
)

var (
	logInfo = log.New(os.Stdout, "Saying ", log.LstdFlags)
)

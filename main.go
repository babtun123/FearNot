package main

import (
	"FearNot/internal/orchestrator"
	"FearNot/internal/scripture"
	"FearNot/internal/verses"
	"log"
	"os"
)

func main() {
	// 1. Open the file (or create it if it doesn't exist)
	// os.O_APPEND: Add new logs to the end
	// os.O_CREATE: Create the file if it's missing
	// os.O_WRONLY: Open for writing only
	// 0666: Standard file permissions
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Application wide logger
	logger := log.New(file, "", log.LstdFlags|log.Lmicroseconds)

	// Get Verse of the Day
	verseOfTheDay, err := verses.GetVerseOfTheDay(logger)
	if err != nil {
		logger.Println(err)
	}

	ScriptureText, err := scripture.GetScripture(logger, verseOfTheDay)
	if err != nil {
		logger.Println(err)
	}

	NewOrchestrator := orchestrator.NewOrchestrator(verseOfTheDay, ScriptureText)
	NewOrchestrator.Run()
}

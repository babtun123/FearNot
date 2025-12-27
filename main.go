package main

import (
	"FearNot/internal/email"
	"FearNot/internal/orchestrator"
	"FearNot/internal/scripture"
	"FearNot/internal/verses"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// 1. Open the file (or create it if it doesn't exist)
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

	// Get the scripture text of the day
	ScriptureText, err := scripture.GetScripture(logger, verseOfTheDay)
	if err != nil {
		logger.Println(err)
	}

	// Print verse and scripture text
	fmt.Println("Verse of the day:", verseOfTheDay)
	fmt.Println("Scripture text: ")
	fmt.Println(ScriptureText)

	// Read from email list file
	emailList := openEmailSendList(logger)

	// Send the email
	err = email.GenerateEmail(logger, emailList, verseOfTheDay, ScriptureText)
	if err != nil {
		logger.Println(err)
	}

	NewOrchestrator := orchestrator.NewOrchestrator(verseOfTheDay, ScriptureText)
	NewOrchestrator.Run()
}

// openEmailSendList opens and reads email list file
func openEmailSendList(logger *log.Logger) []string {
	ret := make([]string, 0)
	file, err := os.Open("email_list.txt")
	if err != nil {
		logger.Println("Could not open email list file.", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ret = append(ret, line)
	}
	if scanner.Err() != nil {
		logger.Println("Could not open email list file.", err)
	}
	return ret
}

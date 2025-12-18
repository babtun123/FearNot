package orchestrator

import (
	"FearNot/internal/verses"
	"fmt"
)

type Orchestrator struct {
	VerseOfTheDay string
}

func NewOrchestrator() error {
	var orchestrator Orchestrator

	orchestrator.VerseOfTheDay = orchestrator.getVerseOfTheDay()

	return nil
}

func (orchestrator *Orchestrator) getVerseOfTheDay() string {
	// Get the verse of the day
	verse, err := verses.GetVerseOfTheDay()
	if err != nil {
		// in the future log an error here
		fmt.Println("could not get verse of day")
	}
	return verse
}

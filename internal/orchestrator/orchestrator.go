package orchestrator

import "fmt"

type Orchestrator struct {
	VerseOfTheDay     string
	ScriptureOfTheDay string
}

func NewOrchestrator(verseOfTheDay string, ScriptureText string) *Orchestrator {
	return &Orchestrator{
		VerseOfTheDay:     verseOfTheDay,
		ScriptureOfTheDay: ScriptureText,
	}
}

func (o *Orchestrator) Run() {
	fmt.Println("Verse of the day:", o.VerseOfTheDay)
	fmt.Println("Scripture text: ")
	fmt.Println(o.ScriptureOfTheDay)
}

//func (orchestrator *Orchestrator) getVerseOfTheDay() string {
//	// Get the verse of the day
//	verse, err := verses.GetVerseOfTheDay()
//	if err != nil {
//		// in the future log an error here
//		fmt.Println("could not get verse of day")
//	}
//	return verse
//}
//
//func (orchestrator *Orchestrator) fetchScriptureOfTheDay() string {
//	// Get the scripture of the day
//	scr, err := scripture.Run(orchestrator.VerseOfTheDay)
//	if err != nil {
//		fmt.Println("could not fetch scripture of day")
//	}
//	return scr
//}

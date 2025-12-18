/*
This is where I will handle the verse of the day.
*/

package verses

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type VerseState struct {
	Verses       []string `json:"verses"`
	CurrentIndex int      `json:"current_index"`
}

const (
	stateFile  = "verse_state.json"
	versesFile = "verses.txt"
)

func GetVerseOfTheDay() (string, error) {
	// Load or initialize state
	state, err := loadOrInitialize()
	if err != nil {
		return "", fmt.Errorf("could not load or initialize state: %w", err)
	}
	//state.reset()
	// Get today's verse
	verse := state.getNextVerse()
	fmt.Printf("Today's verse (%d/%d):\n%s\n",
		state.CurrentIndex, len(state.Verses), verse)

	// Save state
	if err := state.save(); err != nil {
		return "", fmt.Errorf("could not save state: %w", err)
	}

	return verse, nil
}

// loadOrInitialize loads existing state or creates new shuffled state
func loadOrInitialize() (*VerseState, error) {
	// try to load existing states, which is a json file. the first time should return an error
	//because the file should not exist.
	data, err := os.ReadFile(stateFile)
	if err == nil {
		var state VerseState
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, fmt.Errorf("error parsing the state file: %w", err)
		}
		fmt.Println("Loaded existing state")
		return &state, nil
	}

	// If the state file does not exist
	// In this case, the first time running the application
	fmt.Println("Loaded new state")
	verses, err := loadVersesFromFile()
	if err != nil {
		return nil, fmt.Errorf("error loading verses: %w", err)
	}
	state := &VerseState{
		Verses:       verses,
		CurrentIndex: 0,
	}

	// shuffle the verses. Not needed but for fun
	rand.Shuffle(len(state.Verses), func(i, j int) {
		state.Verses[i], state.Verses[j] = state.Verses[j], state.Verses[i]
	})

	fmt.Println("Loaded new verses")
	return state, nil
}

// loadVersesFromFile reads from a text file where the verses are found.
func loadVersesFromFile() ([]string, error) {
	file, err := os.Open(versesFile)
	if err != nil {
		return nil, fmt.Errorf("error opening verses file: %w", err)
	}
	defer file.Close()

	var verses []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		verses = append(verses, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading verses file: %w", err)
	}

	if len(verses) == 0 {
		return nil, fmt.Errorf("no verses found in file")
	}

	return verses, nil
}

// getNextVerse returns the current verse and advances to the next one
func (s *VerseState) getNextVerse() string {
	verse := s.Verses[s.CurrentIndex]
	s.CurrentIndex++

	// check if we've used all verse - reshuffle for next cycle
	if s.CurrentIndex >= len(s.Verses) {
		fmt.Println("Completed full cycle! Reshuffling ")
		rand.Shuffle(len(s.Verses), func(i, j int) {
			s.Verses[i], s.Verses[j] = s.Verses[j], s.Verses[i]
		})
		s.CurrentIndex = 0
	}

	return verse
}

// Save writes the current state to JSON file
func (s *VerseState) save() error {
	data, err := json.MarshalIndent(s, "", "	")
	if err != nil {
		return fmt.Errorf("error serializing verses: %w", err)
	}

	if err := os.WriteFile(stateFile, data, 0644); err != nil {
		return fmt.Errorf("error writing to state file: %w", err)
	}
	return nil
}

// getRemainingCount returns how many verses are left in current cycle
func (s *VerseState) getRemainingCount() int {
	return len(s.Verses) - s.CurrentIndex
}

// reset resets to the beginning and reshuffles
func (s *VerseState) reset() {
	rand.Shuffle(len(s.Verses), func(i, j int) {
		s.Verses[i], s.Verses[j] = s.Verses[j], s.Verses[i]
	})
	s.CurrentIndex = 0
}

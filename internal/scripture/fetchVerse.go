package scripture

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

type Scripture struct {
	logger *log.Logger
}
type ChapterIds []string

type Passages struct {
	BibleId    string     `json:"bibleId"`
	BookId     string     `json:"bookId"`
	ChapterIds ChapterIds `json:"chapterIds"`
	Content    string     `json:"content"`
	Copyright  string     `json:"copyright"`
	Id         string     `json:"id"`
	OrgId      string     `json:"orgId"`
	Reference  string     `json:"reference"`
	VerseCount int        `json:"verseCount"`
}
type Data struct {
	Passages []Passages `json:"passages"`
}

type Meta struct {
	FumsToken string `json:"fumsToken"`
}

type jsonData struct {
	Data `json:"data"`
	Meta `json:"meta"`
}

const bibleID = "de4e12af7f28f599-02"

func NewScripture(logger *log.Logger) *Scripture {
	return &Scripture{logger: logger}
}

// getBibleVerse fetches the bible the verses of the day from "api Bible."
func (s *Scripture) getBibleVerse(reference string, apiKey string) (string, error) {
	encodedRef := url.QueryEscape(reference)

	apiURL := fmt.Sprintf("https://rest.api.bible/v1/bibles/%s/search?query=%s",
		bibleID, encodedRef)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		s.logger.Println("Error creating request to get Bible verse:", err)
		return "", err
	}

	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Println("Client Response Error:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// stripHTML removes HTML tags from text
func stripHTML(html string) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	text := re.ReplaceAllString(html, " ")

	// Decode HTML entities like \u003c
	text = strings.ReplaceAll(text, "\\u003c", "<")
	text = strings.ReplaceAll(text, "\\u003e", ">")
	text = strings.ReplaceAll(text, "¶", "")

	// Clean up extra whitespace
	text = strings.TrimSpace(text)

	// Replace multiple spaces with single space
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	return text
}

// extractVerseText extracts and cleans verse text from API response
func (s *Scripture) extractVerseText(jsonResponse string) (string, error) {
	var data jsonData

	err := json.Unmarshal([]byte(jsonResponse), &data)
	if err != nil {
		s.logger.Println("Error unmarshalling json:", err)
		return "", err
	}

	var cleanText strings.Builder

	for _, v := range data.Passages {
		content := v.Content

		// Clean the HTML
		cleaned := stripHTML(content)
		cleanText.WriteString(cleaned)
		cleanText.WriteString(" ")
	}

	return strings.TrimSpace(cleanText.String()), nil
}

// GetScripture will be exported to orchestrator
func GetScripture(log *log.Logger, verseOfTheDay string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("API_KEY")

	// Get the scripture through external API
	sc := NewScripture(log)
	jsonResponse, err := sc.getBibleVerse(verseOfTheDay, apiKey)
	if err != nil {
		sc.logger.Println(err)
		return "Could not fetch verse of the day", err
	}

	// Extract and clean the verse text
	cleanVerse, err := sc.extractVerseText(jsonResponse)
	if err != nil {
		sc.logger.Println(err)
		return "Could not extract verse of the day", err
	}
	return cleanVerse, nil
}

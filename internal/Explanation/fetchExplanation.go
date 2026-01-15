package Explanation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type AnthropicRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicResponse struct {
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func GetVerseExplanation(verse string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("CLAUDE_API_KEY")

	prompt := fmt.Sprintf(`You are a thoughtful Bible teacher who provides encouraging, practical explanations of Scripture.

Here is today's Bible verse:
%s

Please provide:
1. A brief explanation of what this verse means (2-3 sentences)
2. A practical way someone can apply this to their life today (2-3 sentences)
3. An encouraging closing thought (1-2 sentences)
4. Make sure you use the verse context for a better explanation

Keep the total response under 150 words. Write in a warm, conversational tone.`, verse)

	// Create request body
	reqBody := map[string]interface{}{
		"model":      "claude-3-haiku-20240307", // Changed model name
		"max_tokens": 500,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	content := result["content"].([]interface{})
	if len(content) > 0 {
		textBlock := content[0].(map[string]interface{})
		return textBlock["text"].(string), nil
	}

	return "", fmt.Errorf("no response")
}

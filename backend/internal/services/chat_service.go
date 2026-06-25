package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ChatService struct {
	apiKey string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func NewChatService(apiKey string) *ChatService {
	return &ChatService{apiKey: apiKey}
}

func (service *ChatService) Chat(ctx context.Context, messages []ChatMessage) (string, error) {
	if service.apiKey == "" {
		return service.getMockResponse(messages), nil
	}

	// Format system prompt and history
	var prompt string
	prompt += "You are a helpful and polite smart transit assistant for Addis Ababa. "
	prompt += "Use your knowledge of Addis Ababa's geography, roads, bus stops, and taxi stands (e.g. Mexico Square, Meskel Square, Bole, Piazza) to explain routes and assist commuters.\n\n"

	for _, msg := range messages {
		prompt += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
	}
	prompt += "Assistant:"

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal gemini request: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", service.apiKey)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", fmt.Errorf("create gemini http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return service.getMockResponse(messages), nil // Fallback on request error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return service.getMockResponse(messages), nil // Fallback on non-200 status code
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read gemini response body: %w", err)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(respBytes, &geminiResp); err != nil {
		return "", fmt.Errorf("unmarshal gemini response: %w", err)
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return service.getMockResponse(messages), nil
}

func (service *ChatService) getMockResponse(messages []ChatMessage) string {
	if len(messages) == 0 {
		return "Hello! I am your STOPS Addis Ababa Transit assistant. How can I help you travel today?"
	}

	lastUserMessage := ""
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserMessage = messages[i].Content
			break
		}
	}

	// simple responses for demo/offline fallback mode
	return fmt.Sprintf("[Offline/Fallback Mode] I received your message: \"%s\". I recommend taking the Bole to Piazza taxi or the Mexico Square bus stop options for Addis Ababa transit.", lastUserMessage)
}

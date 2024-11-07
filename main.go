package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Constants for API endpoints and headers
const (
	statusURL         = "https://duckduckgo.com/duckchat/v1/status"
	chatURL           = "https://duckduckgo.com/duckchat/v1/chat"
	statusHeaders     = "1"
	termsOfServiceURL = "https://duckduckgo.com/aichat/privacy-terms" // Not used in the API version
)

// Model represents the AI model used for chat
type Model string

// ModelAlias represents a user-friendly alias for the AI model
type ModelAlias string

// Define available models and their aliases
const (
	GPT4Mini Model = "gpt-4o-mini"
	Claude3  Model = "claude-3-haiku-20240307"
	Llama    Model = "meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo"
	Mixtral  Model = "mistralai/Mixtral-8x7B-Instruct-v0.1"

	GPT4MiniAlias ModelAlias = "gpt-4o-mini"
	Claude3Alias  ModelAlias = "claude-3-haiku"
	LlamaAlias    ModelAlias = "llama"
	MixtralAlias  ModelAlias = "mixtral"
)

// Map model aliases to their corresponding Model values
var modelMap = map[ModelAlias]Model{
	GPT4MiniAlias: GPT4Mini,
	Claude3Alias:  Claude3,
	LlamaAlias:    Llama,
	MixtralAlias:  Mixtral,
}

// Message represents a chat message
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

// ChatPayload represents the payload sent to the chat API
type ChatPayload struct {
	Model    Model     `json:"model"`
	Messages []Message `json:"messages"`
}

// Chat represents a chat session
type Chat struct {
	OldVqd   string
	NewVqd   string
	Model    Model
	Messages []Message
	Client   *http.Client
}

// NewChat creates a new Chat instance
func NewChat(vqd string, model Model) *Chat {
	return &Chat{
		OldVqd:   vqd,
		NewVqd:   vqd,
		Model:    model,
		Messages: []Message{},
		Client:   &http.Client{},
	}
}

// Fetch sends a chat message and returns the response
func (c *Chat) Fetch(content string) (*http.Response, error) {
	c.Messages = append(c.Messages, Message{Content: content, Role: "user"})
	payload := ChatPayload{
		Model:    c.Model,
		Messages: c.Messages,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %v", err)
	}

	req, err := http.NewRequest("POST", chatURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-vqd-4", c.NewVqd)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%d: Failed to send message. %s. Body: %s", resp.StatusCode, resp.Status, string(body))
	}

	return resp, nil
}

// FetchStream sends a chat message and returns a channel for streaming the response
func (c *Chat) FetchStream(content string) (<-chan string, error) {
	resp, err := c.Fetch(content)
	if err != nil {
		return nil, err
	}

	stream := make(chan string)
	go func() {
		defer resp.Body.Close()
		defer close(stream)

		var text strings.Builder
		scanner := bufio.NewScanner(resp.Body)

		for scanner.Scan() {
			line := scanner.Text()

			if line == "data: [DONE]" {
				break
			}

			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimPrefix(line, "data: ")
				var messageData struct {
					Message string `json:"message"`
				}
				if err := json.Unmarshal([]byte(data), &messageData); err != nil {
					log.Printf("Error unmarshaling data: %v\n", err)
					continue // Handle error appropriately
				}

				if messageData.Message != "" {
					text.WriteString(messageData.Message)
					stream <- messageData.Message
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading response body: %v\n", err) // Handle error
		}

		c.OldVqd = c.NewVqd
		c.NewVqd = resp.Header.Get("x-vqd-4")
		c.Messages = append(c.Messages, Message{Content: text.String(), Role: "assistant"})
	}()

	return stream, nil
}

// Redo resets the chat to the previous state (not used in the API version, but kept for potential future use)
func (c *Chat) Redo() {
	c.NewVqd = c.OldVqd
	if len(c.Messages) >= 2 {
		c.Messages = c.Messages[:len(c.Messages)-2]
	}
}

// InitChat initializes a new chat session
func InitChat(model ModelAlias) (*Chat, error) {
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-vqd-accept", statusHeaders)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: Failed to initialize chat. %s", resp.StatusCode, resp.Status)
	}

	vqd := resp.Header.Get("x-vqd-4")
	if vqd == "" {
		return nil, fmt.Errorf("failed to get VQD from response headers")
	}

	return NewChat(vqd, modelMap[model]), nil
}

// CORSMiddleware adds CORS headers to the response
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, User-ID, x-vqd-accept, x-vqd-4, x-vqd-5")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	// Initialize a map to store chat sessions per user (replace with a proper database in production)
	chatSessions := make(map[string]*Chat)

	router.POST("/chat/:model", func(c *gin.Context) {
		modelAlias := ModelAlias(c.Param("model"))
		if _, ok := modelMap[modelAlias]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model"})
			return
		}

		userID := c.GetHeader("User-ID") // Get User-ID from header. You'll need a way to assign and manage these.
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User-ID header is required"})
			return
		}

		chat, ok := chatSessions[userID]
		if !ok { // Create new session if none exists
			newChat, err := InitChat(modelAlias)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			chat = newChat
			chatSessions[userID] = chat
		}

		var reqBody struct {
			Message string `json:"message"`
		}
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		stream, err := chat.FetchStream(reqBody.Message)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := make([]string, 0)
		for chunk := range stream {
			response = append(response, chunk)
		}
		c.JSON(http.StatusOK, gin.H{"response": strings.Join(response, "")}) // Send entire response as one string.

	})

	router.DELETE("/chat/:model", func(c *gin.Context) {
		// modelAlias := ModelAlias(c.Param("model")) // not really *needed*, but good to potentially validate it

		userID := c.GetHeader("User-ID")
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User-ID header is required"})
			return
		}

		if _, ok := chatSessions[userID]; ok {
			delete(chatSessions, userID)
			c.JSON(http.StatusOK, gin.H{"message": "Chat session deleted"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chat session not found"})
		}
	})
	router.GET("/health", healthCheck())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

// HEALTH CHECK
func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

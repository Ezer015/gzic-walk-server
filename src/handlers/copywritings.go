package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

const CopywritingExpiry = 5 * time.Minute

func (s *Resolver) CreateCopywriting(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateCopywriting called")
	name := r.FormValue("name")
	if name == "" {
		log.Println("Name cannot be empty")
		http.Error(w, "Name cannot be empty", http.StatusBadRequest)
		return
	}
	description := r.FormValue("description")
	if description == "" {
		log.Println("Description cannot be empty")
		http.Error(w, "Description cannot be empty", http.StatusBadRequest)
		return
	}
	prompt := r.FormValue("prompt")
	if prompt == "" {
		log.Println("Prompt cannot be empty")
		http.Error(w, "Prompt cannot be empty", http.StatusBadRequest)
		return
	}

	// Register the copywriting in the cache
	copywritingID := s.Caches.CopywritingCache.Set("", CopywritingExpiry)

	// Generate the copywriting in the background
	go func() {
		// Prepare the JSON payload for the request
		jsonData, err := json.Marshal(map[string]any{
			"messages": []map[string]any{
				{
					"role":    "system",
					"content": fmt.Sprintf("%s\n%s", prompt, s.Config.LLMConfig.SystemPrompt),
				},
				{
					"role":    "user",
					"content": fmt.Sprintf("%s: %s", name, description),
				},
			},
			"stream": false,
			"model":  s.Config.LLMConfig.Model,
		})
		if err != nil {
			s.Caches.CopywritingCache.Remove(copywritingID)
			log.Println("Error marshaling JSON:", err)
			return
		}

		// Send the request to the LLM service
		req, err := http.NewRequest("POST", s.Config.LLMConfig.Url, bytes.NewBuffer(jsonData))
		if err != nil {
			s.Caches.CopywritingCache.Remove(copywritingID)
			log.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		if s.Config.LLMConfig.Token != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.Config.LLMConfig.Token))
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			s.Caches.CopywritingCache.Remove(copywritingID)
			log.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			s.Caches.CopywritingCache.Remove(copywritingID)
			log.Println("Error reading response body:", err)
			return
		}
		// Unmarshal the response JSON
		var response map[string]any
		err = json.Unmarshal(body, &response)
		if err != nil {
			s.Caches.CopywritingCache.Remove(copywritingID)
			log.Println("Error unmarshaling JSON:", err)
			return
		}
		// Get the choices from the response
		var content string
		switch s.Config.LLMConfig.Format {
		case "openai":
			// OpenAI Format: {"choices": [{"message": {"content": "Hello, how can I help you?"}}]}
			choices, ok := response["choices"].([]any)
			if !ok {
				log.Println("Error: 'choices' is not a valid []any")
				return
			}
			if len(choices) == 0 {
				log.Println("Error: 'choices' is empty")
				return
			}
			choice, ok := choices[0].(map[string]any)
			if !ok {
				log.Println("Error: 'choice' is not a valid map[string]any")
				return
			}
			message, ok := choice["message"].(map[string]any)
			if !ok {
				log.Println("Error: 'message' is not a valid map[string]any")
				return
			}
			content, ok = message["content"].(string)
			if !ok {
				log.Println("Error: 'content' is not a valid string")
				return
			}
		// Ollama Format: {"message": {"content": "Hello, how can I help you?"}}
		case "ollama":
			message, ok := response["message"].(map[string]any)
			if !ok {
				log.Println("Error: 'message' is not a valid map[string]any")
				return
			}
			content, ok = message["content"].(string)
			if !ok {
				log.Println("Error: 'content' is not a valid string")
				return
			}
		default:
			log.Println("Error: invalid format")
			return
		}

		// Update the copywriting in the cache with the response content
		s.Caches.CopywritingCache.Update(copywritingID, content, CopywritingExpiry)
		log.Println("CreateCopywriting successful")
	}()

	// Return the copywriting ID as JSON
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]any{
		"copywriting_id": copywritingID,
	})
}

func (s *Resolver) GetCopywriting(w http.ResponseWriter, r *http.Request) {
	log.Println("GetCopywriting called")
	// Get the copywriting ID from the URL
	copywritingID, err := strconv.Atoi(r.PathValue("copywriting_id"))
	if err != nil {
		log.Println("Invalid copywriting ID:", err)
		http.Error(w, "Invalid copywriting ID", http.StatusBadRequest)
		return
	}

	// Check the copywriting cache
	cachedCopywriting, found := s.Caches.CopywritingCache.Get(copywritingID)
	if !found {
		log.Println("Copywriting not found")
		http.Error(w, "Copywriting not found", http.StatusNotFound)
		return
	}
	if cachedCopywriting == "" {
		log.Println("Copywriting not ready")
		http.Error(w, "Copywriting not ready", http.StatusAccepted)
		return
	}

	// Return the copywriting as JSON
	w.WriteHeader(http.StatusOK)
	log.Println("GetCopywriting successful")
	json.NewEncoder(w).Encode(map[string]any{
		"copywriting": cachedCopywriting,
	})
}

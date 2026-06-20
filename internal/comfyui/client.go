package comfyui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 30 * time.Second},
	}
}

type queueResponse struct {
	PromptID string `json:"prompt_id"`
}

type ImageOutput struct {
	Filename  string `json:"filename"`
	Subfolder string `json:"subfolder"`
	Type      string `json:"type"`
}

type HistoryEntry struct {
	Outputs map[string]struct {
		Images []ImageOutput `json:"images"`
	} `json:"outputs"`
	Status struct {
		StatusStr string `json:"status_str"`
		Completed bool   `json:"completed"`
	} `json:"status"`
}

func (c *Client) QueuePrompt(workflow map[string]any) (string, error) {
	body := map[string]any{"prompt": workflow}
	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	resp, err := c.http.Post(c.baseURL+"/prompt", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("queue prompt: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("queue prompt %d: %s", resp.StatusCode, b)
	}
	var qr queueResponse
	return qr.PromptID, json.NewDecoder(resp.Body).Decode(&qr)
}

// GetHistory returns nil if the job is not yet in history (still queued/running).
func (c *Client) GetHistory(promptID string) (*HistoryEntry, error) {
	resp, err := c.http.Get(c.baseURL + "/history/" + promptID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var history map[string]*HistoryEntry
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, err
	}
	return history[promptID], nil
}

func (c *Client) DownloadImage(filename, subfolder, imgType string) ([]byte, error) {
	params := url.Values{}
	params.Set("filename", filename)
	params.Set("subfolder", subfolder)
	params.Set("type", imgType)
	resp, err := c.http.Get(c.baseURL + "/view?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("view %d: %s", resp.StatusCode, b)
	}
	return io.ReadAll(resp.Body)
}

func (c *Client) ListCheckpoints() ([]string, error) {
	resp, err := c.http.Get(c.baseURL + "/object_info/CheckpointLoaderSimple")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info map[string]struct {
		Input struct {
			Required map[string][][]any `json:"required"`
		} `json:"input"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	node, ok := info["CheckpointLoaderSimple"]
	if !ok {
		return nil, fmt.Errorf("CheckpointLoaderSimple not found in object_info")
	}
	names := node.Input.Required["ckpt_name"]
	if len(names) == 0 || len(names[0]) == 0 {
		return nil, nil
	}
	result := make([]string, 0, len(names[0]))
	for _, v := range names[0] {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}

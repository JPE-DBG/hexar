package comfyui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func (c *Client) ListCheckpoints() ([]string, error) {	resp, err := c.http.Get(c.baseURL + "/object_info/CheckpointLoaderSimple")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	// Navigate: info["CheckpointLoaderSimple"]["input"]["required"]["ckpt_name"][0] -> []string
	node, _ := info["CheckpointLoaderSimple"].(map[string]any)
	input, _ := node["input"].(map[string]any)
	required, _ := input["required"].(map[string]any)
	ckptName, _ := required["ckpt_name"].([]any)
	if len(ckptName) == 0 {
		return nil, nil
	}
	names, _ := ckptName[0].([]any)
	result := make([]string, 0, len(names))
	for _, v := range names {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}

// UploadImage uploads a local image file to ComfyUI's input directory.
// Returns the filename ComfyUI assigned it (use in LoadImage nodes).
func (c *Client) UploadImage(localPath string) (string, error) {
	f, err := os.Open(localPath)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", localPath, err)
	}
	defer f.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, err := w.CreateFormFile("image", filepath.Base(localPath))
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, f); err != nil {
		return "", err
	}
	w.Close()

	resp, err := c.http.Post(c.baseURL+"/upload/image", w.FormDataContentType(), body)
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload %d: %s", resp.StatusCode, b)
	}
	var result struct {
		Name string `json:"name"`
	}
	return result.Name, json.NewDecoder(resp.Body).Decode(&result)
}

// ListControlNets returns filenames of installed ControlNet models.
func (c *Client) ListControlNets() ([]string, error) {
	resp, err := c.http.Get(c.baseURL + "/object_info/ControlNetLoader")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	node, _ := info["ControlNetLoader"].(map[string]any)
	input, _ := node["input"].(map[string]any)
	required, _ := input["required"].(map[string]any)
	cnName, _ := required["control_net_name"].([]any)
	if len(cnName) == 0 {
		return nil, nil
	}
	names, _ := cnName[0].([]any)
	result := make([]string, 0, len(names))
	for _, v := range names {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result, nil
}

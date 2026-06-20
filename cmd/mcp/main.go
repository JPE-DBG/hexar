package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"hexar/internal/comfyui"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	comfyURL := flag.String("comfy", "http://127.0.0.1:8188", "ComfyUI base URL")
	projectRoot := flag.String("root", ".", "project root for saving assets")
	flag.Parse()

	client := comfyui.New(*comfyURL)
	root := *projectRoot

	s := server.NewMCPServer("hexar-comfyui", "1.0.0",
		server.WithToolCapabilities(false),
	)

	// ── comfyui_models ────────────────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("comfyui_models",
			mcp.WithDescription("List checkpoint models installed in ComfyUI. Call this first to know which checkpoint names are valid."),
		),
		func(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			models, err := client.ListCheckpoints()
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("list models: %v", err)), nil
			}
			result, _ := json.Marshal(map[string]any{"checkpoints": models})
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// ── comfyui_generate ──────────────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("comfyui_generate",
			mcp.WithDescription("Queue a text-to-image generation in ComfyUI. Returns a job_id; use comfyui_wait or comfyui_status to poll."),
			mcp.WithString("prompt", mcp.Required(), mcp.Description("Positive prompt describing the image")),
			mcp.WithString("negative_prompt", mcp.Description("Negative prompt (what to avoid)")),
			mcp.WithString("output_name", mcp.Required(), mcp.Description("Filename prefix for the output, e.g. 'hex-tile' or 'soldier-unit'")),
			mcp.WithString("checkpoint", mcp.Required(), mcp.Description("Checkpoint model filename from comfyui_models, e.g. 'v1-5-pruned-emaonly.safetensors'")),
			mcp.WithNumber("width", mcp.Description("Image width in pixels (default 512)")),
			mcp.WithNumber("height", mcp.Description("Image height in pixels (default 512)")),
			mcp.WithNumber("steps", mcp.Description("Sampling steps (default 20)")),
			mcp.WithNumber("cfg", mcp.Description("CFG scale (default 7.0)")),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			params := comfyui.Txt2ImgParams{
				Prompt:         strArg(args, "prompt", ""),
				NegativePrompt: strArg(args, "negative_prompt", ""),
				OutputPrefix:   strArg(args, "output_name", "hexar"),
				Checkpoint:     strArg(args, "checkpoint", ""),
				Width:          int(floatArg(args, "width", 512)),
				Height:         int(floatArg(args, "height", 512)),
				Steps:          int(floatArg(args, "steps", 20)),
				CFG:            floatArg(args, "cfg", 7.0),
			}
			jobID, err := client.QueuePrompt(comfyui.BasicTxt2Img(params))
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("queue prompt: %v", err)), nil
			}
			result, _ := json.Marshal(map[string]any{
				"job_id":  jobID,
				"message": "queued — use comfyui_wait or comfyui_status to poll",
			})
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// ── comfyui_status ────────────────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("comfyui_status",
			mcp.WithDescription("Check the status of a ComfyUI generation job. Returns 'pending', a status string, or a list of output images when complete."),
			mcp.WithString("job_id", mcp.Required(), mcp.Description("Job ID returned by comfyui_generate")),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			jobID := strArg(args, "job_id", "")
			entry, err := client.GetHistory(jobID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("history lookup: %v", err)), nil
			}
			if entry == nil {
				return mcp.NewToolResultText(`{"status":"pending"}`), nil
			}
			if !entry.Status.Completed {
				result, _ := json.Marshal(map[string]any{"status": entry.Status.StatusStr})
				return mcp.NewToolResultText(string(result)), nil
			}
			images := collectImages(entry)
			result, _ := json.Marshal(map[string]any{"status": "completed", "images": images})
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// ── comfyui_wait ─────────────────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("comfyui_wait",
			mcp.WithDescription("Block until a ComfyUI job completes, polling every 2 seconds. Returns output image list. Prefer this over manual polling."),
			mcp.WithString("job_id", mcp.Required(), mcp.Description("Job ID returned by comfyui_generate")),
			mcp.WithNumber("timeout_seconds", mcp.Description("Max seconds to wait (default 120)")),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			jobID := strArg(args, "job_id", "")
			timeout := time.Duration(floatArg(args, "timeout_seconds", 120)) * time.Second
			deadline := time.Now().Add(timeout)

			for time.Now().Before(deadline) {
				entry, err := client.GetHistory(jobID)
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("history lookup: %v", err)), nil
				}
				if entry != nil && entry.Status.Completed {
					images := collectImages(entry)
					result, _ := json.Marshal(map[string]any{"status": "completed", "images": images})
					return mcp.NewToolResultText(string(result)), nil
				}
				time.Sleep(2 * time.Second)
			}
			result, _ := json.Marshal(map[string]any{"status": "timeout", "job_id": jobID})
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	// ── comfyui_save ─────────────────────────────────────────────────────────
	s.AddTool(
		mcp.NewTool("comfyui_save",
			mcp.WithDescription("Download a completed ComfyUI image and save it into the project. Use image info from comfyui_wait or comfyui_status."),
			mcp.WithString("filename", mcp.Required(), mcp.Description("Filename from the images list")),
			mcp.WithString("subfolder", mcp.Description("Subfolder from the images list (may be empty)")),
			mcp.WithString("type", mcp.Description("Image type from the images list (usually 'output')")),
			mcp.WithString("asset_path", mcp.Required(), mcp.Description("Relative path within project to save to, e.g. 'client/public/assets/hex-tile.png'")),
		),
		func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			args := req.GetArguments()
			filename := strArg(args, "filename", "")
			subfolder := strArg(args, "subfolder", "")
			imgType := strArg(args, "type", "output")
			assetPath := strArg(args, "asset_path", "")

			data, err := client.DownloadImage(filename, subfolder, imgType)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("download: %v", err)), nil
			}
			dest := filepath.Join(root, filepath.FromSlash(assetPath))
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("mkdir: %v", err)), nil
			}
			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("write: %v", err)), nil
			}
			result, _ := json.Marshal(map[string]any{"saved": dest, "bytes": len(data)})
			return mcp.NewToolResultText(string(result)), nil
		},
	)

	log.Printf("hexar-comfyui MCP server starting (comfy=%s, root=%s)", *comfyURL, root)
	if err := server.ServeStdio(s); err != nil && !strings.Contains(err.Error(), "EOF") {
		log.Fatal(err)
	}
}

// ── helpers ───────────────────────────────────────────────────────────────────

func strArg(args map[string]any, key, def string) string {
	if v, ok := args[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}

func floatArg(args map[string]any, key string, def float64) float64 {
	if v, ok := args[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return def
}

func collectImages(entry *comfyui.HistoryEntry) []comfyui.ImageOutput {
	var images []comfyui.ImageOutput
	for _, output := range entry.Outputs {
		images = append(images, output.Images...)
	}
	return images
}

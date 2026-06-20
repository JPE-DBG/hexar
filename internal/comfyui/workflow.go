package comfyui

import "math/rand"

// Txt2ImgParams configures a basic text-to-image generation.
type Txt2ImgParams struct {
	Checkpoint     string
	Prompt         string
	NegativePrompt string
	Width          int
	Height         int
	Steps          int
	CFG            float64
	Seed           int64
	OutputPrefix   string
}

// BasicTxt2Img builds a minimal KSampler workflow for text-to-image generation.
// Nodes 3–9 are the standard ComfyUI default graph layout.
func BasicTxt2Img(p Txt2ImgParams) map[string]any {
	if p.Steps == 0 {
		p.Steps = 20
	}
	if p.CFG == 0 {
		p.CFG = 7.0
	}
	if p.Width == 0 {
		p.Width = 512
	}
	if p.Height == 0 {
		p.Height = 512
	}
	if p.Seed == 0 {
		p.Seed = rand.Int63()
	}
	if p.OutputPrefix == "" {
		p.OutputPrefix = "hexar"
	}

	return map[string]any{
		"4": map[string]any{
			"class_type": "CheckpointLoaderSimple",
			"inputs":     map[string]any{"ckpt_name": p.Checkpoint},
		},
		"5": map[string]any{
			"class_type": "EmptyLatentImage",
			"inputs":     map[string]any{"batch_size": 1, "width": p.Width, "height": p.Height},
		},
		"6": map[string]any{
			"class_type": "CLIPTextEncode",
			"inputs":     map[string]any{"clip": []any{"4", 1}, "text": p.Prompt},
		},
		"7": map[string]any{
			"class_type": "CLIPTextEncode",
			"inputs":     map[string]any{"clip": []any{"4", 1}, "text": p.NegativePrompt},
		},
		"3": map[string]any{
			"class_type": "KSampler",
			"inputs": map[string]any{
				"cfg":          p.CFG,
				"denoise":      1.0,
				"latent_image": []any{"5", 0},
				"model":        []any{"4", 0},
				"negative":     []any{"7", 0},
				"positive":     []any{"6", 0},
				"sampler_name": "euler",
				"scheduler":    "normal",
				"seed":         p.Seed,
				"steps":        p.Steps,
			},
		},
		"8": map[string]any{
			"class_type": "VAEDecode",
			"inputs":     map[string]any{"samples": []any{"3", 0}, "vae": []any{"4", 2}},
		},
		"9": map[string]any{
			"class_type": "SaveImage",
			"inputs":     map[string]any{"filename_prefix": p.OutputPrefix, "images": []any{"8", 0}},
		},
	}
}

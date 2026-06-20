package comfyui

import "math/rand"

// CompositeParams configures an image-on-image composite operation.
type CompositeParams struct {
	BaseImageName    string // uploaded name of the background image
	OverlayImageName string // uploaded name of the image to place on top
	MaskImageName    string // uploaded name of a mask PNG (optional; white=show overlay)
	OutputPrefix     string
	X                int  // x pixel offset for overlay
	Y                int  // y pixel offset for overlay
	ResizeSource     bool // resize overlay to match base dimensions
}

// CompositeImages builds a workflow that places one image on top of another.
// If MaskImageName is set, its red channel is used as the overlay's alpha mask.
func CompositeImages(p CompositeParams) map[string]any {
	if p.OutputPrefix == "" {
		p.OutputPrefix = "hexar-composite"
	}

	wf := map[string]any{
		"1": map[string]any{
			"class_type": "LoadImage",
			"inputs":     map[string]any{"image": p.BaseImageName, "upload": "image"},
		},
		"2": map[string]any{
			"class_type": "LoadImage",
			"inputs":     map[string]any{"image": p.OverlayImageName, "upload": "image"},
		},
	}

	if p.MaskImageName != "" {
		wf["3"] = map[string]any{
			"class_type": "LoadImage",
			"inputs":     map[string]any{"image": p.MaskImageName, "upload": "image"},
		}
		wf["4"] = map[string]any{
			"class_type": "ImageToMask",
			"inputs":     map[string]any{"image": []any{"3", 0}, "channel": "red"},
		}
		wf["5"] = map[string]any{
			"class_type": "ImageCompositeMasked",
			"inputs": map[string]any{
				"destination":   []any{"1", 0},
				"source":        []any{"2", 0},
				"mask":          []any{"4", 0},
				"x":             p.X,
				"y":             p.Y,
				"resize_source": p.ResizeSource,
			},
		}
		wf["6"] = map[string]any{
			"class_type": "SaveImage",
			"inputs":     map[string]any{"filename_prefix": p.OutputPrefix, "images": []any{"5", 0}},
		}
	} else {
		wf["3"] = map[string]any{
			"class_type": "ImageCompositeMasked",
			"inputs": map[string]any{
				"destination":   []any{"1", 0},
				"source":        []any{"2", 0},
				"x":             p.X,
				"y":             p.Y,
				"resize_source": p.ResizeSource,
			},
		}
		wf["4"] = map[string]any{
			"class_type": "SaveImage",
			"inputs":     map[string]any{"filename_prefix": p.OutputPrefix, "images": []any{"3", 0}},
		}
	}

	return wf
}

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

// ControlNetParams configures a ControlNet-guided text-to-image generation.
type ControlNetParams struct {
	Txt2ImgParams
	ControlNetModel    string  // filename in ComfyUI's controlnet folder
	ControlImageName   string  // filename as returned by Client.UploadImage
	ControlNetStrength float64 // 0.0–2.0, typically 1.0
}

// ControlNetTxt2Img builds a KSampler workflow with ControlNetApplyAdvanced.
// The control image must already be uploaded via Client.UploadImage.
func ControlNetTxt2Img(p ControlNetParams) map[string]any {
	if p.Steps == 0 {
		p.Steps = 25
	}
	if p.CFG == 0 {
		p.CFG = 7.0
	}
	if p.Width == 0 {
		p.Width = 1024
	}
	if p.Height == 0 {
		p.Height = 1024
	}
	if p.Seed == 0 {
		p.Seed = rand.Int63()
	}
	if p.OutputPrefix == "" {
		p.OutputPrefix = "hexar-cn"
	}
	if p.ControlNetStrength == 0 {
		p.ControlNetStrength = 1.0
	}

	return map[string]any{
		// 1: checkpoint
		"1": map[string]any{
			"class_type": "CheckpointLoaderSimple",
			"inputs":     map[string]any{"ckpt_name": p.Checkpoint},
		},
		// 2: positive prompt
		"2": map[string]any{
			"class_type": "CLIPTextEncode",
			"inputs":     map[string]any{"clip": []any{"1", 1}, "text": p.Prompt},
		},
		// 3: negative prompt
		"3": map[string]any{
			"class_type": "CLIPTextEncode",
			"inputs":     map[string]any{"clip": []any{"1", 1}, "text": p.NegativePrompt},
		},
		// 4: latent canvas
		"4": map[string]any{
			"class_type": "EmptyLatentImage",
			"inputs":     map[string]any{"batch_size": 1, "width": p.Width, "height": p.Height},
		},
		// 5: control image (uploaded hex mask / canny edges)
		"5": map[string]any{
			"class_type": "LoadImage",
			"inputs":     map[string]any{"image": p.ControlImageName, "upload": "image"},
		},
		// 6: controlnet model
		"6": map[string]any{
			"class_type": "ControlNetLoader",
			"inputs":     map[string]any{"control_net_name": p.ControlNetModel},
		},
		// 7: apply controlnet — modifies positive/negative conditioning
		"7": map[string]any{
			"class_type": "ControlNetApplyAdvanced",
			"inputs": map[string]any{
				"control_net":   []any{"6", 0},
				"image":         []any{"5", 0},
				"positive":      []any{"2", 0},
				"negative":      []any{"3", 0},
				"strength":      p.ControlNetStrength,
				"start_percent": 0.0,
				"end_percent":   1.0,
			},
		},
		// 8: sampler
		"8": map[string]any{
			"class_type": "KSampler",
			"inputs": map[string]any{
				"cfg":          p.CFG,
				"denoise":      1.0,
				"latent_image": []any{"4", 0},
				"model":        []any{"1", 0},
				"negative":     []any{"7", 1},
				"positive":     []any{"7", 0},
				"sampler_name": "euler",
				"scheduler":    "normal",
				"seed":         p.Seed,
				"steps":        p.Steps,
			},
		},
		// 9: decode
		"9": map[string]any{
			"class_type": "VAEDecode",
			"inputs":     map[string]any{"samples": []any{"8", 0}, "vae": []any{"1", 2}},
		},
		// 11: extract hex mask as alpha (white inside hex = 1.0, black outside = 0.0)
		"11": map[string]any{
			"class_type": "ImageToMask",
			"inputs":     map[string]any{"image": []any{"5", 0}, "channel": "red"},
		},
		// 12: composite generated image with hex mask as alpha → RGBA
		"12": map[string]any{
			"class_type": "JoinImageWithAlpha",
			"inputs":     map[string]any{"image": []any{"9", 0}, "alpha": []any{"11", 0}},
		},
		// 10: save RGBA PNG (transparent outside hex)
		"10": map[string]any{
			"class_type": "SaveImage",
			"inputs":     map[string]any{"filename_prefix": p.OutputPrefix, "images": []any{"12", 0}},
		},
	}
}

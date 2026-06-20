---
name: asset-gen
description: "Use when: generating game art assets for Hexar via ComfyUI — hex tiles, unit tokens, buildings, capitals. Chains comfyui MCP tools into a complete generation workflow. Requires ComfyUI running at localhost:8188."
type: skill
---

# Asset Generation Skill — Hexar ComfyUI Pipeline

> **APPROVAL RULE — READ FIRST**
> Everything in this skill (art direction, color palette, prompt templates, asset catalog, filenames) is a **draft scaffold, not approved content**.
> Before generating any asset:
> 1. Present the proposed prompt, style keywords, and target filename to the user
> 2. Wait for explicit approval
> 3. Only then call the ComfyUI MCP tools
>
> Never generate and save an asset speculatively. The user decides what gets created.

Drives the full asset generation chain: reads design specs → builds prompts → calls ComfyUI MCP tools → saves PNGs to `client/public/assets/`.

**Prerequisites:**
- ComfyUI running at `localhost:8188`
- Hex mask at `client/public/assets/hex-mask.png` (used for all hex-shaped assets)
- Finalized art theme at `.claude/design/asset-theme.md` (produced by `/visual-designer`) — **do not generate assets without this file**

---

## Art Direction (DRAFT — not approved, see asset-theme.md when finalized)

**Visual tone:** Dark, strategic, competitive. Not cartoon. References: Polytopia (clean flat hex style), Into the Breach (grid clarity), Clash Royale (high readability at small size).

**Color palette to reinforce in prompts:**
```
Background:     #1a1a2e  deep navy
Player 1:       #4ecdc4  teal
Player 2:       #ff6b6b  red-coral
Unclaimed hex:  #2d3748  dark slate
Capital accent: gold ring (#f7c948) on owner color
Bar/cost:       #f7c948  amber
```

**Style keywords to always include:** `flat design, top-down strategy game, crisp edges, no gradients, limited palette, game asset, hex tile`

**Style keywords to always exclude (negative prompt):** `photorealistic, 3d render, blurry, watermark, signature, text, letters, busy background, noisy`

---

## Tool Chain

Always follow this sequence:

```
1. comfyui_models          → list checkpoints, pick one
2. comfyui_list_controlnets → list controlnets, pick canny/depth for hex shape
3. comfyui_generate_controlnet  OR  comfyui_generate (for units without hex mask)
4. comfyui_wait            → block until job completes
5. comfyui_save            → download and write to client/public/assets/<name>.png
```

**For hex-shaped assets** (tiles, buildings on hexes): always use `comfyui_generate_controlnet` with the hex mask.
**For unit tokens** (soldiers, etc.): use `comfyui_generate` — circular tokens, no hex constraint.
**Batch generation:** queue all jobs first (steps 3 for each asset), then wait+save each one.

### Hex mask path
Always pass the absolute path for `control_image_path`. Derive it from project root:
`<project_root>/client/public/assets/hex-mask.png`

The project root is `c:/Dev/git/hexar` unless otherwise configured.

### Recommended settings
| Setting | Hex tiles (controlnet) | Units (txt2img) |
|---|---|---|
| width / height | 512 × 512 | 256 × 256 |
| steps | 25 | 20 |
| cfg | 7.0 | 7.0 |
| controlnet_strength | 0.85–1.0 | — |

---

## Asset Catalog

The canonical asset filenames under `client/public/assets/`:

### Hex Tiles (use controlnet + hex mask)

| Asset | Filename | Variants |
|---|---|---|
| Unclaimed hex | `hex-unclaimed-v{n}.png` | Generate 3 variants (different seeds) |
| Player 1 claimed hex | `hex-p1-v{n}.png` | 3 variants |
| Player 2 claimed hex | `hex-p2-v{n}.png` | 3 variants |
| Player 1 capital | `hex-capital-p1-v{n}.png` | 2 variants |
| Player 2 capital | `hex-capital-p2-v{n}.png` | 2 variants |

### Units (no hex mask)

| Asset | Filename | Notes |
|---|---|---|
| Soldier (generic) | `soldier-v{n}.png` | 3 variants; renderer tints by player color |

### Buildings (future — when shop is implemented)

| Asset | Filename |
|---|---|
| Tower | `building-tower-v{n}.png` |
| Spawner | `building-spawner-v{n}.png` |
| Bar speed boost | `building-bar-boost-v{n}.png` |
| Aside slot | `building-aside-slot-v{n}.png` |

---

## Prompt Templates

Use these as starting points. Always include the style keywords above.

### Unclaimed hex tile
```
Prompt: top-down hex tile, dark slate terrain, stone texture, unclaimed territory,
        flat design, strategy game asset, hex shape, muted dark colors #2d3748,
        clean edges, no background
Negative: player colors, teal, red, gold, bright, photorealistic, text, watermark
```

### Player 1 hex tile (teal)
```
Prompt: top-down hex tile, teal territory #4ecdc4, claimed land, flat design,
        strategy game asset, glowing teal border, dark interior, hex shape,
        top-down view, clean crisp edges
Negative: red, orange, photorealistic, 3d, text, watermark, noisy
```

### Player 2 hex tile (red)
```
Prompt: top-down hex tile, red-coral territory #ff6b6b, claimed land, flat design,
        strategy game asset, glowing red border, dark interior, hex shape,
        top-down view, clean crisp edges
Negative: teal, blue, photorealistic, 3d, text, watermark, noisy
```

### Capital hex (P1)
```
Prompt: top-down hex tile, teal capital city #4ecdc4, fortress hex, gold ring border
        #f7c948, glowing amber accent, strategic importance, flat design,
        strategy game asset, top-down view, throne room tile
Negative: red, photorealistic, 3d, text, busy
```

### Capital hex (P2)
```
Prompt: top-down hex tile, red capital city #ff6b6b, fortress hex, gold ring border
        #f7c948, glowing amber accent, strategic importance, flat design,
        strategy game asset, top-down view, throne room tile
Negative: teal, photorealistic, 3d, text, busy
```

### Soldier unit token
```
Prompt: small circular game token, soldier icon, top-down view, flat design,
        simple silhouette, strategy game unit, clean icon, dark outline,
        256x256, no background, game sprite
Negative: hex shape, background color, photorealistic, complex scene, text
```

---

## Transparency Rule

The `ControlNetTxt2Img` workflow in `internal/comfyui/workflow.go` automatically applies the hex mask as an alpha channel (node 11: ImageToMask → node 12: JoinImageWithAlpha). All hex tile PNGs saved via `comfyui_save` will have transparent backgrounds outside the hex shape. **No post-processing needed.**

For units (txt2img path), the model must generate on a transparent or solid-color background. If the background is solid, note it for manual removal — or use a ControlNet inpaint approach.

---

## Workflow Example (single asset)

```
1. comfyui_models → e.g. "dreamshaper_8.safetensors"
2. comfyui_list_controlnets → e.g. "control_v11p_sd15_canny.pth"
3. comfyui_generate_controlnet(
     prompt="top-down hex tile, teal territory...",
     negative_prompt="photorealistic, 3d, text...",
     output_name="hex-p1",
     checkpoint="dreamshaper_8.safetensors",
     controlnet_model="control_v11p_sd15_canny.pth",
     control_image_path="c:/Dev/git/hexar/client/public/assets/hex-mask.png",
     width=512, height=512, steps=25, controlnet_strength=0.9
   ) → job_id
4. comfyui_wait(job_id) → {images: [{filename, subfolder, type}]}
5. comfyui_save(filename, subfolder, type,
     asset_path="client/public/assets/hex-p1-v1.png")
```

## Batch workflow (all hex tiles)

Queue all jobs first (step 3 for each), collect job IDs, then wait+save in sequence. This minimizes total wall time since ComfyUI queues jobs internally.

---

## After Generation

> **NEVER read, display, or analyze generated PNG files with the Read tool.**
> Images are large binaries; base64-encoding them into the request body causes 413 errors.
> Trust `comfyui_save` success as confirmation the file was written correctly.

- Report which asset filenames were written so the renderer can reference them
- Trust the hex mask workflow for transparency — no need to open the file to verify
- If a variant looks wrong, regenerate with a different seed (omit seed param to randomize)
- Use `Glob` to confirm a file exists if needed — never `Read` a PNG

# Hexar — Asset Theme
<!-- Written by /visual-designer. Only confirmed decisions are recorded here. -->

## Status key
- `approved` — user explicitly confirmed in conversation
- `draft` — working direction, not yet confirmed

---

## Art Style — `approved`

Medieval fantasy illustration, board game token aesthetic. Casual deck-building game feel (not gritty/realistic). Cards are the primary visual element; hex map is the battlefield backdrop.

**ComfyUI style keywords:** `medieval fantasy illustration, board game art style, flat even lighting, muted colors, clean shapes, subtle inner shading, simple terrain texture`

**Negative keywords (all generations):** `photorealistic, 3d render, anime, watermark, text, letters, blurry, heavy ornaments, complex busy patterns, 3d background, noisy, oversaturated, shadow`

**Reference games:** Dominion (card aesthetic), Polytopia (readable hex map), board game token style

---

## Overall Theme Direction — `approved`

**Direction A — Illustrated Parchment (Dominion-style)**

Warm cream/parchment card backgrounds, hand-drawn illustration feel, medieval fantasy. Cards look like physical objects — the primary visual hero. Hex map is earthy terrain (grass, stone, dirt) serving as the battle backdrop. Cozy and tactile.

**ComfyUI scene keywords:** `top-down hex strategy card game scene, medieval fantasy, warm parchment colors, illustrated board game style, hex map terrain with grass and stone, game cards with ornate borders visible on screen, cozy tavern lighting, game UI concept art`

**Preview generated:** `client/public/assets/previews/preview-A-parchment.png`

---

## Color Palette — `approved`

**Scope:** CSS-only variables for UI chrome and dynamic elements. Hex tiles, card art, unit tokens, buildings are generated image assets — their colors come from the art, not CSS.

```css
:root {
  /* Screen background — dark walnut table */
  --bg:           #2c1e0f;

  /* UI panels (HUD, deck panel, shop panel) */
  --chrome:       #f0e6d0;   /* warm parchment */
  --chrome-dark:  #e0d0b8;   /* panel subdivisions */
  --ui-border:    #b8956a;   /* warm sienna */
  --ui-accent:    #c8a020;   /* gold — card frame borders, decorative accents */

  /* Text on parchment panels */
  --text-primary: #2c1810;   /* dark brown ink */
  --text-muted:   #8a6840;   /* medium brown */

  /* Bar meter */
  --bar-fill:     #d4860a;   /* deep amber */
  --bar-empty:    #c4aa88;   /* pale tan */

  /* Cost badge on cards */
  --cost:         #d4860a;   /* amber */

  /* Status */
  --battle:       #d4860a;   /* amber pulse ring */
  --danger:       #c03030;   /* capital low HP */
  --win:          #c8a020;   /* gold */
}
```

**Asset-driven (not CSS):**
| Role | Status | Asset file |
|---|---|---|
| Player 1 territory | approved | `client/public/assets/hex-tile-p1.png` |
| Player 2 territory | approved | `client/public/assets/hex-tile-p2.png` |
| Unclaimed hex tile | draft | `client/public/assets/hex-tile-unclaimed.png` |

**Note:** 10 hex color variants generated (`previews/hex-v3-*.png`, `previews/hex-v4-*.png`). Forest green + royal blue are P1/P2 defaults.

---

## UI Chrome Textures — `draft`

UI panels use generated texture assets layered under CSS text/controls, not flat CSS colors.

| Asset | Description | Status |
|---|---|---|
| Card frame | Warm brown border, cream parchment center, name banner at top, portrait orientation | **approved** — `client/public/assets/card-frame.png` |
| Panel background | Parchment/aged paper texture, seamless or stretched | draft |
| Bar meter trough | Stone or wood channel, landscape orientation | **approved** — `client/public/assets/bar-meter-trough.png` |
| Board frame border | Decorative parchment trim around hex canvas | draft |

**Approach:** Generated PNG placed as CSS `background-image` on DOM panels. Dynamic elements (bar fill %, button states, text) float on top via CSS/DOM.

---

## Hex Tiles — `approved`

**Generation approach:** ControlNet with hex mask (`client/public/assets/hex-mask.png`), DreamShaper XL Lightning, CFG 2.0, 6 steps, 512×512, ControlNet promax strength 0.9.

**Approved prompt pattern (v3 style):**
```
top-down hex tile, [COLOR] territory, subtle inner shading, simple terrain texture,
light surface decoration, medieval fantasy illustration, muted [COLOR] land,
clean hex shape, board game art style, soft gradient inside
```

**Negative:**
```
photorealistic, 3d render, anime, watermark, text, letters, blurry, heavy ornaments,
complex busy patterns, relief carving, embossed, 3d background, noisy, oversaturated,
border decorations
```

---

## Capital Hex — `approved` (approach), `draft` (final asset)

**Approach:** Generate building sprite separately on white background → composite onto player-colored hex tile at generation time (one capital asset per player color).

**Why separate sprite:** ControlNet can't reliably generate hex tile + building in consistent style per-generation without LoRA. Separate building + composite is more stable.

**Capital building direction (user-provided, needs prompt refinement):**
> A 2D stylized vector art isometric game asset of a medieval capital castle and fortress hub designed for a hex tile grid. Stone-masonry walls, timber-framed keep towers, iron-gated entryways. Layered slate-shingle roofs, small battlements with wooden hoardings, royal heraldic banners. Tightly contained within the hex silhouette. Clean, bold illustrative look suitable for a high-quality indie game.

**Issue:** Current prompt generates graphics that are too simple. Needs more detail/richness while keeping the isometric game asset style.

**Good reference variants generated:** `capital-var-01.png` (round keep), `capital-var-07.png` (wizard tower) — strong silhouette, readable at hex scale.

**Sprite generation settings:** DreamShaper XL Lightning, CFG 3.0, 8 steps, 512×512, white background, no ControlNet.

---

## Other Assets — `draft`

| Asset | Description | Status |
|---|---|---|
| Unclaimed hex tile | Neutral color, no player ownership | **approved** — `client/public/assets/hex-tile-unclaimed.png` |
| Basic Soldier unit | Small token on hex, moves hex-by-hex | draft |
| Tower building | Defensive building placed on hex | draft |
| Bar Speed building | Economy building on hex | draft |
| Spawner building | Military building on hex | draft |

---

## Generation Pipeline

1. Hex tiles: ControlNet + hex mask → transparent PNG (alpha from mask)
2. Capital: building sprite (white bg) → composite over player hex tile → final capital PNG
3. Units/buildings: sprite on white bg → key out background → composite on hex at render time (Canvas renderer)

import { hexToPixel } from '../hexmath';
import { HEX_SIZE, COLORS, BUILDING_POWER, CAPITAL_POWER, FORTIFY_DURATION, TECH_IRON_GRIP, COUNTER_SPEND_CAP } from '../constants';
import { GameState, HexDTO, BattleDTO } from '../state/state';

const BUILDING_LABELS: Record<number, string> = { 1: 'G', 2: 'P', 3: 'R' };

function darken(color: string, amount: number): string {
  const n = parseInt(color.slice(1), 16);
  const r = Math.max(0, (n >> 16) - Math.round(amount * 255));
  const g = Math.max(0, ((n >> 8) & 0xff) - Math.round(amount * 255));
  const b = Math.max(0, (n & 0xff) - Math.round(amount * 255));
  return '#' + [r, g, b].map(v => v.toString(16).padStart(2, '0')).join('');
}

function lighten(color: string, amount: number): string {
  const n = parseInt(color.slice(1), 16);
  const r = Math.min(255, (n >> 16) + Math.round(amount * 255));
  const g = Math.min(255, ((n >> 8) & 0xff) + Math.round(amount * 255));
  const b = Math.min(255, (n & 0xff) + Math.round(amount * 255));
  return '#' + [r, g, b].map(v => v.toString(16).padStart(2, '0')).join('');
}

interface Effect {
  type: 'flash' | 'floater';
  q: number;
  r: number;
  color: string;
  text?: string;
  startTime: number;
  duration: number;
}

export class Renderer {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private offsetX = 0;
  private offsetY = 0;
  private dpr = 1;
  private selectedHex: { q: number; r: number } | null = null;
  private state: GameState | null = null;
  private dropMap = new Set<string>();
  private effects: Effect[] = [];
  private lastFloaterTime = new Map<string, number>();

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d')!;
    this.resize();
    window.addEventListener('resize', () => this.resize());
    this.startLoop();
  }

  private startLoop() {
    const loop = () => {
      if (this.state) this.render(this.state);
      requestAnimationFrame(loop);
    };
    requestAnimationFrame(loop);
  }

  private resize() {
    this.dpr = window.devicePixelRatio || 1;
    const w = window.innerWidth;
    const h = window.innerHeight;
    this.canvas.width = w * this.dpr;
    this.canvas.height = h * this.dpr;
    this.canvas.style.width = w + 'px';
    this.canvas.style.height = h + 'px';
    this.ctx.scale(this.dpr, this.dpr);
    this.offsetX = w / 2;
    this.offsetY = h / 2;
  }

  getOffset(): { x: number; y: number } {
    return { x: this.offsetX, y: this.offsetY };
  }

  pan(dx: number, dy: number) {
    this.offsetX += dx;
    this.offsetY += dy;
  }

  setSelected(hex: { q: number; r: number } | null) {
    this.selectedHex = hex;
  }

  setState(s: GameState) {
    this.state = s;
  }

  setDropMap(keys: Set<string>) {
    this.dropMap = keys;
  }

  addCaptureFlash(q: number, r: number, color: string) {
    this.effects.push({ type: 'flash', q, r, color, startTime: Date.now(), duration: 350 });
  }

  addFloater(q: number, r: number, text: string) {
    const key = `${q},${r}`;
    const now = Date.now();
    if (now - (this.lastFloaterTime.get(key) ?? 0) < 2000) return;
    this.lastFloaterTime.set(key, now);
    this.effects.push({ type: 'floater', q, r, color: COLORS.capital, text, startTime: now, duration: 900 });
  }

  render(state: GameState) {
    const ctx = this.ctx;
    ctx.fillStyle = COLORS.background;
    ctx.fillRect(0, 0, this.canvas.width / this.dpr, this.canvas.height / this.dpr);

    // Pass 1: fills + grid borders
    for (const [, hex] of state.hexes) {
      this.drawHexFill(hex);
    }

    // Pass 2: rings on top of all fills
    for (const [, hex] of state.hexes) {
      this.drawHexRings(hex);
    }

    for (const battle of state.battles) {
      this.drawBattle(battle);
    }

    // Pass 3: labels always on top of fills, rings, and battles
    for (const [, hex] of state.hexes) {
      this.drawHexLabel(hex);
    }

    // Effects always last — on top of everything
    this.drawEffects();
  }

  private hexPath(px: number, py: number) {
    const ctx = this.ctx;
    ctx.beginPath();
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 180) * (60 * i - 30);
      const hx = px + HEX_SIZE * Math.cos(angle);
      const hy = py + HEX_SIZE * Math.sin(angle);
      if (i === 0) ctx.moveTo(hx, hy);
      else ctx.lineTo(hx, hy);
    }
    ctx.closePath();
  }

  private drawHexFill(hex: HexDTO) {
    const ctx = this.ctx;
    ctx.save();
    const { x, y } = hexToPixel({ q: hex.q, r: hex.r });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    this.hexPath(px, py);

    const baseColor = this.hexColor(hex);
    if (hex.owner > 0) {
      const grad = ctx.createRadialGradient(px, py, 0, px, py, HEX_SIZE * 0.85);
      grad.addColorStop(0, lighten(baseColor, hex.capital ? 0.12 : 0.22));
      grad.addColorStop(1, darken(baseColor, hex.capital ? 0.45 : 0.35));
      ctx.fillStyle = grad;
      ctx.fill();
    } else {
      const grad = ctx.createRadialGradient(px, py, 0, px, py, HEX_SIZE * 0.85);
      grad.addColorStop(0, lighten(baseColor, 0.04));
      grad.addColorStop(1, baseColor);
      ctx.fillStyle = grad;
      ctx.fill();
    }
    ctx.strokeStyle = COLORS.grid;
    ctx.lineWidth = 1;
    ctx.stroke();
    ctx.restore();
  }

  private drawHexLabel(hex: HexDTO) {
    const ctx = this.ctx;
    ctx.save();
    const { x, y } = hexToPixel({ q: hex.q, r: hex.r });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    if (hex.capital && hex.building === 0) {
      ctx.beginPath();
      ctx.moveTo(px, py - 8);
      ctx.lineTo(px + 5, py - 3);
      ctx.lineTo(px, py + 2);
      ctx.lineTo(px - 5, py - 3);
      ctx.closePath();
      ctx.fillStyle = COLORS.capital;
      ctx.fill();
    }

    if (hex.building > 0) {
      const label = BUILDING_LABELS[hex.building] ?? '?';
      ctx.font = 'bold 12px monospace';
      ctx.fillStyle = '#ffffff';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      if (hex.building === BUILDING_POWER && hex.owner > 0 && this.state) {
        const ownerPlayer = this.state.players.get(String(hex.owner));
        const capitalBonus = hex.capital ? CAPITAL_POWER : 0;
        const effectivePower = capitalBonus + hex.level + (ownerPlayer?.tech?.[TECH_IRON_GRIP] ? 1 : 0);
        ctx.fillText(`${label}${effectivePower}`, px, py);
      } else {
        ctx.fillText(`${label}${hex.level}`, px, py);
      }
    }

    if (hex.owner > 0 && hex.building !== BUILDING_POWER && this.state) {
      const ownerPlayer = this.state.players.get(String(hex.owner));
      let powerBadge = 0;
      if (hex.capital) powerBadge = CAPITAL_POWER;
      if (ownerPlayer?.tech?.[TECH_IRON_GRIP]) powerBadge++;
      if (powerBadge > 0) {
        ctx.font = 'bold 9px monospace';
        ctx.fillStyle = '#ffdd44';
        ctx.textAlign = 'center';
        ctx.textBaseline = 'bottom';
        ctx.fillText(`${powerBadge}`, px, py + HEX_SIZE * 0.45);
      }
    }

    ctx.restore();
  }

  private drawHexRings(hex: HexDTO) {
    const ctx = this.ctx;
    ctx.save();
    const { x, y } = hexToPixel({ q: hex.q, r: hex.r });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    if (hex.capital && hex.owner !== 0) {
      this.hexPath(px, py);
      ctx.strokeStyle = COLORS.capital;
      ctx.lineWidth = 3;
      ctx.stroke();
    }

    const isSelected = this.selectedHex && this.selectedHex.q === hex.q && this.selectedHex.r === hex.r;
    if (isSelected) {
      this.hexPath(px, py);
      ctx.strokeStyle = '#ffffff';
      ctx.lineWidth = 2.5;
      ctx.stroke();
    }

    // Red pulse for drop-map hexes (crisis auto-drop candidates)
    const key = `${hex.q},${hex.r}`;
    if (this.dropMap.has(key)) {
      const pulse = 0.5 + 0.5 * Math.sin(Date.now() / 300);
      this.hexPath(px, py);
      ctx.strokeStyle = `rgba(255, 60, 60, ${pulse})`;
      ctx.lineWidth = 3;
      ctx.stroke();
    }

    // Fortify timer: shrinking lime border segments (clockwise from top)
    if (hex.fortifyTimer > 0) {
      const hasBattle = this.state?.battles.some(b => b.dq === hex.q && b.dr === hex.r) ?? false;
      this.drawFortifySegments(px, py, hex.fortifyTimer / FORTIFY_DURATION, hasBattle);
    }
    ctx.restore();
  }

  private drawFortifySegments(px: number, py: number, fraction: number, dimmed = false) {
    if (fraction <= 0) return;
    const ctx = this.ctx;
    ctx.save();

    const corners: { x: number; y: number }[] = [];
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 180) * (60 * i - 30);
      corners.push({ x: px + HEX_SIZE * Math.cos(angle), y: py + HEX_SIZE * Math.sin(angle) });
    }

    const cwOrder = [5, 0, 1, 2, 3, 4];
    const clamped = Math.min(fraction, 1);
    const totalSides = clamped * 6;
    const fullSides = Math.floor(totalSides);
    const partial = totalSides - fullSides;

    ctx.strokeStyle = '#c8ff70';
    ctx.lineWidth = dimmed ? 1.5 : 3;
    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
    ctx.globalAlpha = dimmed ? 0.4 : 1.0;

    ctx.beginPath();
    ctx.moveTo(corners[cwOrder[0]].x, corners[cwOrder[0]].y);
    for (let i = 0; i < fullSides && i < 6; i++) {
      const to = corners[cwOrder[(i + 1) % 6]];
      ctx.lineTo(to.x, to.y);
    }
    if (partial > 0 && fullSides < 6) {
      const from = corners[cwOrder[fullSides]];
      const to = corners[cwOrder[(fullSides + 1) % 6]];
      ctx.lineTo(
        from.x + (to.x - from.x) * partial,
        from.y + (to.y - from.y) * partial,
      );
    }
    ctx.stroke();

    ctx.restore();
  }

  private drawBattle(battle: BattleDTO) {
    const ctx = this.ctx;
    ctx.save();
    const { x, y } = hexToPixel({ q: battle.dq, r: battle.dr });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    const cap = Math.min(COUNTER_SPEND_CAP, Math.floor(battle.timeLeft));
    const canBoost = battle.counterBoost < cap;

    if (canBoost) {
      const pulse = 0.5 + 0.5 * Math.sin(Date.now() / 200);
      this.hexPath(px, py);
      ctx.strokeStyle = `rgba(240, 147, 43, ${pulse})`;
      ctx.lineWidth = 3;
      ctx.stroke();
    } else {
      this.hexPath(px, py);
      ctx.strokeStyle = 'rgba(240, 147, 43, 0.25)';
      ctx.lineWidth = 2;
      ctx.stroke();
    }

    // Timer
    ctx.font = 'bold 11px monospace';
    ctx.fillStyle = COLORS.capital;
    ctx.textAlign = 'center';
    ctx.textBaseline = 'bottom';
    ctx.fillText(`${battle.timeLeft.toFixed(1)}s`, px, py - 12);

    // 3 boost dots
    const dotR = 4, spacing = 11;
    for (let i = 0; i < 3; i++) {
      const dx = px + (i - 1) * spacing;
      const dy = py + 16;
      ctx.beginPath();
      ctx.arc(dx, dy, dotR, 0, Math.PI * 2);
      if (i < battle.counterBoost) {
        ctx.fillStyle = COLORS.battle;
        ctx.fill();
      } else if (i < cap) {
        ctx.strokeStyle = COLORS.battle;
        ctx.lineWidth = 1.5;
        ctx.stroke();
      } else {
        ctx.strokeStyle = '#555';
        ctx.lineWidth = 1.5;
        ctx.stroke();
      }
    }
    ctx.restore();
  }

  private drawEffects() {
    const ctx = this.ctx;
    const now = Date.now();
    this.effects = this.effects.filter(e => now - e.startTime < e.duration);

    for (const e of this.effects) {
      const t = (now - e.startTime) / e.duration;
      const { x, y } = hexToPixel({ q: e.q, r: e.r });
      const px = x + this.offsetX;
      const py = y + this.offsetY;

      if (e.type === 'flash') {
        const alpha = (1 - t) * 0.55;
        const alphaHex = Math.round(alpha * 255).toString(16).padStart(2, '0');
        ctx.save();
        this.hexPath(px, py);
        ctx.fillStyle = e.color + alphaHex;
        ctx.fill();
        ctx.restore();
      } else if (e.type === 'floater' && e.text) {
        const alpha = 1 - t;
        const dy = -22 * t;
        ctx.save();
        ctx.globalAlpha = alpha;
        ctx.font = 'bold 11px monospace';
        ctx.fillStyle = e.color;
        ctx.textAlign = 'center';
        ctx.textBaseline = 'middle';
        ctx.fillText(e.text, px, py + dy);
        ctx.restore();
      }
    }
  }

  private hexColor(hex: HexDTO): string {
    switch (hex.owner) {
      case 1: return COLORS.player1;
      case 2: return COLORS.player2;
      default: return COLORS.unclaimed;
    }
  }
}

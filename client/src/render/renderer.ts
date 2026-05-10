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

export class Renderer {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private offsetX = 0;
  private offsetY = 0;
  private selectedHex: { q: number; r: number } | null = null;
  private state: GameState | null = null;
  private dropMap = new Set<string>();

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
    this.canvas.width = window.innerWidth;
    this.canvas.height = window.innerHeight;
    this.offsetX = this.canvas.width / 2;
    this.offsetY = this.canvas.height / 2;
  }

  getOffset(): { x: number; y: number } {
    return { x: this.offsetX, y: this.offsetY };
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

  render(state: GameState) {
    const ctx = this.ctx;
    ctx.fillStyle = COLORS.background;
    ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

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
    const { x, y } = hexToPixel({ q: hex.q, r: hex.r });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    this.hexPath(px, py);
    ctx.fillStyle = hex.capital ? darken(this.hexColor(hex), 0.25) : this.hexColor(hex);
    ctx.fill();
    ctx.strokeStyle = COLORS.grid;
    ctx.lineWidth = 1;
    ctx.stroke();

    if (hex.building > 0) {
      const label = BUILDING_LABELS[hex.building] ?? '?';
      ctx.font = 'bold 12px monospace';
      ctx.fillStyle = '#ffffff';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';

      if (hex.building === BUILDING_POWER && hex.owner > 0 && this.state) {
        // Show effective power: capital innate + building level + Iron Grip
        const ownerPlayer = this.state.players.get(String(hex.owner));
        const capitalBonus = hex.capital ? CAPITAL_POWER : 0;
        const effectivePower = capitalBonus + hex.level + (ownerPlayer?.tech?.[TECH_IRON_GRIP] ? 1 : 0);
        ctx.fillText(`${label}${effectivePower}`, px, py);
      } else {
        ctx.fillText(`${label}${hex.level}`, px, py);
      }
    }

    // Small yellow power badge for non-Power owned hexes that have power > 0
    // (capital innate power, or Iron Grip on economy/research/empty hexes)
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
  }

  private drawHexRings(hex: HexDTO) {
    const ctx = this.ctx;
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
    // Dimmed when a battle is active on this hex — battle takes visual priority
    if (hex.fortifyTimer > 0) {
      const hasBattle = this.state?.battles.some(b => b.dq === hex.q && b.dr === hex.r) ?? false;
      this.drawFortifySegments(px, py, hex.fortifyTimer / FORTIFY_DURATION, hasBattle);
    }
  }

  private drawFortifySegments(px: number, py: number, fraction: number, dimmed = false) {
    if (fraction <= 0) return;
    const ctx = this.ctx;

    // 6 corners: angle = (60*i - 30)° for i=0..5
    // i=0: -30° upper-right, i=5: 270° top
    const corners: { x: number; y: number }[] = [];
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 180) * (60 * i - 30);
      corners.push({ x: px + HEX_SIZE * Math.cos(angle), y: py + HEX_SIZE * Math.sin(angle) });
    }

    // Clockwise from top (i=5): 5→0→1→2→3→4
    const cwOrder = [5, 0, 1, 2, 3, 4];
    const totalSides = fraction * 6;
    const fullSides = Math.floor(totalSides);
    const partial = totalSides - fullSides;

    ctx.strokeStyle = '#c8ff70'; // Bright lime — readable on both P1 teal and P2 red
    ctx.lineWidth = dimmed ? 1.5 : 3;
    ctx.lineCap = 'round';
    if (dimmed) ctx.globalAlpha = 0.4;

    for (let i = 0; i <= fullSides && i < 6; i++) {
      const from = corners[cwOrder[i]];
      const to = corners[cwOrder[(i + 1) % 6]];
      ctx.beginPath();
      ctx.moveTo(from.x, from.y);
      if (i < fullSides) {
        ctx.lineTo(to.x, to.y);
      } else if (partial > 0) {
        ctx.lineTo(from.x + (to.x - from.x) * partial, from.y + (to.y - from.y) * partial);
      }
      ctx.stroke();
    }

    if (dimmed) ctx.globalAlpha = 1.0;
  }

  private drawBattle(battle: BattleDTO) {
    const ctx = this.ctx;
    const { x, y } = hexToPixel({ q: battle.dq, r: battle.dr });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    const cap = Math.min(COUNTER_SPEND_CAP, Math.floor(battle.timeLeft));
    const canBoost = battle.counterBoost < cap;

    // Amber pulse only while counter-spend is still actionable; dim static ring when capped
    if (canBoost) {
      const pulse = 0.5 + 0.5 * Math.sin(Date.now() / 200);
      this.hexPath(px, py);
      ctx.strokeStyle = `rgba(255, 200, 0, ${pulse})`;
      ctx.lineWidth = 3;
      ctx.stroke();
    } else {
      this.hexPath(px, py);
      ctx.strokeStyle = 'rgba(255, 200, 0, 0.25)';
      ctx.lineWidth = 2;
      ctx.stroke();
    }

    // Timer
    ctx.font = 'bold 11px monospace';
    ctx.fillStyle = '#ffd93d';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'bottom';
    ctx.fillText(`${battle.timeLeft.toFixed(1)}s`, px, py - 12);

    // 3 boost dots: filled amber = used, outline amber = available, gray = time-capped
    const dotR = 4, spacing = 11;
    for (let i = 0; i < 3; i++) {
      const dx = px + (i - 1) * spacing;
      const dy = py + 16;
      ctx.beginPath();
      ctx.arc(dx, dy, dotR, 0, Math.PI * 2);
      if (i < battle.counterBoost) {
        ctx.fillStyle = '#ffd93d';
        ctx.fill();
      } else if (i < cap) {
        ctx.strokeStyle = '#ffd93d';
        ctx.lineWidth = 1.5;
        ctx.stroke();
      } else {
        ctx.strokeStyle = '#555';
        ctx.lineWidth = 1.5;
        ctx.stroke();
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

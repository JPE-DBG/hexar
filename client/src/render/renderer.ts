import { hexToPixel } from '../hexmath';
import { HEX_SIZE, COLORS } from '../constants';
import { GameState, HexDTO, BattleDTO } from '../state/state';

const BUILDING_LABELS: Record<number, string> = { 1: 'E', 2: 'D', 3: 'R' };

export class Renderer {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private offsetX = 0;
  private offsetY = 0;
  private selectedHex: { q: number; r: number } | null = null;

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d')!;
    this.resize();
    window.addEventListener('resize', () => this.resize());
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

  render(state: GameState) {
    const ctx = this.ctx;
    ctx.fillStyle = COLORS.background;
    ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

    for (const [, hex] of state.hexes) {
      this.drawHex(hex);
    }

    for (const battle of state.battles) {
      this.drawBattle(battle);
    }
  }

  private drawHex(hex: HexDTO) {
    const ctx = this.ctx;
    const { x, y } = hexToPixel({ q: hex.q, r: hex.r });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    ctx.beginPath();
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 180) * (60 * i - 30);
      const hx = px + HEX_SIZE * Math.cos(angle);
      const hy = py + HEX_SIZE * Math.sin(angle);
      if (i === 0) ctx.moveTo(hx, hy);
      else ctx.lineTo(hx, hy);
    }
    ctx.closePath();

    ctx.fillStyle = this.hexColor(hex);
    ctx.fill();

    const isSelected = this.selectedHex && this.selectedHex.q === hex.q && this.selectedHex.r === hex.r;
    ctx.strokeStyle = isSelected ? '#ffffff' : COLORS.grid;
    ctx.lineWidth = isSelected ? 2.5 : 1;
    ctx.stroke();

    if (hex.capital && hex.owner !== 0) {
      ctx.beginPath();
      ctx.arc(px, py, 4, 0, Math.PI * 2);
      ctx.fillStyle = COLORS.capital;
      ctx.fill();
    }

    if (hex.building > 0) {
      const label = BUILDING_LABELS[hex.building] ?? '?';
      ctx.font = 'bold 12px monospace';
      ctx.fillStyle = '#ffffff';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      const textY = hex.capital ? py + 10 : py;
      ctx.fillText(`${label}${hex.level}`, px, textY);
    }
  }

  private drawBattle(battle: BattleDTO) {
    const ctx = this.ctx;
    const { x, y } = hexToPixel({ q: battle.dq, r: battle.dr });
    const px = x + this.offsetX;
    const py = y + this.offsetY;

    const pulse = 0.5 + 0.5 * Math.sin(Date.now() / 200);
    ctx.beginPath();
    for (let i = 0; i < 6; i++) {
      const angle = (Math.PI / 180) * (60 * i - 30);
      const hx = px + HEX_SIZE * Math.cos(angle);
      const hy = py + HEX_SIZE * Math.sin(angle);
      if (i === 0) ctx.moveTo(hx, hy);
      else ctx.lineTo(hx, hy);
    }
    ctx.closePath();
    ctx.strokeStyle = `rgba(255, 200, 0, ${pulse})`;
    ctx.lineWidth = 3;
    ctx.stroke();

    ctx.font = 'bold 11px monospace';
    ctx.fillStyle = '#ffd93d';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'bottom';
    ctx.fillText(`${battle.timeLeft.toFixed(1)}s`, px, py - 12);
  }

  private hexColor(hex: HexDTO): string {
    switch (hex.owner) {
      case 1: return COLORS.player1;
      case 2: return COLORS.player2;
      default: return COLORS.unclaimed;
    }
  }
}

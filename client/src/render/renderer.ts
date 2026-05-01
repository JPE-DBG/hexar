import { hexToPixel } from '../hexmath';
import { HEX_SIZE, COLORS } from '../constants';
import { GameState, HexDTO } from '../state/state';

export class Renderer {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private offsetX = 0;
  private offsetY = 0;

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

  render(state: GameState) {
    const ctx = this.ctx;
    ctx.fillStyle = COLORS.background;
    ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

    for (const [, hex] of state.hexes) {
      this.drawHex(hex);
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
    ctx.strokeStyle = COLORS.grid;
    ctx.lineWidth = 1;
    ctx.stroke();

    if (hex.capital && hex.owner !== 0) {
      ctx.beginPath();
      ctx.arc(px, py, 5, 0, Math.PI * 2);
      ctx.fillStyle = COLORS.capital;
      ctx.fill();
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

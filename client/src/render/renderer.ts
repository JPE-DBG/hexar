import { hexToPixel } from '../hexmath.js';
import { HEX_SIZE } from '../constants.js';

export interface MockHex {
  q: number;
  r: number;
  owner: number;
  capital: boolean;
  unit?: number;
}

const SQRT3 = Math.sqrt(3);
const HEX_W = SQRT3 * HEX_SIZE;
const HEX_H = 2 * HEX_SIZE;

export class Renderer {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private dpr = 1;
  private offsetX = 0;
  private offsetY = 0;
  private imgs: Record<string, HTMLImageElement> = {};
  private imgsLoaded = 0;
  private readonly totalImgs = 4;
  private hexes: MockHex[] = [];

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d')!;
    this.loadImages();
    window.addEventListener('resize', () => this.resize());
    this.resize();
  }

  private loadImages() {
    const names = ['hex-tile-unclaimed', 'hex-tile-p1', 'hex-tile-p2', 'hex-tile-capital'];
    for (const name of names) {
      const img = new Image();
      img.src = `/assets/${name}.png`;
      img.onload = () => {
        this.imgsLoaded++;
        if (this.imgsLoaded === this.totalImgs) this.render();
      };
      this.imgs[name] = img;
    }
  }

  private resize() {
    this.dpr = window.devicePixelRatio || 1;
    const w = this.canvas.clientWidth;
    const h = this.canvas.clientHeight;
    if (w === 0 || h === 0) return;
    this.canvas.width = w * this.dpr;
    this.canvas.height = h * this.dpr;
    this.ctx.scale(this.dpr, this.dpr);
    this.offsetX = w / 2;
    this.offsetY = h / 2;
    this.render();
  }

  setHexes(hexes: MockHex[]) {
    this.hexes = hexes;
    this.render();
  }

  private render() {
    const w = this.canvas.width / this.dpr;
    const h = this.canvas.height / this.dpr;
    this.ctx.clearRect(0, 0, w, h);
    this.ctx.fillStyle = '#2c1e0f';
    this.ctx.fillRect(0, 0, w, h);

    if (this.imgsLoaded < this.totalImgs) return;

    for (const hex of this.hexes) {
      const { x, y } = hexToPixel(hex);
      const px = x + this.offsetX;
      const py = y + this.offsetY;

      const imgName = hex.capital ? 'hex-tile-capital'
        : hex.owner === 1 ? 'hex-tile-p1'
        : hex.owner === 2 ? 'hex-tile-p2'
        : 'hex-tile-unclaimed';

      // Clip to hex shape so square PNG edges don't bleed into neighbours
      this.ctx.save();
      this.ctx.beginPath();
      for (let i = 0; i < 6; i++) {
        const angle = (Math.PI / 180) * (60 * i - 30);
        const hx = px + HEX_SIZE * Math.cos(angle);
        const hy = py + HEX_SIZE * Math.sin(angle);
        if (i === 0) this.ctx.moveTo(hx, hy); else this.ctx.lineTo(hx, hy);
      }
      this.ctx.closePath();
      this.ctx.clip();
      this.ctx.translate(px, py);
      this.ctx.rotate(Math.PI / 6);
      this.ctx.drawImage(this.imgs[imgName], -HEX_H / 2, -HEX_W / 2, HEX_H, HEX_W);
      this.ctx.restore();

      if (hex.unit !== undefined) {
        const r = HEX_SIZE * 0.28;
        this.ctx.beginPath();
        this.ctx.arc(px, py, r, 0, Math.PI * 2);
        this.ctx.fillStyle = hex.unit === 1 ? 'rgba(78,205,196,0.9)' : 'rgba(255,107,107,0.9)';
        this.ctx.fill();
        this.ctx.strokeStyle = hex.unit === 1 ? '#1a5f5a' : '#7a1f1f';
        this.ctx.lineWidth = 2;
        this.ctx.stroke();
        this.ctx.fillStyle = '#fff';
        this.ctx.font = `bold ${Math.round(r * 1.2)}px sans-serif`;
        this.ctx.textAlign = 'center';
        this.ctx.textBaseline = 'middle';
        this.ctx.fillText('S', px, py);
      }
    }
  }
}

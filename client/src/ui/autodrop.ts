import { PlayerDTO } from '../state/state';

export class AutoDropPanel {
  private el: HTMLElement;

  constructor(parent: HTMLElement) {
    this.el = document.createElement('div');
    this.el.id = 'auto-drop';
    this.el.style.cssText = `
      position: fixed; top: 50px; right: 16px;
      background: rgba(80,20,20,0.95); border: 1px solid #ff4444; border-radius: 8px;
      padding: 12px; display: none; flex-direction: column; gap: 6px;
      font-family: monospace; color: #e0e0e0; font-size: 12px; min-width: 180px; z-index: 90;
    `;
    parent.appendChild(this.el);
  }

  update(player: PlayerDTO | null) {
    if (!player?.autoDropActive) {
      this.el.style.display = 'none';
      return;
    }
    this.el.innerHTML = `
      <div style="color:#ff8888;font-weight:bold">Income Negative</div>
      <div style="color:#ffcccc">Auto-drop in ${(player.autoDropGrace ?? 0).toFixed(1)}s</div>
      <div style="color:#aaa;font-size:11px">Click a red hex to drop it</div>`;
    this.el.style.display = 'flex';
  }
}

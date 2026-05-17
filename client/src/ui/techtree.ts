import { PlayerDTO } from '../state/state';
import { TECH_DEFS, TechDef } from '../constants';

const ARCHETYPE_COLORS: Record<string, string> = {
  Aggressor: '#ff6b6b',
  Defender:  '#4ecdc4',
  Builder:   '#ffd93d',
  Territorial: '#a8d8a8',
};

export class TechTreePanel {
  private el: HTMLElement;
  private visible = false;
  private onUnlock: (techId: number) => void;

  constructor(parent: HTMLElement, onUnlock: (techId: number) => void) {
    this.onUnlock = onUnlock;
    this.el = document.createElement('div');
    this.el.id = 'tech-tree';
    parent.appendChild(this.el);

    this.el.addEventListener('pointerdown', (e) => {
      if ((e.target as HTMLElement).id === 'tech-tree-close') {
        e.preventDefault();
        this.toggle();
        return;
      }
      const btn = (e.target as HTMLElement).closest('button[data-tech-id]') as HTMLElement | null;
      if (!btn || btn.hasAttribute('disabled')) return;
      e.preventDefault();
      const id = parseInt(btn.getAttribute('data-tech-id') ?? '-1', 10);
      if (id >= 0) this.onUnlock(id);
    });
  }

  toggle() {
    this.visible = !this.visible;
    this.el.classList.toggle('visible', this.visible);
  }

  update(player: PlayerDTO) {
    if (!this.visible) return;
    const tp = player.tp;
    let html = `<div style="position:relative;font-size:14px;font-weight:bold;margin-bottom:4px">Tech Tree <span style="color:#aaa;font-size:11px">[T to close]</span> &nbsp; Research: ${tp.toFixed(0)}<button id="tech-tree-close" style="position:absolute;top:0;right:0;background:none;border:none;color:#aaa;font-size:16px;cursor:pointer;line-height:1;padding:0 2px">✕</button></div>`;
    html += '<div class="tech-grid">';
    for (const def of TECH_DEFS) {
      html += this.renderCard(def, player.tech?.[def.id] ?? false, tp >= def.cost);
    }
    html += '</div>';
    this.el.innerHTML = html;
  }

  private renderCard(def: TechDef, owned: boolean, canAfford: boolean): string {
    const archetypeColor = ARCHETYPE_COLORS[def.archetype] ?? '#aaa';
    const bg = owned ? '#1a3a1a' : '#1a1a2e';
    const border = owned ? '#4a8a4a' : '#3a3a5a';
    let btnHtml = '';
    if (owned) {
      btnHtml = `<div style="color:#6a8a6a;font-size:11px;margin-top:4px">Owned — active</div>`;
    } else {
      const btnStyle = canAfford
        ? 'background:#4a4a6a;color:#fff;border:1px solid #6a6a8a;padding:3px 8px;border-radius:4px;cursor:pointer;font-family:monospace;font-size:11px'
        : 'background:#2a2a3a;color:#555;border:1px solid #3a3a4a;padding:3px 8px;border-radius:4px;font-family:monospace;font-size:11px';
      btnHtml = `<button style="${btnStyle}" data-tech-id="${def.id}" ${canAfford ? '' : 'disabled'}>Unlock (${def.cost} Research)</button>`;
    }
    return `
      <div style="background:${bg};border:1px solid ${border};border-radius:6px;padding:8px">
        <div style="display:flex;justify-content:space-between;align-items:center">
          <span style="font-weight:bold">${def.name}</span>
          <span style="color:${archetypeColor};font-size:10px">${def.archetype}</span>
        </div>
        <div style="color:#aaa;margin:3px 0 5px">${def.description}</div>
        ${btnHtml}
      </div>`;
  }
}

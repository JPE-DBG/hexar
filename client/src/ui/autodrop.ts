import { PlayerDTO, HexDTO, BattleDTO } from '../state/state';
import { BASE_INCOME_PER_SEC, GOLD_PER_LEVEL, GOLD_BONUS_MULTIPLIER, BUILDING_GOLD } from '../constants';

function calcHexIncome(hex: HexDTO): number {
  if (hex.building === BUILDING_GOLD) {
    return (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * hex.level) * GOLD_BONUS_MULTIPLIER;
  }
  return BASE_INCOME_PER_SEC;
}

export class AutoDropPanel {
  private el: HTMLElement;
  private onDrop: (q: number, r: number) => void;

  constructor(parent: HTMLElement, onDrop: (q: number, r: number) => void) {
    this.onDrop = onDrop;
    this.el = document.createElement('div');
    this.el.id = 'auto-drop';
    this.el.style.cssText = `
      position: fixed; top: 50px; right: 16px;
      background: rgba(80,20,20,0.95); border: 1px solid #ff4444; border-radius: 8px;
      padding: 12px; display: none; flex-direction: column; gap: 6px;
      font-family: monospace; color: #e0e0e0; font-size: 12px; min-width: 200px; z-index: 90;
    `;
    parent.appendChild(this.el);

    this.el.addEventListener('pointerdown', (e) => {
      const row = (e.target as HTMLElement).closest('[data-drop-q]') as HTMLElement | null;
      if (!row) return;
      e.preventDefault();
      const q = parseInt(row.getAttribute('data-drop-q') ?? '0', 10);
      const r = parseInt(row.getAttribute('data-drop-r') ?? '0', 10);
      this.onDrop(q, r);
    });
  }

  update(player: PlayerDTO | null, hexes: Map<string, HexDTO>, battles: BattleDTO[]) {
    if (!player?.autoDropActive) {
      this.el.style.display = 'none';
      return;
    }

    const droppable = [...hexes.values()].filter(h =>
      h.owner === player.id &&
      !h.capital &&
      !battles.some(b => b.dq === h.q && b.dr === h.r)
    ).sort((a, b) => calcHexIncome(a) - calcHexIncome(b));

    let html = `<div style="color:#ff8888;font-weight:bold">Auto-Drop Warning</div>`;
    html += `<div style="color:#ffaaaa">Income negative! Drop a hex.</div>`;
    html += `<div style="color:#ffcccc">Auto-drop in ${player.autoDropGrace.toFixed(1)}s</div>`;
    html += `<div style="color:#aaa;margin-top:4px;font-size:11px">Click to drop (lowest income first):</div>`;

    for (const hex of droppable) {
      const income = calcHexIncome(hex).toFixed(1);
      const bldg = hex.building === 0 ? '' : ` [${['','G','P','R'][hex.building]}${hex.level}]`;
      html += `
        <div data-drop-q="${hex.q}" data-drop-r="${hex.r}"
          style="cursor:pointer;padding:4px 6px;border-radius:4px;background:#3a1010;border:1px solid #5a2020;
            display:flex;justify-content:space-between"
          onmouseover="this.style.background='#5a1818'" onmouseout="this.style.background='#3a1010'">
          <span>[${hex.q},${hex.r}]${bldg}</span>
          <span style="color:#aaa">${income}/s</span>
        </div>`;
    }

    this.el.innerHTML = html;
    this.el.style.display = 'flex';
  }
}

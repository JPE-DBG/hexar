import { CardType } from '../constants.js';

export interface ShopCard {
  name: string;
  type: CardType;
  cost: number;
  qty: number;
}

const TYPE_ICON: Record<CardType, string> = {
  unit: '⚔',
  building: '🏗',
  bar: '⚡',
};

export function renderShop(container: HTMLElement, cards: ShopCard[]) {
  container.innerHTML = '';

  for (const card of cards) {
    const el = document.createElement('div');
    el.className = 'shop-card';
    el.innerHTML = `
      <span class="cost-badge">${card.cost}</span>
      <span class="shop-type-icon">${TYPE_ICON[card.type]}</span>
      <div class="shop-card-info">
        <div class="shop-card-name">${card.name}</div>
        <div class="shop-card-meta">×${card.qty} left</div>
      </div>
      <button class="btn-buy">BUY</button>
    `;
    container.appendChild(el);
  }

  const remove = document.createElement('div');
  remove.className = 'shop-remove';
  remove.innerHTML = `
    <div style="font-size:16px">✂</div>
    <div style="font-weight:bold;font-size:10px">REMOVE</div>
    <span class="cost-badge">5</span>
    <div style="font-size:9px">top · aside</div>
  `;
  container.appendChild(remove);
}

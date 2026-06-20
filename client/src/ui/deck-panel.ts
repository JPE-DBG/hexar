import { CardType } from '../constants.js';

export interface MockCard {
  name: string;
  type: CardType;
  cost: number;
  effect: string;
}

const TYPE_ICON: Record<CardType, string> = {
  unit: '⚔',
  building: '🏗',
  bar: '⚡',
};

function buildCard(card: MockCard, size: 'top' | 'aside'): HTMLElement {
  const el = document.createElement('div');
  el.className = `card card-${size === 'top' ? 'top-card' : 'aside'}`;

  const actions = size === 'top'
    ? `<div class="card-actions">
        <button class="btn-play">▶ PLAY</button>
        <button class="btn-push">→ PUSH</button>
       </div>`
    : `<div class="card-actions"><button class="btn-play">▶ PLAY</button></div>`;

  el.innerHTML = `
    <span class="card-type-icon">${TYPE_ICON[card.type]}</span>
    <span class="cost-badge">${card.cost}</span>
    <div class="card-name">${card.name}</div>
    <div class="card-effect">${card.effect}</div>
    ${actions}
  `;
  return el;
}

export function renderDeckPanel(container: HTMLElement, deck: MockCard[], aside: MockCard | null, deckCount: number) {
  container.innerHTML = '';

  // Top card
  if (deck.length > 0) {
    container.appendChild(buildCard(deck[0], 'top'));
  }

  // Aside slot
  const asideWrapper = document.createElement('div');
  asideWrapper.style.display = 'flex';
  asideWrapper.style.flexDirection = 'column';
  asideWrapper.style.alignItems = 'center';
  asideWrapper.style.gap = '2px';
  const asideLabel = document.createElement('div');
  asideLabel.className = 'card-aside-label';
  asideLabel.textContent = 'ASIDE';
  asideWrapper.appendChild(asideLabel);

  if (aside) {
    asideWrapper.appendChild(buildCard(aside, 'aside'));
  } else {
    const empty = document.createElement('div');
    empty.className = 'aside-empty';
    empty.innerHTML = `<div class="aside-empty-label">empty</div><div style="font-size:18px;color:var(--ui-border)">→</div>`;
    asideWrapper.appendChild(empty);
  }
  container.appendChild(asideWrapper);

  // Card backs (deck depth)
  const backs = document.createElement('div');
  backs.className = 'deck-backs';
  const sizes = [{ w: 90, h: 110 }, { w: 75, h: 92 }, { w: 62, h: 76 }];
  sizes.forEach(({ w, h }, i) => {
    const back = document.createElement('div');
    back.className = 'card-back';
    back.style.width = `${w}px`;
    back.style.height = `${h}px`;
    back.style.left = `${i * 8}px`;
    back.style.top = `${i * 4}px`;
    back.style.zIndex = `${3 - i}`;
    backs.appendChild(back);
  });
  backs.style.height = `${sizes[0].h + 8}px`;

  const backsWithLabel = document.createElement('div');
  backsWithLabel.style.display = 'flex';
  backsWithLabel.style.flexDirection = 'column';
  backsWithLabel.style.alignItems = 'flex-start';
  backsWithLabel.style.gap = '4px';
  backsWithLabel.appendChild(backs);
  const countLabel = document.createElement('div');
  countLabel.className = 'deck-count-label';
  countLabel.textContent = `${deckCount} cards`;
  backsWithLabel.appendChild(countLabel);
  container.appendChild(backsWithLabel);
}

export function renderBarMeter(container: HTMLElement, bar: number, max = 10) {
  container.innerHTML = '';

  const segments = document.createElement('div');
  segments.className = 'bar-segments';

  for (let i = 0; i < max; i++) {
    const seg = document.createElement('div');
    seg.className = 'bar-seg';
    const filled = bar >= i + 1;
    const partial = !filled && bar > i;
    if (filled) {
      seg.classList.add('filled');
    } else if (partial) {
      seg.classList.add('partial');
      const pct = Math.round((bar - i) * 100);
      seg.style.setProperty('--fill-pct', `${pct}%`);
    }
    segments.appendChild(seg);
  }

  const label = document.createElement('div');
  label.className = 'bar-label';
  label.textContent = `${bar.toFixed(1)} / ${max}`;

  container.appendChild(segments);
  container.appendChild(label);
}

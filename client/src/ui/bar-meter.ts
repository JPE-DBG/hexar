export function renderBarMeter(container: HTMLElement, bar: number, max = 10) {
  container.innerHTML = '';

  const whole = Math.min(Math.floor(bar), max);
  const frac = bar - Math.floor(bar);

  const label = document.createElement('span');
  label.className = 'bar-value-label';
  label.textContent = `${whole}`;
  container.appendChild(label);

  const seg = document.createElement('div');
  seg.className = 'bar-seg';
  if (frac > 0) {
    seg.classList.add('partial');
    seg.style.setProperty('--fill-pct', `${Math.round(frac * 100)}%`);
  }
  container.appendChild(seg);
}

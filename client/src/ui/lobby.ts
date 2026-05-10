import { COLORS } from '../constants';

export interface LobbyResult {
  code: string;
  token: string;
  playerId: number;
}

export class LobbyUI {
  private el: HTMLElement;
  private onReady: (result: LobbyResult) => void;

  constructor(container: HTMLElement, onReady: (result: LobbyResult) => void) {
    this.onReady = onReady;
    this.el = document.createElement('div');
    this.el.className = 'overlay';
    this.el.style.background = COLORS.background;
    this.el.innerHTML = `
      <h1 class="lobby-title" style="font-size:48px;color:${COLORS.accent}">HEXAR</h1>
      <div id="lobby-main" class="lobby-main">
        <button id="createBtn" class="btn" style="background:${COLORS.accent};color:${COLORS.background}">Create Game</button>
        <div style="color:#555;margin:4px 0">— or —</div>
        <div class="lobby-row">
          <input id="codeInput" class="code-input" placeholder="ROOM CODE" maxlength="4">
          <button id="joinBtn" class="btn" style="background:${COLORS.player2};color:#fff">Join</button>
        </div>
        <div id="lobbyStatus" class="lobby-status" style="color:${COLORS.textMuted}"></div>
      </div>
      <div id="lobby-waiting" class="lobby-waiting">
        <div style="color:${COLORS.textMuted};font-size:14px">Share this code with your opponent:</div>
        <div id="displayCode" class="lobby-code" style="font-size:64px;color:${COLORS.accent}"></div>
        <button id="enterBtn" class="btn" style="background:${COLORS.accent};color:${COLORS.background};margin-top:8px">Enter Game</button>
      </div>
    `;
    container.appendChild(this.el);
    this.bind();
  }

  private bind() {
    const createBtn = this.el.querySelector('#createBtn') as HTMLButtonElement;
    const joinBtn = this.el.querySelector('#joinBtn') as HTMLButtonElement;
    const codeInput = this.el.querySelector('#codeInput') as HTMLInputElement;
    const status = this.el.querySelector('#lobbyStatus') as HTMLElement;

    createBtn.addEventListener('click', async () => {
      createBtn.disabled = true;
      status.textContent = 'Creating room...';
      try {
        const res = await fetch('/lobby/create', { method: 'POST' });
        if (!res.ok) throw new Error('Server error');
        const data = await res.json() as LobbyResult;
        this.showWaiting(data);
      } catch {
        status.textContent = 'Error creating room. Try again.';
        createBtn.disabled = false;
      }
    });

    joinBtn.addEventListener('click', () => this.doJoin(codeInput, joinBtn, status));
    codeInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') this.doJoin(codeInput, joinBtn, status);
    });
    codeInput.addEventListener('input', () => {
      codeInput.value = codeInput.value.toUpperCase();
    });
  }

  private showWaiting(result: LobbyResult) {
    const main = this.el.querySelector('#lobby-main') as HTMLElement;
    const waiting = this.el.querySelector('#lobby-waiting') as HTMLElement;
    const displayCode = this.el.querySelector('#displayCode') as HTMLElement;
    const enterBtn = this.el.querySelector('#enterBtn') as HTMLButtonElement;

    main.style.display = 'none';
    waiting.style.display = 'flex';
    displayCode.textContent = result.code;

    enterBtn.addEventListener('click', () => {
      this.hide();
      this.onReady(result);
    });
  }

  private async doJoin(codeInput: HTMLInputElement, joinBtn: HTMLButtonElement, status: HTMLElement) {
    const code = codeInput.value.toUpperCase().trim();
    if (!code) { status.textContent = 'Enter a room code'; return; }
    joinBtn.disabled = true;
    status.textContent = 'Joining...';
    try {
      const res = await fetch('/lobby/join', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code }),
      });
      if (!res.ok) {
        const txt = await res.text();
        throw new Error(txt.trim() || 'Failed to join');
      }
      const data = await res.json() as { token: string; playerId: number };
      this.hide();
      this.onReady({ code, ...data });
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : 'Error joining room';
      status.textContent = msg;
      joinBtn.disabled = false;
    }
  }

  hide() {
    this.el.style.display = 'none';
  }
}

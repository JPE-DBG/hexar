export interface LobbyResult {
  code: string;
  token: string;
  playerId: number;
}

const STYLE = `
  position:fixed;top:0;left:0;width:100%;height:100%;
  display:flex;flex-direction:column;align-items:center;justify-content:center;
  background:#1a1a2e;z-index:1000;font-family:monospace;color:#e0e0e0;
`;

const BTN = `padding:12px 32px;font-size:18px;border:none;border-radius:4px;cursor:pointer;font-family:monospace;`;

export class LobbyUI {
  private el: HTMLElement;
  private onReady: (result: LobbyResult) => void;

  constructor(container: HTMLElement, onReady: (result: LobbyResult) => void) {
    this.onReady = onReady;
    this.el = document.createElement('div');
    this.el.style.cssText = STYLE;
    this.el.innerHTML = `
      <h1 style="font-size:48px;color:#4ecdc4;margin-bottom:40px;letter-spacing:4px">HEXAR</h1>
      <div id="lobby-main" style="display:flex;flex-direction:column;gap:16px;align-items:center">
        <button id="createBtn" style="${BTN}background:#4ecdc4;color:#1a1a2e">Create Game</button>
        <div style="color:#555;margin:4px 0">— or —</div>
        <div style="display:flex;gap:8px">
          <input id="codeInput" placeholder="ROOM CODE" maxlength="4"
            style="padding:12px;font-size:18px;background:#2a2a4e;color:#e0e0e0;border:1px solid #444;
                   border-radius:4px;width:150px;font-family:monospace;text-transform:uppercase;text-align:center">
          <button id="joinBtn" style="${BTN}background:#ff6b6b;color:#fff">Join</button>
        </div>
        <div id="lobbyStatus" style="color:#aaa;margin-top:8px;min-height:28px;text-align:center"></div>
      </div>
      <div id="lobby-waiting" style="display:none;flex-direction:column;align-items:center;gap:16px">
        <div style="color:#aaa;font-size:14px">Share this code with your opponent:</div>
        <div id="displayCode" style="font-size:64px;color:#4ecdc4;letter-spacing:12px;font-weight:bold"></div>
        <button id="enterBtn" style="${BTN}background:#4ecdc4;color:#1a1a2e;margin-top:8px">Enter Game</button>
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

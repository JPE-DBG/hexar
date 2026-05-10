import { SnapshotMsg, DeltaMsg } from '../state/state';

export interface WelcomeMsg {
  type: 'welcome';
  playerId: number;
}

type ServerMsg = SnapshotMsg | DeltaMsg | WelcomeMsg;

export interface ConnectionHandlers {
  onSnapshot: (msg: SnapshotMsg) => void;
  onDelta: (msg: DeltaMsg) => void;
  onWelcome: (msg: WelcomeMsg) => void;
  onDisconnect?: () => void;
}

export class Connection {
  private ws: WebSocket | null = null;
  private handlers: ConnectionHandlers;
  private code: string;
  private token: string;
  private reconnectAttempts = 0;
  private readonly maxReconnectAttempts = 5;
  private closed = false;

  constructor(code: string, token: string, handlers: ConnectionHandlers) {
    this.code = code;
    this.token = token;
    this.handlers = handlers;
    this.connect();
  }

  send(msg: object) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg));
    }
  }

  close() {
    this.closed = true;
    this.ws?.close();
  }

  private buildUrl(): string {
    return `ws://${window.location.host}/ws?code=${this.code}&token=${this.token}`;
  }

  private connect() {
    if (this.closed) return;
    this.ws = new WebSocket(this.buildUrl());

    this.ws.onopen = () => {
      console.log('connected');
      this.reconnectAttempts = 0;
    };

    this.ws.onmessage = (ev) => {
      const msg = JSON.parse(ev.data) as ServerMsg;
      switch (msg.type) {
        case 'welcome':
          this.handlers.onWelcome(msg);
          break;
        case 'snapshot':
          this.handlers.onSnapshot(msg);
          break;
        case 'delta':
          this.handlers.onDelta(msg);
          break;
      }
    };

    this.ws.onclose = () => {
      if (this.closed) return;
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 16000);
        this.reconnectAttempts++;
        console.log(`disconnected, reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        setTimeout(() => this.connect(), delay);
      } else {
        console.log('max reconnect attempts reached');
        this.handlers.onDisconnect?.();
      }
    };

    this.ws.onerror = () => {
      this.ws?.close();
    };
  }
}

import { SnapshotMsg, DeltaMsg } from '../state/state';

const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_INITIAL_DELAY = 1000; // ms, doubles each attempt
const RECONNECT_MAX_DELAY = 16000;    // ms cap on backoff

export interface WelcomeMsg {
  type: 'welcome';
  playerId: number;
}

type ServerMsg = SnapshotMsg | DeltaMsg | WelcomeMsg;

export interface ConnectionHandlers {
  onSnapshot: (msg: SnapshotMsg) => void;
  onDelta: (msg: DeltaMsg) => void;
  onWelcome: (msg: WelcomeMsg) => void;
  onReconnecting?: (attempt: number, max: number) => void;
  onAlreadyConnected?: () => void;
  onDisconnect?: () => void;
}

export class Connection {
  private ws: WebSocket | null = null;
  private handlers: ConnectionHandlers;
  private code: string;
  private token: string;
  private reconnectAttempts = 0;
  private readonly maxReconnectAttempts = MAX_RECONNECT_ATTEMPTS;
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

    this.ws.onclose = (event: CloseEvent) => {
      if (this.closed) return;
      if (event.code === 4001) {
        this.closed = true;
        this.handlers.onAlreadyConnected?.();
        return;
      }
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        const delay = Math.min(RECONNECT_INITIAL_DELAY * Math.pow(2, this.reconnectAttempts), RECONNECT_MAX_DELAY);
        this.reconnectAttempts++;
        console.log(`disconnected, reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
        this.handlers.onReconnecting?.(this.reconnectAttempts, this.maxReconnectAttempts);
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

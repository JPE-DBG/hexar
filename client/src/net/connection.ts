import { SnapshotMsg } from '../state/state';

export interface WelcomeMsg {
  type: 'welcome';
  playerId: number;
}

type ServerMsg = SnapshotMsg | WelcomeMsg;

export interface ConnectionHandlers {
  onSnapshot: (msg: SnapshotMsg) => void;
  onWelcome: (msg: WelcomeMsg) => void;
}

export class Connection {
  private ws: WebSocket | null = null;
  private handlers: ConnectionHandlers;
  private url: string;

  constructor(url: string, handlers: ConnectionHandlers) {
    this.url = url;
    this.handlers = handlers;
    this.connect();
  }

  send(msg: object) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(msg));
    }
  }

  private connect() {
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log('connected');
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
      }
    };

    this.ws.onclose = () => {
      console.log('disconnected, reconnecting...');
      setTimeout(() => this.connect(), 1000);
    };

    this.ws.onerror = () => {
      this.ws?.close();
    };
  }
}

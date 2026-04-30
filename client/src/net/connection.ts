import { SnapshotMsg } from '../state/state';

export type MessageHandler = (msg: SnapshotMsg) => void;

export class Connection {
  private ws: WebSocket | null = null;
  private handler: MessageHandler;
  private url: string;

  constructor(url: string, handler: MessageHandler) {
    this.url = url;
    this.handler = handler;
    this.connect();
  }

  private connect() {
    this.ws = new WebSocket(this.url);

    this.ws.onopen = () => {
      console.log('connected');
    };

    this.ws.onmessage = (ev) => {
      const msg = JSON.parse(ev.data) as SnapshotMsg;
      this.handler(msg);
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

export interface PingResult {
  ip: string;
  status: 'online' | 'down';
  latency: number;
  ttl: number;
  seq: number;
  timestamp: string;
}

export interface WSMessage {
  type: 'ping_result' | 'status_change';
  data: PingResult;
}

export interface MonitorRequest {
  ip: string;
}

export interface APIResponse {
  success: boolean;
  message: string;
}

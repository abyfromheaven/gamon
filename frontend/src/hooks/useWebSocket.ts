import { useCallback, useEffect, useRef, useState } from 'react';
import type { MonitorResult, StatusChange } from '../types';

function websocketURL(): string {
  const configured = import.meta.env.VITE_WS_URL;
  if (configured) return configured;
  const apiURL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080';
  return `${apiURL.replace(/^http/, 'ws').replace(/\/$/, '')}/ws`;
}

interface InitialStateItem {
  device_id: number;
  ip: string;
  method: string;
  status: MonitorResult['status'];
  latency_ms: number;
  last_check: string;
}

export function useWebSocket(onReconnect?: () => void) {
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [monitorResults, setMonitorResults] = useState<Map<number, MonitorResult>>(new Map());
  const [lastStatusChange, setLastStatusChange] = useState<StatusChange | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  const connect = useCallback(() => {
    if (socketRef.current?.readyState === WebSocket.OPEN) return;
    const socket = new WebSocket(websocketURL());

    socket.onopen = () => {
      setIsConnected(true);
      onReconnect?.();
    };
    socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data) as { type: string; data: unknown };
        if (message.type === 'initial_state' && Array.isArray(message.data)) {
          const results = message.data.map((item: InitialStateItem) => ({
            device_id: item.device_id,
            ip: item.ip,
            status: item.status,
            latency_ms: item.latency_ms,
            ttl: 0,
            seq: 0,
            timestamp: item.last_check,
            method: item.method,
            details: {},
          }));
          setMonitorResults(new Map(results.map((item) => [item.device_id, item])));
        }
        if (message.type === 'check_result') {
          const result = message.data as MonitorResult;
          setMonitorResults((current) => new Map(current).set(result.device_id, result));
        }
        if (message.type === 'status_change') {
          setLastStatusChange(message.data as StatusChange);
        }
      } catch {
        // A malformed message must not disconnect the monitoring view.
      }
    };
    socket.onclose = () => {
      setIsConnected(false);
      reconnectTimer.current = setTimeout(connect, 3000);
    };
    socket.onerror = () => socket.close();
    socketRef.current = socket;
  }, [onReconnect]);

  useEffect(() => {
    connect();
    return () => {
      if (reconnectTimer.current) clearTimeout(reconnectTimer.current);
      socketRef.current?.close();
    };
  }, [connect]);

  return { isConnected, monitorResults, lastStatusChange };
}

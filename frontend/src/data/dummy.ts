import type { DashboardData } from '../types';

export const dashboardData: DashboardData = {
  summary: {
    totalDevices: 47,
    online: 42,
    offline: 3,
    warnings: 2,
  },
  deviceBreakdown: [
    { type: 'Server', count: 12, online: 11 },
    { type: 'Router', count: 8, online: 7 },
    { type: 'Switch', count: 15, online: 14 },
    { type: 'Access Point', count: 10, online: 9 },
    { type: 'Website', count: 2, online: 1 },
  ],
  latestAlerts: [
    {
      id: 1,
      device: 'Router-03',
      message: 'Device unreachable — no response for 3 check intervals',
      severity: 'critical',
      time: '2 min ago',
      timestamp: Date.now() - 120000,
    },
    {
      id: 2,
      device: 'Switch-12',
      message: 'Latency above threshold — 245ms avg (limit: 100ms)',
      severity: 'warning',
      time: '5 min ago',
      timestamp: Date.now() - 300000,
    },
    {
      id: 3,
      device: 'AP-07',
      message: 'Connection restored after 4 failed checks',
      severity: 'info',
      time: '8 min ago',
      timestamp: Date.now() - 480000,
    },
    {
      id: 4,
      device: 'Web-01',
      message: 'HTTP status 503 — service unavailable',
      severity: 'critical',
      time: '12 min ago',
      timestamp: Date.now() - 720000,
    },
    {
      id: 5,
      device: 'Server-05',
      message: 'Disk usage above 90% — 92.4% used',
      severity: 'warning',
      time: '18 min ago',
      timestamp: Date.now() - 1080000,
    },
  ],
  systemStatus: {
    monitoring: 'Running',
    checkInterval: '3s',
    lastScan: '14:32:04',
    notifications: 'Active',
  },
};

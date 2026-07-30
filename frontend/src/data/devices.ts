import type { DeviceType } from '../types';

export interface Device {
  id: number;
  name: string;
  type: DeviceType;
  ip: string;
  method: 'ICMP Ping';
  port: number | null;
  location: string;
  status: 'active' | 'inactive';
  lastSeen: string | null;
}

export const devices: Device[] = [
  { id: 1, name: 'Core Router-01', type: 'Router', ip: '192.168.1.1', method: 'ICMP Ping', port: null, location: 'Ruang Server Utama', status: 'active', lastSeen: '2 min ago' },
  { id: 2, name: 'Core Router-02', type: 'Router', ip: '192.168.1.2', method: 'ICMP Ping', port: null, location: 'Ruang Server Utama', status: 'active', lastSeen: '1 min ago' },
  { id: 3, name: 'Web Server-01', type: 'Server', ip: '192.168.1.10', method: 'ICMP Ping', port: 443, location: 'Ruang Server Utama', status: 'active', lastSeen: '1 min ago' },
  { id: 4, name: 'DB Server-01', type: 'Server', ip: '192.168.1.11', method: 'ICMP Ping', port: 3306, location: 'Ruang Server Utama', status: 'active', lastSeen: '30 sec ago' },
  { id: 5, name: 'DB Server-02', type: 'Server', ip: '192.168.1.12', method: 'ICMP Ping', port: 5432, location: 'Ruang Server Utama', status: 'inactive', lastSeen: '25 min ago' },
  { id: 6, name: 'Switch-Lantai1', type: 'Switch', ip: '192.168.2.1', method: 'ICMP Ping', port: null, location: 'Lantai 1 - Reception', status: 'active', lastSeen: '45 sec ago' },
  { id: 7, name: 'Switch-Lantai2', type: 'Switch', ip: '192.168.2.2', method: 'ICMP Ping', port: null, location: 'Lantai 2 - Ruang Rapat', status: 'active', lastSeen: '1 min ago' },
  { id: 8, name: 'Switch-Lantai3', type: 'Switch', ip: '192.168.2.3', method: 'ICMP Ping', port: null, location: 'Lantai 3 - Kantor', status: 'active', lastSeen: '3 min ago' },
  { id: 9, name: 'AP Lantai 1', type: 'Access Point', ip: '192.168.3.10', method: 'ICMP Ping', port: null, location: 'Lantai 1 - Lobby', status: 'active', lastSeen: '30 sec ago' },
  { id: 10, name: 'AP Lantai 2', type: 'Access Point', ip: '192.168.3.11', method: 'ICMP Ping', port: null, location: 'Lantai 2 - Ruang Rapat', status: 'active', lastSeen: '20 sec ago' },
  { id: 11, name: 'AP Lantai 3', type: 'Access Point', ip: '192.168.3.12', method: 'ICMP Ping', port: null, location: 'Lantai 3 - Kantor', status: 'inactive', lastSeen: '1 jam ago' },
  { id: 12, name: 'AP Gudang', type: 'Access Point', ip: '192.168.3.13', method: 'ICMP Ping', port: null, location: 'Gudang', status: 'active', lastSeen: '1 min ago' },
  { id: 13, name: 'Website Corp', type: 'Website', ip: '103.25.48.12', method: 'ICMP Ping', port: 443, location: 'External', status: 'active', lastSeen: '2 min ago' },
  { id: 14, name: 'Website Portal', type: 'Website', ip: '103.25.48.13', method: 'ICMP Ping', port: 80, location: 'External', status: 'active', lastSeen: '1 min ago' },
  { id: 15, name: 'Backup Router', type: 'Router', ip: '192.168.1.3', method: 'ICMP Ping', port: null, location: 'Ruang Server Utama', status: 'inactive', lastSeen: '3 hari lalu' },
];

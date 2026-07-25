import { useState } from 'react';
import type { APIResponse } from '../types';

interface Props {
  onMonitor: (ip: string) => void;
}

export function InputForm({ onMonitor }: Props) {
  const [ip, setIp] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    const trimmed = ip.trim();
    if (!trimmed) {
      setError('IP address wajib diisi');
      return;
    }

    const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/;
    const domainRegex = /^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$/;
    if (!ipRegex.test(trimmed) && !domainRegex.test(trimmed)) {
      setError('Format IP/domain tidak valid');
      return;
    }

    setLoading(true);
    try {
      const res = await fetch('http://localhost:8080/api/monitor', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ip: trimmed }),
      });
      const data: APIResponse = await res.json();
      if (data.success) {
        onMonitor(trimmed);
        setIp('');
      } else {
        setError(data.message);
      }
    } catch {
      setError('Gagal menghubungi server');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-3">
      <div className="flex gap-2">
        <input
          type="text"
          value={ip}
          onChange={(e) => setIp(e.target.value)}
          placeholder="Masukkan IP (contoh: 192.168.1.1)"
          className="flex-1 px-4 py-3 bg-gray-800 border border-gray-600 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-green-500 transition-colors"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-3 bg-green-600 hover:bg-green-700 disabled:bg-gray-600 text-white font-semibold rounded-lg transition-colors cursor-pointer"
        >
          {loading ? 'Memproses...' : 'Monitor'}
        </button>
      </div>
      {error && (
        <p className="text-red-400 text-sm">{error}</p>
      )}
    </form>
  );
}

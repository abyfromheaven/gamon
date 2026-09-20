import { useEffect, useMemo, useRef, useState } from 'react';
import type { Device, MonitoringRecord, PingHistoryRecord } from '../lib/api';
import { fetchDeviceHistory, fetchDevices, fetchMonitoring } from '../lib/api';
import type { MonitorResult } from '../types';
import { PageHeader } from '../components/PageHeader';
import { StatusIndicator } from '../components/StatusIndicator';
import { LatencyChart } from '../components/LatencyChart';

// ============================================================
// HALAMAN MONITORING
// ============================================================
// Halaman ini menampilkan status real-time semua perangkat yang dipantau.
//
// Fitur utama:
// 1. Tabel daftar perangkat dengan status, latency, dan uptime
// 2. Panel detail untuk melihat info lengkap satu perangkat
// 3. Grafik latency (waktu respons) dari 50 data terakhir
// 4. Badge uptime dengan warna yang menunjukkan kondisi perangkat
//
// Data diperoleh dari 2 sumber:
// - REST API (fetchMonitoring) -> Data awal dari database
// - WebSocket (monitorResults) -> Update real-time setiap kali ada pengecekan
// ============================================================

export function MonitoringPage({ monitorResults, reconnectKey }: { monitorResults: Map<number, MonitorResult>; onViewAlerts: () => void; reconnectKey: number }) {
  // ============================================================
  // STATE (DATA YANG BERUBAH-UBAH)
  // ============================================================

  /** Daftar semua perangkat dari database */
  const [devices, setDevices] = useState<Device[]>([]);

  /** Data status monitoring dari server */
  const [records, setRecords] = useState<MonitoringRecord[]>([]);

  /** Perangkat yang sedang dipilih (untuk ditampilkan detailnya) */
  const [selected, setSelected] = useState<MonitoringRecord | null>(null);

  /** Riwayat pengecekan perangkat yang dipilih (untuk grafik) */
  const [history, setHistory] = useState<PingHistoryRecord[]>([]);

  /** Kata kunci pencarian perangkat */
  const [search, setSearch] = useState('');

  /** Filter berdasarkan status (semua/online/offline) */
  const [statusFilter, setStatusFilter] = useState<'all' | MonitoringRecord['status']>('all');

  /** Filter berdasarkan jenis perangkat */
  const [typeFilter, setTypeFilter] = useState<'all' | Device['type']>('all');

  /** Pesan error jika gagal memuat data */
  const [error, setError] = useState('');

  /** Cache riwayat pencekan agar tidak perlu fetch ulang */
  const historyCache = useRef<Map<number, PingHistoryRecord[]>>(new Map());

  // ============================================================
  // MENGAMBIL DATA DARI SERVER
  // ============================================================

  // Ambil data perangkat dan monitoring saat halaman pertama kali dimuat
  useEffect(() => {
    void Promise.all([fetchDevices(), fetchMonitoring()])
      .then(([nextDevices, nextRecords]) => { setDevices(nextDevices); setRecords(nextRecords); })
      .catch((err: unknown) => setError(err instanceof Error ? err.message : 'Gagal memuat monitoring.'));
  }, [reconnectKey]);

  // ============================================================
  // MENGGABUNGKAN DATA REAL-TIME DENGAN DATA AWAL
  // ============================================================

  // Gabungkan data dari database dengan update real-time dari WebSocket
  // Ini membuat tampilan selalu up-to-date tanpa perlu refresh
  const live = useMemo(
    () => records.map((record) => {
      const result = monitorResults.get(record.device_id);
      if (!result) return record;

      // Ambil last_online dari WebSocket details jika ada
      const detailsLastOnline = result.details?.last_online as string | undefined;
      const lastOnline = detailsLastOnline || record.last_online;

      return {
        ...record,
        status: result.status,
        latency_ms: result.latency_ms,
        last_check: result.timestamp,
        last_online: lastOnline || record.last_online,
      };
    }),
    [records, monitorResults]
  );

  // ============================================================
  // FILTER DATA
  // ============================================================

  // Filter data berdasarkan pencarian, status, dan jenis perangkat
  const filtered = live.filter(
    (record) =>
      (statusFilter === 'all' || record.status === statusFilter) &&
      (typeFilter === 'all' || record.device_type === typeFilter) &&
      (!search || `${record.device_name} ${record.ip}`.toLowerCase().includes(search.toLowerCase()))
  );

  // ============================================================
  // PEMILIHAN PERANGKAT DAN RIWAYAT
  // ============================================================

  // Fungsi saat perangkat dipilih di tabel
  const choose = async (record: MonitoringRecord) => {
    setSelected(record);

    // Cek apakah riwayat sudah ada di cache
    const cached = historyCache.current.get(record.device_id);
    if (cached) {
      setHistory(cached);
      return;
    }

    // Jika belum ada di cache, ambil dari server
    try {
      const data = await fetchDeviceHistory(record.device_id);
      const sliced = data.slice(0, 50); // Ambil 50 data terakhir
      historyCache.current.set(record.device_id, sliced);
      setHistory(sliced);
    } catch {
      setHistory([]);
    }
  };

  // Update riwayat dengan data real-time dari WebSocket
  useEffect(() => {
    if (!selected) return;
    const deviceId = selected.device_id;
    const liveResult = monitorResults.get(deviceId);
    if (!liveResult) return;

    // Buat record baru dari hasil pengecekan terakhir
    const newRecord: PingHistoryRecord = {
      id: Date.now(),
      status: liveResult.status,
      latency_ms: liveResult.latency_ms,
      ttl: liveResult.ttl,
      seq: liveResult.seq,
      details: JSON.stringify(liveResult.details),
      timestamp: liveResult.timestamp,
    };

    // Tambahkan ke riwayat dan batasi hanya 50 data terakhir
    const cached = historyCache.current.get(deviceId) || [];
    const updated = [...cached, newRecord].slice(-50);
    historyCache.current.set(deviceId, updated);
    setHistory(updated);
  }, [monitorResults, selected?.device_id]);

  // ============================================================
  // PERSIAPAN DATA UNTUK GRAFIK
  // ============================================================

  // Ubah data riwayat menjadi format yang bisa dibaca oleh komponen grafik
  const chartData = useMemo(
    () => history
      .slice(-50) // Ambil 50 data terakhir
      .reverse()  // Balik urutan (dari lama ke baru)
      .map((item) => {
        const date = new Date(item.timestamp);
        const time = `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}`;
        return { time, latency: item.latency_ms };
      }),
    [history]
  );

  // ============================================================
  // HITUNG JUMLAH PERANGKAT
  // ============================================================

  const counts = {
    total: live.length,
    online: live.filter((item) => item.status === 'online').length,
    offline: live.filter((item) => item.status === 'offline').length,
  };

  // Cari data perangkat yang sedang dipilih
  const detail = selected && devices.find((device) => device.id === selected.device_id);

  // ============================================================
  // FUNGSI UNTUK BADGE UPTIME
  // ============================================================

  // Membuat badge uptime dengan warna berdasarkan durasi status
  // - Hijau: Online lama (lebih dari 5 menit)
  // - Biru: Online 1-5 menit
  // - Kuning: Online kurang dari 1 menit (baru pulih)
  // - Merah: Offline
  const getUptimeBadge = (status: string, duration: string) => {
    if (!duration) return null;

    // Ubah string durasi (contoh: "2j 15m 30d") menjadi total detik
    const parts = duration.split(' ');
    let totalSeconds = 0;
    for (const part of parts) {
      if (part.endsWith('j')) {
        totalSeconds += parseInt(part) * 3600; // jam -> detik
      } else if (part.endsWith('m')) {
        totalSeconds += parseInt(part) * 60; // menit -> detik
      } else if (part.endsWith('d')) {
        totalSeconds += parseInt(part); // sudah detik
      }
    }

    // Tentukan warna badge berdasarkan durasi dan status
    let bgColor = 'bg-success/10 text-success'; // Default: hijau
    let borderColor = 'border-success/30';

    if (status === 'offline') {
      // Offline: merah
      bgColor = 'bg-danger/10 text-danger';
      borderColor = 'border-danger/30';
    } else if (totalSeconds < 60) {
      // Online kurang dari 1 menit: kuning (baru pulih)
      bgColor = 'bg-warning/10 text-warning';
      borderColor = 'border-warning/30';
    } else if (totalSeconds < 300) {
      // Online 1-5 menit: biru
      bgColor = 'bg-accent/10 text-accent';
      borderColor = 'border-accent/30';
    }
    // Online lebih dari 5 menit: hijau (default)

    return (
      <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border ${bgColor} ${borderColor}`}>
        {status === 'online' ? '↑' : '↓'} {duration}
      </span>
    );
  };

  // ============================================================
  // TAMPILAN (RENDER)
  // ============================================================

  return (
    <div className="max-w-[1200px] mx-auto px-4 lg:px-8 py-6 lg:py-8 space-y-5">
      {/* Header halaman */}
      <PageHeader
        title="Monitoring"
        subtitle={`${counts.total} devices · ${counts.online} online · ${counts.offline} offline`}
      />

      {/* Pesan error jika ada */}
      {error && <p className="rounded bg-danger-muted p-3 text-sm text-danger">{error}</p>}

      {/* Filter pencarian dan dropdown */}
      <div className="flex flex-wrap gap-3">
        <input
          className="min-w-56 flex-1 rounded-lg border border-border bg-surface p-3 text-sm"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder="Search device or IP..."
        />
        <select
          className="rounded-lg border border-border bg-surface px-3 text-sm"
          value={statusFilter}
          onChange={(event) => setStatusFilter(event.target.value as typeof statusFilter)}
        >
          <option value="all">All status</option>
          <option value="online">Online</option>
          <option value="offline">Offline</option>
          <option value="unknown">Unknown</option>
        </select>
        <select
          className="rounded-lg border border-border bg-surface px-3 text-sm"
          value={typeFilter}
          onChange={(event) => setTypeFilter(event.target.value as typeof typeFilter)}
        >
          <option value="all">All device types</option>
          {['Server', 'Router', 'Switch', 'Access Point', 'Website'].map((type) => (
            <option key={type}>{type}</option>
          ))}
        </select>
      </div>

      {/* Konten utama: Tabel + Panel Detail */}
      <div className="grid lg:grid-cols-3 gap-5">
        {/* ============================================================ */}
        {/* TABEL DAFTAR PERANGKAT */}
        {/* ============================================================ */}
        <div className="lg:col-span-2 rounded-xl border border-border bg-surface overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-text-muted">
                <th className="p-4">Device</th>
                <th>Status</th>
                <th>Latency</th>
                <th>Uptime</th>
                <th>Last Check</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((record) => (
                <tr
                  key={record.device_id}
                  onClick={() => void choose(record)}
                  className={`cursor-pointer border-t border-border/50 hover:bg-surface-elevated transition-colors ${
                    selected?.device_id === record.device_id ? 'bg-surface-elevated' : ''
                  }`}
                >
                  {/* Nama dan IP perangkat */}
                  <td className="p-4">
                    <p className="font-medium">{record.device_name}</p>
                    <p className="font-mono text-xs text-text-muted">{record.ip} · {record.device_type}</p>
                  </td>
                  {/* Indikator status (hijau/merah) */}
                  <td>
                    <StatusIndicator status={record.status} />
                  </td>
                  {/* Waktu respons */}
                  <td className="font-mono text-xs">{record.latency_ms} ms</td>
                  {/* Badge uptime */}
                  <td>
                    {getUptimeBadge(record.status, record.current_status_duration)}
                  </td>
                  {/* Kapan terakhir kali dicek */}
                  <td className="text-xs text-text-muted">
                    {record.last_check ? new Date(record.last_check).toLocaleTimeString('id-ID') : '—'}
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={5} className="p-8 text-center text-text-muted">Tidak ada device</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        {/* ============================================================ */}
        {/* PANEL DETAIL PERANGKAT */}
        {/* ============================================================ */}
        <div className="rounded-xl border border-border bg-surface p-5 space-y-4">
          {detail && selected ? (
            <>
              {/* Nama dan jenis perangkat */}
              <div>
                <h3 className="text-lg font-semibold">{detail.name}</h3>
                <p className="text-xs text-text-muted">{detail.type} · {detail.method}</p>
              </div>

              {/* Informasi dasar perangkat */}
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-text-muted">IP</span>
                  <span className="font-mono">{detail.ip}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-text-muted">Status</span>
                  <StatusIndicator status={selected.status} />
                </div>
                <div className="flex justify-between">
                  <span className="text-text-muted">Latency</span>
                  <span className="font-mono">{selected.latency_ms} ms</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-text-muted">Interval</span>
                  <span>{selected.interval}s</span>
                </div>

                {/* ============================================================ */}
                {/* INFO UPTIME / STATUS */}
                {/* ============================================================ */}
                <div className="pt-2 border-t border-border/50">
                  <p className="text-text-muted text-xs mb-2">Status Info</p>
                  <div className="space-y-2">
                    {/* Status saat ini dengan badge */}
                    <div className="flex justify-between items-center">
                      <span className="text-text-muted text-xs">Status saat ini</span>
                      {getUptimeBadge(selected.status, selected.current_status_duration)}
                    </div>

                    {/* Sejak kapan status ini aktif */}
                    {selected.current_status_since && (
                      <div className="flex justify-between items-center">
                        <span className="text-text-muted text-xs">
                          {selected.status === 'online' ? 'Online sejak' : 'Offline sejak'}
                        </span>
                        <span className="font-mono text-xs">{selected.current_status_since}</span>
                      </div>
                    )}

                    {/* Kapan terakhir kali perangkat online */}
                    {selected.last_online && (
                      <div className="flex justify-between items-center">
                        <span className="text-text-muted text-xs">Last online</span>
                        <span className="font-mono text-xs">
                          {new Date(selected.last_online).toLocaleTimeString('id-ID')}
                        </span>
                      </div>
                    )}
                  </div>
                </div>

                {/* Lokasi perangkat */}
                {detail.location && (
                  <div className="flex justify-between">
                    <span className="text-text-muted">Lokasi</span>
                    <span>{detail.location}</span>
                  </div>
                )}

                {/* Deskripsi perangkat */}
                {detail.description && (
                  <div className="pt-2 border-t border-border/50">
                    <p className="text-text-muted text-xs mb-1">Deskripsi</p>
                    <p className="text-sm">{detail.description}</p>
                  </div>
                )}
              </div>

              {/* ============================================================ */}
              {/* GRAFIK LATENCY */}
              {/* ============================================================ */}
              <div className="pt-3 border-t border-border/50">
                <p className="text-xs text-text-muted mb-3">Latency History (50 data terakhir)</p>
                <LatencyChart data={chartData} />
              </div>
            </>
          ) : (
            /* Pesan jika belum ada perangkat yang dipilih */
            <div className="flex items-center justify-center h-full min-h-[300px] text-text-muted text-sm">
              Pilih device untuk melihat detail
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

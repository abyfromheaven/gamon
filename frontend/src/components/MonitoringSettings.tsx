import { useCallback, useEffect, useState } from 'react';
import { getSettings, updateSettings, type AppSettings } from '../lib/api';

export function MonitoringSettings() {
  const [settings, setSettings] = useState<AppSettings | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const loadSettings = useCallback(async () => {
    try {
      setError('');
      const data = await getSettings();
      setSettings(data);
    } catch {
      setError('Gagal memuat settings');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { void loadSettings(); }, [loadSettings]);

  const handleSave = async () => {
    if (!settings) return;
    try {
      setSaving(true);
      setError('');
      setSuccess('');
      await updateSettings(settings);
      setSuccess('Settings berhasil disimpan');
      setTimeout(() => setSuccess(''), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal menyimpan settings');
    } finally {
      setSaving(false);
    }
  };

  if (loading || !settings) {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="w-5 h-5 border-2 border-accent border-t-transparent rounded-full animate-spin" />
        <span className="ml-3 text-sm text-text-muted">Loading...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {error && (
        <div className="rounded-lg bg-danger/10 border border-danger/20 p-3 text-sm text-danger">
          {error}
        </div>
      )}
      {success && (
        <div className="rounded-lg bg-success/10 border border-success/20 p-3 text-sm text-success">
          {success}
        </div>
      )}

      {/* Failure Threshold */}
      <div className="space-y-2">
        <label className="block text-sm font-medium text-text-primary">
          Failure Threshold
        </label>
        <p className="text-xs text-text-muted">
          Berapa kali gagal berturut-turut sebelum device dianggap offline
        </p>
        <div className="flex items-center gap-3">
          <input
            type="range"
            min={1}
            max={10}
            value={settings.failure_threshold}
            onChange={(e) => setSettings({ ...settings, failure_threshold: Number(e.target.value) })}
            className="flex-1 h-2 bg-surface-elevated rounded-lg appearance-none cursor-pointer accent-accent"
          />
          <span className="w-12 text-center text-sm font-mono font-semibold text-accent">
            {settings.failure_threshold}x
          </span>
        </div>
        <p className="text-[11px] text-text-muted">
          {settings.failure_threshold === 1 && 'Sangat sensitif — langsung alert saat gagal 1x'}
          {settings.failure_threshold === 2 && 'Cukup sensitif — alert saat gagal 2x berturut-turut'}
          {settings.failure_threshold === 3 && 'Default — alert saat gagal 3x berturut-turut'}
          {settings.failure_threshold >= 4 && settings.failure_threshold <= 5 && 'Sedikit sabar — alert saat gagal beberapa kali'}
          {settings.failure_threshold > 5 && 'Santai — butuh banyak kegagalan sebelum alert'}
        </p>
      </div>

      {/* Check Interval */}
      <div className="space-y-2">
        <label className="block text-sm font-medium text-text-primary">
          Default Check Interval
        </label>
        <p className="text-xs text-text-muted">
          Interval default untuk device baru (dalam detik)
        </p>
        <div className="flex items-center gap-3">
          <input
            type="range"
            min={1}
            max={30}
            value={settings.check_interval}
            onChange={(e) => setSettings({ ...settings, check_interval: Number(e.target.value) })}
            className="flex-1 h-2 bg-surface-elevated rounded-lg appearance-none cursor-pointer accent-accent"
          />
          <span className="w-14 text-center text-sm font-mono font-semibold text-accent">
            {settings.check_interval}s
          </span>
        </div>
        <p className="text-[11px] text-text-muted">
          {settings.check_interval <= 2 && 'Sangat cepat — real-time tapi lebih berat'}
          {settings.check_interval >= 3 && settings.check_interval <= 5 && 'Default — balance antara responsif dan ringan'}
          {settings.check_interval >= 6 && settings.check_interval <= 10 && 'Sedang — lebih ringan untuk jaringan lambat'}
          {settings.check_interval > 10 && 'Lambat — sangat ringan, cocok untuk banyak device'}
        </p>
      </div>

      {/* Save Button */}
      <div className="flex items-center gap-3 pt-2">
        <button
          onClick={() => void handleSave()}
          disabled={saving}
          className="px-5 py-2.5 bg-accent text-white text-sm font-medium rounded-lg hover:bg-accent/90 transition-colors disabled:opacity-50 cursor-pointer"
        >
          {saving ? 'Menyimpan...' : 'Simpan Settings'}
        </button>
        <button
          onClick={() => setSettings({ failure_threshold: 3, check_interval: 3, notifications_enabled: true })}
          disabled={saving}
          className="px-4 py-2.5 bg-surface border border-border text-text-secondary text-sm font-medium rounded-lg hover:bg-surface-elevated transition-colors disabled:opacity-50 cursor-pointer"
        >
          Default
        </button>
      </div>
    </div>
  );
}

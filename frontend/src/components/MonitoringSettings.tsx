import { useCallback, useEffect, useState } from 'react';
import { getSettings, updateSettings, type AppSettings } from '../lib/api';
import { audioAlert } from '../lib/audioAlert';
import { getDesktopNotificationPermissionStatus, requestDesktopNotificationPermission } from '../lib/desktopNotify';
import { getNotificationClientSettings, saveNotificationClientSettings, type NotificationClientSettings } from '../lib/notificationSettings';

export function MonitoringSettings() {
  const [settings, setSettings] = useState<AppSettings | null>(null);
  const [clientSettings, setClientSettings] = useState<NotificationClientSettings>(getNotificationClientSettings());
  const [desktopPermission, setDesktopPermission] = useState(getDesktopNotificationPermissionStatus());
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [testingAudio, setTestingAudio] = useState(false);
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

  useEffect(() => {
    void loadSettings();
  }, [loadSettings]);

  const handleClientSettingChange = (key: keyof NotificationClientSettings, value: boolean | number) => {
    const updated = saveNotificationClientSettings({ [key]: value });
    setClientSettings(updated);
    if (key === 'volume') {
      audioAlert.setVolume(value as number);
    }
  };

  const handleTestAudio = async () => {
    setTestingAudio(true);
    await audioAlert.playTestSound();
    setTestingAudio(false);
  };

  const handleRequestPermission = async () => {
    const perm = await requestDesktopNotificationPermission();
    setDesktopPermission(perm);
  };

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

      <hr className="border-border/50" />

      {/* Emergency Alert Settings */}
      <div className="space-y-4">
        <div>
          <h3 className="text-sm font-semibold text-text-primary flex items-center gap-2">
            <span>🚨 Peringatan Darurat & Audio Alert</span>
          </h3>
          <p className="text-xs text-text-muted mt-1">
            Pengaturan notifikasi darurat langsung ketika status perangkat berubah menjadi Offline.
          </p>
        </div>

        {/* Toggle Block Screen Modal */}
        <div className="flex items-center justify-between p-3 rounded-lg bg-bg/50 border border-border/50">
          <div>
            <span className="text-sm font-medium text-text-primary block">Block Screen Modal Darurat</span>
            <span className="text-xs text-text-muted block">Tampilkan pop-up dialog merah di tengah layar (Default: Aktif)</span>
          </div>
          <button
            onClick={() => handleClientSettingChange('blockScreenEnabled', !clientSettings.blockScreenEnabled)}
            className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out ${
              clientSettings.blockScreenEnabled ? 'bg-danger' : 'bg-surface-elevated'
            }`}
          >
            <span
              className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                clientSettings.blockScreenEnabled ? 'translate-x-5' : 'translate-x-0'
              }`}
            />
          </button>
        </div>

        {/* Toggle Siren Audio */}
        <div className="flex items-center justify-between p-3 rounded-lg bg-bg/50 border border-border/50">
          <div>
            <span className="text-sm font-medium text-text-primary block">Suara Alarm Sirine</span>
            <span className="text-xs text-text-muted block">Bunyikan suara sirine darurat berulang sampai di-acknowledge (Default: Aktif)</span>
          </div>
          <button
            onClick={() => handleClientSettingChange('soundAlarmEnabled', !clientSettings.soundAlarmEnabled)}
            className={`relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out ${
              clientSettings.soundAlarmEnabled ? 'bg-danger' : 'bg-surface-elevated'
            }`}
          >
            <span
              className={`pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out ${
                clientSettings.soundAlarmEnabled ? 'translate-x-5' : 'translate-x-0'
              }`}
            />
          </button>
        </div>

        {/* Test Sound & Audio Volume */}
        {clientSettings.soundAlarmEnabled && (
          <div className="p-3 rounded-lg bg-surface/50 border border-border/50 space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-xs font-medium text-text-secondary">Volume Audio:</span>
              <input
                type="range"
                min={0.1}
                max={1}
                step={0.1}
                value={clientSettings.volume}
                onChange={(e) => handleClientSettingChange('volume', parseFloat(e.target.value))}
                className="w-32 h-1.5 bg-surface-elevated rounded-lg appearance-none cursor-pointer accent-accent"
              />
            </div>
            <button
              onClick={() => void handleTestAudio()}
              disabled={testingAudio}
              className="w-full py-2 px-3 rounded-md bg-amber-500/20 hover:bg-amber-500/30 text-amber-400 border border-amber-500/30 text-xs font-semibold transition-colors cursor-pointer flex items-center justify-center gap-2"
            >
              {testingAudio ? '🔔 Membunyikan Tes Suara...' : '🔊 Tes Suara Alarm (Tes Audio)'}
            </button>
          </div>
        )}

        {/* Desktop Push Notification Status (Mandatory Alert) */}
        <div className="p-3 rounded-lg bg-bg/50 border border-border/50 flex flex-col gap-2">
          <div className="flex items-center justify-between">
            <div>
              <span className="text-sm font-medium text-text-primary block">Desktop Push Notification (OS)</span>
              <span className="text-xs text-text-muted block">Alert Minimal Wajib (Notifikasi OS Windows/Linux/macOS)</span>
            </div>
            <span
              className={`px-2.5 py-1 rounded-full text-xs font-bold border ${
                desktopPermission === 'granted'
                  ? 'bg-success/15 text-success border-success/30'
                  : desktopPermission === 'denied'
                    ? 'bg-danger/15 text-danger border-danger/30'
                    : 'bg-warning/15 text-warning border-warning/30'
              }`}
            >
              {desktopPermission === 'granted'
                ? 'Aktif / Izinkan'
                : desktopPermission === 'denied'
                  ? 'Ditolak Browser'
                  : 'Belum Diizinkan'}
            </span>
          </div>

          {desktopPermission !== 'granted' && (
            <button
              onClick={() => void handleRequestPermission()}
              className="mt-1 w-full py-2 px-3 rounded-md bg-accent/20 hover:bg-accent/30 text-accent border border-accent/30 text-xs font-semibold transition-colors cursor-pointer"
            >
              🖥️ Minta Izin Desktop Push Notification Browser
            </button>
          )}
        </div>
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
          onClick={() => setSettings({ failure_threshold: 3, notifications_enabled: true })}
          disabled={saving}
          className="px-4 py-2.5 bg-surface border border-border text-text-secondary text-sm font-medium rounded-lg hover:bg-surface-elevated transition-colors disabled:opacity-50 cursor-pointer"
        >
          Default
        </button>
      </div>
    </div>
  );
}

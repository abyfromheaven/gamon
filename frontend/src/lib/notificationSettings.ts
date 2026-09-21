export interface NotificationClientSettings {
  soundAlarmEnabled: boolean;
  blockScreenEnabled: boolean;
  volume: number;
}

const STORAGE_KEY = 'gamon_notification_settings';

const DEFAULT_SETTINGS: NotificationClientSettings = {
  soundAlarmEnabled: true,
  blockScreenEnabled: true,
  volume: 0.7,
};

export function getNotificationClientSettings(): NotificationClientSettings {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      return { ...DEFAULT_SETTINGS, ...parsed };
    }
  } catch {
    // fallback to defaults
  }
  return DEFAULT_SETTINGS;
}

export function saveNotificationClientSettings(settings: Partial<NotificationClientSettings>): NotificationClientSettings {
  const current = getNotificationClientSettings();
  const updated = { ...current, ...settings };
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(updated));
  } catch {
    // ignore localstorage errors
  }
  return updated;
}

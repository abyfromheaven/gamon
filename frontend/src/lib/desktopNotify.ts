// Native Desktop Push Notification & Tab Title Flasher Utility

let titleInterval: ReturnType<typeof setInterval> | null = null;
let originalTitle = document.title || 'Gamon Monitoring';

export function getDesktopNotificationPermissionStatus(): NotificationPermission | 'unsupported' {
  if (!('Notification' in window)) return 'unsupported';
  return Notification.permission;
}

export async function requestDesktopNotificationPermission(): Promise<NotificationPermission | 'unsupported'> {
  if (!('Notification' in window)) {
    console.warn('Desktop notifications are not supported in this browser.');
    return 'unsupported';
  }
  try {
    const permission = await Notification.requestPermission();
    return permission;
  } catch (err) {
    console.error('Error requesting desktop notification permission:', err);
    return Notification.permission;
  }
}

export function sendDesktopNotification(title: string, body: string): boolean {
  if (!('Notification' in window)) return false;

  if (Notification.permission === 'granted') {
    try {
      const notification = new Notification(title, {
        body,
        icon: '/favicon.ico',
        tag: 'gamon-critical-alert',
        requireInteraction: true, // Keep notification on screen until user dismisses/clicks
      });

      notification.onclick = () => {
        window.focus();
        notification.close();
      };
      return true;
    } catch (err) {
      console.warn('Failed to dispatch desktop notification:', err);
      return false;
    }
  }
  return false;
}

export function startTabTitleFlash(alertMessage: string) {
  if (titleInterval) return;
  originalTitle = document.title;
  let toggle = false;

  titleInterval = setInterval(() => {
    document.title = toggle ? `🚨 [PERINGANTAN] ${alertMessage}` : `⚠️ PERIKSA GAMON DASHBOARD!`;
    toggle = !toggle;
  }, 1000);
}

export function stopTabTitleFlash() {
  if (titleInterval) {
    clearInterval(titleInterval);
    titleInterval = null;
    document.title = originalTitle;
  }
}

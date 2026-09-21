import type { AlertBannerData } from './AlertBanner';

interface EmergencyAlertModalProps {
  data: AlertBannerData;
  onAcknowledge: () => void;
}

export function EmergencyAlertModal({ data, onAcknowledge }: EmergencyAlertModalProps) {
  const isRecovery = data.new_status === 'online' && data.old_status === 'offline';

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6">
      {/* Red pulse backdrop vignette glow */}
      <div
        className="absolute inset-0 bg-red-950/80 backdrop-blur-md animate-pulse"
        style={{
          boxShadow: 'inset 0 0 100px rgba(220, 38, 38, 0.6)',
        }}
      />

      {/* Main Alert Card */}
      <div className="relative w-full max-w-lg bg-surface border-2 border-danger/80 rounded-2xl shadow-[0_0_50px_rgba(220,38,38,0.5)] overflow-hidden animate-fade-in-scale">
        {/* Top Hazard Stripe Bar */}
        <div className="h-3 bg-gradient-to-r from-danger via-amber-500 to-danger animate-pulse" />

        <div className="p-6 sm:p-8">
          {/* Header Icon & Title */}
          <div className="flex items-center gap-4 mb-5">
            <div className="relative flex items-center justify-center w-16 h-16 rounded-2xl bg-danger/20 border border-danger/40 text-danger shrink-0">
              <span className="absolute inset-0 rounded-2xl bg-danger/30 animate-ping opacity-75" />
              <svg className="w-10 h-10 relative z-10" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.2}>
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z"
                />
              </svg>
            </div>

            <div>
              <div className="inline-flex items-center gap-2 px-2.5 py-0.5 rounded-full text-xs font-bold bg-danger/20 text-danger border border-danger/30 uppercase tracking-wide mb-1">
                🚨 Critical Alert Event
              </div>
              <h2 className="text-xl font-extrabold text-text-primary">
                {isRecovery ? 'PERANGKAT PULIH' : 'PERANGKAT DOWN / OFFLINE'}
              </h2>
            </div>
          </div>

          {/* Alert Device Box */}
          <div className="bg-bg/80 border border-border rounded-xl p-4 space-y-3 mb-6">
            <div className="flex justify-between items-center border-b border-border/60 pb-2">
              <span className="text-xs text-text-muted">Nama Perangkat:</span>
              <span className="text-sm font-bold text-text-primary">{data.device_name}</span>
            </div>
            {data.device_ip && (
              <div className="flex justify-between items-center border-b border-border/60 pb-2">
                <span className="text-xs text-text-muted">Alamat IP:</span>
                <span className="text-sm font-mono font-semibold text-accent">{data.device_ip}</span>
              </div>
            )}
            <div className="flex justify-between items-center border-b border-border/60 pb-2">
              <span className="text-xs text-text-muted">Perubahan Status:</span>
              <span className="text-sm font-bold text-danger uppercase">
                {data.old_status} ➔ {data.new_status}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-xs text-text-muted">Waktu Kejadian:</span>
              <span className="text-xs font-mono text-text-secondary">{new Date(data.timestamp).toLocaleString('id-ID')}</span>
            </div>
          </div>

          {/* Sound Alarm Indicator */}
          <div className="flex items-center gap-2 text-xs text-amber-400 bg-amber-500/10 border border-amber-500/20 rounded-lg p-3 mb-6">
            <svg className="w-5 h-5 shrink-0 animate-bounce text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.536 8.464a5 5 0 010 7.072m2.828-9.9a9 9 0 010 12.728M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
            </svg>
            <span>
              Alarm sirine suara sedang berbunyi. Tekan <strong>Acknowledge</strong> untuk mematikan alarm.
            </span>
          </div>

          {/* Big Action Button */}
          <button
            onClick={onAcknowledge}
            className="w-full py-4 px-6 rounded-xl bg-danger hover:bg-red-700 active:scale-[0.98] text-white font-extrabold text-base tracking-wide shadow-lg shadow-danger/40 transition-all cursor-pointer flex items-center justify-center gap-2"
          >
            <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            ACKNOWLEDGE / TANGGAPI ALERT
          </button>
        </div>
      </div>
    </div>
  );
}

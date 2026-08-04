import { useCallback, useEffect, useState } from 'react';
import { generatePairingToken, getTelegramStatus, disconnectTelegram } from '../lib/api';

type PairingStep = 'idle' | 'waiting' | 'connected';

export function TelegramSettings() {
  const [status, setStatus] = useState<'loading' | 'connected' | 'disconnected'>('loading');
  const [chatID, setChatID] = useState('');
  const [pairedAt, setPairedAt] = useState<string | null>(null);
  const [step, setStep] = useState<PairingStep>('idle');
  const [token, setToken] = useState('');
  const [error, setError] = useState('');
  const [isCopied, setIsCopied] = useState(false);

  const loadStatus = useCallback(async () => {
    try {
      setError('');
      const data = await getTelegramStatus();
      if (data.status === 'connected') {
        setStatus('connected');
        setChatID(data.chat_id);
        setPairedAt(data.paired_at);
        setStep('connected');
      } else {
        setStatus('disconnected');
        setStep('idle');
      }
    } catch {
      setStatus('disconnected');
    }
  }, []);

  useEffect(() => { void loadStatus(); }, [loadStatus]);

  const handleConnect = async () => {
    try {
      setError('');
      const data = await generatePairingToken();
      setToken(data.token);
      setStep('waiting');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal generate token');
    }
  };

  const handleDisconnect = async () => {
    try {
      setError('');
      await disconnectTelegram();
      setStatus('disconnected');
      setStep('idle');
      setChatID('');
      setPairedAt(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal disconnect');
    }
  };

  const handleCopyToken = async () => {
    try {
      await navigator.clipboard.writeText(token);
      setIsCopied(true);
      setTimeout(() => setIsCopied(false), 2000);
    } catch {
      // Fallback: select text
      const el = document.getElementById('pairing-token');
      if (el) {
        const range = document.createRange();
        range.selectNodeContents(el);
        const sel = window.getSelection();
        sel?.removeAllRanges();
        sel?.addRange(range);
      }
    }
  };

  // Auto-refresh when waiting for pairing
  useEffect(() => {
    if (step !== 'waiting') return;
    const interval = setInterval(async () => {
      try {
        const data = await getTelegramStatus();
        if (data.status === 'connected') {
          setStatus('connected');
          setChatID(data.chat_id);
          setPairedAt(data.paired_at);
          setStep('connected');
        }
      } catch {
        // ignore
      }
    }, 2000);
    return () => clearInterval(interval);
  }, [step]);

  if (status === 'loading') {
    return (
      <div className="flex items-center justify-center py-8">
        <div className="w-5 h-5 border-2 border-accent border-t-transparent rounded-full animate-spin" />
        <span className="ml-3 text-sm text-text-muted">Loading...</span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {error && (
        <div className="rounded-lg bg-danger/10 border border-danger/20 p-3 text-sm text-danger">
          {error}
        </div>
      )}

      {/* Status */}
      <div className="flex items-center justify-between py-3 border-b border-border/30">
        <span className="text-sm text-text-muted">Status</span>
        {step === 'connected' ? (
          <span className="inline-flex items-center gap-2 text-sm font-medium text-success">
            <span className="w-2 h-2 rounded-full bg-success animate-breathe" />
            Connected
          </span>
        ) : (
          <span className="inline-flex items-center gap-2 text-sm font-medium text-text-muted">
            <span className="w-2 h-2 rounded-full bg-text-muted" />
            Not Connected
          </span>
        )}
      </div>

      {/* Connected Info */}
      {step === 'connected' && (
        <>
          <div className="flex items-center justify-between py-3 border-b border-border/30">
            <span className="text-sm text-text-muted">Chat ID</span>
            <span className="text-sm font-mono text-text-primary">{chatID}</span>
          </div>
          {pairedAt && (
            <div className="flex items-center justify-between py-3 border-b border-border/30">
              <span className="text-sm text-text-muted">Connected since</span>
              <span className="text-sm text-text-primary">{pairedAt}</span>
            </div>
          )}
          <button
            onClick={() => void handleDisconnect()}
            className="w-full px-4 py-2.5 bg-danger/10 border border-danger/30 text-danger text-sm font-medium rounded-lg hover:bg-danger/20 transition-colors cursor-pointer"
          >
            Disconnect Telegram
          </button>
        </>
      )}

      {/* Waiting for pairing */}
      {step === 'waiting' && (
        <div className="space-y-4">
          <div className="py-4 text-center">
            <div className="inline-flex items-center gap-2 text-sm text-accent">
              <div className="w-4 h-4 border-2 border-accent border-t-transparent rounded-full animate-spin" />
              Waiting for pairing...
            </div>
          </div>

          {/* Token Display */}
          <div>
            <label className="block text-[11px] text-text-muted uppercase tracking-wider font-semibold mb-2">
              Pairing Token
            </label>
            <div className="flex items-center gap-2">
              <div
                id="pairing-token"
                className="flex-1 px-4 py-3 bg-bg border border-border rounded-lg font-mono text-lg text-accent font-bold tracking-wider"
              >
                {token}
              </div>
              <button
                onClick={() => void handleCopyToken()}
                className="px-4 py-3 bg-surface-elevated border border-border rounded-lg text-sm font-medium text-text-primary hover:bg-surface transition-colors cursor-pointer"
              >
                {isCopied ? '✓' : 'Copy'}
              </button>
            </div>
          </div>

          {/* Instructions */}
          <div className="bg-bg/50 border border-border/30 rounded-lg p-4 space-y-3">
            <p className="text-sm font-semibold text-text-primary">Cara Pairing:</p>
            <ol className="space-y-2 text-sm text-text-secondary list-decimal list-inside">
              <li>Buka Bot Telegram @GardaMonitoringBot</li>
              <li>Kirim perintah: <code className="px-1.5 py-0.5 bg-surface-elevated rounded text-accent font-mono text-xs">/pair {token}</code></li>
              <li>Tunggu konfirmasi dari bot</li>
            </ol>
          </div>

          <p className="text-xs text-text-muted text-center">
            Token berlaku selama 30 hari
          </p>

          <button
            onClick={() => { setStep('idle'); setToken(''); }}
            className="w-full px-4 py-2.5 bg-surface border border-border text-text-secondary text-sm font-medium rounded-lg hover:bg-surface-elevated transition-colors cursor-pointer"
          >
            Cancel
          </button>
        </div>
      )}

      {/* Idle: Not connected */}
      {step === 'idle' && (
        <div className="py-6 text-center space-y-4">
          <p className="text-sm text-text-muted">
            Hubungkan Telegram kamu untuk menerima notifikasi alert secara real-time.
          </p>
          <button
            onClick={() => void handleConnect()}
            className="px-6 py-2.5 bg-accent text-white text-sm font-medium rounded-lg hover:bg-accent/90 transition-colors cursor-pointer"
          >
            Connect Telegram
          </button>
        </div>
      )}
    </div>
  );
}

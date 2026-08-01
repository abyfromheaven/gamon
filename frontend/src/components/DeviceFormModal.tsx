import { useEffect, useState } from 'react';
import type { Device, DeviceInput, DeviceType } from '../lib/api';

interface DeviceFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (device: DeviceInput) => Promise<void>;
  editDevice?: Device | null;
  isSaving: boolean;
}

const deviceTypes: DeviceType[] = ['Server', 'Router', 'Switch', 'Access Point', 'Website'];

interface FormData {
  name: string;
  type: DeviceType;
  ip: string;
  location: string;
  description: string;
  checkInterval: string;
}

const initialForm: FormData = { name: '', type: 'Server', ip: '', location: '', description: '', checkInterval: '3' };

export function DeviceFormModal({ isOpen, onClose, onSave, editDevice, isSaving }: DeviceFormModalProps) {
  const [form, setForm] = useState<FormData>(initialForm);
  const [error, setError] = useState('');

  useEffect(() => {
    setForm(editDevice ? {
      name: editDevice.name, type: editDevice.type, ip: editDevice.ip, location: editDevice.location,
      description: editDevice.description, checkInterval: String(editDevice.check_interval),
    } : initialForm);
    setError('');
  }, [editDevice, isOpen]);

  if (!isOpen) return null;

  const update = (field: keyof FormData, value: string) => setForm((current) => ({ ...current, [field]: value }));
  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    const interval = Number(form.checkInterval);
    if (!form.name.trim() || !form.ip.trim()) return setError('Nama device dan alamat IP wajib diisi.');
    if (!Number.isInteger(interval) || interval < 1) return setError('Interval harus berupa angka minimal 1 detik.');
    setError('');
    await onSave({ name: form.name.trim(), type: form.type, ip: form.ip.trim(), method: 'ICMP Ping', location: form.location.trim(), description: form.description.trim(), check_interval: interval });
  };
  const input = 'w-full px-3 py-2.5 bg-bg border border-border rounded-lg text-sm text-text-primary focus:outline-none focus:border-accent/50';

  return <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
    <button aria-label="Close" className="absolute inset-0 bg-bg/60 backdrop-blur-sm" onClick={onClose} />
    <form onSubmit={submit} className="relative w-full max-w-lg bg-surface border border-border rounded-xl shadow-2xl p-6 space-y-4">
      <div className="flex items-center justify-between"><h2 className="text-lg font-semibold">{editDevice ? 'Edit Device' : 'Add New Device'}</h2><span className="text-xs text-text-muted">ICMP Ping</span></div>
      {error && <p className="rounded bg-danger-muted p-3 text-sm text-danger">{error}</p>}
      <label className="block text-xs text-text-secondary">Device Name<input className={input} value={form.name} onChange={(e) => update('name', e.target.value)} /></label>
      <label className="block text-xs text-text-secondary">Device Type<select className={input} value={form.type} onChange={(e) => update('type', e.target.value)}>{deviceTypes.map((type) => <option key={type}>{type}</option>)}</select></label>
      <label className="block text-xs text-text-secondary">IP Address<input className={`${input} font-mono`} value={form.ip} onChange={(e) => update('ip', e.target.value)} placeholder="192.168.1.1" /></label>
      <div className="grid grid-cols-2 gap-4"><label className="block text-xs text-text-secondary">Check Interval (seconds)<input className={input} type="number" min="1" value={form.checkInterval} onChange={(e) => update('checkInterval', e.target.value)} /></label><label className="block text-xs text-text-secondary">Location (optional)<input className={input} value={form.location} onChange={(e) => update('location', e.target.value)} /></label></div>
      <label className="block text-xs text-text-secondary">Description (optional)<textarea className={input} value={form.description} onChange={(e) => update('description', e.target.value)} rows={3} /></label>
      <div className="flex justify-end gap-3"><button type="button" onClick={onClose} disabled={isSaving} className="px-4 py-2 text-sm text-text-secondary">Cancel</button><button disabled={isSaving} className="rounded-lg bg-accent px-5 py-2.5 text-sm font-semibold text-white disabled:opacity-60">{isSaving ? 'Saving...' : editDevice ? 'Save Changes' : 'Add Device'}</button></div>
    </form>
  </div>;
}

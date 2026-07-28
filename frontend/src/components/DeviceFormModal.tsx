import { useState, useEffect } from 'react';
import type { DeviceType } from '../types';
import type { Device } from '../data/devices';

interface DeviceFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (device: Omit<Device, 'id' | 'lastSeen'>) => void;
  editDevice?: Device | null;
}

const deviceTypes: DeviceType[] = ['Server', 'Router', 'Switch', 'Access Point', 'Website'];

interface FormData {
  name: string;
  type: DeviceType;
  ip: string;
  port: string;
  location: string;
  status: 'active' | 'inactive';
}

const initialFormData: FormData = {
  name: '',
  type: 'Server',
  ip: '',
  port: '',
  location: '',
  status: 'active',
};

export function DeviceFormModal({ isOpen, onClose, onSave, editDevice }: DeviceFormModalProps) {
  const [form, setForm] = useState<FormData>(initialFormData);
  const [errors, setErrors] = useState<Partial<Record<keyof FormData, string>>>({});

  useEffect(() => {
    if (editDevice) {
      setForm({
        name: editDevice.name,
        type: editDevice.type,
        ip: editDevice.ip,
        port: editDevice.port?.toString() ?? '',
        location: editDevice.location,
        status: editDevice.status,
      });
    } else {
      setForm(initialFormData);
    }
    setErrors({});
  }, [editDevice, isOpen]);

  const validate = (): boolean => {
    const newErrors: Partial<Record<keyof FormData, string>> = {};

    if (!form.name.trim()) {
      newErrors.name = 'Name is required';
    }

    if (!form.ip.trim()) {
      newErrors.ip = 'IP address is required';
    } else {
      const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/;
      const domainRegex = /^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$/;
      if (!ipRegex.test(form.ip.trim()) && !domainRegex.test(form.ip.trim())) {
        newErrors.ip = 'Invalid IP or domain format';
      }
    }

    if (!form.location.trim()) {
      newErrors.location = 'Location is required';
    }

    if (form.port && (isNaN(Number(form.port)) || Number(form.port) < 1 || Number(form.port) > 65535)) {
      newErrors.port = 'Port must be 1-65535';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    onSave({
      name: form.name.trim(),
      type: form.type,
      ip: form.ip.trim(),
      method: 'ICMP Ping',
      port: form.port ? Number(form.port) : null,
      location: form.location.trim(),
      status: form.status,
    });
  };

  const updateField = (field: keyof FormData, value: string) => {
    setForm((prev) => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: undefined }));
    }
  };

  if (!isOpen) return null;

  const inputBase = 'w-full px-3 py-2.5 bg-bg border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-accent/50 transition-colors';

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="absolute inset-0 bg-bg/60 backdrop-blur-sm animate-fade-in"
        onClick={onClose}
      />

      <div className="relative w-full max-w-lg bg-surface border border-border rounded-xl shadow-2xl animate-fade-in-scale">
        <div className="flex items-center justify-between px-6 py-4 border-b border-border/50">
          <h2 className="text-lg font-semibold text-text-primary">
            {editDevice ? 'Edit Device' : 'Add New Device'}
          </h2>
          <button
            onClick={onClose}
            className="p-1 rounded-md text-text-muted hover:text-text-primary hover:bg-surface-elevated transition-colors cursor-pointer"
          >
            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-xs font-medium text-text-secondary mb-1.5">Device Name</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => updateField('name', e.target.value)}
              placeholder="e.g. Core Router-01"
              className={`${inputBase} ${errors.name ? 'border-danger' : 'border-border'}`}
            />
            {errors.name && <p className="mt-1 text-xs text-danger">{errors.name}</p>}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1.5">Device Type</label>
              <select
                value={form.type}
                onChange={(e) => updateField('type', e.target.value)}
                className={`${inputBase} cursor-pointer`}
              >
                {deviceTypes.map((t) => (
                  <option key={t} value={t}>{t}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1.5">Status</label>
              <select
                value={form.status}
                onChange={(e) => updateField('status', e.target.value)}
                className={`${inputBase} cursor-pointer`}
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
              </select>
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-text-secondary mb-1.5">IP Address</label>
            <input
              type="text"
              value={form.ip}
              onChange={(e) => updateField('ip', e.target.value)}
              placeholder="e.g. 192.168.1.1"
              className={`${inputBase} font-mono ${errors.ip ? 'border-danger' : 'border-border'}`}
            />
            {errors.ip && <p className="mt-1 text-xs text-danger">{errors.ip}</p>}
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1.5">Method</label>
              <input
                type="text"
                value="ICMP Ping"
                disabled
                className={`${inputBase} font-mono text-text-muted border-border cursor-not-allowed opacity-60`}
              />
            </div>
            <div>
              <label className="block text-xs font-medium text-text-secondary mb-1.5">Port <span className="text-text-muted">(optional)</span></label>
              <input
                type="number"
                value={form.port}
                onChange={(e) => updateField('port', e.target.value)}
                placeholder="e.g. 443"
                min={1}
                max={65535}
                className={`${inputBase} font-mono ${errors.port ? 'border-danger' : 'border-border'}`}
              />
              {errors.port && <p className="mt-1 text-xs text-danger">{errors.port}</p>}
            </div>
          </div>

          <div>
            <label className="block text-xs font-medium text-text-secondary mb-1.5">Location / Ruangan</label>
            <input
              type="text"
              value={form.location}
              onChange={(e) => updateField('location', e.target.value)}
              placeholder="e.g. Ruang Server Utama"
              className={`${inputBase} ${errors.location ? 'border-danger' : 'border-border'}`}
            />
            {errors.location && <p className="mt-1 text-xs text-danger">{errors.location}</p>}
          </div>

          <div className="flex items-center justify-end gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2.5 rounded-lg text-sm font-medium text-text-secondary hover:text-text-primary hover:bg-surface-elevated transition-colors cursor-pointer"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-5 py-2.5 bg-accent hover:bg-accent/90 text-white text-sm font-semibold rounded-lg transition-colors cursor-pointer"
            >
              {editDevice ? 'Save Changes' : 'Add Device'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

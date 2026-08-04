import { useEffect, useMemo, useState } from 'react';
import type { Device, DeviceInput, DeviceType } from '../lib/api';
import { createDevice, deleteDevice, fetchDevices, updateDevice, toggleDeviceStatus } from '../lib/api';
import { PageHeader } from '../components/PageHeader';
import { SearchBar } from '../components/SearchBar';
import { DeviceTable } from '../components/DeviceTable';
import { DeviceFormModal } from '../components/DeviceFormModal';
import { ConfirmDialog } from '../components/ConfirmDialog';

type FilterType = 'All' | DeviceType;

export function DeviceManagementPage() {
  const [devices, setDevices] = useState<Device[]>([]);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<FilterType>('All');
  const [editing, setEditing] = useState<Device | null>(null);
  const [deleting, setDeleting] = useState<Device | null>(null);
  const [isFormOpen, setFormOpen] = useState(false);
  const [isSaving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [toast, setToast] = useState<string | null>(null);

  const load = async () => {
    try {
      setError('');
      setDevices(await fetchDevices());
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal memuat device.');
    }
  };

  useEffect(() => { void load(); }, []);

  const filtered = useMemo(
    () => devices.filter(
      (device) =>
        (filter === 'All' || device.type === filter) &&
        (!search || device.name.toLowerCase().includes(search.toLowerCase()) || device.ip.includes(search))
    ),
    [devices, filter, search]
  );

  const save = async (input: DeviceInput) => {
    setSaving(true);
    try {
      const saved = editing
        ? await updateDevice(editing.id, input)
        : await createDevice(input);
      setDevices((current) =>
        editing
          ? current.map((device) => device.id === saved.id ? saved : device)
          : [saved, ...current]
      );
      setFormOpen(false);
      setEditing(null);
      setToast(editing ? `Device "${saved.name}" berhasil diupdate` : `Device "${saved.name}" berhasil ditambahkan`);
      setTimeout(() => setToast(null), 3000);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal menyimpan device.');
    } finally {
      setSaving(false);
    }
  };

  const confirmDelete = async () => {
    if (!deleting) return;
    try {
      await deleteDevice(deleting.id);
      setDevices((current) => current.filter((device) => device.id !== deleting.id));
      setDeleting(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal menghapus device.');
    }
  };

  const toggleDevice = async (device: Device) => {
    const newStatus = device.status === 'active' ? 'inactive' : 'active';
    try {
      await toggleDeviceStatus(device.id, newStatus);
      setDevices((current) =>
        current.map((d) => d.id === device.id ? { ...d, status: newStatus } : d)
      );
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Gagal mengubah status device.');
    }
  };

  const active = devices.filter((device) => device.status === 'active').length;

  return (
    <div className="min-h-full">
      <div className="max-w-[1200px] mx-auto px-4 lg:px-8 py-6 lg:py-8">
        <PageHeader
          title="Device Management"
          subtitle={`${devices.length} devices registered · ${active} active`}
          actionLabel="Add Device"
          onAction={() => { setEditing(null); setFormOpen(true); }}
        />
        {error && <p className="mb-4 rounded bg-danger-muted p-3 text-sm text-danger">{error}</p>}
        {toast && (
          <div className="mb-4 rounded bg-surface border border-border p-3 text-sm text-text-primary animate-fade-in">
            {toast}
          </div>
        )}
        <SearchBar search={search} onSearchChange={setSearch} activeFilter={filter} onFilterChange={setFilter} />
        <DeviceTable
          devices={filtered}
          onEdit={(device) => { setEditing(device); setFormOpen(true); }}
          onDelete={setDeleting}
          onToggle={toggleDevice}
        />
      </div>
      <DeviceFormModal
        isOpen={isFormOpen}
        onClose={() => { setFormOpen(false); setEditing(null); }}
        onSave={save}
        editDevice={editing}
        isSaving={isSaving}
      />
      <ConfirmDialog
        isOpen={!!deleting}
        onClose={() => setDeleting(null)}
        onConfirm={() => void confirmDelete()}
        title="Delete Device"
        message={`Are you sure you want to delete "${deleting?.name}"? This action cannot be undone.`}
        confirmLabel="Delete"
      />
    </div>
  );
}

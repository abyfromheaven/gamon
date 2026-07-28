import { useState, useMemo } from 'react';
import type { DeviceType } from '../types';
import { devices as initialDevices } from '../data/devices';
import type { Device } from '../data/devices';
import { PageHeader } from '../components/PageHeader';
import { SearchBar } from '../components/SearchBar';
import { DeviceTable } from '../components/DeviceTable';
import { DeviceFormModal } from '../components/DeviceFormModal';
import { ConfirmDialog } from '../components/ConfirmDialog';

type FilterType = 'All' | DeviceType;

export function DeviceManagementPage() {
  const [devices, setDevices] = useState<Device[]>(initialDevices);
  const [search, setSearch] = useState('');
  const [filter, setFilter] = useState<FilterType>('All');
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingDevice, setEditingDevice] = useState<Device | null>(null);
  const [deletingDevice, setDeletingDevice] = useState<Device | null>(null);

  const filteredDevices = useMemo(() => {
    return devices.filter((d) => {
      const matchesSearch = search === '' ||
        d.name.toLowerCase().includes(search.toLowerCase()) ||
        d.ip.includes(search);
      const matchesFilter = filter === 'All' || d.type === filter;
      return matchesSearch && matchesFilter;
    });
  }, [devices, search, filter]);

  const handleAdd = () => {
    setEditingDevice(null);
    setIsFormOpen(true);
  };

  const handleEdit = (device: Device) => {
    setEditingDevice(device);
    setIsFormOpen(true);
  };

  const handleDelete = (device: Device) => {
    setDeletingDevice(device);
  };

  const handleSave = (data: Omit<Device, 'id' | 'lastSeen'>) => {
    if (editingDevice) {
      setDevices((prev) =>
        prev.map((d) => (d.id === editingDevice.id ? { ...d, ...data } : d))
      );
    } else {
      const newDevice: Device = {
        ...data,
        id: Math.max(...devices.map((d) => d.id), 0) + 1,
        lastSeen: null,
      };
      setDevices((prev) => [newDevice, ...prev]);
    }
    setIsFormOpen(false);
    setEditingDevice(null);
  };

  const handleConfirmDelete = () => {
    if (deletingDevice) {
      setDevices((prev) => prev.filter((d) => d.id !== deletingDevice.id));
      setDeletingDevice(null);
    }
  };

  const activeCount = devices.filter((d) => d.status === 'active').length;

  return (
    <div className="min-h-full">
      <div className="max-w-[1200px] mx-auto px-4 lg:px-8 py-6 lg:py-8">
        <PageHeader
          title="Device Management"
          subtitle={`${devices.length} devices registered · ${activeCount} active`}
          actionLabel="Add Device"
          onAction={handleAdd}
        />

        <SearchBar
          search={search}
          onSearchChange={setSearch}
          activeFilter={filter}
          onFilterChange={setFilter}
        />

        <DeviceTable
          devices={filteredDevices}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      </div>

      <DeviceFormModal
        isOpen={isFormOpen}
        onClose={() => { setIsFormOpen(false); setEditingDevice(null); }}
        onSave={handleSave}
        editDevice={editingDevice}
      />

      <ConfirmDialog
        isOpen={!!deletingDevice}
        onClose={() => setDeletingDevice(null)}
        onConfirm={handleConfirmDelete}
        title="Delete Device"
        message={`Are you sure you want to delete "${deletingDevice?.name}"? This action cannot be undone.`}
        confirmLabel="Delete"
      />
    </div>
  );
}

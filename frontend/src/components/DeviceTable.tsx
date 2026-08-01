import type { Device } from '../lib/api';
import { StatusBadge } from './StatusBadge';

interface DeviceTableProps {
  devices: Device[];
  onEdit: (device: Device) => void;
  onDelete: (device: Device) => void;
  onToggle: (device: Device) => void;
}

export function DeviceTable({ devices, onEdit, onDelete, onToggle }: DeviceTableProps) {
  if (devices.length === 0) {
    return (
      <div className="animate-fade-in bg-surface border border-border rounded-xl p-12 text-center">
        <div className="text-text-muted text-4xl mb-3">○</div>
        <p className="text-text-secondary font-medium">No devices found</p>
        <p className="text-text-muted text-sm mt-1">Try adjusting your search or filter</p>
      </div>
    );
  }

  return (
    <div className="animate-fade-in-up anim-delay-2 bg-surface border border-border rounded-xl overflow-hidden">
      {/* Desktop Table */}
      <div className="hidden lg:block overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border/50">
              <th className="text-left text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">Status</th>
              <th className="text-left text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">Name</th>
              <th className="text-left text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">Type</th>
              <th className="text-left text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">IP Address</th>
              <th className="text-left text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">Method</th>
              <th className="text-left text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">Port</th>
              <th className="text-left text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">Location</th>
              <th className="text-right text-[10px] uppercase tracking-widest text-text-muted font-medium px-5 py-3">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/30">
            {devices.map((device) => (
              <tr
                key={device.id}
                className="hover:bg-surface-elevated/50 transition-colors duration-150"
              >
                <td className="px-5 py-3.5">
                  <StatusBadge status={device.status} />
                </td>
                <td className="px-5 py-3.5">
                  <span className="text-sm font-medium text-text-primary">{device.name}</span>
                </td>
                <td className="px-5 py-3.5">
                  <span className="inline-flex px-2 py-0.5 rounded text-[11px] font-mono font-medium bg-bg/50 text-text-secondary">
                    {device.type}
                  </span>
                </td>
                <td className="px-5 py-3.5">
                  <span className="font-mono text-sm text-text-primary">{device.ip}</span>
                </td>
                <td className="px-5 py-3.5">
                  <span className="text-sm text-text-muted">{device.method}</span>
                </td>
                <td className="px-5 py-3.5">
                  <span className="font-mono text-sm text-text-muted">
                    {device.port ?? '—'}
                  </span>
                </td>
                <td className="px-5 py-3.5">
                  <span className="text-sm text-text-secondary">{device.location}</span>
                </td>
                <td className="px-5 py-3.5">
                  <div className="flex items-center justify-end gap-1">
                    <button
                      onClick={() => onToggle(device)}
                      className={`p-1.5 rounded-md transition-colors duration-150 cursor-pointer ${device.status === 'active' ? 'text-success hover:bg-success-muted' : 'text-text-muted hover:bg-surface-elevated'}`}
                      title={device.status === 'active' ? 'Nonaktifkan' : 'Aktifkan'}
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M5.636 5.636a9 9 0 1012.728 0M12 3v9" />
                      </svg>
                    </button>
                    <button
                      onClick={() => onEdit(device)}
                      className="p-1.5 rounded-md text-text-muted hover:text-accent hover:bg-accent-muted transition-colors duration-150 cursor-pointer"
                      title="Edit device"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                      </svg>
                    </button>
                    <button
                      onClick={() => onDelete(device)}
                      className="p-1.5 rounded-md text-text-muted hover:text-danger hover:bg-danger-muted transition-colors duration-150 cursor-pointer"
                      title="Delete device"
                    >
                      <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                      </svg>
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Mobile Cards */}
      <div className="lg:hidden divide-y divide-border/30">
        {devices.map((device) => (
          <div key={device.id} className="p-4 space-y-3">
            <div className="flex items-start justify-between">
              <div className="flex items-center gap-2">
                <StatusBadge status={device.status} />
                <span className="text-sm font-medium text-text-primary">{device.name}</span>
              </div>
              <div className="flex items-center gap-1">
                <button
                  onClick={() => onToggle(device)}
                  className={`p-1.5 rounded-md transition-colors cursor-pointer ${device.status === 'active' ? 'text-success hover:bg-success-muted' : 'text-text-muted hover:bg-surface-elevated'}`}
                  title={device.status === 'active' ? 'Nonaktifkan' : 'Aktifkan'}
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M5.636 5.636a9 9 0 1012.728 0M12 3v9" />
                  </svg>
                </button>
                <button
                  onClick={() => onEdit(device)}
                  className="p-1.5 rounded-md text-text-muted hover:text-accent hover:bg-accent-muted transition-colors cursor-pointer"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                  </svg>
                </button>
                <button
                  onClick={() => onDelete(device)}
                  className="p-1.5 rounded-md text-text-muted hover:text-danger hover:bg-danger-muted transition-colors cursor-pointer"
                >
                  <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                  </svg>
                </button>
              </div>
            </div>
            <div className="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
              <div>
                <span className="text-text-muted">Type</span>
                <span className="ml-2 font-mono text-text-secondary">{device.type}</span>
              </div>
              <div>
                <span className="text-text-muted">IP</span>
                <span className="ml-2 font-mono text-text-primary">{device.ip}</span>
              </div>
              <div>
                <span className="text-text-muted">Method</span>
                <span className="ml-2 text-text-secondary">{device.method}</span>
              </div>
              <div>
                <span className="text-text-muted">Port</span>
                <span className="ml-2 font-mono text-text-secondary">{device.port ?? '—'}</span>
              </div>
              <div className="col-span-2">
                <span className="text-text-muted">Location</span>
                <span className="ml-2 text-text-secondary">{device.location}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

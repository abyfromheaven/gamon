import type { AlertStatus, AlertSeverity, DeviceType } from '../types';

interface AlertFiltersProps {
  searchQuery: string;
  onSearchChange: (query: string) => void;
  statusFilter: AlertStatus | 'all';
  onStatusFilterChange: (status: AlertStatus | 'all') => void;
  severityFilter: AlertSeverity | 'all';
  onSeverityFilterChange: (severity: AlertSeverity | 'all') => void;
  deviceTypeFilter: DeviceType | 'all';
  onDeviceTypeFilterChange: (type: DeviceType | 'all') => void;
}

const statusOptions: { value: AlertStatus | 'all'; label: string }[] = [
  { value: 'all', label: 'All Status' },
  { value: 'ongoing', label: 'Ongoing' },
  { value: 'resolved', label: 'Resolved' },
];

const severityOptions: { value: AlertSeverity | 'all'; label: string }[] = [
  { value: 'all', label: 'All Severity' },
  { value: 'critical', label: 'Critical' },
];

const deviceTypeOptions: { value: DeviceType | 'all'; label: string }[] = [
  { value: 'all', label: 'All Devices' },
  { value: 'Server', label: 'Server' },
  { value: 'Router', label: 'Router' },
  { value: 'Switch', label: 'Switch' },
  { value: 'Access Point', label: 'Access Point' },
  { value: 'Website', label: 'Website' },
];

export function AlertFilters({
  searchQuery,
  onSearchChange,
  statusFilter,
  onStatusFilterChange,
  severityFilter,
  onSeverityFilterChange,
  deviceTypeFilter,
  onDeviceTypeFilterChange,
}: AlertFiltersProps) {
  return (
    <div className="animate-fade-in-up anim-delay-4 flex flex-col lg:flex-row gap-3">
      {/* Search */}
      <div className="relative flex-1">
        <svg
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={1.5}
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" />
        </svg>
        <input
          type="text"
          placeholder="Search alerts..."
          value={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 bg-surface border border-border rounded-lg text-sm text-text-primary placeholder:text-text-muted focus:outline-none focus:border-accent/50 transition-colors"
        />
      </div>

      {/* Filters */}
      <div className="flex gap-2 flex-wrap">
        <select
          value={statusFilter}
          onChange={(e) => onStatusFilterChange(e.target.value as AlertStatus | 'all')}
          className="px-3 py-2.5 bg-surface border border-border rounded-lg text-sm text-text-secondary focus:outline-none focus:border-accent/50 transition-colors cursor-pointer"
        >
          {statusOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>{opt.label}</option>
          ))}
        </select>

        <select
          value={severityFilter}
          onChange={(e) => onSeverityFilterChange(e.target.value as AlertSeverity | 'all')}
          className="px-3 py-2.5 bg-surface border border-border rounded-lg text-sm text-text-secondary focus:outline-none focus:border-accent/50 transition-colors cursor-pointer"
        >
          {severityOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>{opt.label}</option>
          ))}
        </select>

        <select
          value={deviceTypeFilter}
          onChange={(e) => onDeviceTypeFilterChange(e.target.value as DeviceType | 'all')}
          className="px-3 py-2.5 bg-surface border border-border rounded-lg text-sm text-text-secondary focus:outline-none focus:border-accent/50 transition-colors cursor-pointer"
        >
          {deviceTypeOptions.map((opt) => (
            <option key={opt.value} value={opt.value}>{opt.label}</option>
          ))}
        </select>
      </div>
    </div>
  );
}

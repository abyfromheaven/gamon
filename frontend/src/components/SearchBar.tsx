type FilterType = 'All' | string;

interface SearchBarProps {
  search: string;
  onSearchChange: (value: string) => void;
  activeFilter: FilterType;
  onFilterChange: (filter: FilterType) => void;
  availableTypes: string[];
}

export function SearchBar({ search, onSearchChange, activeFilter, onFilterChange, availableTypes }: SearchBarProps) {
  return (
    <div className="animate-fade-in-up anim-delay-1 flex flex-col sm:flex-row gap-3 mb-6">
      {/* Search Input */}
      <div className="relative flex-1">
        <svg
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          placeholder="Search by name or IP..."
          className="w-full pl-10 pr-4 py-2.5 bg-surface border border-border rounded-lg text-sm text-text-primary placeholder-text-muted focus:outline-none focus:border-accent/50 transition-colors"
        />
      </div>

      {/* Filter Dropdown */}
      <select
        value={activeFilter}
        onChange={(e) => onFilterChange(e.target.value as FilterType)}
        className="px-3 py-2.5 bg-surface border border-border rounded-lg text-sm text-text-secondary focus:outline-none focus:border-accent/50 transition-colors cursor-pointer"
      >
        <option value="All">All Types</option>
        {availableTypes.map((type) => (
          <option key={type} value={type}>{type}</option>
        ))}
      </select>
    </div>
  );
}

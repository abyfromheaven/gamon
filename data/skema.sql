CREATE TABLE devices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			ip TEXT NOT NULL,
			url TEXT DEFAULT '',
			port INTEGER,
			method TEXT NOT NULL DEFAULT 'ICMP Ping',
			location TEXT DEFAULT '',
			check_interval INTEGER DEFAULT 3,
			status TEXT DEFAULT 'active',
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
CREATE TABLE sqlite_sequence(name,seq);
CREATE TABLE ping_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			status TEXT NOT NULL,
			latency_ms REAL DEFAULT 0,
			ttl INTEGER DEFAULT 0,
			seq INTEGER DEFAULT 0,
			details TEXT DEFAULT '{}',
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
		);
CREATE TABLE alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'ongoing',
			severity TEXT NOT NULL DEFAULT 'low',
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			resolved_at DATETIME,
			description TEXT DEFAULT '', alert_type TEXT DEFAULT 'critical', acknowledged BOOLEAN DEFAULT FALSE, acknowledged_at DATETIME,
			FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
		);
CREATE INDEX idx_ping_history_device_id ON ping_history(device_id);
CREATE INDEX idx_ping_history_timestamp ON ping_history(timestamp);
CREATE INDEX idx_alerts_device_id ON alerts(device_id);
CREATE INDEX idx_alerts_status ON alerts(status);
CREATE TABLE telegram_pairing (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			token TEXT NOT NULL UNIQUE,
			chat_id TEXT DEFAULT '',
			status TEXT DEFAULT 'pending',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME NOT NULL,
			paired_at DATETIME
		);
CREATE TABLE settings (
			key TEXT PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

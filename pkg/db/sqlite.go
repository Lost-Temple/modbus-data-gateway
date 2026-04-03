package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gateway/pkg/models"
	"encoding/json"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS history_data (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id TEXT NOT NULL,
		timestamp INTEGER NOT NULL,
		flow_instant_m3h REAL DEFAULT 0,
		forward_total_m3 REAL DEFAULT 0,
		reverse_total_m3 REAL DEFAULT 0,
		pressure REAL DEFAULT 0,
		power_comm_v REAL DEFAULT 0,
		power_meter_v REAL DEFAULT 0,
		st_word INTEGER DEFAULT 0,
		payload_json TEXT,
		last_ack_ts INTEGER DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_history_device_ts ON history_data(device_id, timestamp);
	CREATE INDEX IF NOT EXISTS idx_history_ack ON history_data(last_ack_ts);
	`
	_, err := db.Exec(schema)
	return err
}

func (s *Store) Save(data models.NormalizedData) error {
	query := `
		INSERT INTO history_data (
			device_id, timestamp, flow_instant_m3h, forward_total_m3,
			reverse_total_m3, pressure, power_comm_v, power_meter_v,
			st_word, payload_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	payload, err := json.Marshal(data.Tags)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	_, err = s.db.Exec(query,
		data.DeviceID,
		data.Timestamp.Unix(),
		data.FlowInstantM3H,
		data.ForwardTotalM3,
		data.ReverseTotalM3,
		data.Pressure,
		data.PowerCommV,
		data.PowerMeterV,
		data.StatusWord,
		string(payload),
	)

	return err
}

func (s *Store) Close() error {
	return s.db.Close()
}

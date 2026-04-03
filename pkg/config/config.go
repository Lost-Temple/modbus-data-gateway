package config

// Config represents the root configuration
type Config struct {
	Ingest  IngestConfig  `yaml:"ingest"`
	Device  DeviceConfig  `yaml:"device"`
	Hzws    HzwsConfig    `yaml:"hzws"`
	Routes  []RouteConfig `yaml:"routes"`
	MQTT    *MQTTConfig   `yaml:"mqtt,omitempty"`
}

type IngestConfig struct {
	Mode               string                   `yaml:"mode"`
	SpcModbusTcpSlave  *SpcModbusTcpSlaveConfig `yaml:"spc_modbus_tcp_slave,omitempty"`
	SpcHttpPush        *SpcHttpPushConfig       `yaml:"spc_http_push,omitempty"`
	DirectModbusRtu    *DirectModbusRtuConfig   `yaml:"direct_modbus_rtu,omitempty"`
	DirectModbusTcp    *DirectModbusTcpConfig   `yaml:"direct_modbus_tcp,omitempty"`
	SpcSqlite          *SpcSqliteConfig         `yaml:"spc_sqlite,omitempty"`
}

type SpcModbusTcpSlaveConfig struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	UnitID int    `yaml:"unit_id"`
}

type SpcHttpPushConfig struct {
	ListenHost string `yaml:"listen_host"`
	ListenPort int    `yaml:"listen_port"`
	Path       string `yaml:"path"`
}

type DirectModbusRtuConfig struct {
	Port     string `yaml:"port"`
	Baudrate int    `yaml:"baudrate"`
	Databits int    `yaml:"databits"`
	Parity   string `yaml:"parity"`
	Stopbits int    `yaml:"stopbits"`
}

type DirectModbusTcpConfig struct {
	Devices []struct {
		Name                string `yaml:"name"`
		Host                string `yaml:"host"`
		Port                int    `yaml:"port"`
		UnitID              int    `yaml:"unit_id"`
		PollIntervalSeconds int    `yaml:"poll_interval_seconds"`
		PointsRef           string `yaml:"points_ref"`
	} `yaml:"devices"`
}

type SpcSqliteConfig struct {
	DbPath string `yaml:"db_path"`
}

type DeviceConfig struct {
	IntervalSeconds       int `yaml:"interval_seconds"`
	ReportIntervalMinutes int `yaml:"report_interval_minutes"`
}

type HzwsConfig struct {
	ServerHost   string `yaml:"server_host"`
	ServerPort   int    `yaml:"server_port"`
	MeterIdBcd8  string `yaml:"meter_id_bcd8"`
	SimIdType    int    `yaml:"sim_id_type"`
	SimAscii     string `yaml:"sim_ascii"`
}

type RouteConfig struct {
	Name      string `yaml:"name"`
	Match     struct {
		DeviceTags []string `yaml:"device_tags"`
	} `yaml:"match"`
	Transform struct {
		Profile string `yaml:"profile"`
	} `yaml:"transform"`
	Outputs []struct {
		Type   string `yaml:"type"`
		Target string `yaml:"target"`
	} `yaml:"outputs"`
}

type MQTTConfig struct {
	Broker    string `yaml:"broker"`
	Port      int    `yaml:"port"`
	FuelTopic string `yaml:"fuel_topic"`
}

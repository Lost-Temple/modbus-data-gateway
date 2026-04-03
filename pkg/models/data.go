package models

import "time"

// NormalizedData 代表归一化后的统一数据模型
type NormalizedData struct {
	Timestamp      time.Time         `json:"ts"`
	DeviceID       string            `json:"device_id"`
	FlowInstantM3H float64           `json:"flow_instant_m3h"`
	ForwardTotalM3 float64           `json:"forward_total_m3"`
	ReverseTotalM3 float64           `json:"reverse_total_m3"`
	Pressure       float64           `json:"pressure"`
	PowerCommV     float64           `json:"power_comm_v"`
	PowerMeterV    float64           `json:"power_meter_v"`
	StatusWord     uint32            `json:"st_word"`
	Tags           map[string]string `json:"tags"`
}

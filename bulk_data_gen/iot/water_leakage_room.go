package iot

import (
	. "github.com/taosdata/timeseriesdatabase-comparisons/bulk_data_gen/common"
	"time"
)

var (
	WaterLeakageRoomByteString = []byte("water_leakage_room") // heap optimization
)

var (
	// Field keys for 'air condition indoor' points.
	WaterLeakageRoomFieldKeys = [][]byte{
		[]byte("leakage"),
		[]byte("battery_voltage"),
	}
)

type WaterLeakageRoomMeasurement struct {
	sensorId      []byte
	roomId        []byte
	timestamp     time.Time
	distributions []Distribution
}

func NewWaterLeakageRoomMeasurement(start time.Time, roomId []byte, sensorId []byte) *WaterLeakageRoomMeasurement {
	distributions := make([]Distribution, len(WaterLeakageRoomFieldKeys))
	//state
	distributions[0] = TSD(0, 1, 0)
	//battery_voltage
	distributions[1] = MUDWD(ND(0.01, 0.005), 1, 3.2, 3.2)

	return &WaterLeakageRoomMeasurement{
		timestamp:     start,
		distributions: distributions,
		sensorId:      sensorId,
		roomId:        roomId,
	}
}

func (m *WaterLeakageRoomMeasurement) Tick(d time.Duration) {
	m.timestamp = m.timestamp.Add(d)
	for i := range m.distributions {
		m.distributions[i].Advance()
	}
}

func (m *WaterLeakageRoomMeasurement) ToPoint(p *Point) bool {
	p.SetMeasurementName(WaterLeakageRoomByteString)
	p.SetTimestamp(&m.timestamp)
	p.AppendTag(SensorHomeTagKeys[0], m.sensorId)
	p.AppendTag(RoomTagKey, m.roomId)
	for i := range m.distributions {
		p.AppendField(WaterLeakageRoomFieldKeys[i], m.distributions[i].Get())
	}
	return true
}

package devops

import (
	. "github.com/taosdata/timeseriesdatabase-comparisons/bulk_data_gen/common"
	"math/rand"
	"time"
)

var (
	KernelByteString   = []byte("kernel") // heap optimization
	BootTimeByteString = []byte("boot_time")
	KernelFields       = []LabeledDistributionMaker{
		{[]byte("interrupts"), func() Distribution { return MWD(ND(5, 1), 0) }},
		{[]byte("context_switches"), func() Distribution { return MWD(ND(5, 1), 0) }},
		{[]byte("processes_forked"), func() Distribution { return MWD(ND(5, 1), 0) }},
		{[]byte("disk_pages_in"), func() Distribution { return MWD(ND(5, 1), 0) }},
		{[]byte("disk_pages_out"), func() Distribution { return MWD(ND(5, 1), 0) }},
	}
)

type KernelMeasurement struct {
	timestamp time.Time

	bootTime      int64
	uptime        time.Duration
	distributions []Distribution
}

func NewKernelMeasurement(start time.Time) *KernelMeasurement {
	distributions := make([]Distribution, len(KernelFields))
	for i := range KernelFields {
		distributions[i] = KernelFields[i].DistributionMaker()
	}

	bootTime := rand.Int63n(240)
	return &KernelMeasurement{
		bootTime: bootTime,

		timestamp:     start,
		distributions: distributions,
	}
}

func (m *KernelMeasurement) Tick(d time.Duration) {
	m.timestamp = m.timestamp.Add(d)

	for i := range m.distributions {
		m.distributions[i].Advance()
	}
}

func (m *KernelMeasurement) ToPoint(p *Point) bool {
	p.SetMeasurementName(KernelByteString)
	p.SetTimestamp(&m.timestamp)

	p.AppendField(BootTimeByteString, m.bootTime)
	for i := range m.distributions {
		p.AppendField(KernelFields[i].Label, int64(m.distributions[i].Get()))
	}
	return true
}

package ph

import (
	"fmt"
	"github.com/reef-pi/reef-pi/controller/utils"
)

const ReadingsBucket = "ph_readings"

type Measurement struct {
	Ph   float64        `json:"pH"`
	Time utils.TeleTime `json:"time"`
	sum  float64
	len  int
}

func (m1 Measurement) Rollup(mx utils.Metric) (utils.Metric, bool) {
	m2 := mx.(Measurement)
	m := Measurement{Time: m1.Time, Ph: m1.Ph, sum: m1.sum, len: m1.len}
	if m1.Time.Hour() == m2.Time.Hour() {
		m.sum += m2.Ph
		m.len += 1
		m.Ph = m.sum / float64(m.len)
		return m, false
	}
	return m2, true
}

func (m1 Measurement) Before(mx utils.Metric) bool {
	m2 := mx.(Measurement)
	return m1.Time.Before(m2.Time)
}

func notifyIfNeeded(t *utils.Telemetry, name string, n Notify, reading float64) {
	if !n.Enable {
		return
	}
	subject := "[Reef-Pi ALERT] ph out of range"
	format := "Current ph value from probe '%s' (%f) is out of acceptable range ( %f -%f )"
	body := fmt.Sprintf(format, reading, name, n.Min, n.Max)
	if reading >= n.Max {
		t.Alert(subject, "Tank ph is high. "+body)
		return
	}
	if reading <= n.Min {
		t.Alert(subject, "Tank ph is low. "+body)
		return
	}
}

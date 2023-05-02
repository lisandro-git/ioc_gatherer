package cowrieParser

import (
	"github.com/prometheus/client_golang/prometheus"
)

type IPConnection struct {
	connections *prometheus.Desc
}

// NewIPConnection returns a new CpuUsage collector
func NewIPConnection() *IPConnection {
	return &IPConnection{
		connections: prometheus.NewDesc("connectedIP", "IPs used to connect to the honeypots", []string{"honeypot", "src_ip", "proto", "timestamp", "event_id", "latitude", "longitude", "country"}, nil),
	}
}

func (c *IPConnection) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.connections
}

func (c *IPConnection) Collect(ch chan<- prometheus.Metric) {
	//var data *IOCs = ReadFile("./data_files/cowrie.json")

	/*for i := 0; i < len(data.IocGroup); i++ {
		ch <- prometheus.MustNewConstMetric(
			c.connections,
			prometheus.CounterValue,
			float64(data.IocGroup[i].HitCount),
			"cowrie",
			data.IocGroup[i].Src_ip,
			data.IocGroup[i].Proto,
			data.IocGroup[i].Timestamp,
			data.IocGroup[i].Event_id,
			data.IocGroup[i].Latitude,
			data.IocGroup[i].Longitude,
			data.IocGroup[i].IOCCountryOrigin)
	}*/
}

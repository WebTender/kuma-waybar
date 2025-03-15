package kuma

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Kuma struct {
	baseUrl string
	apiKey string
	req *http.Request

	mutex sync.Mutex
	currentMetrics *[]Metric
	currentMonitors *[]*Monitor
	LastUpdated int64
}

type MonitorStatus uint8
const (
	Down MonitorStatus = iota
	Up
	Pending
	Maintenance 
)
func ParseMonitorStatus(value string) (MonitorStatus, error) {
	num, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return Pending, err
	}

	status := MonitorStatus(num)
	if status < Down || status > Maintenance {
		return Pending, errors.New("invalid MonitorStatus value")
	}

	return status, nil
}

type MonitorType string
const (
	HTTP MonitorType = "http"
	TCP MonitorType = "port"
	PUSH MonitorType = "push"
	GROUP MonitorType = "group"
	PING MonitorType = "ping"
)

func ParseMonitorType(value string) MonitorType {
	switch value {
	case "http":
		return HTTP
	case "port":
		return TCP
	case "push":
		return PUSH
	case "group":
		return GROUP
	case "ping":
		return PING
	default:
		println("WARN Unknown monitor type: " + value)
		return ""
	}
}

type Metric struct {
	Key string
	Value string
	Labels map[string]string
}

type Monitor struct {
	Status MonitorStatus
	Type MonitorType
	Name string
	Url string
	Hostname string
	Port uint16
	ResponseTime uint64
}

func New(baseUrl string, apiKey string) (*Kuma, error) {
	if (baseUrl == "") {
		return nil, errors.New("baseUrl is required")
	}
	if (apiKey == "") {
		return nil, errors.New("apiKey is required")
	}
	baseUrl = strings.TrimRight(baseUrl, "/")

	req, err := http.NewRequest("GET", baseUrl + "/metrics", nil)
	if (err != nil) {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(":" + apiKey)))

	return &Kuma{baseUrl: baseUrl, apiKey: apiKey, req: req}, nil
}

func (kuma *Kuma) GetMetrics() ([]Metric, []*Monitor, error) {
	resp, err := http.DefaultClient.Do(kuma.req)
	if (err != nil) {
		return nil, nil, err
	}
	if (resp.StatusCode != 200) {
		return nil, nil, errors.New("failed to get metrics: " + resp.Status)
	}
	
	body, err := io.ReadAll(resp.Body)
	if (err != nil) {
		return nil, nil, err
	}
	
	lines := strings.Split(string(body), "\n")

	newMetrics := []Metric{}
	newMonitors := []*Monitor{}

	for _, line := range lines {
		if (line == "") {
			continue
		}

		parts := strings.Split(line, "{")
		if (len(parts) == 1) {
			parts2 := strings.Split(line, " ")
			if (len(parts2) != 2) {
				continue
			}
			key := strings.TrimSpace(parts2[0])
			value := strings.TrimSpace(parts2[1])
			newMetric := parseMetric(key, "", value)
			newMetrics = append(newMetrics, newMetric)
			continue;
		}
		if (len(parts) != 2) {
			continue
		}
		key := strings.TrimSpace(parts[0])

		parts = strings.Split(parts[1], "}")
		if (len(parts) != 2) {
			continue
		}
		mapStr := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		newMetric := parseMetric(key, mapStr, value)
		newMetrics = append(newMetrics, newMetric)

		newMonitor, isNew := matchOrNewMonitorFromMetric(&newMetric,&newMonitors)
		if (isNew) {
			newMonitors = append(newMonitors, newMonitor)
		}
	}

	kuma.mutex.Lock()
	kuma.LastUpdated = time.Now().Unix()
	kuma.currentMetrics = &newMetrics
	kuma.currentMonitors = &newMonitors
	kuma.mutex.Unlock()

	return newMetrics, newMonitors, nil
}

// Opens the dashboard in the default browser
func (kuma *Kuma) Open() {
	err := exec.Command("open", kuma.baseUrl + "/dashboard").Run()
	if (err != nil) {
		fmt.Println("Failed to open URL:")
		os.Exit(1)
	}
}

func parseMetric(key string, keyValueMapString string, value string) Metric {
	labels := map[string]string{}
	if (keyValueMapString != "") {
		mapParts := strings.Split(keyValueMapString, ",")
		for _, mapPart := range mapParts {
			mapPart = strings.TrimSpace(mapPart)
			mapPartParts := strings.Split(mapPart, "=")
			if (len(mapPartParts) != 2) {
				continue
			}
			mapKey := strings.TrimSpace(mapPartParts[0])
			mapValue := strings.TrimSpace(mapPartParts[1])
			if (mapValue[0] == '"' && mapValue[len(mapValue) - 1] == '"') {
				mapValue = mapValue[1:len(mapValue) - 1]
			}
			labels[mapKey] = mapValue
		}
	}
	return Metric{
		key,
		value,
		labels,
	}
}

var ErrorNotAMonitor = errors.New("metric is not a monitor")

func matchOrNewMonitorFromMetric(metric *Metric, monitors *[]*Monitor) (*Monitor, bool) {
	name := metric.Labels["monitor_name"]
	if name == "" {
		return nil, false
	}
	for _, monitor := range *monitors {
		if (monitor.Name == name) {
			applyMetricValue(metric, monitor)
			return monitor, false
		}
	}
	monitor, err := newMonitorFromMetric(metric)
	if (err != nil) {
		return nil, false
	}

	return &monitor, true
}

func newMonitorFromMetric(metric *Metric) (Monitor, error) {
	status := Pending
	
	port, _ := strconv.ParseUint(metric.Labels["monitor_port"], 10, 16)

	var monitor = Monitor{
		Name: metric.Labels["monitor_name"],
		Status: status,
		Type: ParseMonitorType(metric.Labels["monitor_type"]),
		Port: uint16(port),
		Hostname: metric.Labels["monitor_hostname"],
		Url: metric.Labels["monitor_url"],
	}

	applyMetricValue(metric, &monitor)

	return monitor, nil
}

func applyMetricValue(metric *Metric, monitor *Monitor) {
	switch(metric.Key) {
		case "monitor_status":
			status, err := ParseMonitorStatus(metric.Value)
			if (err == nil) {
				monitor.Status = status
			}
		case "monitor_response_time":
			time, err := strconv.ParseUint(metric.Value, 10, 64)
			if (err == nil) {
				monitor.ResponseTime = time
			}
	}
}
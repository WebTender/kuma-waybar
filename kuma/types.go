package kuma

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

var ErrorNotAMonitor = errors.New("metric is not a monitor")

type Kuma struct {
	baseUrl string
	apiKey string
	req *http.Request

	LastUpdated int64
}

type MonitorStatus uint8
const (
	Down MonitorStatus = iota
	Up
	Pending
	Maintenance 
	Paused
)
func ParseMonitorStatus(value string) (MonitorStatus, error) {
	num, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return Paused, err
	}

	status := MonitorStatus(num)
	if status < Down || status > Paused {
		return Paused, errors.New("invalid MonitorStatus value")
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
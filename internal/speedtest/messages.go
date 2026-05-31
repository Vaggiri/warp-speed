package speedtest

import "time"

// Messages for the speedtest engine to send to the UI

type ServerLocatingMsg struct{}
type ServerFoundMsg struct {
	Name     string
	Location string
	Host     string
}

type PingResultMsg struct {
	Latency time.Duration
	Jitter  time.Duration
}

type DownloadProgressMsg struct {
	Progress float64 // 0.0 to 1.0
	SpeedMBs float64 // Current speed in MB/s
}

type UploadProgressMsg struct {
	Progress float64 // 0.0 to 1.0
	SpeedMBs float64 // Current speed in MB/s
}

type TestCompleteMsg struct{}
type ErrorMsg struct {
	Err error
}

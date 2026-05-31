package monitor

import (
	"net"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	psnet "github.com/shirou/gopsutil/v3/net"
)

type PingStatusMsg struct {
	CloudflareLatency time.Duration
	GoogleLatency     time.Duration
	SpeedtestLatency  time.Duration
	DownloadMbps      float64
	UploadMbps        float64
}

type MonitorTickMsg time.Time

var (
	lastRx   uint64
	lastTx   uint64
	lastTime time.Time
)

func init() {
	// Initialize counters
	stats, err := psnet.IOCounters(false)
	if err == nil && len(stats) > 0 {
		lastRx = stats[0].BytesRecv
		lastTx = stats[0].BytesSent
	}
	lastTime = time.Now()
}

func Tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return MonitorTickMsg(t)
	})
}

func MeasureLatency(speedtestHost string) tea.Cmd {
	return func() tea.Msg {
		cLat := measureTCP("1.1.1.1:53")
		gLat := measureTCP("8.8.8.8:53")
		
		sLat := time.Duration(0)
		if speedtestHost != "" {
			sLat = measureTCP(speedtestHost)
		}

		// Calculate bandwidth
		now := time.Now()
		elapsed := now.Sub(lastTime).Seconds()
		dlMbps := 0.0
		ulMbps := 0.0

		stats, err := psnet.IOCounters(false)
		if err == nil && len(stats) > 0 && elapsed > 0 {
			currRx := stats[0].BytesRecv
			currTx := stats[0].BytesSent
			
			// Bytes to bits (x8) then to Megabits (/1e6)
			dlMbps = (float64(currRx-lastRx) * 8 / 1000000) / elapsed
			ulMbps = (float64(currTx-lastTx) * 8 / 1000000) / elapsed

			lastRx = currRx
			lastTx = currTx
		}
		lastTime = now

		return PingStatusMsg{
			CloudflareLatency: cLat,
			GoogleLatency:     gLat,
			SpeedtestLatency:  sLat,
			DownloadMbps:      dlMbps,
			UploadMbps:        ulMbps,
		}
	}
}

func measureTCP(address string) time.Duration {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return 0 // 0 means offline/timeout
	}
	defer conn.Close()
	return time.Since(start)
}

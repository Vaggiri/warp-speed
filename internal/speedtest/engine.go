package speedtest

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/showwin/speedtest-go/speedtest"
)

var (
	client      *speedtest.Speedtest
	testServer  *speedtest.Server
	allServers  speedtest.Servers
)

func init() {
	client = speedtest.New()
}

func StartLocatingServer() tea.Cmd {
	return func() tea.Msg {
		if testServer != nil {
			return ServerFoundMsg{
				Name:     testServer.Name,
				Location: testServer.Country,
				Host:     testServer.Host,
			}
		}

		_, err := client.FetchUserInfo()
		if err != nil {
			return ErrorMsg{err}
		}

		serverList, err := client.FetchServers()
		if err != nil {
			return ErrorMsg{err}
		}

		targets, err := serverList.FindServer([]int{})
		if err != nil || len(targets) == 0 {
			return ErrorMsg{err}
		}

		testServer = targets[0]
		return ServerFoundMsg{
			Name:     testServer.Name,
			Location: testServer.Country,
			Host:     testServer.Host,
		}
	}
}

type ServerListMsg struct {
	Servers speedtest.Servers
}

func FetchAllServersCmd() tea.Cmd {
	return func() tea.Msg {
		_, err := client.FetchUserInfo()
		if err != nil {
			return ErrorMsg{err}
		}
		allServers, err = client.FetchServers()
		if err != nil {
			return ErrorMsg{err}
		}
		return ServerListMsg{Servers: allServers}
	}
}

func SetTargetServer(s *speedtest.Server) {
	testServer = s
}

func ClearTargetServer() {
	testServer = nil
}

func GetTargetServer() *speedtest.Server {
	return testServer
}

func StartPing(host string) tea.Cmd {
	return func() tea.Msg {
		err := testServer.PingTest(nil)
		if err != nil {
			return ErrorMsg{err}
		}

		return PingResultMsg{
			Latency: testServer.Latency,
			Jitter:  testServer.Jitter,
		}
	}
}

func StartDownload() tea.Cmd {
	return func() tea.Msg {
		err := testServer.DownloadTest()
		if err != nil {
			return ErrorMsg{err}
		}
		
		speedMBs := (float64(testServer.DLSpeed) * 8) / 1024 / 1024 // Mbps
		return DownloadProgressMsg{Progress: 1.0, SpeedMBs: speedMBs}
	}
}

func StartUpload() tea.Cmd {
	return func() tea.Msg {
		err := testServer.UploadTest()
		if err != nil {
			return ErrorMsg{err}
		}
		
		speedMBs := (float64(testServer.ULSpeed) * 8) / 1024 / 1024 // Mbps
		return UploadProgressMsg{Progress: 1.0, SpeedMBs: speedMBs}
	}
}

func FinalizeTest() tea.Cmd {
	return func() tea.Msg {
		return TestCompleteMsg{}
	}
}

func GetCurrentDownloadSpeed() float64 {
	if testServer == nil {
		return 0
	}
	return (float64(testServer.DLSpeed) * 8) / 1024 / 1024
}

func GetCurrentUploadSpeed() float64 {
	if testServer == nil {
		return 0
	}
	return (float64(testServer.ULSpeed) * 8) / 1024 / 1024
}

// Tick message for progress updates
type TickMsg time.Time

func Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

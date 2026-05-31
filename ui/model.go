package ui

import (
	"fmt"
	"time"
	"warp-speed/internal/history"
	"warp-speed/internal/monitor"
	"warp-speed/internal/speedtest"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	speedtestgo "github.com/showwin/speedtest-go/speedtest"
)

type Screen int

const (
	ScreenMenu Screen = iota
	ScreenSpeedTest
	ScreenMonitor
	ScreenHistory
	ScreenServerSelect
)

type SpeedTestState int

const (
	StateInit SpeedTestState = iota
	StateLocating
	StatePing
	StateDownload
	StateUpload
	StateDone
	StateError
)

// List Item
type serverItem struct {
	server *speedtestgo.Server
}
func (i serverItem) Title() string       { return i.server.Name }
func (i serverItem) Description() string { return fmt.Sprintf("%s (%s) - %s", i.server.Sponsor, i.server.Country, i.server.Host) }
func (i serverItem) FilterValue() string { return i.server.Name + " " + i.server.Country + " " + i.server.Sponsor }

type Model struct {
	currentScreen Screen
	
	// Menu State
	menuIndex int
	menuItems []string

	// SpeedTest State
	stState     SpeedTestState
	spinner     spinner.Model
	progressBar progress.Model
	serverName  string
	serverLoc   string
	serverHost  string
	ping        time.Duration
	jitter      time.Duration
	dlSpeed     float64
	ulSpeed     float64
	progressVal float64
	err         error

	// Monitor State
	monCloudflare time.Duration
	monGoogle     time.Duration
	monSpeedtest  time.Duration
	monDownload   float64
	monUpload     float64
	
	// Graph Data
	histCloudflare []float64
	histGoogle     []float64
	histSpeedtest  []float64
	histDownload   []float64
	histUpload     []float64

	// History State
	historyTable table.Model
	
	// Server Select State
	serverList list.Model
	fetchingServers bool
}

func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Line
	s.Style = SpinnerStyle
	
	p := progress.New(
		progress.WithSolidFill(string(ColorAccent)),
		progress.WithoutPercentage(),
	)

	// Setup Server List
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Speedtest Server"
	l.SetShowStatusBar(false)

	// Setup History Table
	columns := []table.Column{
		{Title: "Date", Width: 20},
		{Title: "Server", Width: 30},
		{Title: "Ping", Width: 10},
		{Title: "Down (Mbps)", Width: 15},
		{Title: "Up (Mbps)", Width: 15},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	ts := table.DefaultStyles()
	ts.Header = ts.Header.BorderStyle(lipgloss.NormalBorder()).BorderForeground(ColorBorder).BorderBottom(true).Bold(false)
	ts.Selected = ts.Selected.Foreground(ColorBg).Background(ColorAccent).Bold(false)
	t.SetStyles(ts)

	return Model{
		currentScreen: ScreenMenu,
		menuIndex:     0,
		menuItems:     []string{"Speed Test", "Real-Time Status Monitor", "Select Server", "View History"},
		stState:       StateInit,
		spinner:       s,
		progressBar:   p,
		serverList:    l,
		historyTable:  t,
	}
}

func (m *Model) pushGraphData(history *[]float64, val float64) {
	*history = append(*history, val)
	if len(*history) > 40 {
		*history = (*history)[1:]
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			tea.SetWindowTitle("Warp Speed CLI")
			return nil
		},
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := AppStyle.GetFrameSize()
		m.serverList.SetSize(msg.Width-h, msg.Height-v-4)
		
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || (msg.String() == "q" && m.currentScreen != ScreenServerSelect) {
			return m, tea.Quit
		}

		if m.currentScreen == ScreenMenu {
			switch msg.String() {
			case "up", "k":
				if m.menuIndex > 0 {
					m.menuIndex--
				}
			case "down", "j":
				if m.menuIndex < len(m.menuItems)-1 {
					m.menuIndex++
				}
			case "enter":
				switch m.menuIndex {
				case 0: // Speed Test
					if speedtest.GetTargetServer() == nil {
						// Redirect to server selection
						m.currentScreen = ScreenServerSelect
						if len(m.serverList.Items()) == 0 {
							m.fetchingServers = true
							cmds = append(cmds, speedtest.FetchAllServersCmd())
						}
					} else {
						m.currentScreen = ScreenSpeedTest
						cmds = append(cmds, speedtest.Tick(), func() tea.Msg { return speedtest.ServerLocatingMsg{} })
					}
				case 1: // Monitor
					m.currentScreen = ScreenMonitor
					cmds = append(cmds, func() tea.Msg { return speedtest.ServerLocatingMsg{} })
				case 2: // Select Server
					m.currentScreen = ScreenServerSelect
					if len(m.serverList.Items()) == 0 {
						m.fetchingServers = true
						cmds = append(cmds, speedtest.FetchAllServersCmd())
					}
				case 3: // View History
					m.currentScreen = ScreenHistory
					// Load records
					records, _ := history.LoadRecords()
					rows := []table.Row{}
					for i := len(records) - 1; i >= 0; i-- {
						r := records[i]
						rows = append(rows, table.Row{
							r.Date.Format("2006-01-02 15:04:05"),
							r.Server,
							fmt.Sprintf("%.2f ms", r.Latency),
							fmt.Sprintf("%.2f", r.Download),
							fmt.Sprintf("%.2f", r.Upload),
						})
					}
					m.historyTable.SetRows(rows)
				}
			}
			return m, tea.Batch(cmds...)
		}
		
		if msg.String() == "esc" {
			if m.currentScreen == ScreenServerSelect && m.serverList.FilterState() == list.Filtering {
				// Don't go back to menu, let list handle esc to cancel filter
			} else {
				m.currentScreen = ScreenMenu
				m.stState = StateInit
				m.progressVal = 0
				return m, nil
			}
		}

		if m.currentScreen == ScreenServerSelect {
			if msg.String() == "enter" {
				if i, ok := m.serverList.SelectedItem().(serverItem); ok {
					speedtest.SetTargetServer(i.server)
					m.currentScreen = ScreenMenu // go back to menu
					return m, nil
				}
			}
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progressBar.Update(msg)
		m.progressBar = progressModel.(progress.Model)
		cmds = append(cmds, cmd)
		
	// Server List Messages
	case speedtest.ServerListMsg:
		m.fetchingServers = false
		items := make([]list.Item, len(msg.Servers))
		for i, s := range msg.Servers {
			items[i] = serverItem{server: s}
		}
		cmd = m.serverList.SetItems(items)
		cmds = append(cmds, cmd)

	// SpeedTest & Monitor Messages
	case speedtest.ServerLocatingMsg:
		if m.currentScreen == ScreenSpeedTest {
			m.stState = StateLocating
		}
		cmds = append(cmds, speedtest.StartLocatingServer())
		
	case speedtest.ServerFoundMsg:
		m.serverName = msg.Name
		m.serverLoc = msg.Location
		m.serverHost = msg.Host
		
		if m.currentScreen == ScreenSpeedTest {
			m.stState = StatePing
			cmds = append(cmds, speedtest.StartPing(msg.Host))
		} else if m.currentScreen == ScreenMonitor {
			cmds = append(cmds, monitor.MeasureLatency(m.serverHost), monitor.Tick())
		}
		
	case speedtest.PingResultMsg:
		if m.currentScreen == ScreenSpeedTest {
			m.stState = StateDownload
			m.ping = msg.Latency
			m.jitter = msg.Jitter
			m.progressVal = 0
			cmds = append(cmds, speedtest.StartDownload())
		}
		
	case speedtest.DownloadProgressMsg:
		if m.currentScreen == ScreenSpeedTest {
			m.dlSpeed = msg.SpeedMBs
			m.progressVal = msg.Progress
			cmd = m.progressBar.SetPercent(msg.Progress)
			cmds = append(cmds, cmd)
			if msg.Progress >= 1.0 {
				m.stState = StateUpload
				m.progressVal = 0
				cmds = append(cmds, speedtest.StartUpload())
			}
		}
		
	case speedtest.UploadProgressMsg:
		if m.currentScreen == ScreenSpeedTest {
			m.ulSpeed = msg.SpeedMBs
			m.progressVal = msg.Progress
			cmd = m.progressBar.SetPercent(msg.Progress)
			cmds = append(cmds, cmd)
			if msg.Progress >= 1.0 {
				m.stState = StateDone
				cmds = append(cmds, speedtest.FinalizeTest())
			}
		}
		
	case speedtest.TestCompleteMsg:
		if m.currentScreen == ScreenSpeedTest {
			m.stState = StateDone
			// Save History
			history.SaveRecord(history.Record{
				Date:     time.Now(),
				Server:   fmt.Sprintf("%s (%s)", m.serverName, m.serverLoc),
				Latency:  float64(m.ping.Milliseconds()),
				Download: m.dlSpeed,
				Upload:   m.ulSpeed,
			})
		}
		
	case speedtest.ErrorMsg:
		if m.currentScreen == ScreenSpeedTest {
			m.stState = StateError
			m.err = msg.Err
		}

	case speedtest.TickMsg:
		if m.currentScreen == ScreenSpeedTest {
			if m.stState == StateDownload {
				m.dlSpeed = speedtest.GetCurrentDownloadSpeed()
				m.progressVal += 0.01
				if m.progressVal > 0.99 {
					m.progressVal = 0.99
				}
				cmd = m.progressBar.SetPercent(m.progressVal)
				cmds = append(cmds, cmd)
			} else if m.stState == StateUpload {
				m.ulSpeed = speedtest.GetCurrentUploadSpeed()
				m.progressVal += 0.01
				if m.progressVal > 0.99 {
					m.progressVal = 0.99
				}
				cmd = m.progressBar.SetPercent(m.progressVal)
				cmds = append(cmds, cmd)
			}
			cmds = append(cmds, speedtest.Tick())
		}

	case monitor.MonitorTickMsg:
		if m.currentScreen == ScreenMonitor {
			cmds = append(cmds, monitor.MeasureLatency(m.serverHost), monitor.Tick())
		}

	case monitor.PingStatusMsg:
		if m.currentScreen == ScreenMonitor {
			m.monCloudflare = msg.CloudflareLatency
			m.monGoogle = msg.GoogleLatency
			m.monSpeedtest = msg.SpeedtestLatency
			m.monDownload = msg.DownloadMbps
			m.monUpload = msg.UploadMbps

			m.pushGraphData(&m.histCloudflare, float64(msg.CloudflareLatency.Milliseconds()))
			m.pushGraphData(&m.histGoogle, float64(msg.GoogleLatency.Milliseconds()))
			m.pushGraphData(&m.histSpeedtest, float64(msg.SpeedtestLatency.Milliseconds()))
			m.pushGraphData(&m.histDownload, msg.DownloadMbps)
			m.pushGraphData(&m.histUpload, msg.UploadMbps)
		}
	}

	// Update components based on screen
	if m.currentScreen == ScreenServerSelect {
		m.serverList, cmd = m.serverList.Update(msg)
		cmds = append(cmds, cmd)
	} else if m.currentScreen == ScreenHistory {
		m.historyTable, cmd = m.historyTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

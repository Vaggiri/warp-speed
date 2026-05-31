package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
)

func (m Model) View() string {
	if m.err != nil {
		return AppStyle.Render(ErrorStyle.Render(fmt.Sprintf("Critical Error: %v\nPress q to quit.", m.err)))
	}

	var sb strings.Builder

	// Global Header
	sb.WriteString(HeaderStyle.Render(" WARP-SPEED // Enterprise Edge CLI "))
	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(ColorBorder).Render(" Developed by Girisudhan V | VAG CREATIONS "))
	sb.WriteString("\n\n")

	// Screen Routing
	switch m.currentScreen {
	case ScreenMenu:
		sb.WriteString(m.viewMenu())
	case ScreenSpeedTest:
		sb.WriteString(m.viewSpeedTest())
	case ScreenMonitor:
		sb.WriteString(m.viewMonitor())
	case ScreenHistory:
		sb.WriteString(m.viewHistory())
	case ScreenServerSelect:
		sb.WriteString(m.viewServerSelect())
	}

	// Global Footer
	sb.WriteString("\n\n")
	if m.currentScreen == ScreenMenu {
		sb.WriteString(FooterStyle.Render("Use ↑/↓ to navigate • Enter to select • q to quit"))
	} else if m.currentScreen == ScreenServerSelect {
		sb.WriteString(FooterStyle.Render("Use ↑/↓ to navigate • / to filter • Enter to select • Esc to go back"))
	} else {
		sb.WriteString(FooterStyle.Render("Esc: Return to Menu • q: Quit"))
	}

	return AppStyle.Render(sb.String())
}

func (m Model) viewMenu() string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("Select Operation Mode"))
	sb.WriteString("\n\n")

	for i, item := range m.menuItems {
		if i == m.menuIndex {
			sb.WriteString(SelectedItemStyle.Render(fmt.Sprintf("▶ %s", item)))
		} else {
			sb.WriteString(ItemStyle.Render(fmt.Sprintf("  %s", item)))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (m Model) viewSpeedTest() string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("Diagnostic: Throughput & Latency"))
	sb.WriteString("\n\n")

	if m.stState > StateLocating {
		serverInfo := fmt.Sprintf("Target Node: %s (%s)", m.serverName, m.serverLoc)
		sb.WriteString(SubtitleStyle.Render(serverInfo))
		sb.WriteString("\n\n")
	} else if m.stState == StateLocating {
		sb.WriteString(fmt.Sprintf("%s Establishing edge connection...\n", m.spinner.View()))
		return sb.String()
	}

	// Metrics Row
	metricsRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderMetricBox("Latency", fmt.Sprintf("%d", m.ping.Milliseconds()), "ms"),
		m.renderMetricBox("Jitter", fmt.Sprintf("%d", m.jitter.Milliseconds()), "ms"),
		m.renderMetricBox("Download", fmt.Sprintf("%.2f", m.dlSpeed), "Mbps"),
		m.renderMetricBox("Upload", fmt.Sprintf("%.2f", m.ulSpeed), "Mbps"),
	)
	sb.WriteString(metricsRow)
	sb.WriteString("\n\n")

	// Progress Area
	sb.WriteString(m.progressBar.View())
	sb.WriteString("\n")

	if m.stState == StatePing {
		sb.WriteString(fmt.Sprintf("%s Measuring ICMP Latency...", m.spinner.View()))
	} else if m.stState == StateDownload {
		sb.WriteString(fmt.Sprintf("%s Executing Ingress Test...", m.spinner.View()))
	} else if m.stState == StateUpload {
		sb.WriteString(fmt.Sprintf("%s Executing Egress Test...", m.spinner.View()))
	} else if m.stState == StateDone {
		sb.WriteString(SuccessStyle.Render("✓ Diagnostic Complete. All systems operational. Results saved to history."))
	}

	return sb.String()
}

func (m Model) viewMonitor() string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("Diagnostic: Real-Time Edge Telemetry"))
	sb.WriteString("\n\n")

	if m.serverName == "" {
		sb.WriteString(fmt.Sprintf("%s Initializing telemetry endpoints...", m.spinner.View()))
		return sb.String()
	}

	sb.WriteString(SubtitleStyle.Render("Continuously monitoring network health across diverse routes."))
	sb.WriteString("\n\n")

	// Render the three ping targets
	row := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderStatusBox("Cloudflare (1.1.1.1)", m.monCloudflare, m.histCloudflare),
		m.renderStatusBox("Google (8.8.8.8)", m.monGoogle, m.histGoogle),
		m.renderStatusBox(fmt.Sprintf("Speedtest (%s)", m.serverName), m.monSpeedtest, m.histSpeedtest),
	)

	bandwidthRow := lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderGraphBox("Current Ingress (Mbps)", fmt.Sprintf("%.2f", m.monDownload), m.histDownload),
		m.renderGraphBox("Current Egress (Mbps)", fmt.Sprintf("%.2f", m.monUpload), m.histUpload),
	)

	sb.WriteString(row)
	sb.WriteString("\n")
	sb.WriteString(bandwidthRow)
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("%s Polling...", m.spinner.View()))

	return sb.String()
}

func (m Model) viewHistory() string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("Test History Archive"))
	sb.WriteString("\n\n")
	sb.WriteString(m.historyTable.View())
	return sb.String()
}

func (m Model) viewServerSelect() string {
	if m.fetchingServers {
		return fmt.Sprintf("%s Fetching global edge nodes...", m.spinner.View())
	}
	return m.serverList.View()
}

func (m Model) renderMetricBox(label, value, unit string) string {
	content := fmt.Sprintf("%s\n%s %s",
		MetricLabelStyle.Render(strings.ToUpper(label)),
		MetricValueStyle.Render(value),
		MetricUnitStyle.Render(unit),
	)
	return BoxStyle.Render(content)
}

func (m Model) renderStatusBox(target string, latency time.Duration, hist []float64) string {
	status := SuccessStyle.Render("● ONLINE")
	val := fmt.Sprintf("%d ms", latency.Milliseconds())
	
	if latency == 0 {
		status = ErrorStyle.Render("● OFFLINE")
		val = "ERR"
	} else if latency > 150*time.Millisecond {
		status = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true).Render("● DEGRADED")
	}

	graph := ""
	if len(hist) > 0 {
		graph = asciigraph.Plot(hist, asciigraph.Height(4), asciigraph.Width(20))
	}

	content := fmt.Sprintf("%s\n%s\n\n%s\n\n%s",
		MetricLabelStyle.Render(target),
		MetricValueStyle.Render(val),
		status,
		graph,
	)
	return BoxStyle.Width(35).Height(12).Render(content)
}

func (m Model) renderGraphBox(label string, value string, hist []float64) string {
	graph := ""
	if len(hist) > 0 {
		graph = asciigraph.Plot(hist, asciigraph.Height(4), asciigraph.Width(30))
	}

	content := fmt.Sprintf("%s\n%s\n\n%s",
		MetricLabelStyle.Render(label),
		MetricValueStyle.Render(value),
		graph,
	)
	return BoxStyle.Width(45).Height(10).Render(content)
}

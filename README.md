# Warp-Speed CLI 🚀

`warp-speed` is a high-performance, enterprise-grade network diagnostic and telemetry CLI tool built in Go. It leverages the power of Cobra for command routing and the Charm.sh ecosystem (Bubble Tea, Lipgloss, Bubbles) to deliver a stunning, interactive Terminal User Interface (TUI).

---

## 🌟 Key Features

1. **Enterprise Dashboard UI**: A sleek, SaaS-inspired dark mode aesthetic with crisp borders, professional accent colors, and a persistent global layout.
2. **Comprehensive Speed Diagnostics**: Accurately measure ICMP Latency, Jitter, Ingress (Download) throughput, and Egress (Upload) throughput.
3. **Real-Time Telemetry Monitor**: A dedicated dashboard that uses asynchronous TCP polling and `gopsutil` to plot live, rolling ASCII sparkline graphs of:
    *   Cloudflare Edge (1.1.1.1) Latency
    *   Google DNS (8.8.8.8) Latency
    *   Selected Node Latency
    *   Current Network Interface Ingress & Egress Bandwidth
4. **Advanced Target Selection**: Fetch a global list of edge nodes and use fuzzy-searching to pinpoint the exact server you want to test against.
5. **Local History Archive**: Automatically persists the results of your last 100 speed tests to a local JSON database (`~/.warp-speed/history.json`), viewable directly in the terminal via an interactive data table.

---

## 🛠 Architecture & Tech Stack

*   **Language**: Go 1.21+
*   **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
*   **Terminal UI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (The Elm Architecture for the terminal)
*   **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss) (CSS for the terminal)
*   **UI Components**: [Bubbles](https://github.com/charmbracelet/bubbles) (Tables, Lists, Spinners, Progress Bars)
*   **Graphing Engine**: [Asciigraph](https://github.com/guptarohit/asciigraph)
*   **System Telemetry**: [Gopsutil](https://github.com/shirou/gopsutil)
*   **Speedtest Backend**: [Speedtest-go](https://github.com/showwin/speedtest-go)

### Directory Structure

```text
warp-speed/
├── main.go                     # Application entry point
├── cmd/
│   └── root.go                 # Cobra root command setup
├── internal/
│   ├── history/
│   │   └── history.go          # JSON persistence logic
│   ├── monitor/
│   │   └── engine.go           # TCP pinging and bandwidth polling
│   └── speedtest/
│       ├── engine.go           # Speedtest-go wrapper & server fetching
│       └── messages.go         # Bubble Tea message definitions
└── ui/
    ├── model.go                # Central state machine & screen router
    ├── style.go                # Enterprise design tokens and Lipgloss styles
    └── view.go                 # Component rendering and layout
```

---

## 🚀 Installation & Deployment

### Quick Start (Local Compilation)
1. Clone the repository.
2. Ensure you have Go installed.
3. Run the following command in the project root:
   ```bash
   go build -o warp-speed
   ```
4. Run the executable: `./warp-speed` (or `.\warp-speed.exe` on Windows).

### Global Installation (Developer Standard)
If you want to run `warp-speed` from anywhere on your system, we provide professional installation scripts that automatically inject the binary into your system `PATH`.

**For Windows (PowerShell):**
```powershell
.\install.ps1
```
*This creates a `%USERPROFILE%\.warp-speed\bin` directory, copies the executable, and updates your User PATH.*

**For macOS / Linux (Bash/Zsh):**
```bash
./install.sh
```
*This places the binary in `~/.local/bin` and updates your `.bashrc` or `.zshrc`.*

*(For more details on `PATH` management and symlinking vs. copying, see `INSTALL.md`).*

---

## 🕹 Usage Guide

When you launch `warp-speed`, you are presented with the Main Menu. Use the `Up` and `Down` arrow keys to navigate, and hit `Enter` to select. Press `Esc` from any screen to return to the menu, and `q` to quit the application.

### 1. Speed Test
*   If you haven't selected a server yet, selecting this will intelligently redirect you to the **Select Server** screen.
*   Once a server is selected, it will run a sequential diagnostic (Ping -> Download -> Upload) and display a progress bar.
*   Upon completion, the results are automatically saved to your local history.

### 2. Real-Time Status Monitor
*   This mode continuously polls your active network connections.
*   It measures TCP latency to Cloudflare, Google, and your selected server every second.
*   It calculates the delta of your network interface's I/O counters to display current upload/download speeds.
*   All data is mapped to live sparkline graphs directly in the terminal.

### 3. Select Server
*   Fetches thousands of global edge nodes.
*   Press `/` to enter filter mode. You can fuzzy-search by City, Country, or Sponsor Name.
*   Hit `Enter` to lock in your target.

### 4. View History
*   Opens an interactive data table displaying your archived test runs.
*   Use the arrow keys to scroll through your history.

---

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! 
If you plan to add new screens, please follow the Elm Architecture pattern established in `ui/model.go` and define your distinct views in `ui/view.go`. 

*Design changes should be strictly confined to `ui/style.go` to maintain the enterprise aesthetic.*

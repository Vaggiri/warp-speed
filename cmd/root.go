package cmd

import (
	"fmt"
	"warp-speed/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "warp-speed",
	Short: "A sleek, high-performance CLI internet speed test",
	Long:  `warp-speed is a terminal-based internet speed test tool featuring a beautiful Pro aesthetic UI built with Bubble Tea.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("error running UI: %w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Flags and configuration settings can go here
}

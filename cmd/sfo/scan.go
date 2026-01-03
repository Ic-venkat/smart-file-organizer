package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"smart-organizer/internal/scanner"
)

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan directory and show file analytics",
	Long:  `Recursively scans a directory to count files and group them by extension.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		// Verify path exists
		info, err := os.Stat(targetPath)
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", targetPath, err)
			os.Exit(1)
		}
		if !info.IsDir() {
			fmt.Printf("Path %s is not a directory\n", targetPath)
			os.Exit(1)
		}

		fmt.Printf("Scanning directory: %s\n", targetPath)
		
		stats, err := scanner.Scan(targetPath)
		if err != nil {
			fmt.Printf("Scan failed: %v\n", err)
			os.Exit(1)
		}

		runTUI(stats)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

// TUI Model
type model struct {
	table     table.Model
	stats     *scanner.Analytics
	quitting  bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		headerHeight := 4 // Title + Padding
		footerHeight := 2 // Help text
		m.table.SetHeight(msg.Height - headerHeight - footerHeight)
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	totalStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(1, 0)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	baseView := lipgloss.JoinVertical(lipgloss.Left,
		totalStyle.Render(fmt.Sprintf("Total Files: %d • Total Size: %s", m.stats.TotalFiles, formatBytes(m.stats.TotalSize))),
		m.table.View(),
		helpStyle.Render("Press 'q' to quit • ↑/↓ to scroll"),
	)
	
	// Add padding
	return lipgloss.NewStyle().Margin(1, 1).Render(baseView)
}

func runTUI(stats *scanner.Analytics) {
	columns := []table.Column{
		{Title: "Extension", Width: 12},
		{Title: "Count", Width: 8},
		{Title: "Count %", Width: 8},
		{Title: "Size", Width: 10},
		{Title: "Size %", Width: 8},
	}

	// Prepare rows
	type kv struct {
		Key   string
		Value int
		Size  int64
	}
	var ss []kv
	for k, v := range stats.Extensions {
		ss = append(ss, kv{k, v, stats.Sizes[k]})
	}
	// Sort by Size by default (more interesting usually)
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Size > ss[j].Size
	})

	var rows []table.Row
	for _, kv := range ss {
		countPct := float64(kv.Value) / float64(stats.TotalFiles) * 100
		sizePct := float64(kv.Size) / float64(stats.TotalSize) * 100
		
		// Avoid NaN if total size is 0
		if stats.TotalSize == 0 {
			sizePct = 0
		}

		rows = append(rows, table.Row{
			kv.Key,
			fmt.Sprintf("%d", kv.Value),
			fmt.Sprintf("%.1f%%", countPct),
			formatBytes(kv.Size),
			fmt.Sprintf("%.1f%%", sizePct),
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20), // Increased height for desktop
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true).
		Foreground(lipgloss.Color("231")). // White text
		Background(lipgloss.Color("63"))   // Purple background
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{table: t, stats: stats}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
	
	// Print summary after TUI exits
	fmt.Printf("\nScan complete. Scanned %d files. Total Size: %s\n", stats.TotalFiles, formatBytes(stats.TotalSize))
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

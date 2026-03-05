package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#3B82F6")).
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3B82F6")).
			Padding(1, 2).
			Width(60)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3B82F6")).
			Padding(1, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#3B82F6")).
				Bold(true)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	dimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))
)

// MenuItem represents a menu item in the command palette
type MenuItem struct {
	title       string
	description string
	command     string
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

// Model represents the TUI state
type tuiModel struct {
	textInput    textinput.Model
	list         list.Model
	state        string // "input", "menu", "connecting", "executing", "result"
	connected    bool
	errorMsg     string
	resultMsg    string
	width        int
	height       int
	selectedCmd  string
	quitting     bool
	filterText   string // Text used to filter menu items
}

// Messages
type connectMsg struct {
	success bool
	err     error
}

type executeMsg struct {
	success bool
	output  string
	err     error
}

// Initialize the TUI
func initialTuiModel() tuiModel {
	ti := textinput.New()
	ti.Placeholder = "Type / to open command palette or 'exit' to quit"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	// Define menu items
	items := []list.Item{
		MenuItem{title: "Status", description: "Show database status and metrics", command: "status"},
		MenuItem{title: "Schema Discovery", description: "List public base tables on primary DB", command: "schema-discovery"},
		MenuItem{title: "Schema Check", description: "Validate schema consistency between DBs", command: "schema-check"},
		MenuItem{title: "Sync Health Check", description: "Show per-table counts and sync status", command: "sync-health-check"},
		MenuItem{title: "Print Schema", description: "Print detailed schema information", command: "print-schema"},
		MenuItem{title: "Print Table", description: "Print detailed table information", command: "print-table"},
		MenuItem{title: "Print Tables", description: "Print detailed info for all tables", command: "print-tables"},
		MenuItem{title: "Print All", description: "Print all schemas and tables", command: "print-all"},
		MenuItem{title: "Print Table Data", description: "Print actual table data/entries", command: "print-table-data"},
		MenuItem{title: "Sync", description: "Synchronize databases", command: "sync"},
		MenuItem{title: "Migrate", description: "Run database migrations", command: "migrate"},
		MenuItem{title: "CSV Backup", description: "Export table(s) to CSV files", command: "csv-backup"},
		MenuItem{title: "CSV Seed", description: "Import CSV files into database", command: "csv-seed"},
		MenuItem{title: "SQL Query", description: "Execute direct SQL query", command: "sql"},
		MenuItem{title: "Sizes", description: "Show database sizes", command: "sizes"},
		MenuItem{title: "Exit", description: "Exit the application", command: "exit"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "📋 DBManager Command Palette"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	return tuiModel{
		textInput: ti,
		list:      l,
		state:     "input", // Start in input state (not connecting)
		connected: false,
	}
}

func (m tuiModel) Init() tea.Cmd {
	// Start with input focused, user can type commands
	return textinput.Blink
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width - 4)
		m.list.SetHeight(msg.Height - 10)
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case "input":
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "enter":
				input := strings.TrimSpace(m.textInput.Value())
				if input == "" {
					return m, nil
				}
				if input == "exit" || input == "quit" {
					m.quitting = true
					return m, tea.Quit
				}
				
				// Execute the command with arguments
				m.selectedCmd = input
				m.state = "executing"
				m.textInput.SetValue("")
				return m, func() tea.Msg {
					return executeCommandWithArgs(input)
				}
			
			// Check if user types "/" to trigger menu immediately
			default:
				// First update the text input
				m.textInput, cmd = m.textInput.Update(msg)
				input := m.textInput.Value()
				
				// If user just typed "/" (single character), immediately show menu
				if input == "/" {
					// Connect if not already connected
					if !m.connected {
						m.state = "connecting"
						m.textInput.SetValue("")
						return m, func() tea.Msg {
							return connectToDatabase()
						}
					} else {
						// Already connected, show menu directly
						m.state = "menu"
						m.textInput.SetValue("")
						return m, nil
					}
				}
				
				return m, cmd
			}

		case "menu":
			switch msg.String() {
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
				
			case "esc":
				// Go back to input without inserting anything
				m.state = "input"
				m.textInput.SetValue("")
				return m, nil

			case "enter":
				// Insert selected command into input field
				if item, ok := m.list.SelectedItem().(MenuItem); ok {
					if item.command == "exit" {
						m.quitting = true
						return m, tea.Quit
					}
					
					// Insert command into input field and return to input mode
					m.textInput.SetValue(item.command + " ")
					m.textInput.CursorEnd()
					m.state = "input"
					return m, nil
				}
			}

		case "result":
			switch msg.String() {
			case "enter", "esc":
				m.state = "input" // Return to input mode
				m.errorMsg = ""
				m.resultMsg = ""
				m.textInput.SetValue("")
				return m, nil
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
		}

	case connectMsg:
		if msg.success {
			m.connected = true
			m.state = "menu" // Show menu after successful connection
			m.errorMsg = ""
		} else {
			m.state = "input"
			m.errorMsg = fmt.Sprintf("Failed to connect to database: %v", msg.err)
		}
		return m, nil

	case executeMsg:
		m.state = "result"
		if msg.success {
			m.resultMsg = msg.output
			m.errorMsg = ""
		} else {
			m.errorMsg = msg.output // executeMsg.output contains error message
			m.resultMsg = ""
		}
		return m, nil
	}

	// Update the active component
	switch m.state {
	case "input":
		m.textInput, cmd = m.textInput.Update(msg)
	case "menu":
		m.list, cmd = m.list.Update(msg)
	}

	return m, cmd
}

func (m tuiModel) View() string {
	if m.quitting {
		return successStyle.Render("👋 Goodbye! Thanks for using DBManager.\n")
	}

	// Ensure minimum dimensions if not set yet
	if m.width == 0 {
		m.width = 80
	}
	if m.height == 0 {
		m.height = 24
	}

	var content string

	// Title
	title := titleStyle.Render("🗄️  DBManager Interactive CLI")

	switch m.state {
	case "input":
		// Input state
		var statusLine string
		if m.connected {
			statusLine = successStyle.Render("✅ Connected to database")
		} else {
			statusLine = dimmedStyle.Render("Not connected (type '/' to connect)")
		}
		
		help := helpStyle.Render("Type command or '/' for menu • 'exit' to quit • Ctrl+C to exit")
		inputBox := inputStyle.Render(m.textInput.View())

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			statusLine,
			"",
			inputBox,
			"",
			help,
		)

		if m.errorMsg != "" {
			content += "\n\n" + errorStyle.Render("❌ "+m.errorMsg)
		}

	case "connecting":
		// Connecting state
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			dimmedStyle.Render("🔌 Connecting to database..."),
		)

	case "menu":
		// Menu state - ensure list is properly sized
		if m.width > 0 && m.height > 0 {
			listHeight := m.height - 10 // Reserve space for title, status, help
			if listHeight < 5 {
				listHeight = 5
			}
			listWidth := m.width - 4
			if listWidth < 40 {
				listWidth = 40
			}
			m.list.SetSize(listWidth, listHeight)
		}
		
		connectionStatus := successStyle.Render("✅ Connected")
		help := helpStyle.Render("↑/↓: Navigate • Enter: Select • Esc: Back • Ctrl+C: Exit")

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			connectionStatus,
			"",
			m.list.View(),
			"",
			help,
		)

	case "executing":
		// Executing state
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			dimmedStyle.Render(fmt.Sprintf("⚙️  Executing command: %s...", m.selectedCmd)),
		)

	case "result":
		// Result state - truncate long output to prevent overflow
		help := helpStyle.Render("Press Enter or Esc to continue")

		var result string
		if m.errorMsg != "" {
			result = errorStyle.Render("❌ " + m.errorMsg)
		} else {
			// Limit result lines to fit screen
			maxResultLines := m.height - 8
			if maxResultLines < 5 {
				maxResultLines = 5
			}
			lines := strings.Split(m.resultMsg, "\n")
			if len(lines) > maxResultLines {
				lines = lines[:maxResultLines]
				lines = append(lines, dimmedStyle.Render("... (output truncated, use CLI for full output)"))
			}
			result = strings.Join(lines, "\n")
		}

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			result,
			"",
			help,
		)

	}

	// Always render at top-left with consistent sizing
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Left,
			lipgloss.Top,
			content,
		)
	}

	return content
}

// connectToDatabase connects to the database in the background
func connectToDatabase() connectMsg {
	// Initialize database connection silently
	err := initializeForSQLSilent("database.yaml")
	if err != nil {
		return connectMsg{success: false, err: err}
	}
	return connectMsg{success: true}
}

// executeCommandWithArgs executes a command string with arguments
func executeCommandWithArgs(input string) executeMsg {
	var output strings.Builder
	
	// Parse command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return executeMsg{
			success: false,
			output:  "❌ No command provided",
		}
	}
	
	cmd := parts[0]
	args := parts[1:]
	
	// Parse arguments into a map
	argMap := make(map[string]string)
	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			kv := strings.SplitN(arg[2:], "=", 2)
			if len(kv) == 2 {
				argMap[kv[0]] = kv[1]
			} else {
				argMap[kv[0]] = "true" // Flag without value
			}
		}
	}
	
	// Execute based on command
	switch cmd {
	case "status":
		handleStatus()
		output.WriteString("✅ Status command completed successfully\n")
		output.WriteString("Check the terminal output above for detailed status.")
		
	case "schema-discovery":
		handleSchemaDiscovery()
		output.WriteString("✅ Schema discovery completed successfully\n")
		output.WriteString("Check the terminal output above for schema list.")
		
	case "schema-check":
		handleSchemaCheck()
		output.WriteString("✅ Schema check completed successfully\n")
		output.WriteString("Check the terminal output above for schema consistency report.")
		
	case "sync-health-check":
		handleSyncHealthCheck()
		output.WriteString("✅ Sync health check completed successfully\n")
		output.WriteString("Check the terminal output above for sync status.")
		
	case "print-schema":
		schemaName := argMap["schema"]
		target := argMap["target"]
		if target == "" {
			target = "primary"
		}
		if schemaName == "" {
			output.WriteString("❌ Error: --schema parameter is required\n\n")
			output.WriteString("Usage: print-schema --schema=<name> --target=<primary|backup|signal_db>\n")
			output.WriteString("Example: print-schema --schema=auth --target=primary")
		} else {
			handlePrintSchema(schemaName, target)
			output.WriteString(fmt.Sprintf("✅ Schema '%s' printed successfully from %s\n", schemaName, target))
			output.WriteString("Check the terminal output above for detailed schema info.")
		}
		
	case "print-table":
		tableName := argMap["table"]
		target := argMap["target"]
		if target == "" {
			target = "primary"
		}
		if tableName == "" {
			output.WriteString("❌ Error: --table parameter is required\n\n")
			output.WriteString("Usage: print-table --table=<name> --target=<primary|backup|signal_db>\n")
			output.WriteString("Example: print-table --table=users --target=primary")
		} else {
			handlePrintTable(tableName, target)
			output.WriteString(fmt.Sprintf("✅ Table '%s' printed successfully from %s\n", tableName, target))
			output.WriteString("Check the terminal output above for detailed table info.")
		}
		
	case "print-tables":
		schemaName := argMap["schema"]
		target := argMap["target"]
		if target == "" {
			target = "primary"
		}
		if schemaName == "" {
			output.WriteString("❌ Error: --schema parameter is required\n\n")
			output.WriteString("Usage: print-tables --schema=<name> --target=<primary|backup|signal_db>\n")
			output.WriteString("Example: print-tables --schema=auth --target=primary")
		} else {
			handlePrintTables(schemaName, target)
			output.WriteString(fmt.Sprintf("✅ Tables in schema '%s' printed successfully from %s\n", schemaName, target))
			output.WriteString("Check the terminal output above for detailed tables info.")
		}
		
	case "print-all":
		target := argMap["target"]
		if target == "" {
			target = "primary"
		}
		handlePrintAll(target)
		output.WriteString(fmt.Sprintf("✅ All schemas and tables printed successfully from %s\n", target))
		output.WriteString("Check the terminal output above for complete database overview.")
		
	case "print-table-data":
		tableName := argMap["table"]
		target := argMap["target"]
		limitStr := argMap["limit"]
		if target == "" {
			target = "primary"
		}
		limit := 10 // default limit
		if limitStr != "" {
			if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
				limit = parsedLimit
			}
		}
		if tableName == "" {
			output.WriteString("❌ Error: --table parameter is required\n\n")
			output.WriteString("Usage: print-table-data --table=<name> --limit=<num> --target=<primary|backup|signal_db>\n")
			output.WriteString("Example: print-table-data --table=users --limit=10 --target=primary")
		} else {
			handlePrintTableData(tableName, target, limit)
			output.WriteString(fmt.Sprintf("✅ Data from table '%s' printed successfully from %s\n", tableName, target))
			output.WriteString("Check the terminal output above for table data.")
		}
		
	case "sync":
		if argMap["commit"] != "true" {
			output.WriteString("⚠️  Dry-run mode (use --commit to apply changes)\n\n")
		}
		output.WriteString("Note: Sync command execution in TUI is limited.\n")
		output.WriteString("For full sync functionality, please use CLI mode:\n")
		output.WriteString("  dbmanager sync --commit --prune --table=<name>")
		
	case "migrate":
		target := argMap["target"]
		if target == "" {
			target = "primary"
		}
		output.WriteString("Note: Migration command execution in TUI is limited.\n")
		output.WriteString("For full migration functionality, please use CLI mode:\n")
		output.WriteString(fmt.Sprintf("  dbmanager migrate --target=%s", target))
		
	case "csv-backup":
		tableName := argMap["table"]
		source := argMap["source"]
		if source == "" {
			source = "primary"
		}
		if tableName == "" {
			output.WriteString("❌ Error: --table parameter is required\n\n")
			output.WriteString("Usage: csv-backup --table=<name> --source=<primary|backup|signal_db>\n")
			output.WriteString("Example: csv-backup --table=users --source=primary")
		} else {
			output.WriteString("Note: CSV backup command execution in TUI is limited.\n")
			output.WriteString("For full CSV backup functionality, please use CLI mode:\n")
			output.WriteString(fmt.Sprintf("  dbmanager csv-backup --table=%s --source=%s", tableName, source))
		}
		
	case "csv-seed":
		tableName := argMap["table"]
		target := argMap["target"]
		if target == "" {
			target = "primary"
		}
		if tableName == "" {
			output.WriteString("❌ Error: --table parameter is required\n\n")
			output.WriteString("Usage: csv-seed --table=<name> --file=<path> --target=<primary|backup|signal_db>\n")
			output.WriteString("Example: csv-seed --table=users --file=users.csv --target=primary")
		} else {
			output.WriteString("Note: CSV seed command execution in TUI is limited.\n")
			output.WriteString("For full CSV seed functionality, please use CLI mode:\n")
			output.WriteString(fmt.Sprintf("  dbmanager csv-seed --table=%s --target=%s", tableName, target))
		}
		
	case "sql":
		sqlQuery := argMap["sql"]
		target := argMap["target"]
		if target == "" {
			target = "primary"
		}
		if sqlQuery == "" {
			output.WriteString("❌ Error: --sql parameter is required\n\n")
			output.WriteString("Usage: sql --sql=\"<query>\" --target=<primary|backup|signal_db>\n")
			output.WriteString("Example: sql --sql=\"SELECT COUNT(*) FROM users\" --target=primary")
		} else {
			handleSQL(sqlQuery, target)
			output.WriteString(fmt.Sprintf("✅ SQL query executed successfully on %s\n", target))
			output.WriteString("Check the terminal output above for query results.")
		}
		
	case "sizes":
		handleSizes()
		output.WriteString("✅ Database sizes retrieved successfully\n")
		output.WriteString("Check the terminal output above for size information.")
		
	default:
		output.WriteString(fmt.Sprintf("❌ Unknown command: %s\n\n", cmd))
		output.WriteString("Type '/' to see available commands or 'exit' to quit.")
	}
	
	return executeMsg{
		success: true,
		output:  output.String(),
	}
}

// RunInteractiveTUI starts the interactive TUI
func RunInteractiveTUI() error {
	// Use alt screen to avoid terminal pollution
	// WithAltScreen ensures clean state on entry/exit
	p := tea.NewProgram(
		initialTuiModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	
	model, err := p.Run()
	
	// Clean up on exit
	if err != nil {
		return err
	}
	
	// Check if user wants to quit
	if m, ok := model.(tuiModel); ok && m.quitting {
		// Clean exit
		return nil
	}
	
	return err
}

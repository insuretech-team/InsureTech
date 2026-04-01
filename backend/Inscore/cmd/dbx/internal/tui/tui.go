package tui

// ─────────────────────────────────────────────────────────────────────────────
// Double-FSM TUI for DBManager
//
// FSM 1 — DB State Machine (dbState):
//   disconnected → connecting → connected → idle
//   Any state can transition to: error
//
// FSM 2 — TUI Screen Machine (tuiState):
//   input → connecting → menu → executing → pager → input
//
// The two FSMs are orthogonal:
//   dbState  tracks whether we have a live DB connection
//   tuiState tracks which screen/widget is currently rendered
//
// Key fix: all handler output is captured via io.Pipe / os.Stdout redirect
// into a strings.Builder and then displayed in the viewport pager.
// Nothing ever calls fmt.Println to raw stdout while the TUI is running.
// ─────────────────────────────────────────────────────────────────────────────

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── DB FSM states ─────────────────────────────────────────────────────────────

type dbState uint8

const (
	dbDisconnected dbState = iota // no connection attempted
	dbConnecting                  // background goroutine running
	dbConnected                   // healthy connection
	dbError                       // connection failed
)

func (s dbState) String() string {
	return [...]string{"disconnected", "connecting", "connected", "error"}[s]
}

// ── TUI FSM states ────────────────────────────────────────────────────────

type tuiState uint8

const (
	tuiInput      tuiState = iota // main command input prompt
	tuiConnecting                 // spinner while DB connects
	tuiMenu                       // command palette list
	tuiForm                       // form fields for command arguments
	tuiExecuting                  // spinner while command runs
	tuiPager                      // scrollable result viewport
)

func (s tuiState) String() string {
	return [...]string{"input", "connecting", "menu", "form", "executing", "pager"}[s]
}

// ── Handlers ──────────────────────────────────────────────────────────────────

// Handlers provides dbx operations used by the TUI.
// Each handler that produces output accepts an io.Writer so the TUI can
// capture the output instead of letting it leak to raw stdout.
type Handlers struct {
	InitializeForSQLSilent func(configPath string) error

	// Output-capturing handlers — write to w instead of os.Stdout
	Status          func(w io.Writer)
	SchemaDiscovery func(w io.Writer)
	SchemaCheck     func(w io.Writer)
	SyncHealthCheck func(w io.Writer)
	PrintSchema     func(w io.Writer, schemaName, targetDB string)
	PrintTable      func(w io.Writer, tableName, targetDB string)
	PrintTables     func(w io.Writer, schemaName, targetDB string)
	PrintAll        func(w io.Writer, targetDB string)
	PrintTableData  func(w io.Writer, tableName, targetDB string, limit int)
	SQL             func(w io.Writer, query, targetDB string)
	Sizes           func(w io.Writer)

	// Autocomplete providers — called once after connect to populate dropdowns
	GetSchemas func() ([]string, error)              // returns all schema names
	GetTables  func(schema string) ([]string, error) // returns all table names in schema
}

// ── Styles ────────────────────────────────────────────────────────────────────

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

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#3B82F6")).
				Bold(true)

	_ = selectedItemStyle // suppress unused warning (used by list delegate)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	dimmedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262"))

	pagerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3B82F6")).
			Padding(0, 1)
)

// ── Menu items ────────────────────────────────────────────────────────────────

type MenuItem struct {
	title       string
	description string
	command     string
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.description }
func (i MenuItem) FilterValue() string { return i.title }

var menuItems = []list.Item{
	MenuItem{"Status", "Show database status and metrics", "status"},
	MenuItem{"Schema Discovery", "List public base tables on primary DB", "schema-discovery"},
	MenuItem{"Schema Check", "Validate schema consistency between DBs", "schema-check"},
	MenuItem{"Sync Health Check", "Show per-table counts and sync status", "sync-health-check"},
	MenuItem{"Print Schema", "Print detailed schema info (--schema=<name>)", "print-schema"},
	MenuItem{"Print Table", "Print detailed table info (--table=<name>)", "print-table"},
	MenuItem{"Print Tables", "Print all tables in schema (--schema=<name>)", "print-tables"},
	MenuItem{"Print All", "Print all schemas and tables", "print-all"},
	MenuItem{"Print Table Data", "Print table rows (--table=<name> --limit=N)", "print-table-data"},
	MenuItem{"SQL Query", "Execute SQL (--sql=\"<query>\" --target=<db>)", "sql"},
	MenuItem{"Sizes", "Show database sizes", "sizes"},
	MenuItem{"Sync", "Synchronize databases", "sync"},
	MenuItem{"Migrate", "Run database migrations", "migrate"},
	MenuItem{"CSV Backup", "Export table(s) to CSV files", "csv-backup"},
	MenuItem{"CSV Seed", "Import CSV files into database", "csv-seed"},
	MenuItem{"Help", "Show command reference", "help"},
	MenuItem{"Exit", "Exit the application", "exit"},
}

// ── Form types ────────────────────────────────────────────────────────────

type FieldType uint8

const (
	FieldTypeSelect FieldType = iota
	FieldTypeText
)

type FormField struct {
	Name        string // argument name e.g. "schema", "table", "target"
	Label       string // display label e.g. "Select schema:"
	Type        FieldType
	Options     []string                                          // for FieldTypeSelect — populated dynamically
	LoadOptions func(answers map[string]string) ([]string, error) // lazy loader using prior answers
	Optional    bool                                              // if true, show a "(skip)" option at top of list
	Default     string                                            // default value if skipped
}

type CommandForm struct {
	Command string // base command e.g. "print-table"
	Fields  []FormField
}

// ── Tea messages ──────────────────────────────────────────────────────────────

type connectDoneMsg struct{ err error }
type executeDoneMsg struct {
	output string
	err    error
}
type loadOptionsMsg struct {
	options []string
	err     error
}
type cacheLoadedMsg struct {
	schemas []string
	tables  map[string][]string
	err     error
}

// ── Model ─────────────────────────────────────────────────────────────────────

type Model struct {
	// FSM states
	db  dbState
	tui tuiState

	// widgets
	input    textinput.Model
	list     list.Model
	viewport viewport.Model
	spinner  spinner.Model

	// handlers
	handlers Handlers

	// display
	width  int
	height int

	// runtime
	selectedCmd string
	errorMsg    string
	quitting    bool

	// Form state
	currentForm *CommandForm
	formAnswers map[string]string
	formStep    int             // which field we're on
	formList    list.Model      // picker for FieldTypeSelect
	formInput   textinput.Model // for FieldTypeText
	formLoading bool            // true while LoadOptions is running

	// Cached autocomplete data (loaded after connect)
	cachedSchemas []string
	cachedTables  map[string][]string // schema -> tables

	// Pre-built command forms (rebuilt after cache loads so closures see live data)
	forms []CommandForm
}

func initialModel(handlers Handlers) Model {
	// Text input
	ti := textinput.New()
	ti.Placeholder = "Type command or '/' for menu • 'exit' to quit"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	// Command list
	l := list.New(menuItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = "dbx Command Palette"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	// Spinner
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#3B82F6"))

	// Viewport (sized properly on first WindowSizeMsg)
	vp := viewport.New(80, 20)
	vp.Style = pagerStyle

	m := Model{
		db:           dbDisconnected,
		tui:          tuiInput,
		input:        ti,
		list:         l,
		viewport:     vp,
		spinner:      sp,
		handlers:     handlers,
		formAnswers:  make(map[string]string),
		cachedTables: make(map[string][]string),
	}
	m.rebuildForms()
	return m
}

// ── Init ──────────────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// ── Window resize ─────────────────────────────────────────────────────────
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width-4, msg.Height-12)
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10
		return m, nil

	// ── Spinner tick (only when animating) ───────────────────────────────────
	case spinner.TickMsg:
		if m.tui == tuiConnecting || m.tui == tuiExecuting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	// ── DB connect result ─────────────────────────────────────────────────
	case connectDoneMsg:
		if msg.err != nil {
			m.db = dbError
			m.tui = tuiInput
			m.errorMsg = fmt.Sprintf("Connection failed: %v", msg.err)
		} else {
			m.db = dbConnected
			m.tui = tuiMenu
			m.errorMsg = ""
			// Load cache asynchronously
			cmds = append(cmds, m.loadCache())
		}
		return m, tea.Batch(cmds...)

	// ── Cache loaded ───────────────────────────────────────────────────────
	case cacheLoadedMsg:
		if msg.err != nil {
			// Non-fatal: forms still work with fallback "public"
			m.errorMsg = fmt.Sprintf("Schema cache: %v", msg.err)
		} else {
			m.cachedSchemas = msg.schemas
			m.cachedTables = msg.tables
			m.errorMsg = ""
		}
		// Rebuild forms so closures capture the now-populated cache via pointer
		m.rebuildForms()
		return m, nil

	// ── Load options result ────────────────────────────────────────────────
	case loadOptionsMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.formList = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
		} else {
			items := make([]list.Item, 0, len(msg.options)+2)
			items = append(items, MenuItem{"← Back", "Go back to previous field", ""})
			currentField := m.currentForm.Fields[m.formStep]
			if currentField.Optional {
				items = append(items, MenuItem{
					fmt.Sprintf("(use default: %s)", currentField.Default),
					"Skip this field",
					"__default__",
				})
			}
			for _, opt := range msg.options {
				items = append(items, MenuItem{opt, "", opt})
			}
			m.formList = list.New(items, list.NewDefaultDelegate(), 0, 0)
			m.formList.Title = fmt.Sprintf("Step %d/%d: %s", m.formStep+1, len(m.currentForm.Fields), currentField.Label)
			m.formList.SetShowStatusBar(false)
			m.formList.SetFilteringEnabled(false)
		}
		m.formLoading = false
		return m, nil

	// ── Command execute result ────────────────────────────────────────────────
	case executeDoneMsg:
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.viewport.SetContent(errorStyle.Render("Error: " + msg.err.Error()))
		} else {
			m.errorMsg = ""
			m.viewport.SetContent(msg.output)
		}
		m.viewport.GotoTop()
		m.tui = tuiPager
		return m, nil

	// ── Keyboard ──────────────────────────────────────────────────────────────
	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.tui {

		// ── Input screen ──────────────────────────────────────────────────────
		case tuiInput:
			switch msg.String() {
			case "enter":
				raw := strings.TrimSpace(m.input.Value())
				if raw == "" {
					return m, nil
				}
				if raw == "exit" || raw == "quit" {
					m.quitting = true
					return m, tea.Quit
				}
				if raw == "/" {
					return m.openMenu()
				}
				m.selectedCmd = raw
				m.input.SetValue("")
				return m.startExecute(raw)

			default:
				var cmd tea.Cmd
				m.input, cmd = m.input.Update(msg)
				// Auto-open menu when user types "/"
				if m.input.Value() == "/" {
					return m.openMenu()
				}
				cmds = append(cmds, cmd)
			}

		// ── Menu screen ───────────────────────────────────────────────────
		case tuiMenu:
			switch msg.String() {
			case "esc":
				m.tui = tuiInput
				m.input.SetValue("")
				return m, nil
			case "enter":
				if item, ok := m.list.SelectedItem().(MenuItem); ok {
					if item.command == "exit" {
						m.quitting = true
						return m, tea.Quit
					}
					// Check if command needs a form or direct execution
					form := m.findFormForCommand(item.command)
					if form != nil {
						// Start form
						return m.startForm(form)
					}
					// Direct execution (no form needed)
					m.input.SetValue(item.command)
					m.selectedCmd = item.command
					m.input.SetValue("")
					return m.startExecute(item.command)
				}
			default:
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
			}

		// ── Form screen ───────────────────────────────────────────────────────
		case tuiForm:
			switch msg.String() {
			case "esc":
				if m.formStep > 0 {
					m.formStep--
					return m.updateFormStep()
				}
				// Go back to menu if on first step
				m.tui = tuiMenu
				m.currentForm = nil
				m.formAnswers = make(map[string]string)
				m.formStep = 0
				return m, nil
			case "enter":
				currentField := m.currentForm.Fields[m.formStep]
				if currentField.Type == FieldTypeSelect {
					if item, ok := m.formList.SelectedItem().(MenuItem); ok {
						if item.command == "" && item.title == "← Back" {
							// Go back
							if m.formStep > 0 {
								m.formStep--
								return m.updateFormStep()
							}
							// Back to menu
							m.tui = tuiMenu
							m.currentForm = nil
							m.formAnswers = make(map[string]string)
							m.formStep = 0
							return m, nil
						}
						// Record answer
						if item.command == "__default__" {
							m.formAnswers[currentField.Name] = currentField.Default
						} else {
							m.formAnswers[currentField.Name] = item.command
						}
						// Move to next field or execute
						if m.formStep < len(m.currentForm.Fields)-1 {
							m.formStep++
							return m.updateFormStep()
						}
						// All fields filled — execute command
						return m.executeForm()
					}
				} else if currentField.Type == FieldTypeText {
					// Record text input
					val := strings.TrimSpace(m.formInput.Value())
					if val == "" && currentField.Optional {
						m.formAnswers[currentField.Name] = currentField.Default
					} else {
						m.formAnswers[currentField.Name] = val
					}
					m.formInput.SetValue("")
					// Move to next field or execute
					if m.formStep < len(m.currentForm.Fields)-1 {
						m.formStep++
						return m.updateFormStep()
					}
					// All fields filled — execute command
					return m.executeForm()
				}
			default:
				if m.currentForm.Fields[m.formStep].Type == FieldTypeSelect {
					var cmd tea.Cmd
					m.formList, cmd = m.formList.Update(msg)
					cmds = append(cmds, cmd)
				} else {
					var cmd tea.Cmd
					m.formInput, cmd = m.formInput.Update(msg)
					cmds = append(cmds, cmd)
				}
			}

		// ── Pager screen ───────────────────────────────────────────────────
		case tuiPager:
			switch msg.String() {
			case "esc", "q", "enter":
				m.tui = tuiInput
				m.errorMsg = ""
				m.input.SetValue("")
				return m, nil
			default:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				cmds = append(cmds, cmd)
			}

		// ── Connecting / Executing — ignore key input ─────────────────────────
		case tuiConnecting, tuiExecuting:
			// swallow all keys except ctrl+c (handled above)
		}
	}

	return m, tea.Batch(cmds...)
}

// openMenu transitions to the menu, connecting first if needed.
func (m Model) openMenu() (tea.Model, tea.Cmd) {
	m.input.SetValue("")
	if m.db != dbConnected {
		m.db = dbConnecting
		m.tui = tuiConnecting
		return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
			err := m.handlers.InitializeForSQLSilent("database.yaml")
			return connectDoneMsg{err: err}
		})
	}
	m.tui = tuiMenu
	return m, nil
}

// startExecute transitions to executing state and fires the command.
func (m Model) startExecute(input string) (tea.Model, tea.Cmd) {
	// Connect first if needed (for direct commands without going through menu)
	if m.db != dbConnected {
		m.db = dbConnecting
		m.tui = tuiConnecting
		handlers := m.handlers
		return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
			if err := handlers.InitializeForSQLSilent("database.yaml"); err != nil {
				return connectDoneMsg{err: err}
			}
			// After connect, immediately run the command
			out, err := runCommand(input, handlers)
			if err != nil {
				return executeDoneMsg{err: err}
			}
			return executeDoneMsg{output: out}
		})
	}
	m.tui = tuiExecuting
	handlers := m.handlers
	return m, tea.Batch(m.spinner.Tick, func() tea.Msg {
		out, err := runCommand(input, handlers)
		if err != nil {
			return executeDoneMsg{err: err}
		}
		return executeDoneMsg{output: out}
	})
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.quitting {
		return successStyle.Render("Goodbye! Thanks for using dbx.\n")
	}

	w := m.width
	if w == 0 {
		w = 80
	}
	h := m.height
	if h == 0 {
		h = 24
	}

	title := titleStyle.Render("dbx — InsureTech Database eXplorer")
	dbBadge := m.dbBadge()

	switch m.tui {

	case tuiInput:
		help := helpStyle.Render("'/' = menu  •  type command + Enter  •  'exit' to quit  •  Ctrl+C")
		inputBox := inputStyle.Render(m.input.View())
		rows := []string{title, "", dbBadge, "", inputBox, "", help}
		if m.errorMsg != "" {
			rows = append(rows, "", errorStyle.Render("Error: "+m.errorMsg))
		}
		return lipgloss.Place(w, h, lipgloss.Left, lipgloss.Top,
			lipgloss.JoinVertical(lipgloss.Left, rows...))

	case tuiForm:
		return m.renderForm(w, h, title, dbBadge)

	case tuiConnecting:
		body := lipgloss.JoinVertical(lipgloss.Left,
			title, "",
			m.spinner.View()+" Connecting to database...",
		)
		return lipgloss.Place(w, h, lipgloss.Left, lipgloss.Top, body)

	case tuiMenu:
		m.list.SetSize(w-4, h-12)
		help := helpStyle.Render("↑/↓ navigate  •  Enter select  •  Esc back  •  / filter  •  Ctrl+C quit")
		body := lipgloss.JoinVertical(lipgloss.Left,
			title, "", dbBadge, "",
			m.list.View(), "",
			help,
		)
		return lipgloss.Place(w, h, lipgloss.Left, lipgloss.Top, body)

	case tuiExecuting:
		body := lipgloss.JoinVertical(lipgloss.Left,
			title, "",
			m.spinner.View()+" Executing: "+dimmedStyle.Render(m.selectedCmd),
		)
		return lipgloss.Place(w, h, lipgloss.Left, lipgloss.Top, body)

	case tuiPager:
		m.viewport.Width = w - 4
		m.viewport.Height = h - 10
		scrollPct := fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100)
		help := helpStyle.Render("↑/↓ or j/k scroll  •  PgUp/PgDn  •  Enter/Esc/q back  •  Ctrl+C quit")
		scrollInfo := dimmedStyle.Render(scrollPct)
		body := lipgloss.JoinVertical(lipgloss.Left,
			title, "",
			lipgloss.JoinHorizontal(lipgloss.Top,
				dimmedStyle.Render("Command: "+m.selectedCmd),
				strings.Repeat(" ", max(0, w-20-len(m.selectedCmd))),
				scrollInfo,
			),
			"",
			m.viewport.View(),
			"",
			help,
		)
		return lipgloss.Place(w, h, lipgloss.Left, lipgloss.Top, body)
	}

	return ""
}

// dbBadge returns a styled one-line DB connection status.
func (m Model) dbBadge() string {
	switch m.db {
	case dbConnected:
		return successStyle.Render("Connected to database")
	case dbConnecting:
		return dimmedStyle.Render("Connecting...")
	case dbError:
		return errorStyle.Render("Connection error")
	default:
		return dimmedStyle.Render("Not connected (type '/' to connect)")
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ── Output capture helper ─────────────────────────────────────────────────────

// captureStdout redirects os.Stdout to a pipe, calls fn(), then returns
// everything fn wrote to stdout as a string.
// This is the safety net for any handler that still uses fmt.Println internally.
func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fn()
		return ""
	}
	os.Stdout = w

	var wg sync.WaitGroup
	var buf strings.Builder
	wg.Add(1)
	go func() {
		defer wg.Done()
		io.Copy(&buf, r)
	}()

	fn()
	w.Close()
	os.Stdout = old
	wg.Wait()
	r.Close()
	return buf.String()
}

// ── Command dispatcher ────────────────────────────────────────────────────────

// runCommand parses the input string, dispatches to the correct handler,
// and returns the captured output as a string.
func runCommand(input string, h Handlers) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", fmt.Errorf("no command provided")
	}

	cmd := parts[0]
	argMap := parseArgs(parts[1:])

	var buf strings.Builder

	switch cmd {

	case "status":
		if h.Status == nil {
			return "", fmt.Errorf("status handler not configured")
		}
		h.Status(&buf)

	case "schema-discovery":
		if h.SchemaDiscovery == nil {
			return "", fmt.Errorf("schema-discovery handler not configured")
		}
		h.SchemaDiscovery(&buf)

	case "schema-check":
		if h.SchemaCheck == nil {
			return "", fmt.Errorf("schema-check handler not configured")
		}
		h.SchemaCheck(&buf)

	case "sync-health-check":
		if h.SyncHealthCheck == nil {
			return "", fmt.Errorf("sync-health-check handler not configured")
		}
		h.SyncHealthCheck(&buf)

	case "print-schema":
		if h.PrintSchema == nil {
			return "", fmt.Errorf("print-schema handler not configured")
		}
		schema := argMap["schema"]
		target := argOrDefault(argMap, "target", "primary")
		h.PrintSchema(&buf, schema, target)

	case "print-table":
		if h.PrintTable == nil {
			return "", fmt.Errorf("print-table handler not configured")
		}
		table := argMap["table"]
		if table == "" {
			return "", fmt.Errorf("--table is required\nUsage: print-table --table=<name> [--target=primary|backup]")
		}
		target := argOrDefault(argMap, "target", "primary")
		h.PrintTable(&buf, table, target)

	case "print-tables":
		if h.PrintTables == nil {
			return "", fmt.Errorf("print-tables handler not configured")
		}
		schema := argMap["schema"]
		target := argOrDefault(argMap, "target", "primary")
		h.PrintTables(&buf, schema, target)

	case "print-all":
		if h.PrintAll == nil {
			return "", fmt.Errorf("print-all handler not configured")
		}
		target := argOrDefault(argMap, "target", "primary")
		h.PrintAll(&buf, target)

	case "print-table-data":
		if h.PrintTableData == nil {
			return "", fmt.Errorf("print-table-data handler not configured")
		}
		table := argMap["table"]
		if table == "" {
			return "", fmt.Errorf("--table is required\nUsage: print-table-data --table=<name> [--limit=N] [--target=primary|backup]")
		}
		target := argOrDefault(argMap, "target", "primary")
		limit := 10
		if ls := argMap["limit"]; ls != "" {
			if n, err := strconv.Atoi(ls); err == nil && n > 0 {
				limit = n
			}
		}
		h.PrintTableData(&buf, table, target, limit)

	case "sql":
		if h.SQL == nil {
			return "", fmt.Errorf("sql handler not configured")
		}
		query := argMap["sql"]
		if query == "" {
			return "", fmt.Errorf("--sql is required\nUsage: sql --sql=\"<query>\" [--target=primary|backup|both]")
		}
		target := argOrDefault(argMap, "target", "primary")
		h.SQL(&buf, query, target)

	case "sizes":
		if h.Sizes == nil {
			return "", fmt.Errorf("sizes handler not configured")
		}
		h.Sizes(&buf)

	case "help":
		fmt.Fprint(&buf, getHelpReference())

	case "sync", "migrate", "csv-backup", "csv-seed":
		fmt.Fprintf(&buf, "Command '%s' is not supported in the interactive TUI.\n", cmd)
		fmt.Fprintf(&buf, "Use the CLI directly:\n\n")
		fmt.Fprintf(&buf, "  dbx %s %s\n", cmd, strings.Join(parts[1:], " "))

	default:
		fmt.Fprintf(&buf, "Unknown command: %s\n\n", cmd)
		fmt.Fprintf(&buf, "Type '/' to browse available commands.\n")
		fmt.Fprintf(&buf, "Or type a command like: status, print-schema --schema=public\n")
	}

	return buf.String(), nil
}

// parseArgs converts ["--key=value", "--flag"] into map[key]value.
func parseArgs(args []string) map[string]string {
	m := make(map[string]string)
	for _, a := range args {
		if strings.HasPrefix(a, "--") {
			kv := strings.SplitN(a[2:], "=", 2)
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			} else {
				m[kv[0]] = "true"
			}
		}
	}
	return m
}

func argOrDefault(m map[string]string, key, def string) string {
	if v, ok := m[key]; ok && v != "" {
		return v
	}
	return def
}

// ── Help reference ───────────────────────────────────────────────────────

func getHelpReference() string {
	return `═══════════════════════════════════════════════════════════════════════════════
 dbx — InsureTech Database eXplorer
═══════════════════════════════════════════════════════════════════════════════

 CLI:  dbx <command> [flags]
 TUI:  dbx  (no args) → type command or press '/' for guided menu

 STATUS & DIAGNOSTICS
 ─────────────────────────────────────────────────────────────────────────────
  status              Connection status and metrics
  schema-discovery    List all base tables on primary
  schema-check        Validate primary vs backup schema consistency
  sync-health-check   Per-table row counts (primary vs backup)
  sizes               Database and table sizes

 SYNC
 ─────────────────────────────────────────────────────────────────────────────
  sync                Sync primary → backup (authoritative upsert)
    dbx sync                              # dry-run
    dbx sync --commit --prune             # full sync
    dbx sync --table=auth.users --commit  # single table
    dbx sync --commit --fail-on-drift     # CI mode

 MIGRATIONS
 ─────────────────────────────────────────────────────────────────────────────
  migrate             Run proto-driven SQL migrations
    dbx migrate --target=primary
    dbx migrate --target=both --strict --prune

 PRINT / INSPECT    (TUI: guided dropdowns with live schema/table lists)
 ─────────────────────────────────────────────────────────────────────────────
  print-schema        Schema overview
    dbx print-schema --schema=auth --target=primary
    TUI → schema picker (live) → target

  print-tables        All tables in a schema
    dbx print-tables --schema=public --target=primary
    TUI → schema picker → target

  print-table         Single table: columns, FKs, indexes, size
    dbx print-table --table=auth.users --target=primary
    TUI → schema → table (auto-filtered) → target
         produces: --table=auth.users --target=primary

  print-all           Full DB overview (all schemas + tables)
    dbx print-all --target=primary
    TUI → target picker

  print-table-data    View actual table rows
    dbx print-table-data --table=auth.users --limit=20
    TUI → schema → table → limit → target

 SQL
 ─────────────────────────────────────────────────────────────────────────────
  sql                 Execute arbitrary SQL
    dbx sql --sql="SELECT COUNT(*) FROM auth.users" --target=primary
    dbx sql --sql="VACUUM ANALYZE;" --target=both
    TUI → target picker → type query

 MAINTENANCE  (CLI only — not in TUI)
 ─────────────────────────────────────────────────────────────────────────────
  migrate             dbx migrate --target=primary [--strict] [--prune]
  sync                dbx sync --commit --prune [--fail-on-drift]
  sync-repair         dbx sync-repair
  sync-users          dbx sync-users
  failover            dbx failover
  rebuild-backup      dbx rebuild-backup
  csv-backup          dbx csv-backup [--table=X] --source=primary
  csv-seed            dbx csv-seed [--table=X] --target=primary

 TUI NAVIGATION
 ─────────────────────────────────────────────────────────────────────────────
  /           Open command menu      Esc     Back / cancel
  ↑ ↓         Navigate list          j k     Scroll pager
  Enter       Select / confirm       PgUp↓   Page scroll
  help        Show this reference    Ctrl+C  Quit
  exit/quit   Exit application
═══════════════════════════════════════════════════════════════════════════════
`
}

// ── Helper methods for forms ──────────────────────────────────────────────

// CommandForms defines all forms statically — no closures over model pointers.
// Options are resolved at runtime in resolveOptions() using the live model state.
func allCommandForms() []CommandForm {
	targetOptions := []string{"primary", "backup", "both"}
	return []CommandForm{
		{
			Command: "print-schema",
			Fields: []FormField{
				{Name: "schema", Label: "Select schema:", Type: FieldTypeSelect, Optional: true, Default: "public"},
				{Name: "target", Label: "Select target database:", Type: FieldTypeSelect, Options: targetOptions, Default: "primary"},
			},
		},
		{
			Command: "print-tables",
			Fields: []FormField{
				{Name: "schema", Label: "Select schema:", Type: FieldTypeSelect, Optional: true, Default: "public"},
				{Name: "target", Label: "Select target database:", Type: FieldTypeSelect, Options: targetOptions, Default: "primary"},
			},
		},
		{
			Command: "print-table",
			Fields: []FormField{
				{Name: "schema", Label: "Select schema:", Type: FieldTypeSelect, Optional: true, Default: "public"},
				{Name: "table", Label: "Select table:", Type: FieldTypeSelect},
				{Name: "target", Label: "Select target database:", Type: FieldTypeSelect, Options: targetOptions, Default: "primary"},
			},
		},
		{
			Command: "print-table-data",
			Fields: []FormField{
				{Name: "schema", Label: "Select schema:", Type: FieldTypeSelect, Optional: true, Default: "public"},
				{Name: "table", Label: "Select table:", Type: FieldTypeSelect},
				{Name: "limit", Label: "Enter row limit:", Type: FieldTypeText, Optional: true, Default: "10"},
				{Name: "target", Label: "Select target database:", Type: FieldTypeSelect, Options: targetOptions, Default: "primary"},
			},
		},
		{
			Command: "print-all",
			Fields: []FormField{
				{Name: "target", Label: "Select target database:", Type: FieldTypeSelect, Options: targetOptions, Default: "primary"},
			},
		},
		{
			Command: "sql",
			Fields: []FormField{
				{Name: "target", Label: "Select target database:", Type: FieldTypeSelect, Options: targetOptions, Default: "primary"},
				{Name: "sql", Label: "Enter SQL query:", Type: FieldTypeText},
			},
		},
	}
}

func (m *Model) rebuildForms() {
	m.forms = allCommandForms()
}

func (m *Model) findFormForCommand(cmd string) *CommandForm {
	for i := range m.forms {
		if m.forms[i].Command == cmd {
			f := m.forms[i] // copy
			return &f
		}
	}
	return nil
}

// resolveOptions returns the option list for the field at formStep,
// using live model state (cachedSchemas, cachedTables, formAnswers).
// This is called synchronously in updateFormStep — no closures needed.
func (m Model) resolveOptions(field FormField) ([]string, error) {
	switch field.Name {
	case "schema":
		if len(m.cachedSchemas) == 0 {
			return []string{"public"}, nil
		}
		return m.cachedSchemas, nil

	case "table":
		// Get schema from previously answered step
		schema := m.formAnswers["schema"]
		if schema == "" {
			schema = "public"
		}
		tables := m.cachedTables[schema]
		if len(tables) == 0 && m.handlers.GetTables != nil {
			// Live fetch if not cached (blocking but in a tea.Cmd goroutine)
			fetched, err := m.handlers.GetTables(schema)
			if err != nil {
				return nil, fmt.Errorf("could not load tables for schema %q: %v", schema, err)
			}
			return fetched, nil
		}
		if len(tables) == 0 {
			return nil, fmt.Errorf("no tables found for schema %q", schema)
		}
		return tables, nil

	default:
		// Static options already on the field
		return field.Options, nil
	}
}

func (m Model) startForm(form *CommandForm) (tea.Model, tea.Cmd) {
	m.currentForm = form
	m.formAnswers = make(map[string]string)
	m.formStep = 0
	m.tui = tuiForm
	return m.updateFormStep()
}

func (m Model) updateFormStep() (tea.Model, tea.Cmd) {
	if m.formStep >= len(m.currentForm.Fields) {
		return m.executeForm()
	}

	field := m.currentForm.Fields[m.formStep]

	if field.Type == FieldTypeSelect {
		// Deep-copy formAnswers and cachedTables (maps are reference types)
		// so the goroutine snapshot is fully independent of future mutations.
		answersCopy := make(map[string]string, len(m.formAnswers))
		for k, v := range m.formAnswers {
			answersCopy[k] = v
		}
		tablesCopy := make(map[string][]string, len(m.cachedTables))
		for k, v := range m.cachedTables {
			tablesCopy[k] = v
		}
		snap := m
		snap.formAnswers = answersCopy
		snap.cachedTables = tablesCopy
		fieldSnap := field
		getTablesFn := m.handlers.GetTables
		snap.handlers.GetTables = getTablesFn // keep handler ref
		m.formLoading = true
		return m, func() tea.Msg {
			options, err := snap.resolveOptions(fieldSnap)
			return loadOptionsMsg{options: options, err: err}
		}
	}

	// FieldTypeText — show input immediately
	ti := textinput.New()
	ti.Placeholder = "Type value, Enter to confirm"
	if field.Optional && field.Default != "" {
		ti.Placeholder = fmt.Sprintf("Enter to use default: %s", field.Default)
	}
	ti.Focus()
	ti.CharLimit = 512
	ti.Width = 60
	m.formInput = ti
	return m, textinput.Blink
}

func (m Model) executeForm() (tea.Model, tea.Cmd) {
	command := m.currentForm.Command
	answers := m.formAnswers

	var cmdStr string
	switch command {

	case "print-table":
		// table is stored as plain name; prepend schema to form schema.table
		schema := answers["schema"]
		table := answers["table"]
		if schema == "" {
			schema = "public"
		}
		qualifiedTable := schema + "." + table
		target := answers["target"]
		if target == "" {
			target = "primary"
		}
		cmdStr = fmt.Sprintf("print-table --table=%s --target=%s", qualifiedTable, target)

	case "print-table-data":
		schema := answers["schema"]
		table := answers["table"]
		if schema == "" {
			schema = "public"
		}
		qualifiedTable := schema + "." + table
		limit := answers["limit"]
		if limit == "" {
			limit = "10"
		}
		target := answers["target"]
		if target == "" {
			target = "primary"
		}
		cmdStr = fmt.Sprintf("print-table-data --table=%s --limit=%s --target=%s", qualifiedTable, limit, target)

	case "print-schema":
		schema := answers["schema"]
		if schema == "" {
			schema = "public"
		}
		target := answers["target"]
		if target == "" {
			target = "primary"
		}
		cmdStr = fmt.Sprintf("print-schema --schema=%s --target=%s", schema, target)

	case "print-tables":
		schema := answers["schema"]
		if schema == "" {
			schema = "public"
		}
		target := answers["target"]
		if target == "" {
			target = "primary"
		}
		cmdStr = fmt.Sprintf("print-tables --schema=%s --target=%s", schema, target)

	case "print-all":
		target := answers["target"]
		if target == "" {
			target = "primary"
		}
		cmdStr = fmt.Sprintf("print-all --target=%s", target)

	case "sql":
		target := answers["target"]
		if target == "" {
			target = "primary"
		}
		query := answers["sql"]
		cmdStr = fmt.Sprintf("sql --sql=%s --target=%s", query, target)

	default:
		// Generic: append all non-empty answers as flags
		cmdStr = command
		for _, field := range m.currentForm.Fields {
			if val := answers[field.Name]; val != "" {
				cmdStr += " --" + field.Name + "=" + val
			}
		}
	}

	m.selectedCmd = cmdStr
	m.currentForm = nil
	m.formAnswers = make(map[string]string)
	m.formStep = 0
	return m.startExecute(cmdStr)
}

func (m Model) renderForm(w, h int, title, dbBadge string) string {
	if m.currentForm == nil || m.formStep >= len(m.currentForm.Fields) {
		return ""
	}

	field := m.currentForm.Fields[m.formStep]
	stepIndicator := fmt.Sprintf("Step %d/%d: %s", m.formStep+1, len(m.currentForm.Fields), field.Label)

	var content string
	if field.Type == FieldTypeSelect {
		if m.formLoading {
			content = m.spinner.View() + " Loading options..."
		} else {
			m.formList.SetSize(w-4, h-12)
			content = m.formList.View()
		}
	} else {
		content = inputStyle.Render(m.formInput.View())
	}

	help := helpStyle.Render("↑/↓ navigate  •  Enter confirm  •  Esc back  •  Ctrl+C quit")
	body := lipgloss.JoinVertical(lipgloss.Left,
		title, "",
		dbBadge, "",
		stepIndicator, "",
		content, "",
		help,
	)
	return lipgloss.Place(w, h, lipgloss.Left, lipgloss.Top, body)
}

func (m Model) loadCache() tea.Cmd {
	h := m.handlers
	return func() tea.Msg {
		if h.GetSchemas == nil || h.GetTables == nil {
			return cacheLoadedMsg{err: fmt.Errorf("schema/table handlers not configured")}
		}

		schemas, err := h.GetSchemas()
		if err != nil {
			return cacheLoadedMsg{err: fmt.Errorf("failed to load schemas: %v", err)}
		}

		tables := make(map[string][]string, len(schemas))
		for _, schema := range schemas {
			t, err := h.GetTables(schema)
			if err != nil {
				continue // non-fatal: skip schema
			}
			tables[schema] = t
		}

		return cacheLoadedMsg{schemas: schemas, tables: tables}
	}
}

// ── Entry point ───────────────────────────────────────────────────────────────

// RunInteractiveTUI starts the interactive TUI.
func RunInteractiveTUI(handlers Handlers) error {
	p := tea.NewProgram(
		initialModel(handlers),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	_, err := p.Run()
	return err
}

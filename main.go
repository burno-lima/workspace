package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	iconFolder     = "\uf07b"
	iconSearch     = "\uf002"
	iconError      = "\uf00d"
	iconConfig     = "\uf013"
	iconTarget     = "\uf140"
	iconTextEditor = "\uf0f6"
	iconIDE        = "\ue7a5"
	iconArrow      = "\uf054"
	iconArrowRight = "\u25b6"
	iconProject    = "\uf413"
	iconGo         = "\ue626"
	iconPython     = "\ue73c"
	iconJavaScript = "\ue74e"
	iconJava       = "\ue738"
	iconPHP        = "\ue73d"
	iconRuby       = "\ue791"
	iconRust       = "\ue7a8"
	iconCPlusPlus  = "\ue61d"

	// Colors
	primaryColor  = lipgloss.Color("#7D56F4")
	selectedColor = lipgloss.Color("#F780E2")
	normalColor   = lipgloss.Color("#FAFAFA")
	dimColor      = lipgloss.Color("#626262")
	borderColor   = lipgloss.Color("#383838")
	accentColor   = lipgloss.Color("#00D9FF")
	successColor  = lipgloss.Color("#04B575")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(selectedColor).
				Background(lipgloss.Color("#2a2a2a")).
				Bold(true).
				Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(normalColor).
			Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
			Foreground(borderColor).
			Bold(true)

	pathStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(dimColor).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(borderColor).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(dimColor)

	searchStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1)

	searchInputStyle = lipgloss.NewStyle().
				Foreground(normalColor)
)

type project struct {
	name        string
	path        string
	projectType string
}

type Editor struct {
	name    string
	command string
	display string
}

type Config struct {
	DefaultEditor  string
	ProjectEditors map[string]string
}

type ProjectType struct {
	name     string
	ide      string
	patterns []string
}

type model struct {
	projects             []project
	filteredProjects     []project
	cursor               int
	basePath             string
	err                  error
	quitting             bool
	width                int
	height               int
	searchMode           bool
	searchQuery          string
	lastOpened           string
	config               Config
	availableEditors     []Editor
	showConfigMenu       bool
	showEditorSelector   bool
	showProjectSelector  bool
	editorCursor         int
	configMenuCursor     int
	textEditors          []Editor
	jetbrainsEditors     []Editor
	pendingProjectName   string
	pendingProjectPath   string
	pendingProjectType   string
	projectEditorOptions []Editor
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	configDir := filepath.Join(home, ".config", "laucher-path")
	// Create config directory if it doesn't exist
	os.MkdirAll(configDir, 0755)
	return filepath.Join(configDir, "laucher.config")
}

func loadConfig() (Config, error) {
	configPath := getConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{ProjectEditors: make(map[string]string)}, err
	}

	config := Config{ProjectEditors: make(map[string]string)}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key == "default" {
				config.DefaultEditor = value
			} else {
				config.ProjectEditors[key] = value
			}
		}
	}
	return config, nil
}

func saveConfig(config Config) error {
	configPath := getConfigPath()
	var lines []string
	lines = append(lines, "# Launcher Configuration")
	lines = append(lines, "# Format: project-name=editor or default=editor")
	lines = append(lines, "")
	if config.DefaultEditor != "" {
		lines = append(lines, fmt.Sprintf("default=%s", config.DefaultEditor))
		lines = append(lines, "")
	}
	for project, editor := range config.ProjectEditors {
		lines = append(lines, fmt.Sprintf("%s=%s", project, editor))
	}
	data := strings.Join(lines, "\n") + "\n"
	return os.WriteFile(configPath, []byte(data), 0644)
}

func resetConfig() error {
	configPath := getConfigPath()
	// Remove the config file if it exists
	if _, err := os.Stat(configPath); err == nil {
		return os.Remove(configPath)
	}
	return nil
}

func detectAvailableEditors() []Editor {
	editors := []Editor{
		{name: "vscode", command: "code", display: "Visual Studio Code"},
		{name: "cursor", command: "cursor", display: "Cursor"},
		{name: "windsurf", command: "windsurf", display: "Windsurf"},
		{name: "neovim", command: "nvim", display: "Neovim"},
		{name: "vim", command: "vim", display: "Vim"},
		{name: "goland", command: "goland", display: "GoLand"},
		{name: "pycharm", command: "pycharm", display: "PyCharm"},
		{name: "webstorm", command: "webstorm", display: "WebStorm"},
		{name: "intellij", command: "idea", display: "IntelliJ IDEA"},
		{name: "phpstorm", command: "phpstorm", display: "PhpStorm"},
		{name: "rubymine", command: "rubymine", display: "RubyMine"},
		{name: "clion", command: "clion", display: "CLion"},
		{name: "rustrover", command: "rustrover", display: "RustRover"},
		{name: "sublime", command: "subl", display: "Sublime Text"},
		{name: "zed", command: "zed", display: "Zed"},
		{name: "codium", command: "codium", display: "VSCodium"},
		{name: "antigravity", command: "antigravity", display: "Antigravity"},
	}

	var available []Editor
	for _, editor := range editors {
		if isCommandAvailable(editor.command) {
			available = append(available, editor)
		}
	}

	if len(available) == 0 {
		if runtime.GOOS == "darwin" {
			available = append(available, Editor{name: "default", command: "open", display: "Default macOS Editor"})
		} else if runtime.GOOS == "windows" {
			available = append(available, Editor{name: "default", command: "start", display: "Default Windows Editor"})
		} else {
			available = append(available, Editor{name: "default", command: "xdg-open", display: "Default System Editor"})
		}
	}

	return available
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func findJetBrainsIDE(ideName string) string {
	// Simply check if the IDE command is available in PATH
	// This works regardless of installation method (Toolbox, snap, manual, etc.)
	if path, err := exec.LookPath(ideName); err == nil {
		return path
	}
	return ""
}

func getEditorOptionsForProject(projectType string, availableEditors []Editor) []Editor {
	var options []Editor

	// Add general purpose editors first (text editors and IDEs)
	generalEditors := map[string]bool{
		"vscode":      true,
		"cursor":      true,
		"neovim":      true,
		"vim":         true,
		"windsurf":    true,
		"sublime":     true,
		"zed":         true,
		"codium":      true,
		"antigravity": true,
	}

	for _, editor := range availableEditors {
		if generalEditors[editor.name] {
			options = append(options, editor)
		}
	}

	// Add project-specific IDE if available
	var specificIDE string
	switch projectType {
	case "Go":
		specificIDE = "goland"
	case "Python":
		specificIDE = "pycharm"
	case "JavaScript/TypeScript":
		specificIDE = "webstorm"
	case "Java":
		specificIDE = "idea"
	case "PHP":
		specificIDE = "phpstorm"
	case "Ruby":
		specificIDE = "rubymine"
	case "Rust":
		specificIDE = "rustrover"
	case "C/C++":
		specificIDE = "clion"
	}

	if specificIDE != "" {
		for _, editor := range availableEditors {
			if editor.name == specificIDE {
				options = append(options, editor)
				break
			}
		}
	}

	// If no options, return all available editors
	if len(options) == 0 {
		return availableEditors
	}

	return options
}

func getProjectTypes() []ProjectType {
	return []ProjectType{
		{name: "Go", ide: "goland", patterns: []string{"go.mod", "go.sum", "go.work", "go.work.sum"}},
		{name: "Python", ide: "pycharm", patterns: []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"}},
		{name: "JavaScript/TypeScript", ide: "webstorm", patterns: []string{"package.json", "tsconfig.json"}},
		{name: "Java", ide: "idea", patterns: []string{"pom.xml", "build.gradle", "build.gradle.kts"}},
		{name: "PHP", ide: "phpstorm", patterns: []string{"composer.json"}},
		{name: "Ruby", ide: "rubymine", patterns: []string{"Gemfile", "Rakefile"}},
		{name: "Rust", ide: "rustrover", patterns: []string{"Cargo.toml"}},
		{name: "C/C++", ide: "clion", patterns: []string{"CMakeLists.txt", "Makefile"}},
	}
}

func getIconForProjectType(projectType string) string {
	switch projectType {
	case "Go":
		return iconGo
	case "Python":
		return iconPython
	case "JavaScript/TypeScript":
		return iconJavaScript
	case "Java":
		return iconJava
	case "PHP":
		return iconPHP
	case "Ruby":
		return iconRuby
	case "Rust":
		return iconRust
	case "C/C++":
		return iconCPlusPlus
	default:
		return iconProject
	}
}

func detectProjectType(projectPath string) string {
	projectTypes := getProjectTypes()

	// First check the root directory
	for _, pt := range projectTypes {
		for _, pattern := range pt.patterns {
			filePath := filepath.Join(projectPath, pattern)
			if _, err := os.Stat(filePath); err == nil {
				return pt.name
			}
		}
	}

	// If not found in root, search in subdirectories (up to 3 levels deep)
	maxDepth := 3
	foundType := "Unknown"

	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		// Calculate current depth
		relPath, _ := filepath.Rel(projectPath, path)
		depth := len(strings.Split(relPath, string(os.PathSeparator)))

		// Stop if we've gone too deep
		if depth > maxDepth {
			return filepath.SkipDir
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && path != projectPath {
			return filepath.SkipDir
		}

		// Skip common directories that won't have project files
		if info.IsDir() {
			skipDirs := map[string]bool{
				"node_modules": true,
				"vendor":       true,
				"venv":         true,
				".venv":        true,
				"__pycache__":  true,
				"build":        true,
				"dist":         true,
				"target":       true,
			}
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
		}

		// Check if this file matches any project type pattern
		if !info.IsDir() {
			fileName := info.Name()
			for _, pt := range projectTypes {
				for _, pattern := range pt.patterns {
					if fileName == pattern && foundType == "Unknown" {
						foundType = pt.name
						return filepath.SkipAll // Stop searching once found
					}
				}
			}
		}

		return nil
	})

	return foundType
}

func getIDEForProjectType(projectType string) string {
	projectTypes := getProjectTypes()

	for _, pt := range projectTypes {
		if pt.name == projectType {
			// Use the enhanced detection for JetBrains IDEs
			idePath := findJetBrainsIDE(pt.ide)
			if idePath != "" {
				return idePath
			}
			break
		}
	}

	return ""
}

func initialModel(basePath string) model {
	projects, err := loadProjects(basePath)
	config, _ := loadConfig()
	availableEditors := detectAvailableEditors()

	// Separate editors into categories
	var textEditors, jetbrainsEditors []Editor
	for _, editor := range availableEditors {
		if editor.name == "goland" || editor.name == "pycharm" || editor.name == "webstorm" ||
			editor.name == "intellij" || editor.name == "phpstorm" || editor.name == "rubymine" ||
			editor.name == "clion" || editor.name == "rustrover" {
			jetbrainsEditors = append(jetbrainsEditors, editor)
		} else {
			textEditors = append(textEditors, editor)
		}
	}

	return model{
		projects:             projects,
		filteredProjects:     projects,
		cursor:               0,
		basePath:             basePath,
		err:                  err,
		width:                80,
		height:               24,
		searchMode:           false,
		searchQuery:          "",
		lastOpened:           "",
		config:               config,
		availableEditors:     availableEditors,
		textEditors:          textEditors,
		jetbrainsEditors:     jetbrainsEditors,
		showConfigMenu:       false,
		showEditorSelector:   false,
		showProjectSelector:  false,
		editorCursor:         0,
		configMenuCursor:     0,
		projectEditorOptions: []Editor{},
	}
}

func loadProjects(basePath string) ([]project, error) {
	var projects []project

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			projectPath := filepath.Join(basePath, entry.Name())
			projects = append(projects, project{
				name:        entry.Name(),
				path:        projectPath,
				projectType: detectProjectType(projectPath),
			})
		}
	}

	return projects, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case projectOpenedMsg:
		// Don't keep the lastOpened message
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		// Project editor selector modal handling
		if m.showProjectSelector {
			switch msg.String() {
			case "esc", "q":
				m.showProjectSelector = false
				m.editorCursor = 0
				return m, nil
			case "up", "k":
				if m.editorCursor > 0 {
					m.editorCursor--
				}
			case "down", "j":
				if m.editorCursor < len(m.projectEditorOptions)-1 {
					m.editorCursor++
				}
			case "enter", "l":
				selectedEditor := m.projectEditorOptions[m.editorCursor]
				m.config.ProjectEditors[m.pendingProjectName] = selectedEditor.command
				saveConfig(m.config)
				m.showProjectSelector = false
				// Now actually open the project
				cmd := exec.Command(selectedEditor.command, m.pendingProjectPath)
				cmd.Start()
				return m, nil
			}
			return m, nil
		}

		// Config menu modal handling
		if m.showConfigMenu {
			switch msg.String() {
			case "esc", "q":
				m.showConfigMenu = false
				m.configMenuCursor = 0
				return m, nil
			case "up", "k":
				if m.configMenuCursor > 0 {
					m.configMenuCursor--
				}
			case "down", "j":
				if m.configMenuCursor < 1 {
					m.configMenuCursor++
				}
			case "enter", "l":
				if m.configMenuCursor == 0 {
					// Open editor selector
					m.showConfigMenu = false
					m.showEditorSelector = true
					m.editorCursor = 0
					// Populate with all available editors for default editor selection
					m.projectEditorOptions = m.availableEditors
					// Clear pending project info since we're setting default
					m.pendingProjectName = ""
					m.pendingProjectType = ""
					m.configMenuCursor = 0
				} else if m.configMenuCursor == 1 {
					// Reset configuration
					resetConfig()
					// Reload config
					m.config = Config{ProjectEditors: make(map[string]string)}
					m.showConfigMenu = false
					m.configMenuCursor = 0
				}
				return m, nil
			}
			return m, nil
		}

		// Editor selector with categories handling
		if m.showEditorSelector {
			switch msg.String() {
			case "esc", "q", "h":
				m.showEditorSelector = false
				m.editorCursor = 0
				return m, nil
			case "up", "k":
				if m.editorCursor > 0 {
					m.editorCursor--
				}
			case "down", "j":
				if m.editorCursor < len(m.projectEditorOptions)-1 {
					m.editorCursor++
				}
			case "enter", "l":
				selectedEditor := m.projectEditorOptions[m.editorCursor]
				m.config.DefaultEditor = selectedEditor.command
				saveConfig(m.config)
				m.showEditorSelector = false
				m.editorCursor = 0
			}
			return m, nil
		}

		// Search mode handling
		if m.searchMode {
			switch msg.String() {
			case "esc":
				m.searchMode = false
				m.searchQuery = ""
				m.filteredProjects = m.projects
				m.cursor = 0
				return m, nil
			case "enter":
				if len(m.filteredProjects) > 0 {
					m.searchMode = false
					m.searchQuery = ""

					// Check if we need to show modal or open directly
					projectName := m.filteredProjects[m.cursor].name
					projectType := m.filteredProjects[m.cursor].projectType
					projectPath := m.filteredProjects[m.cursor].path

					// Reset filter after opening to show all projects again
					m.filteredProjects = m.projects
					m.cursor = 0

					// Check if project already has a configured editor
					if editorCmd, exists := m.config.ProjectEditors[projectName]; exists {
						cmd := exec.Command(editorCmd, projectPath)
						cmd.Start()
						return m, nil
					}

					// Check if there's a default editor configured
					if m.config.DefaultEditor != "" {
						cmd := exec.Command(m.config.DefaultEditor, projectPath)
						cmd.Start()
						return m, nil
					}

					// No configuration found - show selector modal
					m.pendingProjectName = projectName
					m.pendingProjectPath = projectPath
					m.pendingProjectType = projectType
					m.projectEditorOptions = getEditorOptionsForProject(projectType, m.availableEditors)
					m.showProjectSelector = true
					m.editorCursor = 0
					return m, nil
				}
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.filterProjects()
					if m.cursor >= len(m.filteredProjects) {
						m.cursor = len(m.filteredProjects) - 1
						if m.cursor < 0 {
							m.cursor = 0
						}
					}
				}
			case "up", "ctrl+p":
				cols := m.calculateColumns()
				if m.cursor >= cols {
					m.cursor -= cols
				}
			case "down", "ctrl+n":
				cols := m.calculateColumns()
				if m.cursor+cols < len(m.filteredProjects) {
					m.cursor += cols
				}
			case "left":
				if m.cursor > 0 {
					m.cursor--
				}
			case "right":
				if m.cursor < len(m.filteredProjects)-1 {
					m.cursor++
				}
			default:
				// Add character to search (including h, j, k, l)
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
					m.filterProjects()
					m.cursor = 0
				}
			}
			return m, nil
		}

		// Normal mode handling
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "/":
			m.searchMode = true
			m.searchQuery = ""
			return m, nil

		case "c":
			// Reset search mode when opening config
			m.searchMode = false
			m.searchQuery = ""
			m.showConfigMenu = true
			return m, nil

		case "up", "k":
			cols := m.calculateColumns()
			if m.cursor >= cols {
				m.cursor -= cols
			}

		case "down", "j":
			cols := m.calculateColumns()
			if m.cursor+cols < len(m.filteredProjects) {
				m.cursor += cols
			}

		case "left", "h":
			// Move to previous item (left)
			if m.cursor > 0 {
				m.cursor--
			}

		case "right", "l":
			// Move to next item (right)
			if m.cursor < len(m.filteredProjects)-1 {
				m.cursor++
			}

		case "enter":
			if len(m.filteredProjects) > 0 {
				// Check if we need to show modal or open directly
				projectName := m.filteredProjects[m.cursor].name
				projectType := m.filteredProjects[m.cursor].projectType
				projectPath := m.filteredProjects[m.cursor].path

				// Check if project already has a configured editor
				if editorCmd, exists := m.config.ProjectEditors[projectName]; exists {
					cmd := exec.Command(editorCmd, projectPath)
					cmd.Start()
					return m, nil
				}

				// Check if there's a default editor configured
				if m.config.DefaultEditor != "" {
					cmd := exec.Command(m.config.DefaultEditor, projectPath)
					cmd.Start()
					return m, nil
				}

				// No configuration found - show selector modal
				m.pendingProjectName = projectName
				m.pendingProjectPath = projectPath
				m.pendingProjectType = projectType
				m.projectEditorOptions = getEditorOptionsForProject(projectType, m.availableEditors)
				m.showProjectSelector = true
				m.editorCursor = 0
				return m, nil
			}
		}
	}

	return m, nil
}

func (m *model) filterProjects() {
	if m.searchQuery == "" {
		m.filteredProjects = m.projects
		return
	}

	m.filteredProjects = []project{}
	query := strings.ToLower(m.searchQuery)
	for _, proj := range m.projects {
		if strings.Contains(strings.ToLower(proj.name), query) {
			m.filteredProjects = append(m.filteredProjects, proj)
		}
	}
}

func (m model) calculateColumns() int {
	if len(m.filteredProjects) == 0 {
		return 1
	}

	// Find the longest project name
	maxLen := 0
	for _, proj := range m.filteredProjects {
		if len(proj.name) > maxLen {
			maxLen = len(proj.name)
		}
	}

	// Calculate column width (name + padding + selection indicator)
	columnWidth := maxLen + 8
	availableWidth := m.width - 10 // Account for box borders and padding

	if availableWidth < columnWidth {
		return 1
	}

	cols := availableWidth / columnWidth
	if cols < 1 {
		cols = 1
	}
	if cols > 5 {
		cols = 5 // Max 5 columns
	}

	return cols
}

type projectOpenedMsg struct {
	name string
}

type errMsg struct{ err error }

func (m model) renderPathBox() string {
	boxWidth := m.width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}
	pathContent := pathStyle.Render(fmt.Sprintf("%s %s", iconFolder, m.basePath))

	if m.searchMode {
		return boxStyle.
			Width(boxWidth).
			BorderForeground(accentColor).
			Render(pathContent)
	}

	return boxStyle.
		Width(boxWidth).
		BorderForeground(borderColor).
		Render(pathContent)
}

func (m model) renderSearchBox() string {
	boxWidth := m.width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}

	var searchContent string
	if m.searchMode {
		searchContent = infoStyle.Render(fmt.Sprintf("%s Search: ", iconSearch)) +
			searchInputStyle.Render(m.searchQuery) +
			searchInputStyle.Render("█")
		return boxStyle.
			Width(boxWidth).
			BorderForeground(accentColor).
			Render(searchContent)
	} else {
		searchContent = dimStyle.Render(fmt.Sprintf("%s Search: ", iconSearch)) +
			dimStyle.Render("(press / to search)")
		return boxStyle.
			Width(boxWidth).
			BorderForeground(borderColor).
			Render(searchContent)
	}
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	if m.err != nil {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			boxStyle.Render(fmt.Sprintf("%s Error: %v", iconError, m.err)),
		)
	}

	// If any modal is active, show it
	if m.showProjectSelector {
		return m.renderProjectEditorModal()
	}

	if m.showConfigMenu {
		return m.renderConfigMenu()
	}

	if m.showEditorSelector {
		return m.renderEditorSelector()
	}

	// Otherwise show main view
	return m.renderMainView()
}

func (m model) renderMainView() string {
	boxWidth := m.width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}

	// Render path box (fixed at top)
	pathBox := m.renderPathBox()

	// Render search box (fixed below path)
	searchBox := m.renderSearchBox()

	// Render projects area
	var projectsContent string
	if len(m.filteredProjects) == 0 {
		noResultsMsg := "No projects found"
		if m.searchQuery != "" {
			noResultsMsg = fmt.Sprintf("No projects matching '%s'", m.searchQuery)
		}
		projectsContent = boxStyle.Width(boxWidth).Render(dimStyle.Render(noResultsMsg))
	} else {
		projectsGrid := m.renderProjectsGrid()
		projectsContent = boxStyle.Width(boxWidth).Render(projectsGrid)
	}

	// Render info box for selected project
	var infoBox string
	if len(m.filteredProjects) > 0 {
		selectedProj := m.filteredProjects[m.cursor]
		projectIcon := getIconForProjectType(selectedProj.projectType)

		// Determine which editor will be used
		var editorDisplay string
		if configuredEditor, exists := m.config.ProjectEditors[selectedProj.name]; exists {
			for _, editor := range m.availableEditors {
				if editor.command == configuredEditor {
					editorDisplay = fmt.Sprintf("%s (project specific)", editor.display)
					break
				}
			}
			if editorDisplay == "" {
				editorDisplay = fmt.Sprintf("%s (project specific)", configuredEditor)
			}
		} else if m.config.DefaultEditor != "" {
			for _, editor := range m.availableEditors {
				if editor.command == m.config.DefaultEditor {
					editorDisplay = fmt.Sprintf("%s (default)", editor.display)
					break
				}
			}
			if editorDisplay == "" {
				editorDisplay = fmt.Sprintf("%s (default)", m.config.DefaultEditor)
			}
		} else {
			editorDisplay = "Not configured"
		}

		line1 := infoStyle.Render("Selected: ") + normalItemStyle.Render(fmt.Sprintf("%s %s", projectIcon, selectedProj.name))
		line2 := pathStyle.Render(selectedProj.path)
		line3 := dimStyle.Render(fmt.Sprintf("%s ", iconArrow)) + dimStyle.Render(editorDisplay)

		infoContent := lipgloss.JoinVertical(lipgloss.Left, line1, line2, line3)
		infoBox = boxStyle.Width(boxWidth).Render(infoContent)
	}

	// Render help box (fixed at bottom)
	var helpText string
	if m.searchMode {
		helpText = lipgloss.JoinHorizontal(
			lipgloss.Left,
			dimStyle.Render("↑↓←→/ctrl+n,p: navigate"),
			dimStyle.Render(" • "),
			dimStyle.Render("type to search"),
			dimStyle.Render(" • "),
			statusStyle.Render("enter: open"),
			dimStyle.Render(" • "),
			dimStyle.Render("esc: cancel"),
		)
	} else {
		helpText = lipgloss.JoinHorizontal(
			lipgloss.Left,
			dimStyle.Render("↑↓←→/hjkl: navigate"),
			dimStyle.Render(" • "),
			infoStyle.Render("/: search"),
			dimStyle.Render(" • "),
			infoStyle.Render("c: config"),
			dimStyle.Render(" • "),
			statusStyle.Render("enter: open"),
			dimStyle.Render(" • "),
			dimStyle.Render("q: quit"),
		)
	}
	helpBox := boxStyle.Width(boxWidth).Render(helpText)

	// Build layout
	var sections []string
	sections = append(sections, pathBox)
	sections = append(sections, searchBox)
	sections = append(sections, projectsContent)
	if infoBox != "" {
		sections = append(sections, infoBox)
	}
	sections = append(sections, helpBox)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m model) renderProjectEditorModal() string {
	var modalContent []string

	modalContent = append(modalContent, titleStyle.Render(fmt.Sprintf("%s Select Editor for %s", iconTarget, m.pendingProjectName)))
	if m.pendingProjectType != "Unknown" {
		modalContent = append(modalContent, "\n"+dimStyle.Render("Project type: "+m.pendingProjectType))
	}
	modalContent = append(modalContent, "\n"+dimStyle.Render("Choose your preferred editor for this project:"))
	modalContent = append(modalContent, "")

	for i, editor := range m.projectEditorOptions {
		if i == m.editorCursor {
			modalContent = append(modalContent, selectedItemStyle.Render(fmt.Sprintf("%s %s", iconArrowRight, editor.display)))
		} else {
			modalContent = append(modalContent, normalItemStyle.Render("  "+editor.display))
		}
	}

	modalContent = append(modalContent, "")
	modalContent = append(modalContent, dimStyle.Render("↑/↓: navigate • enter: select • esc: cancel"))

	content := lipgloss.JoinVertical(lipgloss.Left, modalContent...)

	// Responsive width: 80% of terminal width, min 40, max 70
	modalWidth := int(float64(m.width) * 0.8)
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 70 {
		modalWidth = 70
	}

	modal := boxStyle.
		Width(modalWidth).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func (m model) renderConfigMenu() string {
	var modalContent []string

	modalContent = append(modalContent, titleStyle.Render(fmt.Sprintf("%s Configuration", iconConfig)))
	modalContent = append(modalContent, "\n"+dimStyle.Render("Select an option:"))
	modalContent = append(modalContent, "")

	// Option 0: Set Default Editor
	if m.configMenuCursor == 0 {
		modalContent = append(modalContent, selectedItemStyle.Render(fmt.Sprintf("%s Set Default Editor", iconArrowRight)))
	} else {
		modalContent = append(modalContent, normalItemStyle.Render("  Set Default Editor"))
	}

	// Option 1: Reset Configuration
	if m.configMenuCursor == 1 {
		modalContent = append(modalContent, selectedItemStyle.Render(fmt.Sprintf("%s Reset Configuration", iconArrowRight)))
	} else {
		modalContent = append(modalContent, normalItemStyle.Render("  Reset Configuration"))
	}

	modalContent = append(modalContent, "")
	modalContent = append(modalContent, dimStyle.Render("↑/↓: navigate • enter: select • esc: cancel"))

	content := lipgloss.JoinVertical(lipgloss.Left, modalContent...)

	// Responsive width: 60% of terminal width, min 35, max 55
	modalWidth := int(float64(m.width) * 0.6)
	if modalWidth < 35 {
		modalWidth = 35
	}
	if modalWidth > 55 {
		modalWidth = 55
	}

	modal := boxStyle.
		Width(modalWidth).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func (m model) renderEditorSelector() string {
	var modalContent []string

	// Check if we're setting default editor or project-specific editor
	if m.pendingProjectName == "" {
		// Setting default editor - show with categories
		modalContent = append(modalContent, titleStyle.Render(fmt.Sprintf("%s Select Default Editor", iconTarget)))
		modalContent = append(modalContent, "\n"+dimStyle.Render("Choose your default editor:"))
		modalContent = append(modalContent, "")

		// Separate editors by category
		var textEditors, jetbrainsEditors []Editor
		for _, editor := range m.projectEditorOptions {
			if editor.name == "goland" || editor.name == "pycharm" || editor.name == "webstorm" ||
				editor.name == "intellij" || editor.name == "phpstorm" || editor.name == "rubymine" ||
				editor.name == "clion" || editor.name == "rustrover" {
				jetbrainsEditors = append(jetbrainsEditors, editor)
			} else {
				textEditors = append(textEditors, editor)
			}
		}

		// Render Text Editors category
		if len(textEditors) > 0 {
			modalContent = append(modalContent, infoStyle.Render(fmt.Sprintf("%s Text Editors", iconTextEditor)))
			idx := 0
			for _, editor := range textEditors {
				if idx == m.editorCursor {
					modalContent = append(modalContent, selectedItemStyle.Render(fmt.Sprintf("%s %s", iconArrowRight, editor.display)))
				} else {
					modalContent = append(modalContent, normalItemStyle.Render("  "+editor.display))
				}
				idx++
			}
			modalContent = append(modalContent, "")
		}

		// Render JetBrains IDEs category
		if len(jetbrainsEditors) > 0 {
			modalContent = append(modalContent, infoStyle.Render(fmt.Sprintf("%s JetBrains IDEs", iconIDE)))
			idx := len(textEditors)
			for _, editor := range jetbrainsEditors {
				if idx == m.editorCursor {
					modalContent = append(modalContent, selectedItemStyle.Render(fmt.Sprintf("%s %s", iconArrowRight, editor.display)))
				} else {
					modalContent = append(modalContent, normalItemStyle.Render("  "+editor.display))
				}
				idx++
			}
		}
	} else {
		// Setting project-specific editor
		modalContent = append(modalContent, titleStyle.Render(fmt.Sprintf("%s Select Editor for %s", iconTarget, m.pendingProjectName)))
		if m.pendingProjectType != "Unknown" {
			modalContent = append(modalContent, "\n"+dimStyle.Render("Project type: "+m.pendingProjectType))
		}
		modalContent = append(modalContent, "\n"+dimStyle.Render("Choose your preferred editor for this project:"))
		modalContent = append(modalContent, "")

		for i, editor := range m.projectEditorOptions {
			if i == m.editorCursor {
				modalContent = append(modalContent, selectedItemStyle.Render(fmt.Sprintf("%s %s", iconArrowRight, editor.display)))
			} else {
				modalContent = append(modalContent, normalItemStyle.Render("  "+editor.display))
			}
		}
	}

	modalContent = append(modalContent, "")
	modalContent = append(modalContent, dimStyle.Render("↑/↓: navigate • enter: select • esc: cancel"))

	content := lipgloss.JoinVertical(lipgloss.Left, modalContent...)

	// Responsive width: 80% of terminal width, min 40, max 70
	modalWidth := int(float64(m.width) * 0.8)
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 70 {
		modalWidth = 70
	}

	// Create a modal box with appropriate width
	modal := boxStyle.
		Width(modalWidth).
		Render(content)

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func (m model) renderProjectsGrid() string {
	if len(m.filteredProjects) == 0 {
		return ""
	}

	cols := m.calculateColumns()
	rows := (len(m.filteredProjects) + cols - 1) / cols

	// Find max length for column width alignment
	maxLen := 0
	for _, proj := range m.filteredProjects {
		if len(proj.name) > maxLen {
			maxLen = len(proj.name)
		}
	}
	itemWidth := maxLen + 8 // Extra space for icons and padding

	var gridRows []string
	for row := 0; row < rows; row++ {
		var rowItems []string
		for col := 0; col < cols; col++ {
			idx := row*cols + col // Row-major order
			if idx >= len(m.filteredProjects) {
				// Add empty space for alignment
				rowItems = append(rowItems, lipgloss.NewStyle().Width(itemWidth).Render(""))
				continue
			}

			proj := m.filteredProjects[idx]
			projectIcon := getIconForProjectType(proj.projectType)
			var item string
			if idx == m.cursor {
				item = selectedItemStyle.Width(itemWidth).Render(fmt.Sprintf("%s %s %s", iconArrowRight, projectIcon, proj.name))
			} else {
				item = normalItemStyle.Width(itemWidth).Render(fmt.Sprintf("  %s %s", projectIcon, proj.name))
			}
			rowItems = append(rowItems, item)
		}
		if len(rowItems) > 0 {
			gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, rowItems...))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, gridRows...)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: laucher <path>")
		fmt.Println("Example: laucher /home/brunolima/workspace/napp")
		os.Exit(1)
	}

	basePath := os.Args[1]

	// Check if path exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		fmt.Printf("Error: Path '%s' does not exist\n", basePath)
		os.Exit(1)
	}

	p := tea.NewProgram(
		initialModel(basePath),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

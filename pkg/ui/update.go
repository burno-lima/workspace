package ui

import (
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"

	"workspace/pkg/config"
	"workspace/pkg/editor"
	"workspace/pkg/models"
)

func Update(msg tea.Msg, m models.Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case models.ProjectOpenedMsg:
		return m, nil

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		if m.ShowProjectSelector {
			return handleProjectSelector(msg, m)
		}

		if m.ShowConfigMenu {
			return handleConfigMenu(msg, m)
		}

		if m.ShowEditorSelector {
			return handleEditorSelector(msg, m)
		}

		if m.SearchMode {
			return handleSearchMode(msg, m)
		}

		return handleNormalMode(msg, m)
	}

	return m, nil
}

func handleProjectSelector(msg tea.KeyMsg, m models.Model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ShowProjectSelector = false
		m.EditorCursor = 0
		return m, nil
	case "up", "k":
		if m.EditorCursor > 0 {
			m.EditorCursor--
		}
	case "down", "j":
		if m.EditorCursor < len(m.ProjectEditorOptions)-1 {
			m.EditorCursor++
		}
	case "enter", "l":
		selectedEditor := m.ProjectEditorOptions[m.EditorCursor]
		m.Config.ProjectEditors[m.PendingProjectName] = selectedEditor.Command
		config.Save(m.Config)
		m.ShowProjectSelector = false
		cmd := exec.Command(selectedEditor.Command, m.PendingProjectPath)
		cmd.Start()
		return m, nil
	}
	return m, nil
}

func handleConfigMenu(msg tea.KeyMsg, m models.Model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.ShowConfigMenu = false
		m.ConfigMenuCursor = 0
		return m, nil
	case "up", "k":
		if m.ConfigMenuCursor > 0 {
			m.ConfigMenuCursor--
		}
	case "down", "j":
		if m.ConfigMenuCursor < 1 {
			m.ConfigMenuCursor++
		}
	case "enter", "l":
		if m.ConfigMenuCursor == 0 {
			m.ShowConfigMenu = false
			m.ShowEditorSelector = true
			m.EditorCursor = 0
			m.ProjectEditorOptions = m.AvailableEditors
			m.PendingProjectName = ""
			m.PendingProjectType = ""
			m.ConfigMenuCursor = 0
		} else if m.ConfigMenuCursor == 1 {
			config.Reset()
			m.Config = models.Config{ProjectEditors: make(map[string]string)}
			m.ShowConfigMenu = false
			m.ConfigMenuCursor = 0
		}
		return m, nil
	}
	return m, nil
}

func handleEditorSelector(msg tea.KeyMsg, m models.Model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "h":
		m.ShowEditorSelector = false
		m.EditorCursor = 0
		return m, nil
	case "up", "k":
		if m.EditorCursor > 0 {
			m.EditorCursor--
		}
	case "down", "j":
		if m.EditorCursor < len(m.ProjectEditorOptions)-1 {
			m.EditorCursor++
		}
	case "enter", "l":
		selectedEditor := m.ProjectEditorOptions[m.EditorCursor]
		m.Config.DefaultEditor = selectedEditor.Command
		config.Save(m.Config)
		m.ShowEditorSelector = false
		m.EditorCursor = 0
	}
	return m, nil
}

func handleSearchMode(msg tea.KeyMsg, m models.Model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.SearchMode = false
		m.SearchQuery = ""
		m.FilteredProjects = m.Projects
		m.Cursor = 0
		return m, nil
	case "enter":
		if len(m.FilteredProjects) > 0 {
			m.SearchMode = false
			m.SearchQuery = ""

			projectName := m.FilteredProjects[m.Cursor].Name
			projectType := m.FilteredProjects[m.Cursor].ProjectType
			projectPath := m.FilteredProjects[m.Cursor].Path

			m.FilteredProjects = m.Projects
			m.Cursor = 0

			if editorCmd, exists := m.Config.ProjectEditors[projectName]; exists {
				cmd := exec.Command(editorCmd, projectPath)
				cmd.Start()
				return m, nil
			}

			if m.Config.DefaultEditor != "" {
				cmd := exec.Command(m.Config.DefaultEditor, projectPath)
				cmd.Start()
				return m, nil
			}

			m.PendingProjectName = projectName
			m.PendingProjectPath = projectPath
			m.PendingProjectType = projectType
			m.ProjectEditorOptions = editor.GetOptionsForProject(projectType, m.AvailableEditors)
			m.ShowProjectSelector = true
			m.EditorCursor = 0
			return m, nil
		}
	case "backspace":
		if len(m.SearchQuery) > 0 {
			m.SearchQuery = m.SearchQuery[:len(m.SearchQuery)-1]
			FilterProjects(&m)
			if m.Cursor >= len(m.FilteredProjects) {
				m.Cursor = len(m.FilteredProjects) - 1
				if m.Cursor < 0 {
					m.Cursor = 0
				}
			}
		}
	case "up", "ctrl+p":
		cols := CalculateColumns(m)
		if m.Cursor >= cols {
			m.Cursor -= cols
		}
	case "down", "ctrl+n":
		cols := CalculateColumns(m)
		if m.Cursor+cols < len(m.FilteredProjects) {
			m.Cursor += cols
		}
	case "left":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "right":
		if m.Cursor < len(m.FilteredProjects)-1 {
			m.Cursor++
		}
	default:
		if len(msg.String()) == 1 {
			m.SearchQuery += msg.String()
			FilterProjects(&m)
			m.Cursor = 0
		}
	}
	return m, nil
}

func handleNormalMode(msg tea.KeyMsg, m models.Model) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.Quitting = true
		return m, tea.Quit

	case "/":
		m.SearchMode = true
		m.SearchQuery = ""
		return m, nil

	case "c":
		m.SearchMode = false
		m.SearchQuery = ""
		m.ShowConfigMenu = true
		return m, nil

	case "up", "k":
		cols := CalculateColumns(m)
		if m.Cursor >= cols {
			m.Cursor -= cols
		}

	case "down", "j":
		cols := CalculateColumns(m)
		if m.Cursor+cols < len(m.FilteredProjects) {
			m.Cursor += cols
		}

	case "left", "h":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "right", "l":
		if m.Cursor < len(m.FilteredProjects)-1 {
			m.Cursor++
		}

	case "enter":
		if len(m.FilteredProjects) > 0 {
			projectName := m.FilteredProjects[m.Cursor].Name
			projectType := m.FilteredProjects[m.Cursor].ProjectType
			projectPath := m.FilteredProjects[m.Cursor].Path

			if editorCmd, exists := m.Config.ProjectEditors[projectName]; exists {
				cmd := exec.Command(editorCmd, projectPath)
				cmd.Start()
				return m, nil
			}

			if m.Config.DefaultEditor != "" {
				cmd := exec.Command(m.Config.DefaultEditor, projectPath)
				cmd.Start()
				return m, nil
			}

			m.PendingProjectName = projectName
			m.PendingProjectPath = projectPath
			m.PendingProjectType = projectType
			m.ProjectEditorOptions = editor.GetOptionsForProject(projectType, m.AvailableEditors)
			m.ShowProjectSelector = true
			m.EditorCursor = 0
			return m, nil
		}
	}

	return m, nil
}

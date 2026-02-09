package models

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Project struct {
	Name        string
	Path        string
	ProjectType string
}

type Editor struct {
	Name    string
	Command string
	Display string
}

type Config struct {
	DefaultEditor  string
	ProjectEditors map[string]string
}

type ProjectType struct {
	Name     string
	IDE      string
	Patterns []string
}

type Model struct {
	Projects             []Project
	FilteredProjects     []Project
	Cursor               int
	BasePath             string
	Err                  error
	Quitting             bool
	Width                int
	Height               int
	SearchMode           bool
	SearchQuery          string
	LastOpened           string
	Config               Config
	AvailableEditors     []Editor
	ShowConfigMenu       bool
	ShowEditorSelector   bool
	ShowProjectSelector  bool
	EditorCursor         int
	ConfigMenuCursor     int
	TextEditors          []Editor
	JetbrainsEditors     []Editor
	PendingProjectName   string
	PendingProjectPath   string
	PendingProjectType   string
	ProjectEditorOptions []Editor
}

type ProjectOpenedMsg struct {
	Name string
}

type ErrMsg struct {
	Err error
}

func (m Model) Init() tea.Cmd {
	return nil
}

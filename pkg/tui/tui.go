package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"workspace/pkg/models"
	"workspace/pkg/ui"
)

type TUIModel struct {
	Model models.Model
}

func New(basePath string) TUIModel {
	return TUIModel{
		Model: ui.InitialModel(basePath),
	}
}

func (t TUIModel) Init() tea.Cmd {
	return t.Model.Init()
}

func (t TUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	updatedModel, cmd := ui.Update(msg, t.Model)
	t.Model = updatedModel.(models.Model)
	return t, cmd
}

func (t TUIModel) View() string {
	return ui.View(t.Model)
}

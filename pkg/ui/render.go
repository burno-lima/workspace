package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"workspace/pkg/config"
	"workspace/pkg/editor"
	"workspace/pkg/models"
	"workspace/pkg/project"
	str "workspace/pkg/strings"
)

func RenderPathBox(m models.Model) string {
	boxWidth := m.Width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}
	pathContent := PathStyle.Render(fmt.Sprintf("%s %s", str.IconFolder, m.BasePath))

	if m.SearchMode {
		return BoxStyle.
			Width(boxWidth).
			BorderForeground(AccentColor).
			Render(pathContent)
	}

	return BoxStyle.
		Width(boxWidth).
		BorderForeground(BorderColor).
		Render(pathContent)
}

func RenderSearchBox(m models.Model) string {
	boxWidth := m.Width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}

	var searchContent string
	if m.SearchMode {
		searchContent = InfoStyle.Render(fmt.Sprintf("%s %s", str.IconSearch, str.SearchLabel)) +
			SearchInputStyle.Render(m.SearchQuery) +
			SearchInputStyle.Render("█")
		return BoxStyle.
			Width(boxWidth).
			BorderForeground(AccentColor).
			Render(searchContent)
	} else {
		searchContent = DimStyle.Render(fmt.Sprintf("%s %s", str.IconSearch, str.SearchLabel)) +
			DimStyle.Render(str.SearchPlaceholder.String())
		return BoxStyle.
			Width(boxWidth).
			BorderForeground(BorderColor).
			Render(searchContent)
	}
}

func RenderProjectsGrid(m models.Model) string {
	if len(m.FilteredProjects) == 0 {
		return ""
	}

	cols := CalculateColumns(m)
	rows := (len(m.FilteredProjects) + cols - 1) / cols

	maxLen := 0
	for _, proj := range m.FilteredProjects {
		if len(proj.Name) > maxLen {
			maxLen = len(proj.Name)
		}
	}
	itemWidth := maxLen + 8

	var gridRows []string
	for row := 0; row < rows; row++ {
		var rowItems []string
		for col := 0; col < cols; col++ {
			idx := row*cols + col
			if idx >= len(m.FilteredProjects) {
				rowItems = append(rowItems, lipgloss.NewStyle().Width(itemWidth).Render(""))
				continue
			}

			proj := m.FilteredProjects[idx]
			projectIcon := project.GetIconForType(proj.ProjectType)
			var item string
			if idx == m.Cursor {
				item = SelectedItemStyle.Width(itemWidth).Render(fmt.Sprintf("%s %s %s", str.IconArrowRight, projectIcon, proj.Name))
			} else {
				item = NormalItemStyle.Width(itemWidth).Render(fmt.Sprintf("  %s %s", projectIcon, proj.Name))
			}
			rowItems = append(rowItems, item)
		}
		if len(rowItems) > 0 {
			gridRows = append(gridRows, lipgloss.JoinHorizontal(lipgloss.Top, rowItems...))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, gridRows...)
}

func RenderMainView(m models.Model) string {
	boxWidth := m.Width - 4
	if boxWidth < 40 {
		boxWidth = 40
	}

	pathBox := RenderPathBox(m)
	searchBox := RenderSearchBox(m)

	var projectsContent string
	if len(m.FilteredProjects) == 0 {
		noResultsMsg := str.NoProjectsFound.String()
		if m.SearchQuery != "" {
			noResultsMsg = fmt.Sprintf(str.NoProjectsMatching.String(), m.SearchQuery)
		}
		projectsContent = BoxStyle.Width(boxWidth).Render(DimStyle.Render(noResultsMsg))
	} else {
		projectsGrid := RenderProjectsGrid(m)
		projectsContent = BoxStyle.Width(boxWidth).Render(projectsGrid)
	}

	var infoBox string
	if len(m.FilteredProjects) > 0 {
		selectedProj := m.FilteredProjects[m.Cursor]
		projectIcon := project.GetIconForType(selectedProj.ProjectType)

		var editorDisplay string
		if configuredEditor, exists := m.Config.ProjectEditors[selectedProj.Name]; exists {
			for _, ed := range m.AvailableEditors {
				if ed.Command == configuredEditor {
					editorDisplay = fmt.Sprintf(str.EditorProjectSpecific.String(), ed.Display)
					break
				}
			}
			if editorDisplay == "" {
				editorDisplay = fmt.Sprintf(str.EditorProjectSpecific.String(), configuredEditor)
			}
		} else if m.Config.DefaultEditor != "" {
			for _, ed := range m.AvailableEditors {
				if ed.Command == m.Config.DefaultEditor {
					editorDisplay = fmt.Sprintf(str.EditorDefault.String(), ed.Display)
					break
				}
			}
			if editorDisplay == "" {
				editorDisplay = fmt.Sprintf(str.EditorDefault.String(), m.Config.DefaultEditor)
			}
		} else {
			editorDisplay = str.EditorNotConfigured.String()
		}

		line1 := InfoStyle.Render(str.SelectedLabel.String()) + NormalItemStyle.Render(fmt.Sprintf("%s %s", projectIcon, selectedProj.Name))
		line2 := PathStyle.Render(selectedProj.Path)
		line3 := DimStyle.Render(fmt.Sprintf("%s ", str.IconArrow)) + DimStyle.Render(editorDisplay)

		infoContent := lipgloss.JoinVertical(lipgloss.Left, line1, line2, line3)
		infoBox = BoxStyle.Width(boxWidth).Render(infoContent)
	}

	var helpText string
	if m.SearchMode {
		helpText = lipgloss.JoinHorizontal(
			lipgloss.Left,
			DimStyle.Render(str.HelpNavigateArrows.String()),
			DimStyle.Render(" • "),
			DimStyle.Render(str.HelpTypeToSearch.String()),
			DimStyle.Render(" • "),
			StatusStyle.Render(str.HelpOpen.String()),
			DimStyle.Render(" • "),
			DimStyle.Render(str.HelpCancel.String()),
		)
	} else {
		helpText = lipgloss.JoinHorizontal(
			lipgloss.Left,
			DimStyle.Render(str.HelpNavigate.String()),
			DimStyle.Render(" • "),
			InfoStyle.Render(str.HelpSearch.String()),
			DimStyle.Render(" • "),
			InfoStyle.Render(str.HelpConfig.String()),
			DimStyle.Render(" • "),
			StatusStyle.Render(str.HelpOpen.String()),
			DimStyle.Render(" • "),
			DimStyle.Render(str.HelpQuit.String()),
		)
	}
	helpBox := BoxStyle.Width(boxWidth).Render(helpText)

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

func RenderProjectEditorModal(m models.Model) string {
	var modalContent []string

	modalContent = append(modalContent, TitleStyle.Render(fmt.Sprintf("%s %s", str.IconTarget, fmt.Sprintf(str.ModalSelectEditor.String(), m.PendingProjectName))))
	if m.PendingProjectType != str.ProjectTypeUnknown.String() {
		modalContent = append(modalContent, "\n"+DimStyle.Render(str.ProjectTypeLabel.String()+m.PendingProjectType))
	}
	modalContent = append(modalContent, "\n"+DimStyle.Render(str.ModalChooseEditor.String()))
	modalContent = append(modalContent, "")

	for i, ed := range m.ProjectEditorOptions {
		if i == m.EditorCursor {
			modalContent = append(modalContent, SelectedItemStyle.Render(fmt.Sprintf("%s %s", str.IconArrowRight, ed.Display)))
		} else {
			modalContent = append(modalContent, NormalItemStyle.Render("  "+ed.Display))
		}
	}

	modalContent = append(modalContent, "")
	modalContent = append(modalContent, DimStyle.Render(str.HelpNavigateUpDown.String()+" • "+str.HelpSelect.String()+" • "+str.HelpCancel.String()))

	content := lipgloss.JoinVertical(lipgloss.Left, modalContent...)

	modalWidth := int(float64(m.Width) * 0.8)
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 70 {
		modalWidth = 70
	}

	modal := BoxStyle.
		Width(modalWidth).
		Render(content)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func RenderConfigMenu(m models.Model) string {
	var modalContent []string

	modalContent = append(modalContent, TitleStyle.Render(fmt.Sprintf("%s %s", str.IconConfig, str.ModalConfiguration)))
	modalContent = append(modalContent, "\n"+DimStyle.Render(str.ModalSelectOption.String()))
	modalContent = append(modalContent, "")

	if m.ConfigMenuCursor == 0 {
		modalContent = append(modalContent, SelectedItemStyle.Render(fmt.Sprintf("%s %s", str.IconArrowRight, str.ConfigSetDefaultEditor)))
	} else {
		modalContent = append(modalContent, NormalItemStyle.Render("  "+str.ConfigSetDefaultEditor.String()))
	}

	if m.ConfigMenuCursor == 1 {
		modalContent = append(modalContent, SelectedItemStyle.Render(fmt.Sprintf("%s %s", str.IconArrowRight, str.ConfigResetConfig)))
	} else {
		modalContent = append(modalContent, NormalItemStyle.Render("  "+str.ConfigResetConfig.String()))
	}

	modalContent = append(modalContent, "")
	modalContent = append(modalContent, DimStyle.Render(str.HelpNavigateUpDown.String()+" • "+str.HelpSelect.String()+" • "+str.HelpCancel.String()))

	content := lipgloss.JoinVertical(lipgloss.Left, modalContent...)

	modalWidth := int(float64(m.Width) * 0.6)
	if modalWidth < 35 {
		modalWidth = 35
	}
	if modalWidth > 55 {
		modalWidth = 55
	}

	modal := BoxStyle.
		Width(modalWidth).
		Render(content)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func RenderEditorSelector(m models.Model) string {
	var modalContent []string

	if m.PendingProjectName == "" {
		modalContent = append(modalContent, TitleStyle.Render(fmt.Sprintf("%s %s", str.IconTarget, str.ModalSelectDefaultEditor)))
		modalContent = append(modalContent, "\n"+DimStyle.Render(str.ModalChooseDefaultEditor.String()))
		modalContent = append(modalContent, "")

		var textEditors, jetbrainsEditors []models.Editor
		for _, ed := range m.ProjectEditorOptions {
			if ed.Name == str.EditorNameGoLand.String() || ed.Name == str.EditorNamePyCharm.String() || ed.Name == str.EditorNameWebStorm.String() ||
				ed.Name == str.EditorNameIntelliJ.String() || ed.Name == str.EditorNamePhpStorm.String() || ed.Name == str.EditorNameRubyMine.String() ||
				ed.Name == str.EditorNameCLion.String() || ed.Name == str.EditorNameRustRover.String() {
				jetbrainsEditors = append(jetbrainsEditors, ed)
			} else {
				textEditors = append(textEditors, ed)
			}
		}

		if len(textEditors) > 0 {
			modalContent = append(modalContent, InfoStyle.Render(fmt.Sprintf("%s %s", str.IconTextEditor, str.CategoryTextEditors)))
			idx := 0
			for _, ed := range textEditors {
				if idx == m.EditorCursor {
					modalContent = append(modalContent, SelectedItemStyle.Render(fmt.Sprintf("%s %s", str.IconArrowRight, ed.Display)))
				} else {
					modalContent = append(modalContent, NormalItemStyle.Render("  "+ed.Display))
				}
				idx++
			}
			modalContent = append(modalContent, "")
		}

		if len(jetbrainsEditors) > 0 {
			modalContent = append(modalContent, InfoStyle.Render(fmt.Sprintf("%s %s", str.IconIDE, str.CategoryJetBrainsIDEs)))
			idx := len(textEditors)
			for _, ed := range jetbrainsEditors {
				if idx == m.EditorCursor {
					modalContent = append(modalContent, SelectedItemStyle.Render(fmt.Sprintf("%s %s", str.IconArrowRight, ed.Display)))
				} else {
					modalContent = append(modalContent, NormalItemStyle.Render("  "+ed.Display))
				}
				idx++
			}
		}
	} else {
		modalContent = append(modalContent, TitleStyle.Render(fmt.Sprintf("%s %s", str.IconTarget, fmt.Sprintf(str.ModalSelectEditor.String(), m.PendingProjectName))))
		if m.PendingProjectType != str.ProjectTypeUnknown.String() {
			modalContent = append(modalContent, "\n"+DimStyle.Render(str.ProjectTypeLabel.String()+m.PendingProjectType))
		}
		modalContent = append(modalContent, "\n"+DimStyle.Render(str.ModalChooseEditor.String()))
		modalContent = append(modalContent, "")

		for i, ed := range m.ProjectEditorOptions {
			if i == m.EditorCursor {
				modalContent = append(modalContent, SelectedItemStyle.Render(fmt.Sprintf("%s %s", str.IconArrowRight, ed.Display)))
			} else {
				modalContent = append(modalContent, NormalItemStyle.Render("  "+ed.Display))
			}
		}
	}

	modalContent = append(modalContent, "")
	modalContent = append(modalContent, DimStyle.Render(str.HelpNavigateUpDown.String()+" • "+str.HelpSelect.String()+" • "+str.HelpCancel.String()))

	content := lipgloss.JoinVertical(lipgloss.Left, modalContent...)

	modalWidth := int(float64(m.Width) * 0.8)
	if modalWidth < 40 {
		modalWidth = 40
	}
	if modalWidth > 70 {
		modalWidth = 70
	}

	modal := BoxStyle.
		Width(modalWidth).
		Render(content)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
	)
}

func View(m models.Model) string {
	if m.Quitting {
		return ""
	}

	if m.Err != nil {
		return lipgloss.Place(
			m.Width,
			m.Height,
			lipgloss.Center,
			lipgloss.Center,
			BoxStyle.Render(fmt.Sprintf("%s %s", str.IconError, str.ErrorPrefix)),
		)
	}

	if m.ShowProjectSelector {
		return RenderProjectEditorModal(m)
	}

	if m.ShowConfigMenu {
		return RenderConfigMenu(m)
	}

	if m.ShowEditorSelector {
		return RenderEditorSelector(m)
	}

	return RenderMainView(m)
}

func CalculateColumns(m models.Model) int {
	if len(m.FilteredProjects) == 0 {
		return 1
	}

	maxLen := 0
	for _, proj := range m.FilteredProjects {
		if len(proj.Name) > maxLen {
			maxLen = len(proj.Name)
		}
	}

	columnWidth := maxLen + 8
	availableWidth := m.Width - 10

	if availableWidth < columnWidth {
		return 1
	}

	cols := availableWidth / columnWidth
	if cols < 1 {
		cols = 1
	}
	if cols > 5 {
		cols = 5
	}

	return cols
}

func FilterProjects(m *models.Model) {
	if m.SearchQuery == "" {
		m.FilteredProjects = m.Projects
		return
	}

	m.FilteredProjects = []models.Project{}
	query := strings.ToLower(m.SearchQuery)
	for _, proj := range m.Projects {
		if strings.Contains(strings.ToLower(proj.Name), query) {
			m.FilteredProjects = append(m.FilteredProjects, proj)
		}
	}
}

func InitialModel(basePath string) models.Model {
	projects, err := project.Load(basePath)
	cfg, _ := config.Load()
	availableEditors := editor.DetectAvailable()

	textEditors, jetbrainsEditors := editor.SeparateByCategory(availableEditors)

	return models.Model{
		Projects:             projects,
		FilteredProjects:     projects,
		Cursor:               0,
		BasePath:             basePath,
		Err:                  err,
		Width:                80,
		Height:               24,
		SearchMode:           false,
		SearchQuery:          "",
		LastOpened:           "",
		Config:               cfg,
		AvailableEditors:     availableEditors,
		TextEditors:          textEditors,
		JetbrainsEditors:     jetbrainsEditors,
		ShowConfigMenu:       false,
		ShowEditorSelector:   false,
		ShowProjectSelector:  false,
		EditorCursor:         0,
		ConfigMenuCursor:     0,
		ProjectEditorOptions: []models.Editor{},
	}
}

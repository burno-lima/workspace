package editor

import (
	"os/exec"
	"runtime"

	"workspace/pkg/models"
	str "workspace/pkg/strings"
)

func DetectAvailable() []models.Editor {
	editors := []models.Editor{
		{Name: str.EditorNameVSCode.String(), Command: str.EditorCommandVSCode.String(), Display: str.EditorVSCode.String()},
		{Name: str.EditorNameCursor.String(), Command: str.EditorCommandCursor.String(), Display: str.EditorCursor.String()},
		{Name: str.EditorNameWindsurf.String(), Command: str.EditorCommandWindsurf.String(), Display: str.EditorWindsurf.String()},
		{Name: str.EditorNameNeovim.String(), Command: str.EditorCommandNeovim.String(), Display: str.EditorNeovim.String()},
		{Name: str.EditorNameVim.String(), Command: str.EditorCommandVim.String(), Display: str.EditorVim.String()},
		{Name: str.EditorNameGoLand.String(), Command: str.EditorCommandGoLand.String(), Display: str.EditorGoLand.String()},
		{Name: str.EditorNamePyCharm.String(), Command: str.EditorCommandPyCharm.String(), Display: str.EditorPyCharm.String()},
		{Name: str.EditorNameWebStorm.String(), Command: str.EditorCommandWebStorm.String(), Display: str.EditorWebStorm.String()},
		{Name: str.EditorNameIntelliJ.String(), Command: str.EditorCommandIntelliJ.String(), Display: str.EditorIntelliJ.String()},
		{Name: str.EditorNamePhpStorm.String(), Command: str.EditorCommandPhpStorm.String(), Display: str.EditorPhpStorm.String()},
		{Name: str.EditorNameRubyMine.String(), Command: str.EditorCommandRubyMine.String(), Display: str.EditorRubyMine.String()},
		{Name: str.EditorNameCLion.String(), Command: str.EditorCommandCLion.String(), Display: str.EditorCLion.String()},
		{Name: str.EditorNameRustRover.String(), Command: str.EditorCommandRustRover.String(), Display: str.EditorRustRover.String()},
		{Name: str.EditorNameSublime.String(), Command: str.EditorCommandSublime.String(), Display: str.EditorSublime.String()},
		{Name: str.EditorNameZed.String(), Command: str.EditorCommandZed.String(), Display: str.EditorZed.String()},
		{Name: str.EditorNameVSCodium.String(), Command: str.EditorCommandVSCodium.String(), Display: str.EditorVSCodium.String()},
		{Name: str.EditorNameAntigravity.String(), Command: str.EditorCommandAntigravity.String(), Display: str.EditorAntigravity.String()},
	}

	var available []models.Editor
	for _, editor := range editors {
		if isCommandAvailable(editor.Command) {
			available = append(available, editor)
		}
	}

	if len(available) == 0 {
		if runtime.GOOS == str.OSDarwin.String() {
			available = append(available, models.Editor{Name: str.EditorNameDefault.String(), Command: str.EditorCommandOpenMacOS.String(), Display: str.EditorDefaultMacOS.String()})
		} else if runtime.GOOS == str.OSWindows.String() {
			available = append(available, models.Editor{Name: str.EditorNameDefault.String(), Command: str.EditorCommandStartWindows.String(), Display: str.EditorDefaultWindows.String()})
		} else {
			available = append(available, models.Editor{Name: str.EditorNameDefault.String(), Command: str.EditorCommandXdgOpen.String(), Display: str.EditorDefaultSystem.String()})
		}
	}

	return available
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func FindJetBrainsIDE(ideName string) string {
	if path, err := exec.LookPath(ideName); err == nil {
		return path
	}
	return ""
}

func GetOptionsForProject(projectType string, availableEditors []models.Editor) []models.Editor {
	var options []models.Editor

	generalEditors := map[string]bool{
		str.EditorNameVSCode.String():      true,
		str.EditorNameCursor.String():      true,
		str.EditorNameNeovim.String():      true,
		str.EditorNameVim.String():         true,
		str.EditorNameWindsurf.String():    true,
		str.EditorNameSublime.String():     true,
		str.EditorNameZed.String():         true,
		str.EditorNameVSCodium.String():    true,
		str.EditorNameAntigravity.String(): true,
	}

	for _, editor := range availableEditors {
		if generalEditors[editor.Name] {
			options = append(options, editor)
		}
	}

	var specificIDE string
	switch projectType {
	case str.ProjectTypeGo.String():
		specificIDE = str.EditorNameGoLand.String()
	case str.ProjectTypePython.String():
		specificIDE = str.EditorNamePyCharm.String()
	case str.ProjectTypeJavaScript.String():
		specificIDE = str.EditorNameWebStorm.String()
	case str.ProjectTypeJava.String():
		specificIDE = str.EditorNameIntelliJ.String()
	case str.ProjectTypePHP.String():
		specificIDE = str.EditorNamePhpStorm.String()
	case str.ProjectTypeRuby.String():
		specificIDE = str.EditorNameRubyMine.String()
	case str.ProjectTypeRust.String():
		specificIDE = str.EditorNameRustRover.String()
	case str.ProjectTypeCPlusPlus.String():
		specificIDE = str.EditorNameCLion.String()
	}

	if specificIDE != "" {
		for _, editor := range availableEditors {
			if editor.Name == specificIDE {
				options = append(options, editor)
				break
			}
		}
	}

	if len(options) == 0 {
		return availableEditors
	}

	return options
}

func SeparateByCategory(availableEditors []models.Editor) (textEditors, jetbrainsEditors []models.Editor) {
	for _, editor := range availableEditors {
		if editor.Name == str.EditorNameGoLand.String() || editor.Name == str.EditorNamePyCharm.String() || editor.Name == str.EditorNameWebStorm.String() ||
			editor.Name == str.EditorNameIntelliJ.String() || editor.Name == str.EditorNamePhpStorm.String() || editor.Name == str.EditorNameRubyMine.String() ||
			editor.Name == str.EditorNameCLion.String() || editor.Name == str.EditorNameRustRover.String() {
			jetbrainsEditors = append(jetbrainsEditors, editor)
		} else {
			textEditors = append(textEditors, editor)
		}
	}
	return
}

package project

import (
	"os"
	"path/filepath"
	"strings"

	"workspace/pkg/editor"
	"workspace/pkg/models"
	str "workspace/pkg/strings"
)

func GetTypes() []models.ProjectType {
	return []models.ProjectType{
		{Name: str.ProjectTypeGo.String(), IDE: "goland", Patterns: []string{"go.mod", "go.sum", "go.work", "go.work.sum"}},
		{Name: str.ProjectTypePython.String(), IDE: "pycharm", Patterns: []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"}},
		{Name: str.ProjectTypeJavaScript.String(), IDE: "webstorm", Patterns: []string{"package.json", "tsconfig.json"}},
		{Name: str.ProjectTypeJava.String(), IDE: "idea", Patterns: []string{"pom.xml", "build.gradle", "build.gradle.kts"}},
		{Name: str.ProjectTypePHP.String(), IDE: "phpstorm", Patterns: []string{"composer.json"}},
		{Name: str.ProjectTypeRuby.String(), IDE: "rubymine", Patterns: []string{"Gemfile", "Rakefile"}},
		{Name: str.ProjectTypeRust.String(), IDE: "rustrover", Patterns: []string{"Cargo.toml"}},
		{Name: str.ProjectTypeCPlusPlus.String(), IDE: "clion", Patterns: []string{"CMakeLists.txt", "Makefile"}},
	}
}

func GetIconForType(projectType string) string {
	switch projectType {
	case str.ProjectTypeGo.String():
		return str.IconGo.String()
	case str.ProjectTypePython.String():
		return str.IconPython.String()
	case str.ProjectTypeJavaScript.String():
		return str.IconJavaScript.String()
	case str.ProjectTypeJava.String():
		return str.IconJava.String()
	case str.ProjectTypePHP.String():
		return str.IconPHP.String()
	case str.ProjectTypeRuby.String():
		return str.IconRuby.String()
	case str.ProjectTypeRust.String():
		return str.IconRust.String()
	case str.ProjectTypeCPlusPlus.String():
		return str.IconCPlusPlus.String()
	default:
		return str.IconProject.String()
	}
}

func DetectType(projectPath string) string {
	projectTypes := GetTypes()

	for _, pt := range projectTypes {
		for _, pattern := range pt.Patterns {
			filePath := filepath.Join(projectPath, pattern)
			if _, err := os.Stat(filePath); err == nil {
				return pt.Name
			}
		}
	}

	maxDepth := 3
	foundType := str.ProjectTypeUnknown.String()

	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectPath, path)
		depth := len(strings.Split(relPath, string(os.PathSeparator)))

		if depth > maxDepth {
			return filepath.SkipDir
		}

		if info.IsDir() && strings.HasPrefix(info.Name(), ".") && path != projectPath {
			return filepath.SkipDir
		}

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

		if !info.IsDir() {
			fileName := info.Name()
			for _, pt := range projectTypes {
				for _, pattern := range pt.Patterns {
					if fileName == pattern && foundType == str.ProjectTypeUnknown.String() {
						foundType = pt.Name
						return filepath.SkipAll
					}
				}
			}
		}

		return nil
	})

	return foundType
}

func GetIDEForType(projectType string) string {
	projectTypes := GetTypes()

	for _, pt := range projectTypes {
		if pt.Name == projectType {
			idePath := editor.FindJetBrainsIDE(pt.IDE)
			if idePath != "" {
				return idePath
			}
			break
		}
	}

	return ""
}

func Load(basePath string) ([]models.Project, error) {
	var projects []models.Project

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			projectPath := filepath.Join(basePath, entry.Name())
			projects = append(projects, models.Project{
				Name:        entry.Name(),
				Path:        projectPath,
				ProjectType: DetectType(projectPath),
			})
		}
	}

	return projects, nil
}

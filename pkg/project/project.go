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
	var types []models.ProjectType
	for _, pattern := range str.ProjectTypePatterns {
		types = append(types, models.ProjectType{
			Name:     pattern.Type.String(),
			IDE:      pattern.IDE.String(),
			Patterns: pattern.Patterns,
		})
	}
	return types
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
	for _, pattern := range str.ProjectTypePatterns {
		for _, filePattern := range pattern.Patterns {
			filePath := filepath.Join(projectPath, filePattern)
			if _, err := os.Stat(filePath); err == nil {
				return pattern.Type.String()
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
			for _, skipDir := range str.SkipDirectories {
				if info.Name() == skipDir.String() {
					return filepath.SkipDir
				}
			}
		}

		if !info.IsDir() {
			fileName := info.Name()
			for _, pattern := range str.ProjectTypePatterns {
				for _, filePattern := range pattern.Patterns {
					if fileName == filePattern && foundType == str.ProjectTypeUnknown.String() {
						foundType = pattern.Type.String()
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
	for typeKey, pattern := range str.ProjectTypePatterns {
		if typeKey.String() == projectType {
			idePath := editor.FindJetBrainsIDE(pattern.IDE.String())
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

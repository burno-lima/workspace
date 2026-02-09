package strings

type ProjectTypePattern struct {
	Type     UIText
	IDE      IDEName
	Patterns []string
}

var ProjectTypePatterns = map[UIText]ProjectTypePattern{
	ProjectTypeGo: {
		Type:     ProjectTypeGo,
		IDE:      IDENameGoLand,
		Patterns: []string{"go.mod", "go.sum", "go.work", "go.work.sum"},
	},
	ProjectTypePython: {
		Type:     ProjectTypePython,
		IDE:      IDENamePyCharm,
		Patterns: []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"},
	},
	ProjectTypeJavaScript: {
		Type:     ProjectTypeJavaScript,
		IDE:      IDENameWebStorm,
		Patterns: []string{"package.json", "tsconfig.json"},
	},
	ProjectTypeJava: {
		Type:     ProjectTypeJava,
		IDE:      IDENameIntelliJ,
		Patterns: []string{"pom.xml", "build.gradle", "build.gradle.kts"},
	},
	ProjectTypePHP: {
		Type:     ProjectTypePHP,
		IDE:      IDENamePhpStorm,
		Patterns: []string{"composer.json"},
	},
	ProjectTypeRuby: {
		Type:     ProjectTypeRuby,
		IDE:      IDENameRubyMine,
		Patterns: []string{"Gemfile", "Rakefile"},
	},
	ProjectTypeRust: {
		Type:     ProjectTypeRust,
		IDE:      IDENameRustRover,
		Patterns: []string{"Cargo.toml"},
	},
	ProjectTypeCPlusPlus: {
		Type:     ProjectTypeCPlusPlus,
		IDE:      IDENameCLion,
		Patterns: []string{"CMakeLists.txt", "Makefile"},
	},
}

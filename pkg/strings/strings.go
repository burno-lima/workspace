package strings

type UIText string

const (
	AppUsage              UIText = "Usage: workspace <path>"
	AppUsageExample       UIText = "Example: workspace /home/brunolima/workspace/napp"
	ErrorPathNotExist     UIText = "Error: Path '%s' does not exist"
	ErrorPrefix           UIText = "Error: %v"
	NoProjectsFound       UIText = "No projects found"
	NoProjectsMatching    UIText = "No projects matching '%s'"
	ProjectTypeLabel      UIText = "Project type: "
	SelectedLabel         UIText = "Selected: "
	SearchLabel           UIText = "Search: "
	SearchPlaceholder     UIText = "(press / to search)"
	EditorNotConfigured   UIText = "Not configured"
	EditorProjectSpecific UIText = "%s (project specific)"
	EditorDefault         UIText = "%s (default)"

	ModalSelectEditor        UIText = "Select Editor for %s"
	ModalSelectDefaultEditor UIText = "Select Default Editor"
	ModalChooseEditor        UIText = "Choose your preferred editor for this project:"
	ModalChooseDefaultEditor UIText = "Choose your default editor:"
	ModalConfiguration       UIText = "Configuration"
	ModalSelectOption        UIText = "Select an option:"

	ConfigSetDefaultEditor UIText = "Set Default Editor"
	ConfigResetConfig      UIText = "Reset Configuration"
	ConfigHeader           UIText = "# Workspace Configuration"
	ConfigFormat           UIText = "# Format: project-name=editor or default=editor"
	ConfigDefaultPrefix    UIText = "default=%s"

	CategoryTextEditors   UIText = "Text Editors"
	CategoryJetBrainsIDEs UIText = "JetBrains IDEs"

	HelpNavigate       UIText = "↑↓←→/hjkl: navigate"
	HelpNavigateArrows UIText = "↑↓←→/ctrl+n,p: navigate"
	HelpNavigateUpDown UIText = "↑/↓: navigate"
	HelpSearch         UIText = "/: search"
	HelpConfig         UIText = "c: config"
	HelpOpen           UIText = "enter: open"
	HelpSelect         UIText = "enter: select"
	HelpQuit           UIText = "q: quit"
	HelpCancel         UIText = "esc: cancel"
	HelpTypeToSearch   UIText = "type to search"

	EditorVSCode         UIText = "Visual Studio Code"
	EditorCursor         UIText = "Cursor"
	EditorWindsurf       UIText = "Windsurf"
	EditorNeovim         UIText = "Neovim"
	EditorVim            UIText = "Vim"
	EditorGoLand         UIText = "GoLand"
	EditorPyCharm        UIText = "PyCharm"
	EditorWebStorm       UIText = "WebStorm"
	EditorIntelliJ       UIText = "IntelliJ IDEA"
	EditorPhpStorm       UIText = "PhpStorm"
	EditorRubyMine       UIText = "RubyMine"
	EditorCLion          UIText = "CLion"
	EditorRustRover      UIText = "RustRover"
	EditorSublime        UIText = "Sublime Text"
	EditorZed            UIText = "Zed"
	EditorVSCodium       UIText = "VSCodium"
	EditorAntigravity    UIText = "Antigravity"
	EditorDefaultMacOS   UIText = "Default macOS Editor"
	EditorDefaultWindows UIText = "Default Windows Editor"
	EditorDefaultSystem  UIText = "Default System Editor"

	ProjectTypeGo         UIText = "Go"
	ProjectTypePython     UIText = "Python"
	ProjectTypeJavaScript UIText = "JavaScript/TypeScript"
	ProjectTypeJava       UIText = "Java"
	ProjectTypePHP        UIText = "PHP"
	ProjectTypeRuby       UIText = "Ruby"
	ProjectTypeRust       UIText = "Rust"
	ProjectTypeCPlusPlus  UIText = "C/C++"
	ProjectTypeUnknown    UIText = "Unknown"
)

func (t UIText) String() string {
	return string(t)
}

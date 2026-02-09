package strings

type EditorName string

const (
	EditorNameVSCode      EditorName = "vscode"
	EditorNameCursor      EditorName = "cursor"
	EditorNameWindsurf    EditorName = "windsurf"
	EditorNameNeovim      EditorName = "neovim"
	EditorNameVim         EditorName = "vim"
	EditorNameGoLand      EditorName = "goland"
	EditorNamePyCharm     EditorName = "pycharm"
	EditorNameWebStorm    EditorName = "webstorm"
	EditorNameIntelliJ    EditorName = "intellij"
	EditorNamePhpStorm    EditorName = "phpstorm"
	EditorNameRubyMine    EditorName = "rubymine"
	EditorNameCLion       EditorName = "clion"
	EditorNameRustRover   EditorName = "rustrover"
	EditorNameSublime     EditorName = "sublime"
	EditorNameZed         EditorName = "zed"
	EditorNameVSCodium    EditorName = "codium"
	EditorNameAntigravity EditorName = "antigravity"
	EditorNameDefault     EditorName = "default"
)

func (e EditorName) String() string {
	return string(e)
}

type EditorCommand string

const (
	EditorCommandVSCode       EditorCommand = "code"
	EditorCommandCursor       EditorCommand = "cursor"
	EditorCommandWindsurf     EditorCommand = "windsurf"
	EditorCommandNeovim       EditorCommand = "nvim"
	EditorCommandVim          EditorCommand = "vim"
	EditorCommandGoLand       EditorCommand = "goland"
	EditorCommandPyCharm      EditorCommand = "pycharm"
	EditorCommandWebStorm     EditorCommand = "webstorm"
	EditorCommandIntelliJ     EditorCommand = "idea"
	EditorCommandPhpStorm     EditorCommand = "phpstorm"
	EditorCommandRubyMine     EditorCommand = "rubymine"
	EditorCommandCLion        EditorCommand = "clion"
	EditorCommandRustRover    EditorCommand = "rustrover"
	EditorCommandSublime      EditorCommand = "subl"
	EditorCommandZed          EditorCommand = "zed"
	EditorCommandVSCodium     EditorCommand = "codium"
	EditorCommandAntigravity  EditorCommand = "antigravity"
	EditorCommandOpenMacOS    EditorCommand = "open"
	EditorCommandStartWindows EditorCommand = "start"
	EditorCommandXdgOpen      EditorCommand = "xdg-open"
)

func (e EditorCommand) String() string {
	return string(e)
}

type IDEName string

const (
	IDENameGoLand    IDEName = "goland"
	IDENamePyCharm   IDEName = "pycharm"
	IDENameWebStorm  IDEName = "webstorm"
	IDENameIntelliJ  IDEName = "idea"
	IDENamePhpStorm  IDEName = "phpstorm"
	IDENameRubyMine  IDEName = "rubymine"
	IDENameCLion     IDEName = "clion"
	IDENameRustRover IDEName = "rustrover"
)

func (i IDEName) String() string {
	return string(i)
}

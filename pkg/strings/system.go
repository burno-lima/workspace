package strings

type OSName string

const (
	OSDarwin  OSName = "darwin"
	OSWindows OSName = "windows"
	OSLinux   OSName = "linux"
)

func (o OSName) String() string {
	return string(o)
}

type ConfigKey string

const (
	ConfigKeyDefault ConfigKey = "default"
)

func (c ConfigKey) String() string {
	return string(c)
}

type ConfigDir string

const (
	ConfigDirName     ConfigDir = ".config"
	ConfigDirWorkspace ConfigDir = "workspace"
)

func (c ConfigDir) String() string {
	return string(c)
}

type ConfigFile string

const (
	ConfigFileName ConfigFile = "workspace.config"
)

func (c ConfigFile) String() string {
	return string(c)
}

type FilePermission uint32

const (
	PermissionDirDefault  FilePermission = 0755
	PermissionFileDefault FilePermission = 0644
)

type SkipDirectory string

const (
	SkipDirNodeModules SkipDirectory = "node_modules"
	SkipDirVendor      SkipDirectory = "vendor"
	SkipDirVenv        SkipDirectory = "venv"
	SkipDirVenvHidden  SkipDirectory = ".venv"
	SkipDirPycache     SkipDirectory = "__pycache__"
	SkipDirBuild       SkipDirectory = "build"
	SkipDirDist        SkipDirectory = "dist"
	SkipDirTarget      SkipDirectory = "target"
)

func (s SkipDirectory) String() string {
	return string(s)
}

var SkipDirectories = []SkipDirectory{
	SkipDirNodeModules,
	SkipDirVendor,
	SkipDirVenv,
	SkipDirVenvHidden,
	SkipDirPycache,
	SkipDirBuild,
	SkipDirDist,
	SkipDirTarget,
}

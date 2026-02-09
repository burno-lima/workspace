package strings

type Icon string

const (
	IconFolder     Icon = "\uf07b"
	IconSearch     Icon = "\uf002"
	IconError      Icon = "\uf00d"
	IconConfig     Icon = "\uf013"
	IconTarget     Icon = "\uf140"
	IconTextEditor Icon = "\uf0f6"
	IconIDE        Icon = "\ue7a5"
	IconArrow      Icon = "\uf054"
	IconArrowRight Icon = "\u25b6"
	IconProject    Icon = "\uf413"
	IconGo         Icon = "\ue626"
	IconPython     Icon = "\ue73c"
	IconJavaScript Icon = "\ue74e"
	IconJava       Icon = "\ue738"
	IconPHP        Icon = "\ue73d"
	IconRuby       Icon = "\ue791"
	IconRust       Icon = "\ue7a8"
	IconCPlusPlus  Icon = "\ue61d"
)

func (i Icon) String() string {
	return string(i)
}

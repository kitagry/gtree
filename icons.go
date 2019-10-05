package main

import "github.com/fatih/color"

var (
	brown       = "905532"
	aqua        = "3AFFDB"
	blue        = "689FB6"
	darkBlue    = "44788E"
	purple      = "834F79"
	lightPurple = "834F79"
	red         = "AE403F"
	beige       = "F5C06F"
	yellow      = "F09F17"
	orange      = "D4843E"
	darkOrange  = "F16529"
	pink        = "CB6F6F"
	green       = "8FAA54"
	lightGreen  = "31B53E"
	white       = "FFFFFF"
)

type Icon struct {
	Icon  string
	Color color.Attribute
}

var defaultFileIcon = Icon{
	Icon:  "",
	Color: color.FgWhite,
}

var icons map[string]Icon = map[string]Icon{
	"styl":     Icon{Icon: "", Color: color.FgGreen},
	"sass":     Icon{Icon: "", Color: color.FgWhite},
	"scss":     Icon{Icon: "", Color: color.FgMagenta},
	"htm":      Icon{Icon: "", Color: color.FgHiRed},
	"html":     Icon{Icon: "", Color: color.FgHiRed},
	"slim":     Icon{Icon: "", Color: color.FgHiRed},
	"ejs":      Icon{Icon: "", Color: color.FgYellow},
	"css":      Icon{Icon: "", Color: color.FgBlue},
	"less":     Icon{Icon: "", Color: color.FgBlue},
	"md":       Icon{Icon: "", Color: color.FgYellow},
	"markdown": Icon{Icon: "", Color: color.FgYellow},
	"rmd":      Icon{Icon: "", Color: color.FgWhite},
	"json":     Icon{Icon: "", Color: color.FgWhite},
	"js":       Icon{Icon: "", Color: color.FgWhite},
	"mjs":      Icon{Icon: "", Color: color.FgWhite},
	"jsx":      Icon{Icon: "", Color: color.FgBlue},
	"rb":       Icon{Icon: "", Color: color.FgRed},
	"php":      Icon{Icon: "", Color: color.FgMagenta},
	"py":       Icon{Icon: "", Color: color.FgYellow},
	"pyc":      Icon{Icon: "", Color: color.FgYellow},
	"pyo":      Icon{Icon: "", Color: color.FgYellow},
	"pyd":      Icon{Icon: "", Color: color.FgYellow},
	"coffee":   Icon{Icon: "", Color: color.FgYellow},
	"mustache": Icon{Icon: "", Color: color.FgHiRed},
	"hbs":      Icon{Icon: "", Color: color.FgHiRed},
	"conf":     Icon{Icon: "", Color: color.FgWhite},
	"ini":      Icon{Icon: "", Color: color.FgWhite},
	"yml":      Icon{Icon: "", Color: color.FgWhite},
	"yaml":     Icon{Icon: "", Color: color.FgWhite},
	"bat":      Icon{Icon: "", Color: color.FgWhite},
	"jpg":      Icon{Icon: "", Color: color.FgCyan},
	"jpeg":     Icon{Icon: "", Color: color.FgCyan},
	"bmp":      Icon{Icon: "", Color: color.FgCyan},
	"png":      Icon{Icon: "", Color: color.FgCyan},
	"gif":      Icon{Icon: "", Color: color.FgCyan},
	"ico":      Icon{Icon: "", Color: color.FgCyan},
	"twig":     Icon{Icon: "", Color: color.FgGreen},
	"cpp":      Icon{Icon: "", Color: color.FgBlue},
	"cxx":      Icon{Icon: "", Color: color.FgBlue},
	"cc":       Icon{Icon: "", Color: color.FgBlue},
	"cp":       Icon{Icon: "", Color: color.FgBlue},
	"c":        Icon{Icon: "", Color: color.FgBlue},
	"h":        Icon{Icon: "", Color: color.FgWhite},
	"hpp":      Icon{Icon: "", Color: color.FgWhite},
	"hxx":      Icon{Icon: "", Color: color.FgWhite},
	"hs":       Icon{Icon: "", Color: color.FgWhite},
	"lhs":      Icon{Icon: "", Color: color.FgWhite},
	"lua":      Icon{Icon: "", Color: color.FgMagenta},
	"java":     Icon{Icon: "", Color: color.FgMagenta},
	"sh":       Icon{Icon: "", Color: color.FgMagenta},
	"fish":     Icon{Icon: "", Color: color.FgGreen},
	"bash":     Icon{Icon: "", Color: color.FgWhite},
	"zsh":      Icon{Icon: "", Color: color.FgWhite},
	"ksh":      Icon{Icon: "", Color: color.FgWhite},
	"csh":      Icon{Icon: "", Color: color.FgWhite},
	"awk":      Icon{Icon: "", Color: color.FgWhite},
	"ps1":      Icon{Icon: "", Color: color.FgWhite},
	"ml":       Icon{Icon: "λ", Color: color.FgYellow},
	"mli":      Icon{Icon: "λ", Color: color.FgYellow},
	"diff":     Icon{Icon: "", Color: color.FgWhite},
	"db":       Icon{Icon: "", Color: color.FgBlue},
	"sql":      Icon{Icon: "", Color: color.FgBlue},
	"dump":     Icon{Icon: "", Color: color.FgBlue},
	"clj":      Icon{Icon: "", Color: color.FgGreen},
	"cljc":     Icon{Icon: "", Color: color.FgGreen},
	"cljs":     Icon{Icon: "", Color: color.FgGreen},
	"edn":      Icon{Icon: "", Color: color.FgGreen},
	"scala":    Icon{Icon: "", Color: color.FgRed},
	"go":       Icon{Icon: "", Color: color.FgWhite},
	"dart":     Icon{Icon: "", Color: color.FgWhite},
	"xul":      Icon{Icon: "", Color: color.FgHiRed},
	"sln":      Icon{Icon: "", Color: color.FgMagenta},
	"suo":      Icon{Icon: "", Color: color.FgMagenta},
	"pl":       Icon{Icon: "", Color: color.FgBlue},
	"pm":       Icon{Icon: "", Color: color.FgBlue},
	"t":        Icon{Icon: "", Color: color.FgBlue},
	"rss":      Icon{Icon: "", Color: color.FgHiRed},
	"fsscript": Icon{Icon: "", Color: color.FgBlue},
	"fsx":      Icon{Icon: "", Color: color.FgBlue},
	"fs":       Icon{Icon: "", Color: color.FgBlue},
	"fsi":      Icon{Icon: "", Color: color.FgBlue},
	"rs":       Icon{Icon: "", Color: color.FgHiRed},
	"rlib":     Icon{Icon: "", Color: color.FgHiRed},
	"d":        Icon{Icon: "", Color: color.FgRed},
	"erl":      Icon{Icon: "", Color: color.FgMagenta},
	"ex":       Icon{Icon: "", Color: color.FgMagenta},
	"exs":      Icon{Icon: "", Color: color.FgMagenta},
	"eex":      Icon{Icon: "", Color: color.FgMagenta},
	"hrl":      Icon{Icon: "", Color: color.FgMagenta},
	"vim":      Icon{Icon: "", Color: color.FgGreen},
	"ai":       Icon{Icon: "", Color: color.FgHiRed},
	"psd":      Icon{Icon: "", Color: color.FgBlue},
	"psb":      Icon{Icon: "", Color: color.FgBlue},
	"ts":       Icon{Icon: "", Color: color.FgBlue},
	"tsx":      Icon{Icon: "", Color: color.FgWhite},
	"jl":       Icon{Icon: "", Color: color.FgMagenta},
	"pp":       Icon{Icon: "", Color: color.FgWhite},
	"vue":      Icon{Icon: "﵂", Color: color.FgGreen},
}

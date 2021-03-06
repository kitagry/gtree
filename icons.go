package main

import "github.com/gookit/color"

// Icon is a set of icon and color.
type Icon struct {
	Icon  string
	Color color.Color
}

var defaultFolderIcon = Icon{
	Icon:  "",
	Color: color.FgBlue,
}

var defaultFileIcon = Icon{
	Icon:  "",
	Color: color.FgWhite,
}

var icons = map[string]Icon{
	"styl":     {Icon: "", Color: color.FgGreen},
	"sass":     {Icon: "", Color: color.FgWhite},
	"scss":     {Icon: "", Color: color.FgMagenta},
	"htm":      {Icon: "", Color: color.FgLightRed},
	"html":     {Icon: "", Color: color.FgLightRed},
	"slim":     {Icon: "", Color: color.FgLightRed},
	"ejs":      {Icon: "", Color: color.FgYellow},
	"css":      {Icon: "", Color: color.FgBlue},
	"less":     {Icon: "", Color: color.FgBlue},
	"md":       {Icon: "", Color: color.FgYellow},
	"markdown": {Icon: "", Color: color.FgYellow},
	"rmd":      {Icon: "", Color: color.FgWhite},
	"json":     {Icon: "", Color: color.FgWhite},
	"js":       {Icon: "", Color: color.FgWhite},
	"mjs":      {Icon: "", Color: color.FgWhite},
	"jsx":      {Icon: "", Color: color.FgBlue},
	"rb":       {Icon: "", Color: color.FgRed},
	"php":      {Icon: "", Color: color.FgMagenta},
	"py":       {Icon: "", Color: color.FgYellow},
	"pyc":      {Icon: "", Color: color.FgYellow},
	"pyo":      {Icon: "", Color: color.FgYellow},
	"pyd":      {Icon: "", Color: color.FgYellow},
	"coffee":   {Icon: "", Color: color.FgYellow},
	"mustache": {Icon: "", Color: color.FgLightRed},
	"hbs":      {Icon: "", Color: color.FgLightRed},
	"conf":     {Icon: "", Color: color.FgWhite},
	"ini":      {Icon: "", Color: color.FgWhite},
	"yml":      {Icon: "", Color: color.FgWhite},
	"yaml":     {Icon: "", Color: color.FgWhite},
	"bat":      {Icon: "", Color: color.FgWhite},
	"jpg":      {Icon: "", Color: color.FgCyan},
	"jpeg":     {Icon: "", Color: color.FgCyan},
	"bmp":      {Icon: "", Color: color.FgCyan},
	"png":      {Icon: "", Color: color.FgCyan},
	"gif":      {Icon: "", Color: color.FgCyan},
	"ico":      {Icon: "", Color: color.FgCyan},
	"twig":     {Icon: "", Color: color.FgGreen},
	"cpp":      {Icon: "", Color: color.FgBlue},
	"cxx":      {Icon: "", Color: color.FgBlue},
	"cc":       {Icon: "", Color: color.FgBlue},
	"cp":       {Icon: "", Color: color.FgBlue},
	"c":        {Icon: "", Color: color.FgBlue},
	"h":        {Icon: "", Color: color.FgWhite},
	"hpp":      {Icon: "", Color: color.FgWhite},
	"hxx":      {Icon: "", Color: color.FgWhite},
	"hs":       {Icon: "", Color: color.FgWhite},
	"lhs":      {Icon: "", Color: color.FgWhite},
	"lua":      {Icon: "", Color: color.FgMagenta},
	"java":     {Icon: "", Color: color.FgMagenta},
	"sh":       {Icon: "", Color: color.FgMagenta},
	"fish":     {Icon: "", Color: color.FgGreen},
	"bash":     {Icon: "", Color: color.FgWhite},
	"zsh":      {Icon: "", Color: color.FgWhite},
	"ksh":      {Icon: "", Color: color.FgWhite},
	"csh":      {Icon: "", Color: color.FgWhite},
	"awk":      {Icon: "", Color: color.FgWhite},
	"ps1":      {Icon: "", Color: color.FgWhite},
	"ml":       {Icon: "λ", Color: color.FgYellow},
	"mli":      {Icon: "λ", Color: color.FgYellow},
	"diff":     {Icon: "", Color: color.FgWhite},
	"db":       {Icon: "", Color: color.FgBlue},
	"sql":      {Icon: "", Color: color.FgBlue},
	"dump":     {Icon: "", Color: color.FgBlue},
	"clj":      {Icon: "", Color: color.FgGreen},
	"cljc":     {Icon: "", Color: color.FgGreen},
	"cljs":     {Icon: "", Color: color.FgGreen},
	"edn":      {Icon: "", Color: color.FgGreen},
	"scala":    {Icon: "", Color: color.FgRed},
	"go":       {Icon: "", Color: color.FgWhite},
	"dart":     {Icon: "", Color: color.FgWhite},
	"xul":      {Icon: "", Color: color.FgLightRed},
	"sln":      {Icon: "", Color: color.FgMagenta},
	"suo":      {Icon: "", Color: color.FgMagenta},
	"pl":       {Icon: "", Color: color.FgBlue},
	"pm":       {Icon: "", Color: color.FgBlue},
	"t":        {Icon: "", Color: color.FgBlue},
	"rss":      {Icon: "", Color: color.FgLightRed},
	"fsscript": {Icon: "", Color: color.FgBlue},
	"fsx":      {Icon: "", Color: color.FgBlue},
	"fs":       {Icon: "", Color: color.FgBlue},
	"fsi":      {Icon: "", Color: color.FgBlue},
	"rs":       {Icon: "", Color: color.FgLightRed},
	"rlib":     {Icon: "", Color: color.FgLightRed},
	"d":        {Icon: "", Color: color.FgRed},
	"erl":      {Icon: "", Color: color.FgMagenta},
	"ex":       {Icon: "", Color: color.FgMagenta},
	"exs":      {Icon: "", Color: color.FgMagenta},
	"eex":      {Icon: "", Color: color.FgMagenta},
	"hrl":      {Icon: "", Color: color.FgMagenta},
	"vim":      {Icon: "", Color: color.FgGreen},
	"ai":       {Icon: "", Color: color.FgLightRed},
	"psd":      {Icon: "", Color: color.FgBlue},
	"psb":      {Icon: "", Color: color.FgBlue},
	"ts":       {Icon: "", Color: color.FgBlue},
	"tsx":      {Icon: "", Color: color.FgWhite},
	"jl":       {Icon: "", Color: color.FgMagenta},
	"pp":       {Icon: "", Color: color.FgWhite},
	"vue":      {Icon: "﵂", Color: color.FgGreen},
}

func NewIconString(suffix string) string {
	icon, ok := icons[suffix]
	if !ok {
		icon = defaultFileIcon
	}
	return color.New(icon.Color).Sprint(icon.Icon)
}

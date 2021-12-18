package dev

// || TOOLING ||

type ToolingConfig []string

var RequiredTools = ToolingConfig{
	"multipass",
	"kubernetes-cli",
	"krew",
	"yq",
	"helm",
}

// Package main is the entry point for YOLOSINT.
package main

import (
	"log"

	"github.com/Gaurav-Gosain/tuios/pkg/tuios"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/thesavant42/yolosint/internal/style"
)

func main() {
	model := tuios.New(
		tuios.WithTheme(style.DefaultTheme),
		tuios.WithBorderStyle(style.DefaultBorder),
		tuios.WithDockbarPosition(style.DefaultDockPosition),
		tuios.WithWorkspaces(style.DefaultWorkspaces),
		tuios.WithAnimations(style.DefaultAnimations),
		tuios.WithScrollbackLines(style.DefaultScrollback),
	)

	p := tea.NewProgram(model, tuios.ProgramOptions()...)

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}


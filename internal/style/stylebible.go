// Package style defines the visual theme constants for TUIOS windows.
package style

// Theme names supported by TUIOS (via bubbletint).
// Default is Dracula but all 300+ themes are available.
const (
	ThemeDracula    = "dracula"
	ThemeNord       = "nord"
	ThemeTokyoNight = "tokyonight"
	ThemeGruvbox    = "gruvbox"
	ThemeMonokai    = "monokai"
	ThemeSolarized  = "solarized"
	ThemeCatppuccin = "catppuccin"
	ThemeOneDark    = "onedark"
)

// DefaultTheme is the default color theme.
const DefaultTheme = ThemeDracula

// Border style constants for TUIOS windows.
const (
	BorderRounded = "rounded"
	BorderNormal  = "normal"
	BorderThick   = "thick"
	BorderDouble  = "double"
	BorderHidden  = "hidden"
	BorderBlock   = "block"
	BorderASCII   = "ascii"
)

// DefaultBorder is the default window border style.
const DefaultBorder = BorderNormal

// Dockbar position constants.
const (
	DockTop    = "top"
	DockBottom = "bottom"
	DockHidden = "hidden"
)

// DefaultDockPosition is the default dockbar position.
const DefaultDockPosition = DockTop

// Workspace configuration.
const (
	DefaultWorkspaces   = 9
	MinWorkspaces       = 1
	MaxWorkspaces       = 9
	DefaultScrollback   = 10000
	MinScrollback       = 100
	MaxScrollback       = 1000000
)

// UI feature flags.
const (
	DefaultAnimations       = true
	DefaultASCIIOnly        = false // Use Nerd Font icons
	DefaultShowKeys         = false
	DefaultHideWindowButtons = false
)


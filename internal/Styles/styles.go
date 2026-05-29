package Styles

import (
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/charmbracelet/lipgloss"
)

// ═════════════════════════════════════════════════════════════════════════
// 1. TUI/CLI DASHBOARD ENGINE STRUCTURES & LOGIC
// ═════════════════════════════════════════════════════════════════════════

// AppStyles holds all the generated lipgloss styles for a specific theme
type AppStyles struct {
	Title       lipgloss.Style
	Text        lipgloss.Style
	SubText     lipgloss.Style
	Box         lipgloss.Style
	ListLabel   lipgloss.Style
	ListValue   lipgloss.Style
	ListPercent lipgloss.Style

	// Dynamic colors for your bars and graphs
	Color1 lipgloss.Style
	Color2 lipgloss.Style
	Color3 lipgloss.Style
	Color4 lipgloss.Style

	Navbar lipgloss.Style
	Footer lipgloss.Style

	// Timestats
	StatCard      lipgloss.Style
	StatCardTitle lipgloss.Style
	StatCardValue lipgloss.Style
}

// InitStyles acts as a factory. You pass in a ThemeConfig, and it returns
// a full set of Lipgloss styles mapped exactly to those colors.
func InitStyles(theme types.ThemeConfig) AppStyles {
	return AppStyles{
		// Headers & Layout
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(theme.Color1)). 
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Color1)). 
			Padding(0, 1).
			MarginBottom(1),

		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.SubTextColor)).
			Padding(0, 1).
			MarginBottom(1),

		// Base Text
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.TextColor)),

		SubText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.SubTextColor)),

		// Formatted Lists (Perfect for Language/Project stats)
		ListLabel: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.TextColor)).
			Bold(true).
			Width(15), 

		ListValue: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.SubTextColor)).
			Width(10), 

		ListPercent: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Color2)). 
			Italic(true),

		// Raw Colors (Useful for rendering progress bars)
		Color1: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color1)),
		Color2: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color2)),
		Color3: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color3)),
		Color4: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color4)),

		// Navbar
		Navbar: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.BackgroundColor)). 
			Background(lipgloss.Color(theme.Color1)).
			Padding(0, 1).
			MarginBottom(1),

		// Footer
		Footer: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.SubTextColor)).
			MarginTop(1),

		StatCard: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.SubTextColor)).
			Padding(0, 1).
			Align(lipgloss.Center), 

		StatCardTitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.SubTextColor)).
			MarginBottom(1), 

		StatCardValue: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Color1)). 
			Bold(true),
	}
}

// BuildStyles compiles raw hex codes from a theme into Lipgloss styles
func BuildStyles(theme types.ThemeConfig) AppStyles {
	return AppStyles{
		Color1: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color1)),
		Color2: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color2)),
		Color3: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color3)), 
		Color4: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color4)), 

		Title:   lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color1)).Bold(true), 
		Text:    lipgloss.NewStyle().Foreground(lipgloss.Color(theme.TextColor)),
		SubText: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SubTextColor)),

		Navbar: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color2)).Bold(true), 
		Footer: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SubTextColor)),      

		Box: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.BarBackgroundColor)).
			Padding(1, 2),

		ListLabel:   lipgloss.NewStyle().Foreground(lipgloss.Color(theme.TextColor)).Bold(true),
		ListValue:   lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color1)),
		ListPercent: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SubTextColor)).Italic(true),

		StatCard: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.BarBackgroundColor)).
			Padding(0, 1), 
		StatCardTitle: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.SubTextColor)),
		StatCardValue: lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Color1)).Bold(true),
	}
}

// ═════════════════════════════════════════════════════════════════════════
// 2. SVG REPORT RENDERER PALETTES
// ═════════════════════════════════════════════════════════════════════════

// ThemeColors holds every color slot used by the SVG renderer
type ThemeColors struct {
	Background string
	Text       string
	SubText    string
	BarBg      string
	C1         string
	C2         string
	C3         string
	C4         string
}

// GetTheme returns the color palette for the requested theme name.
func GetTheme(name string) ThemeColors {
	switch name {

	case "dark":
		return ThemeColors{
			Background: "#161b22", Text: "#e6edf3", SubText: "#8b949e",
			BarBg: "#30363d", C1: "#58a6ff", C2: "#3fb950",
			C3: "#d29922", C4: "#f85149",
		}

	// ── NEW: Rosé Pine ───────────────────────────────────────────────────────
	case "rosepine":
		return ThemeColors{
			Background: "#191724", // Base   — deep night-sky purple
			Text:       "#e0def4", // Text   — soft lavender white
			SubText:    "#6e6a86", // Muted  — desaturated purple
			BarBg:      "#26233a", // Overlay — dark overlay for empty bars
			C1:         "#eb6f92", // Love   — rose-red   (All Time / highest)
			C2:         "#f6c177", // Gold   — warm amber  (Last 30 Days)
			C3:         "#ebbcba", // Rose   — dusty rose  (Last 7 Days)
			C4:         "#31748f", // Pine   — teal        (Yesterday / lowest)
		}

	default:
		return ThemeColors{
			Background: "#161b22", Text: "#e6edf3", SubText: "#8b949e",
			BarBg: "#30363d", C1: "#58a6ff", C2: "#3fb950",
			C3: "#d29922", C4: "#f85149",
		}
	}
}
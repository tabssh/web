package terminal

import "testing"

func TestCalculateMode(t *testing.T) {
	tests := []struct {
		name string
		cols int
		rows int
		want SizeMode
	}{
		{"micro narrow", 39, 24, SizeModeMicro},
		{"micro short", 100, 9, SizeModeMicro},
		{"minimal low cols", 40, 24, SizeModeMinimal},
		{"minimal low rows", 100, 15, SizeModeMinimal},
		{"compact cols", 60, 24, SizeModeCompact},
		{"compact rows", 100, 23, SizeModeCompact},
		{"standard", 80, 24, SizeModeStandard},
		{"standard upper", 119, 39, SizeModeStandard},
		{"wide", 120, 40, SizeModeWide},
		{"wide upper", 199, 59, SizeModeWide},
		{"ultrawide", 200, 60, SizeModeUltrawide},
		{"ultrawide upper", 399, 79, SizeModeUltrawide},
		{"massive", 400, 80, SizeModeMassive},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateMode(tt.cols, tt.rows); got != tt.want {
				t.Errorf("calculateMode(%d, %d) = %v, want %v", tt.cols, tt.rows, got, tt.want)
			}
		})
	}
}

func TestSizeModeHelpers(t *testing.T) {
	tests := []struct {
		mode    SizeMode
		art     bool
		borders bool
		sidebar bool
		icons   bool
	}{
		{SizeModeMicro, false, false, false, false},
		{SizeModeMinimal, false, false, false, true},
		{SizeModeCompact, false, true, false, true},
		{SizeModeStandard, true, true, false, true},
		{SizeModeWide, true, true, true, true},
		{SizeModeUltrawide, true, true, true, true},
		{SizeModeMassive, true, true, true, true},
	}
	for _, tt := range tests {
		if got := tt.mode.ShowASCIIArt(); got != tt.art {
			t.Errorf("mode %d ShowASCIIArt() = %v, want %v", tt.mode, got, tt.art)
		}
		if got := tt.mode.ShowBorders(); got != tt.borders {
			t.Errorf("mode %d ShowBorders() = %v, want %v", tt.mode, got, tt.borders)
		}
		if got := tt.mode.ShowSidebar(); got != tt.sidebar {
			t.Errorf("mode %d ShowSidebar() = %v, want %v", tt.mode, got, tt.sidebar)
		}
		if got := tt.mode.ShowIcons(); got != tt.icons {
			t.Errorf("mode %d ShowIcons() = %v, want %v", tt.mode, got, tt.icons)
		}
	}
}

func TestGetTerminalSizeDefaults(t *testing.T) {
	size := GetTerminalSize()
	if size.Cols <= 0 || size.Rows <= 0 {
		t.Errorf("GetTerminalSize() = %+v, dimensions must be positive", size)
	}
	if size.Mode != calculateMode(size.Cols, size.Rows) {
		t.Errorf("GetTerminalSize() mode %v inconsistent with dimensions %dx%d", size.Mode, size.Cols, size.Rows)
	}
}

func TestGetSymbols(t *testing.T) {
	uni := GetSymbols(true)
	if uni.OK != "✓" || uni.Horiz != "─" {
		t.Errorf("GetSymbols(true) = %+v, want Unicode set", uni)
	}
	ascii := GetSymbols(false)
	if ascii.OK != "[OK]" || ascii.Error != "[ERROR]" || ascii.Warn != "[WARN]" {
		t.Errorf("GetSymbols(false) = %+v, want ASCII status tags", ascii)
	}
	if ascii.TopLeft != "+" || ascii.Horiz != "-" || ascii.Vert != "|" {
		t.Errorf("GetSymbols(false) = %+v, want +--+ table borders", ascii)
	}
}

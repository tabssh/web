// Symbol sets for terminal output, per AI.md PART 7 (TERM=dumb rules) and
// PART 8 (NO_COLOR). Dumb terminals and NO_COLOR environments get plain
// ASCII status tags and +--+ style tables instead of Unicode symbols.
package terminal

// Symbols is a set of status and box-drawing characters for terminal output.
type Symbols struct {
	OK       string
	Error    string
	Warn     string
	Info     string
	Bullet   string
	Arrow    string
	Horiz    string
	Vert     string
	TopLeft  string
	TopRight string
	BotLeft  string
	BotRight string
	Cross    string
}

// unicodeSymbols is the default symbol set for capable terminals.
var unicodeSymbols = Symbols{
	OK:       "✓",
	Error:    "✗",
	Warn:     "⚠",
	Info:     "ℹ",
	Bullet:   "•",
	Arrow:    "→",
	Horiz:    "─",
	Vert:     "│",
	TopLeft:  "┌",
	TopRight: "┐",
	BotLeft:  "└",
	BotRight: "┘",
	Cross:    "┼",
}

// asciiSymbols is the fallback set for dumb terminals: plain text status
// tags and +--+ table borders, no Unicode.
var asciiSymbols = Symbols{
	OK:       "[OK]",
	Error:    "[ERROR]",
	Warn:     "[WARN]",
	Info:     "[INFO]",
	Bullet:   "*",
	Arrow:    "->",
	Horiz:    "-",
	Vert:     "|",
	TopLeft:  "+",
	TopRight: "+",
	BotLeft:  "+",
	BotRight: "+",
	Cross:    "+",
}

// GetSymbols returns the Unicode symbol set when unicode is true, otherwise
// the plain-ASCII fallback set.
func GetSymbols(unicode bool) Symbols {
	if unicode {
		return unicodeSymbols
	}
	return asciiSymbols
}

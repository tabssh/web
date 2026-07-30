// Shell completions and init commands per AI.md PART 7 (Shell
// Completions): generated live from the running binary for bash, zsh,
// fish, sh, dash, ksh, powershell, and pwsh, auto-detecting the shell
// from $SHELL when omitted. Never shipped as static files.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// supportedShells is the exhaustive list of supported completion shells.
var supportedShells = []string{"bash", "zsh", "fish", "sh", "dash", "ksh", "powershell", "pwsh"}

// serverFlags lists every top-level flag for completion generation.
var serverFlags = []string{
	"--help", "--version", "--status", "--shell", "--mode", "--config",
	"--data", "--cache", "--log", "--backup", "--pid", "--address",
	"--port", "--baseurl", "--daemon", "--debug", "--color", "--lang",
	"--service", "--maintenance", "--update",
}

// detectShell resolves the target shell: explicit argument first, then
// the basename of $SHELL, defaulting to bash (powershell on Windows).
func detectShell(arg string) string {
	if arg != "" {
		return strings.ToLower(arg)
	}
	if sh := os.Getenv("SHELL"); sh != "" {
		return strings.ToLower(filepath.Base(sh))
	}
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	return "bash"
}

// isSupportedShell reports whether the named shell is supported.
func isSupportedShell(name string) bool {
	for _, s := range supportedShells {
		if s == name {
			return true
		}
	}
	return false
}

// completionScript returns the completion script for the given shell and
// binary name. The caller must have validated the shell name.
func completionScript(shell, binName string) string {
	flags := strings.Join(serverFlags, " ")
	switch shell {
	case "bash", "ksh":
		return fmt.Sprintf(`# %[1]s completion for %[2]s
_%[3]s_completions() {
  local cur="${COMP_WORDS[COMP_CWORD]}"
  local prev="${COMP_WORDS[COMP_CWORD-1]}"
  case "$prev" in
    --shell) COMPREPLY=($(compgen -W "completions init help" -- "$cur")); return ;;
    --service) COMPREPLY=($(compgen -W "start stop restart reload --install --disable --uninstall help" -- "$cur")); return ;;
    --maintenance) COMPREPLY=($(compgen -W "backup restore update mode setup help" -- "$cur")); return ;;
    --update) COMPREPLY=($(compgen -W "check yes branch help" -- "$cur")); return ;;
    --mode) COMPREPLY=($(compgen -W "production development" -- "$cur")); return ;;
    --color) COMPREPLY=($(compgen -W "auto yes no" -- "$cur")); return ;;
    --config|--data|--cache|--log|--backup) COMPREPLY=($(compgen -d -- "$cur")); return ;;
    --pid) COMPREPLY=($(compgen -f -- "$cur")); return ;;
  esac
  COMPREPLY=($(compgen -W "%[4]s" -- "$cur"))
}
complete -F _%[3]s_completions %[2]s
`, shell, binName, sanitizeFuncName(binName), flags)
	case "zsh":
		return fmt.Sprintf(`#compdef %[1]s
# zsh completion for %[1]s
_%[2]s() {
  local -a flags
  flags=(%[3]s)
  case "$words[CURRENT-1]" in
    --shell) _values 'shell command' completions init help; return ;;
    --service) _values 'service command' start stop restart reload --install --disable --uninstall help; return ;;
    --maintenance) _values 'maintenance command' backup restore update mode setup help; return ;;
    --update) _values 'update command' check yes branch help; return ;;
    --mode) _values 'mode' production development; return ;;
    --color) _values 'color' auto yes no; return ;;
    --config|--data|--cache|--log|--backup) _directories; return ;;
    --pid) _files; return ;;
  esac
  _describe 'flags' flags
}
compdef _%[2]s %[1]s
`, binName, sanitizeFuncName(binName), flags)
	case "fish":
		var b strings.Builder
		fmt.Fprintf(&b, "# fish completion for %s\n", binName)
		for _, f := range serverFlags {
			fmt.Fprintf(&b, "complete -c %s -l %s\n", binName, strings.TrimPrefix(f, "--"))
		}
		fmt.Fprintf(&b, "complete -c %s -s h -l help\n", binName)
		fmt.Fprintf(&b, "complete -c %s -s v -l version\n", binName)
		return b.String()
	case "sh", "dash":
		return fmt.Sprintf(`# %[1]s does not support programmable completion.
# Available %[2]s flags:
# %[3]s
`, shell, binName, flags)
	case "powershell", "pwsh":
		return fmt.Sprintf(`# PowerShell completion for %[1]s
Register-ArgumentCompleter -Native -CommandName %[1]s -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $flags = @(%[2]s)
    $flags | Where-Object { $_ -like "$wordToComplete*" } | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
`, binName, "'"+strings.Join(serverFlags, "', '")+"'")
	}
	return ""
}

// initCommand returns the shell init line that sources the live
// completion script for the given shell.
func initCommand(shell, binName string) string {
	switch shell {
	case "bash", "zsh", "ksh":
		return fmt.Sprintf("source <(%s --shell completions %s)\n", binName, shell)
	case "fish":
		return fmt.Sprintf("%s --shell completions fish | source\n", binName)
	case "sh", "dash":
		return fmt.Sprintf("# %s does not support programmable completion; no init needed.\n", shell)
	case "powershell", "pwsh":
		return fmt.Sprintf("Invoke-Expression (& %s --shell completions %s | Out-String)\n", binName, shell)
	}
	return ""
}

// sanitizeFuncName converts a binary name into a safe shell function
// name fragment.
func sanitizeFuncName(name string) string {
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

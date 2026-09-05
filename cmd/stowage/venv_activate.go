package main

import (
	"fmt"
	"os"
)

// writeActivationScripts writes one activation script per shell family into
// the venv's bin/ directory. Each is self-contained and, when sourced,
// prepends bin/ to PATH, sets STOWAGE_ENV, tweaks the prompt, and defines a
// "deactivate" that undoes all of that - mirroring python -m venv.
func (v *venv) writeActivationScripts() error {
	bin := v.binDir()
	name := v.name()

	scripts := map[string]string{
		"activate":       fmt.Sprintf(posixActivateTemplate, bin, v.Dir, name),
		"Activate.ps1":   fmt.Sprintf(powershellActivateTemplate, bin, v.Dir, name),
		"activate.bat":   fmt.Sprintf(cmdActivateTemplate, bin, v.Dir, name),
		"deactivate.bat": cmdDeactivateTemplate,
	}

	for filename, content := range scripts {
		path := bin + string(os.PathSeparator) + filename
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}

const posixActivateTemplate = `# This file must be used with "source bin/activate" from bash or zsh -
# you cannot run it directly, and it will not work from fish or csh.

deactivate () {
    if [ -n "${_STOWAGE_OLD_PATH:-}" ] ; then
        PATH="$_STOWAGE_OLD_PATH"
        export PATH
        unset _STOWAGE_OLD_PATH
    fi

    if [ -n "${_STOWAGE_OLD_PS1:-}" ] ; then
        PS1="$_STOWAGE_OLD_PS1"
        export PS1
        unset _STOWAGE_OLD_PS1
    fi

    unset STOWAGE_ENV

    if [ ! "${1:-}" = "nondestructive" ] ; then
        unset -f deactivate
    fi
}

# unset any lingering environment first, in case another stowage venv is
# already active
deactivate nondestructive

_STOWAGE_OLD_PATH="$PATH"
PATH="%s:$PATH"
export PATH

STOWAGE_ENV="%s"
export STOWAGE_ENV

if [ -z "${STOWAGE_ENV_DISABLE_PROMPT:-}" ] ; then
    _STOWAGE_OLD_PS1="${PS1:-}"
    PS1="(%s) ${PS1:-}"
    export PS1
fi
`

const powershellActivateTemplate = `function global:deactivate ([switch]$NonDestructive) {
    if (Test-Path variable:_STOWAGE_OLD_PATH) {
        $env:PATH = $variable:_STOWAGE_OLD_PATH
        Remove-Variable "_STOWAGE_OLD_PATH" -Scope global
    }

    if (Test-Path function:_STOWAGE_OLD_PROMPT) {
        Copy-Item function:_STOWAGE_OLD_PROMPT function:prompt
        Remove-Item function:_STOWAGE_OLD_PROMPT
    }

    Remove-Item Env:STOWAGE_ENV -ErrorAction SilentlyContinue

    if (-not $NonDestructive) {
        Remove-Item function:deactivate
    }
}

# unset any lingering environment first, in case another stowage venv is
# already active
deactivate -NonDestructive

New-Variable -Scope global -Name _STOWAGE_OLD_PATH -Value $env:PATH

$env:PATH = "%s;" + $env:PATH
$env:STOWAGE_ENV = "%s"

if (-not $env:STOWAGE_ENV_DISABLE_PROMPT) {
    Copy-Item function:prompt function:_STOWAGE_OLD_PROMPT
    function global:prompt {
        Write-Host "(%s) " -NoNewline -ForegroundColor Green
        _STOWAGE_OLD_PROMPT
    }
}
`

const cmdActivateTemplate = `@echo off
if defined _STOWAGE_OLD_PATH set "PATH=%%_STOWAGE_OLD_PATH%%"
set "_STOWAGE_OLD_PATH=%%PATH%%"
set "PATH=%s;%%PATH%%"
set "STOWAGE_ENV=%s"
if defined _STOWAGE_OLD_PROMPT set "PROMPT=%%_STOWAGE_OLD_PROMPT%%"
set "_STOWAGE_OLD_PROMPT=%%PROMPT%%"
set "PROMPT=(%s) %%PROMPT%%"
`

const cmdDeactivateTemplate = `@echo off
if defined _STOWAGE_OLD_PATH set "PATH=%_STOWAGE_OLD_PATH%"
if defined _STOWAGE_OLD_PATH set "_STOWAGE_OLD_PATH="
if defined _STOWAGE_OLD_PROMPT set "PROMPT=%_STOWAGE_OLD_PROMPT%"
if defined _STOWAGE_OLD_PROMPT set "_STOWAGE_OLD_PROMPT="
set "STOWAGE_ENV="
`

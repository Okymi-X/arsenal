// Package logx provides minimal, leveled logging without color or emoji.
//
// Output is quiet by default. Verbose mode unlocks debug and info lines.
// Status markers are plain ASCII: [ok], [fail], [warn], ->.
package logx

import (
	"fmt"
	"io"
)

// Level controls how much output a Logger emits.
type Level int

const (
	// LevelQuiet emits warnings and errors only.
	LevelQuiet Level = iota
	// LevelVerbose emits info and debug lines in addition to warnings.
	LevelVerbose
)

// Logger writes leveled, unadorned messages to a pair of writers.
//
// Informational output goes to out; warnings and errors go to errw.
// A Logger holds no global state and is safe to construct per command.
type Logger struct {
	level Level
	out   io.Writer
	errw  io.Writer
}

// New returns a Logger writing normal output to out and diagnostics to errw.
func New(level Level, out, errw io.Writer) *Logger {
	return &Logger{level: level, out: out, errw: errw}
}

// Verbose reports whether the logger emits info and debug lines.
func (l *Logger) Verbose() bool { return l.level >= LevelVerbose }

// Debugf prints a debug line only in verbose mode.
func (l *Logger) Debugf(format string, args ...any) {
	if l.level >= LevelVerbose {
		fmt.Fprintf(l.errw, "debug: "+format+"\n", args...)
	}
}

// Infof prints an informational line only in verbose mode.
func (l *Logger) Infof(format string, args ...any) {
	if l.level >= LevelVerbose {
		fmt.Fprintf(l.out, format+"\n", args...)
	}
}

// Printf prints a line regardless of level (primary command output).
func (l *Logger) Printf(format string, args ...any) {
	fmt.Fprintf(l.out, format+"\n", args...)
}

// Warnf prints a warning line regardless of level.
func (l *Logger) Warnf(format string, args ...any) {
	fmt.Fprintf(l.errw, "[warn] "+format+"\n", args...)
}

// Errorf prints an error line regardless of level.
func (l *Logger) Errorf(format string, args ...any) {
	fmt.Fprintf(l.errw, "[fail] "+format+"\n", args...)
}

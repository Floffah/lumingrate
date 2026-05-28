package util

import "strings"

func NormaliseCommand(command string) string {
	return strings.ToLower(strings.TrimSpace(command))
}

func JoinForDisplay(commands []string) string {
	return strings.Join(commands, ", ")
}

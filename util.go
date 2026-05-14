package main

import "fmt"

const (
	SYSTEM = "SYSTEM"
)

func PrintStringErr(name string, pid int, err string) string {
	return fmt.Sprintf("[%s | %d] : error: %s\n", name, pid, err)
}

func PrintErr(name string, pid int, err string) string {
	return fmt.Sprintf("[%s | %d] : error: %s\n", name, pid, err)
}

func PrintLog(name string, pid int, messages ...string) {
	msg := messages[0]

	if len(messages) > 1 {
		for i := 1; i < len(messages); i++ {
			msg = msg + " | " + messages[i]
		}
	}

	fmt.Printf("[%s | %d] :  %s\n", name, pid, msg)
}

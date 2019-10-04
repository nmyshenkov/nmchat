package message

import (
	"errors"
	"strings"
	"time"
)

// Message - message
type Message struct {
	FromNikname string
	From        int
	ToNikname   string
	To          int
	Body        string
}

// FromStringMessage - get full string answer
func (m *Message) FromStringMessage(additional ...string) string {
	return getChatTime() + m.FromNikname + ": " + m.Body + strings.Join(additional, " ")
}

// FromByteMessage - get full byte answer
func (m *Message) FromByteMessage(additional ...string) []byte {
	return []byte(getChatTime() + m.FromNikname + ": " + m.Body + strings.Join(additional, " "))
}

// GetCommendArg - get command argument from string
func GetCommendArg(command, str string) (string, error) {

	// check string contains
	if !strings.Contains(str, command) {
		return "", errors.New("command not found")
	}

	return strings.TrimSpace(strings.Replace(str, command, "", 1)), nil
}

// TrimText = separate command and text
func TrimText(text string) (bool, string, string) {
	strs := strings.SplitN(strings.TrimSpace(text), " ", 2)
	if len(strs) < 2 {
		return false, "", ""
	}

	strs[1] = strings.TrimSpace(strs[1])

	if strs[0][0] == 47 {
		msgTo := strings.Replace(strs[0], "/", "", 1)
		return true, msgTo, strs[1]
	}

	return false, strs[0], strs[1]
}

func getChatTime() string {
	return "[" + time.Now().Format("3:04:05PM") + "] "
}

package main

import "testing"

func TestCmdText(t *testing.T) {
	td := []struct {
		args      []string
		cmd, text string
	}{
		{},
		{[]string{"start", "coding"}, "start", "coding"},
	}
	for _, d := range td {
		t.Run(d.text, func(t *testing.T) {
			actualCmd, actualText := cmdText(d.args)
			if actualCmd != d.cmd {
				t.Errorf("expected cmd %s got %s", d.cmd, actualCmd)
			}
			if actualText != d.text {
				t.Errorf("expected text %s got %s", d.text, actualText)
			}
		})
	}
}

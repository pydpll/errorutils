package errorutils

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCompareOptions(t *testing.T) {
	//set logrus level to debug
	logrus.SetLevel(logrus.DebugLevel)
	tests := []struct {
		target   Option
		template Option
		expected bool
		name     string
	}{
		{WithMsg("this is the target, it has a message"), WithMsg("test"), true, "same"},
		{WithMsg("this is the target, it has a message"), WithExitCode(1), false, "different"},
		{WithMsg("this is the target, it has a message"), WithAltPrint("test"), false, "different"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := compareOptions(tc.target, tc.template)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

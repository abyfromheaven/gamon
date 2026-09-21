package monitor

import (
	"testing"
)

func TestParsePingOutput(t *testing.T) {
	// We can test parsing logic or mock command output if needed.
	// Let's test checking logic or helper parsing if extracted, or test PingOnce behavior.
	// Since PingOnce executes system ping, let's verify it doesn't crash on invalid IPs.
	res := PingOnce("256.256.256.256", 1)
	if res.Status != StatusOffline {
		t.Errorf("Expected invalid IP to be offline, got %s", res.Status)
	}
}

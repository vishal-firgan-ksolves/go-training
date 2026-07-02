package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestCLIConverter(t *testing.T) {

	// Table-driven test simulating keyboard inputs
	tests := []struct {
		name           string
		simulatedInput string
		expectedOutput string
	}{
		{
			name:           "Celsius to Fahrenheit (Boiling)",
			simulatedInput: "1\n100\n", // Simulates typing "1" [ENTER] "100" [ENTER]
			expectedOutput: "100.00°C = 212.00°F",
		},
		{
			name:           "Fahrenheit to Celsius (Freezing)",
			simulatedInput: "2\n32\n",  // Simulates typing "2" [ENTER] "32" [ENTER]
			expectedOutput: "32.00°F = 0.00°C",
		},
		{
			name:           "Invalid Option",
			simulatedInput: "abc\n",    // Simulates bad user input
			expectedOutput: "Error: Invalid option format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command("go", "run", "main.go")
			
			// Inject the simulated keyboard strokes
			cmd.Stdin = strings.NewReader(tc.simulatedInput)

			// Capture whatever prints to the terminal
			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out // Capture errors too

			// Execute the command
			_ = cmd.Run() 

			// Verify the output matches what we expect
			if !strings.Contains(out.String(), tc.expectedOutput) {
				t.Errorf("Expected output to contain '%s', but got:\n%s", tc.expectedOutput, out.String())
			}
		})
	}
}
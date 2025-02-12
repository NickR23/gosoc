package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestMain sets up and tears down the Docker container for all tests.
func TestMain(m *testing.M) {
    wd, err := os.Getwd()
    if err != nil {
        fmt.Println("Failed to get working directory")
        os.Exit(1)
    }

	// Start the container
    cmd := exec.Command("docker", "run", "--rm", "-d",
                        "--name", "gosoc_e2e_test_server",
                        "-v", filepath.Join(wd, "config")+":/config",
                        "-v", filepath.Join(wd, "reports")+":/reports",
                        "-p", "9001:9001",
                        "crossbario/autobahn-testsuite:latest")
	if err := cmd.Run(); err != nil {
		fmt.Println("Failed to start Docker container:", err)
		os.Exit(1)
	}

	// Wait for server to be ready
	time.Sleep(2 * time.Second) // Adjust based on actual server startup time

	// Run the tests
	exitCode := m.Run()

	// Stop the container after tests complete
	exec.Command("docker", "stop", "test_server").Run()

	os.Exit(exitCode)
}

// Example test that runs against the server
func TestClientRequest(t *testing.T) {
	resp, err := myClient.DoRequest("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

	containerName := "gosoc_e2e_test_server"
	cmd := exec.Command("docker", "run", "--rm", "-d",
		"--name", containerName,
		"-v", filepath.Join(wd, "config")+":/config",
		"-v", filepath.Join(wd, "reports")+":/reports",
		"-p", "9001:9001",
		"crossbario/autobahn-testsuite:latest")
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to start Docker container:", err)
	}

	// Wait for server to be ready
	time.Sleep(2 * time.Second)

	// Run the tests
	exitCode := m.Run()

	// Server cleanup
	exec.Command("docker", "stop", "test_server").Run()
	log.Println("Shutting down container:", containerName)
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

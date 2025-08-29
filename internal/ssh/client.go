package ssh

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
	"github.com/slackpass/slackpass/internal/config"
)

// Client handles SSH operations
type Client struct {
	config *config.Config
}

// NewClient creates a new SSH client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
	}
}

// WaitForConnection waits for SSH to become available on the instance
func (c *Client) WaitForConnection(name string) error {
	// Get instance IP address
	ip, port, err := c.getInstanceAddress(name)
	if err != nil {
		return fmt.Errorf("failed to get instance address: %w", err)
	}

	// Wait for SSH port to be open
	timeout := time.Duration(c.config.SSHTimeout) * time.Second
	start := time.Now()

	for time.Since(start) < timeout {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 5*time.Second)
		if err == nil {
			conn.Close()
			// Wait a bit more for SSH service to be fully ready
			time.Sleep(2 * time.Second)
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout waiting for SSH connection to %s", name)
}

// Shell opens an interactive shell session to the instance
func (c *Client) Shell(name string) error {
	ip, port, err := c.getInstanceAddress(name)
	if err != nil {
		return fmt.Errorf("failed to get instance address: %w", err)
	}

	// Use external SSH client for interactive session
	cmd := exec.Command("ssh",
		"-i", c.config.SSHKeyPath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", fmt.Sprintf("%d", port),
		fmt.Sprintf("%s@%s", c.config.SSHUser, ip),
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Exec executes a command on the instance
func (c *Client) Exec(name, command string) error {
	ip, port, err := c.getInstanceAddress(name)
	if err != nil {
		return fmt.Errorf("failed to get instance address: %w", err)
	}

	// Create SSH client connection
	client, err := c.createSSHClient(ip, port)
	if err != nil {
		return fmt.Errorf("failed to create SSH client: %w", err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up I/O
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Run command
	return session.Run(command)
}

// CopyFile copies a file to the instance
func (c *Client) CopyFile(name, localPath, remotePath string) error {
	ip, port, err := c.getInstanceAddress(name)
	if err != nil {
		return fmt.Errorf("failed to get instance address: %w", err)
	}

	// Use SCP for file transfer
	cmd := exec.Command("scp",
		"-i", c.config.SSHKeyPath,
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-P", fmt.Sprintf("%d", port),
		localPath,
		fmt.Sprintf("%s@%s:%s", c.config.SSHUser, ip, remotePath),
	)

	return cmd.Run()
}

// Helper methods

func (c *Client) getInstanceAddress(name string) (string, int, error) {
	// For now, return localhost with a calculated port
	// TODO: Implement proper IP discovery from QEMU/libvirt
	// This is a simplified implementation that assumes port forwarding
	
	// Calculate port based on instance name hash
	port := 2222 // Base port
	for _, char := range name {
		port += int(char)
	}
	port = 2222 + (port % 1000) // Keep in reasonable range

	return "127.0.0.1", port, nil
}

func (c *Client) createSSHClient(host string, port int) (*ssh.Client, error) {
	// Read private key
	key, err := os.ReadFile(c.config.SSHKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SSH key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SSH key: %w", err)
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: c.config.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// Connect to SSH server
	address := fmt.Sprintf("%s:%d", host, port)
	return ssh.Dial("tcp", address, config)
}

// GenerateSSHKey generates an SSH key pair for slackpass
func (c *Client) GenerateSSHKey() error {
	// Check if key already exists
	if _, err := os.Stat(c.config.SSHKeyPath); err == nil {
		return nil // Key already exists
	}

	// Generate SSH key pair using ssh-keygen
	cmd := exec.Command("ssh-keygen",
		"-t", "rsa",
		"-b", "2048",
		"-f", c.config.SSHKeyPath,
		"-N", "", // No passphrase
		"-C", "slackpass@localhost",
	)

	return cmd.Run()
}

// GetPublicKey returns the public key content
func (c *Client) GetPublicKey() (string, error) {
	pubKeyPath := c.config.SSHKeyPath + ".pub"
	data, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read public key: %w", err)
	}
	return string(data), nil
}
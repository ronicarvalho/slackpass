package vm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/slackpass/slackpass/internal/config"
	"github.com/slackpass/slackpass/internal/kvm"
	"github.com/slackpass/slackpass/internal/ssh"
)

// Manager handles virtual machine operations
type Manager struct {
	config    *config.Config
	kvmClient *kvm.Client
	sshClient *ssh.Client
}

// NewManager creates a new VM manager
func NewManager() *Manager {
	cfg := config.Load()
	return &Manager{
		config:    cfg,
		kvmClient: kvm.NewClient(cfg),
		sshClient: ssh.NewClient(cfg),
	}
}

// Launch creates and starts a new virtual machine
func (m *Manager) Launch(config *LaunchConfig) error {
	// Generate name if not provided
	if config.Name == "" {
		config.Name = generateInstanceName()
	}

	// Validate that instance doesn't already exist
	if m.instanceExists(config.Name) {
		return fmt.Errorf("instance '%s' already exists", config.Name)
	}

	fmt.Printf("Launching %s...\n", config.Name)

	// Create instance directory
	instanceDir := filepath.Join(m.config.InstancesDir, config.Name)
	if err := os.MkdirAll(instanceDir, 0755); err != nil {
		return fmt.Errorf("failed to create instance directory: %w", err)
	}

	// Download or prepare image
	imagePath, err := m.prepareImage(config.Image, instanceDir)
	if err != nil {
		return fmt.Errorf("failed to prepare image: %w", err)
	}

	// Create VM configuration
	vmConfig := &kvm.VMConfig{
		Name:      config.Name,
		ImagePath: imagePath,
		CPUs:      config.CPUs,
		Memory:    config.Memory,
		Disk:      config.Disk,
		CloudInit: config.CloudInit,
	}

	// Create and start the VM
	if err := m.kvmClient.Create(vmConfig); err != nil {
		return fmt.Errorf("failed to create VM: %w", err)
	}

	if err := m.kvmClient.Start(config.Name); err != nil {
		return fmt.Errorf("failed to start VM: %w", err)
	}

	// Wait for SSH to be available
	fmt.Printf("Waiting for %s to be ready...\n", config.Name)
	if err := m.sshClient.WaitForConnection(config.Name); err != nil {
		return fmt.Errorf("failed to establish SSH connection: %w", err)
	}

	fmt.Printf("Launched: %s\n", config.Name)
	return nil
}

// List returns all virtual machine instances
func (m *Manager) List() ([]*kvm.Instance, error) {
	return m.kvmClient.List()
}

// Shell opens an interactive shell to the specified instance
func (m *Manager) Shell(name string) error {
	if name == "" {
		// If no name provided, try to connect to the only running instance
		instances, err := m.List()
		if err != nil {
			return err
		}

		running := make([]*kvm.Instance, 0)
		for _, instance := range instances {
			if instance.State == "Running" {
				running = append(running, instance)
			}
		}

		if len(running) == 0 {
			return fmt.Errorf("no running instances found")
		}
		if len(running) > 1 {
			return fmt.Errorf("multiple instances running, please specify name")
		}

		name = running[0].Name
	}

	return m.sshClient.Shell(name)
}

// Exec executes a command on the specified instance
func (m *Manager) Exec(name, command string) error {
	return m.sshClient.Exec(name, command)
}

// Start starts the specified instance
func (m *Manager) Start(name string) error {
	return m.kvmClient.Start(name)
}

// Stop stops the specified instance
func (m *Manager) Stop(name string, force bool) error {
	return m.kvmClient.Stop(name, force)
}

// Delete deletes the specified instance
func (m *Manager) Delete(name string, purge, force bool) error {
	return m.kvmClient.Delete(name, purge, force)
}

// Info displays detailed information about the specified instance
func (m *Manager) Info(name string) error {
	return m.kvmClient.Info(name)
}

// Helper functions

func (m *Manager) instanceExists(name string) bool {
	instanceDir := filepath.Join(m.config.InstancesDir, name)
	_, err := os.Stat(instanceDir)
	return err == nil
}

func (m *Manager) prepareImage(image, instanceDir string) (string, error) {
	// This will be implemented to download/prepare cloud images
	// For now, return a placeholder path
	imagePath := filepath.Join(instanceDir, "disk.qcow2")
	return imagePath, nil
}

func generateInstanceName() string {
	// Generate a random name like "keen-butterfly"
	adjectives := []string{"keen", "bold", "swift", "bright", "calm", "eager"}
	animals := []string{"butterfly", "dolphin", "eagle", "fox", "hawk", "lion"}
	
	adj := adjectives[len(adjectives)%6]
	animal := animals[len(animals)%6]
	
	return fmt.Sprintf("%s-%s", adj, animal)
}
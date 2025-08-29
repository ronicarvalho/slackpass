package images

import (
	"fmt"
	"strings"
)

// Manager handles image operations
type Manager struct {
	images map[string][]*ImageInfo
}

// NewManager creates a new image manager
func NewManager() *Manager {
	return &Manager{
		images: getBuiltinImages(),
	}
}

// Find searches for available images
func (m *Manager) Find(filter string, remoteOnly bool) ([]*ImageInfo, error) {
	var result []*ImageInfo

	for distro, versions := range m.images {
		// Apply filter if specified
		if filter != "" && !strings.Contains(distro, filter) {
			continue
		}

		for _, img := range versions {
			result = append(result, img)
		}
	}

	return result, nil
}

// Download downloads an image to local cache
func (m *Manager) Download(image string) (string, error) {
	// TODO: Implement image downloading
	return "", fmt.Errorf("image downloading not implemented yet")
}

// GetImagePath returns the local path to an image
func (m *Manager) GetImagePath(image string) (string, error) {
	// TODO: Implement image path resolution
	return "", fmt.Errorf("image path resolution not implemented yet")
}

// getBuiltinImages returns the list of supported images
func getBuiltinImages() map[string][]*ImageInfo {
	return map[string][]*ImageInfo{
		"debian": {
			{
				Distribution: "debian",
				Version:      "bookworm",
				Aliases:      "12, latest",
				Description:  "Debian 12 (Bookworm)",
				Architecture: "amd64",
				URL:          "https://cloud.debian.org/images/cloud/bookworm/latest/debian-12-generic-amd64.qcow2",
			},
			{
				Distribution: "debian",
				Version:      "bullseye",
				Aliases:      "11",
				Description:  "Debian 11 (Bullseye)",
				Architecture: "amd64",
				URL:          "https://cloud.debian.org/images/cloud/bullseye/latest/debian-11-generic-amd64.qcow2",
			},
		},
		"fedora": {
			{
				Distribution: "fedora",
				Version:      "39",
				Aliases:      "latest",
				Description:  "Fedora 39",
				Architecture: "amd64",
				URL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/39/Cloud/x86_64/images/Fedora-Cloud-Base-39-1.5.x86_64.qcow2",
			},
			{
				Distribution: "fedora",
				Version:      "38",
				Aliases:      "",
				Description:  "Fedora 38",
				Architecture: "amd64",
				URL:          "https://download.fedoraproject.org/pub/fedora/linux/releases/38/Cloud/x86_64/images/Fedora-Cloud-Base-38-1.6.x86_64.qcow2",
			},
		},
		"almalinux": {
			{
				Distribution: "almalinux",
				Version:      "9",
				Aliases:      "latest",
				Description:  "AlmaLinux 9",
				Architecture: "amd64",
				URL:          "https://repo.almalinux.org/almalinux/9/cloud/x86_64/images/AlmaLinux-9-GenericCloud-latest.x86_64.qcow2",
			},
			{
				Distribution: "almalinux",
				Version:      "8",
				Aliases:      "",
				Description:  "AlmaLinux 8",
				Architecture: "amd64",
				URL:          "https://repo.almalinux.org/almalinux/8/cloud/x86_64/images/AlmaLinux-8-GenericCloud-latest.x86_64.qcow2",
			},
		},
		"rockylinux": {
			{
				Distribution: "rockylinux",
				Version:      "9",
				Aliases:      "latest",
				Description:  "Rocky Linux 9",
				Architecture: "amd64",
				URL:          "https://download.rockylinux.org/pub/rocky/9/images/x86_64/Rocky-9-GenericCloud-Base.latest.x86_64.qcow2",
			},
			{
				Distribution: "rockylinux",
				Version:      "8",
				Aliases:      "",
				Description:  "Rocky Linux 8",
				Architecture: "amd64",
				URL:          "https://download.rockylinux.org/pub/rocky/8/images/x86_64/Rocky-8-GenericCloud-Base.latest.x86_64.qcow2",
			},
		},
		"centos": {
			{
				Distribution: "centos",
				Version:      "stream9",
				Aliases:      "latest",
				Description:  "CentOS Stream 9",
				Architecture: "amd64",
				URL:          "https://cloud.centos.org/centos/9-stream/x86_64/images/CentOS-Stream-GenericCloud-9-latest.x86_64.qcow2",
			},
			{
				Distribution: "centos",
				Version:      "stream8",
				Aliases:      "",
				Description:  "CentOS Stream 8",
				Architecture: "amd64",
				URL:          "https://cloud.centos.org/centos/8-stream/x86_64/images/CentOS-Stream-GenericCloud-8-latest.x86_64.qcow2",
			},
		},
		"opensuse": {
			{
				Distribution: "opensuse",
				Version:      "tumbleweed",
				Aliases:      "latest",
				Description:  "openSUSE Tumbleweed",
				Architecture: "amd64",
				URL:          "https://download.opensuse.org/repositories/Cloud:/Images:/Leap_15.5/images/openSUSE-Leap-15.5-OpenStack.x86_64.qcow2",
			},
			{
				Distribution: "opensuse",
				Version:      "leap",
				Aliases:      "15.5",
				Description:  "openSUSE Leap 15.5",
				Architecture: "amd64",
				URL:          "https://download.opensuse.org/repositories/Cloud:/Images:/Leap_15.5/images/openSUSE-Leap-15.5-OpenStack.x86_64.qcow2",
			},
		},
		"gentoo": {
			{
				Distribution: "gentoo",
				Version:      "latest",
				Aliases:      "current",
				Description:  "Gentoo Linux (Latest)",
				Architecture: "amd64",
				URL:          "https://bouncer.gentoo.org/fetch/root/all/releases/amd64/autobuilds/current-stage3-amd64-openrc/stage3-amd64-openrc-latest.tar.xz",
			},
		},
		"slackware": {
			{
				Distribution: "slackware",
				Version:      "15.0",
				Aliases:      "latest, current",
				Description:  "Slackware Linux 15.0",
				Architecture: "amd64",
				URL:          "https://mirrors.slackware.com/slackware/slackware64-15.0/",
			},
			{
				Distribution: "slackware",
				Version:      "14.2",
				Aliases:      "",
				Description:  "Slackware Linux 14.2",
				Architecture: "amd64",
				URL:          "https://mirrors.slackware.com/slackware/slackware64-14.2/",
			},
		},
	}
}
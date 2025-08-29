package images

// ImageInfo represents information about a cloud image
type ImageInfo struct {
	Distribution string `json:"distribution"` // e.g., "debian"
	Version      string `json:"version"`      // e.g., "bookworm"
	Aliases      string `json:"aliases"`     // e.g., "12, latest"
	Description  string `json:"description"`  // Human readable description
	Architecture string `json:"architecture"` // e.g., "amd64"
	URL          string `json:"url"`          // Download URL
	Checksum     string `json:"checksum"`     // SHA256 checksum
	Size         int64  `json:"size"`         // File size in bytes
	Cached       bool   `json:"cached"`       // Whether image is cached locally
	LocalPath    string `json:"local_path"`   // Local file path if cached
}

// ImageRepository represents a repository of cloud images
type ImageRepository struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Description string            `json:"description"`
	Images      map[string]string `json:"images"` // distro -> base_url mapping
}

// DownloadProgress represents the progress of an image download
type DownloadProgress struct {
	ImageName     string `json:"image_name"`
	TotalBytes    int64  `json:"total_bytes"`
	Downloaded    int64  `json:"downloaded"`
	Percentage    float64 `json:"percentage"`
	Speed         string  `json:"speed"`
	TimeRemaining string  `json:"time_remaining"`
}
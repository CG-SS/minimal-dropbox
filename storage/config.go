package storage

type System string

const (
	FileSystem System = "file_system"
	Memory     System = "memory"
	Nop        System = "nop"
)

type Config struct {
	System     System `envconfig:"STORAGE_SYSTEM" default:"file_system"`
	ManagedDir string `envconfig:"STORAGE_MANAGED_DIR" default:"./tmp"`
	BufferSize int    `envconfig:"STORAGE_BUFFER_SIZE" default:"1024"`
}

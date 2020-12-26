package homed

import (
	"os"
	"path/filepath"
)

// Homed needs to be used to load data efficiently
type Homed struct {
	basePath string

	data map[FileType]map[string]File
}

// New returns a new Homed
func New(basePath string) *Homed {
	return &Homed{
		basePath: basePath,
		data:     map[FileType]map[string]File{},
	}
}

func (h *Homed) filePath(file File) string {
	return filepath.Join(
		h.basePath,
		string(file.FileType()),
		file.FileName()+".yaml",
	)
}

// Add adds data in the cache
func (h *Homed) Add(file File) error {
	t := file.FileType()
	name := file.FileName()

	if _, ok := h.data[t]; !ok {
		h.data[t] = map[string]File{}
	}

	h.data[t][name] = file
	return nil
}

// Get gets data from the cache
func (h *Homed) Get(fileType FileType, name string) (File, error) {
	var file File
	var err error

	_, ok := h.data[fileType]
	if !ok {
		file, err = h.Load(name, fileType)
	}

	if err != nil {
		return nil, err
	}

	if file == nil {
		file, ok = h.data[fileType][name]
		if !ok {
			file, err = h.Load(name, fileType)
		}
	}

	if err != nil {
		return nil, err
	}

	if err := h.Add(file); err != nil {
		return nil, err
	}

	return file, nil
}

// Load loads the file
func (h *Homed) Load(name string, fileType FileType) (File, error) {
	file, err := NewFile(name, fileType)
	if err != nil {
		return nil, err
	}

	data := file.configFormat()
	path := h.filePath(file)
	err = readFile(path, data)
	if err != nil {
		return nil, err
	}

	// Add the file in the cache to avoid loops
	if err := h.Add(file); err != nil {
		return nil, err
	}

	return file, file.fromConfig(data, h)
}

// Save saves a file
func (h *Homed) Save(file File) error {
	config, err := file.toConfig()
	if err != nil {
		return err
	}

	return writeFile(h.filePath(file), true, config)
}

// ListFiles lists the files in a configuration directory
func (h *Homed) ListFiles(fileType FileType) ([]string, error) {
	roomDir := filepath.Join(h.basePath, string(fileType))
	names := []string{}
	return names, filepath.Walk(roomDir, func(path string, f os.FileInfo, err error) error {
		if f == nil || f.IsDir() {
			return nil
		}

		names = append(names, removeExt(f.Name()))
		return nil
	})
}

// Delete deletes a file
func (h *Homed) Delete(file File) error {
	return deleteFile(h.filePath(file))
}

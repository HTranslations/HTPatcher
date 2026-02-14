// Package rgss3a provides reading functionality for RGSS3A archives used by RPG Maker VX Ace.
package rgss3a

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

const header = "RGSSAD\x00\x03"

// FileEntry represents a file entry in the archive.
type FileEntry struct {
	Offset uint32
	Size   uint32
	Key    uint32
	Name   string
}

// Archive represents an opened RGSS3A archive.
type Archive struct {
	file    *os.File
	entries map[string]*FileEntry // normalized path -> entry
}

// Open opens an RGSS3A archive for reading.
func Open(path string) (*Archive, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open archive: %w", err)
	}

	archive := &Archive{
		file:    file,
		entries: make(map[string]*FileEntry),
	}

	if err := archive.readEntries(); err != nil {
		file.Close()
		return nil, err
	}

	return archive, nil
}

// Close closes the archive.
func (a *Archive) Close() error {
	if a.file != nil {
		return a.file.Close()
	}
	return nil
}

// readEntries reads all file entries from the archive.
func (a *Archive) readEntries() error {
	// Verify header
	headerBytes := make([]byte, 8)
	if _, err := a.file.Read(headerBytes); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	if string(headerBytes) != header {
		return fmt.Errorf("invalid RGSS3A header: %q", headerBytes)
	}

	// Read the base key (at offset 8)
	var baseKey uint32
	if err := binary.Read(a.file, binary.LittleEndian, &baseKey); err != nil {
		return fmt.Errorf("failed to read base key: %w", err)
	}
	// Transform the key: key = key * 9 + 3
	baseKey = baseKey*9 + 3

	// Read all entries
	for {
		// For each entry, use a fresh copy of baseKey
		key := baseKey

		// Read offset (4 bytes, XOR with key)
		var encOffset uint32
		if err := binary.Read(a.file, binary.LittleEndian, &encOffset); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read offset: %w", err)
		}
		offset := encOffset ^ key

		// Offset 0 marks end of entries
		if offset == 0 {
			break
		}

		// Read size
		var encSize uint32
		if err := binary.Read(a.file, binary.LittleEndian, &encSize); err != nil {
			return fmt.Errorf("failed to read size: %w", err)
		}
		size := encSize ^ key

		// Read file key
		var encFileKey uint32
		if err := binary.Read(a.file, binary.LittleEndian, &encFileKey); err != nil {
			return fmt.Errorf("failed to read file key: %w", err)
		}
		fileKey := encFileKey ^ key

		// Read name length
		var encNameLen uint32
		if err := binary.Read(a.file, binary.LittleEndian, &encNameLen); err != nil {
			return fmt.Errorf("failed to read name length: %w", err)
		}
		nameLen := encNameLen ^ key

		// Sanity check
		if nameLen > 1000 {
			return fmt.Errorf("invalid name length: %d", nameLen)
		}

		// Read and decrypt name
		nameBytes := make([]byte, nameLen)
		if _, err := io.ReadFull(a.file, nameBytes); err != nil {
			return fmt.Errorf("failed to read name: %w", err)
		}

		// Decrypt name byte by byte using key
		for i := uint32(0); i < nameLen; i++ {
			nameBytes[i] ^= byte(key >> (8 * (i % 4)))
		}

		// Normalize path: convert Windows separators and lowercase for lookup
		name := strings.ReplaceAll(string(nameBytes), "\\", "/")

		entry := &FileEntry{
			Offset: offset,
			Size:   size,
			Key:    fileKey,
			Name:   name,
		}

		// Store with lowercase key for case-insensitive lookup
		a.entries[strings.ToLower(name)] = entry
	}

	return nil
}

// List returns all file paths in the archive.
func (a *Archive) List() []string {
	paths := make([]string, 0, len(a.entries))
	for _, entry := range a.entries {
		paths = append(paths, entry.Name)
	}
	return paths
}

// ListDir returns all files in a specific directory.
func (a *Archive) ListDir(dir string) []string {
	dir = strings.ToLower(strings.TrimSuffix(dir, "/"))
	if dir != "" {
		dir = dir + "/"
	}

	var paths []string
	for key, entry := range a.entries {
		if strings.HasPrefix(key, dir) {
			paths = append(paths, entry.Name)
		}
	}
	return paths
}

// HasFile checks if a file exists in the archive.
func (a *Archive) HasFile(path string) bool {
	_, ok := a.entries[strings.ToLower(path)]
	return ok
}

// ReadFile reads and decrypts a file from the archive into memory.
func (a *Archive) ReadFile(path string) ([]byte, error) {
	entry, ok := a.entries[strings.ToLower(path)]
	if !ok {
		return nil, fmt.Errorf("file not found in archive: %s", path)
	}

	// Read encrypted data
	data := make([]byte, entry.Size)
	if _, err := a.file.ReadAt(data, int64(entry.Offset)); err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	// Decrypt data using file's own key
	fileKey := entry.Key
	for i := uint32(0); i < entry.Size; i++ {
		data[i] ^= byte(fileKey >> (8 * (i % 4)))
		if (i+1)%4 == 0 {
			fileKey = (fileKey*7 + 3) & 0xFFFFFFFF
		}
	}

	return data, nil
}

// GetEntry returns the entry for a file path, or nil if not found.
func (a *Archive) GetEntry(path string) *FileEntry {
	return a.entries[strings.ToLower(path)]
}

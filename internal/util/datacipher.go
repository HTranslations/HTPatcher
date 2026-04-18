package util

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// DataCipher describes the per-file XOR obfuscation that some games inject
// into rmmz_managers.js. The on-disk data/*.json files are XORed with a
// per-file key derived as: K ^ (sum(charCodeAt(basenameWithoutExt)) & 0xFF).
// XOR is symmetric, so encrypt and decrypt are the same operation.
type DataCipher struct {
	K byte
}

var dataCipherKeyRe = regexp.MustCompile(`window\._K\s*=\s*(\d+)`)

// DetectDataCipher inspects rmmz_managers.js for the injected XOR scheme and
// returns the cipher if found. Returns nil if the file is unreadable or the
// scheme is absent.
func DetectDataCipher(jsPath string) *DataCipher {
	data, err := os.ReadFile(filepath.Join(jsPath, "rmmz_managers.js"))
	if err != nil {
		return nil
	}
	m := dataCipherKeyRe.FindSubmatch(data)
	if m == nil {
		return nil
	}
	k, err := strconv.Atoi(string(m[1]))
	if err != nil || k < 0 || k > 255 {
		return nil
	}
	return &DataCipher{K: byte(k)}
}

// FileKey returns the XOR key for a given json filename (basename or full path).
func (c *DataCipher) FileKey(filename string) byte {
	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	var sum int
	for _, r := range name {
		sum += int(r)
	}
	return c.K ^ byte(sum&0xff)
}

// Transform XORs data in place against the per-file key and returns it.
// The input slice is modified.
func (c *DataCipher) Transform(filename string, data []byte) []byte {
	key := c.FileKey(filename)
	for i := range data {
		data[i] ^= key
	}
	return data
}

// ReadDataFile reads a json data file, decrypting it if cipher is non-nil.
func ReadDataFile(path string, cipher *DataCipher) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if cipher != nil {
		cipher.Transform(path, data)
	}
	return data, nil
}

// WriteDataFile writes a json data file, encrypting it if cipher is non-nil.
// The input slice is not modified.
func WriteDataFile(path string, data []byte, cipher *DataCipher) error {
	out := data
	if cipher != nil {
		out = make([]byte, len(data))
		copy(out, data)
		cipher.Transform(path, out)
	}
	return os.WriteFile(path, out, 0644)
}

// IsUnderDir reports whether path lies inside dir (after cleaning, no symlink
// resolution). Returns false on any error or when paths are unrelated.
func IsUnderDir(path, dir string) bool {
	if dir == "" {
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(absDir, absPath)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

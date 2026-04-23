package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WrapperCipher handles the {uid, bid, data} base64+XOR obfuscation used by
// some RPG Maker MZ games (e.g. RJ01515993).
type WrapperCipher struct {
	// Meta maps basename -> preserved uid/bid from the original files.
	Meta map[string]WrapperMeta
}

// WrapperMeta holds the original wrapper metadata for a file.
type WrapperMeta struct {
	UID string
	BID string
}

// DetectWrapperCipher inspects the data directory for JSON files wrapped in
// {"uid": "...", "bid": "...", "data": "<base64>"}. If any are found, it
// returns a WrapperCipher populated with the metadata from all scanned files.
func DetectWrapperCipher(dataPath string) *WrapperCipher {
	entries, err := os.ReadDir(dataPath)
	if err != nil {
		return nil
	}

	cipher := &WrapperCipher{Meta: make(map[string]WrapperMeta)}
	detected := false

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
			continue
		}

		path := filepath.Join(dataPath, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var wrapper struct {
			UID  string `json:"uid"`
			BID  string `json:"bid"`
			Data string `json:"data"`
		}
		if err := json.Unmarshal(data, &wrapper); err != nil {
			continue
		}
		if wrapper.Data == "" {
			continue
		}

		detected = true
		cipher.Meta[entry.Name()] = WrapperMeta{
			UID: wrapper.UID,
			BID: wrapper.BID,
		}
	}

	if !detected {
		return nil
	}

	return cipher
}

// ReadFile reads a wrapped JSON file, decrypts the payload, and returns the
// inner JSON bytes. If the file is not a valid wrapper, it falls back to
// returning the raw file contents.
func (c *WrapperCipher) ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		UID  string `json:"uid"`
		BID  string `json:"bid"`
		Data string `json:"data"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		// Not a wrapper — return raw data.
		return data, nil
	}
	if wrapper.Data == "" {
		return data, nil
	}

	raw, err := base64.StdEncoding.DecodeString(wrapper.Data)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 in wrapper for %s: %w", filepath.Base(path), err)
	}

	decrypted := wrapperDecryptBytes(raw, filepath.Base(path))
	return decrypted, nil
}

// WriteFile encrypts the given JSON bytes and writes them back in the
// {"uid", "bid", "data"} wrapper format, preserving original metadata.
func (c *WrapperCipher) WriteFile(path string, data []byte) error {
	filename := filepath.Base(path)
	meta, ok := c.Meta[filename]
	if !ok {
		meta = WrapperMeta{UID: "00000000", BID: "1.6.0"}
	}

	encrypted := wrapperEncryptBytes(data, filename)
	wrapper := map[string]interface{}{
		"uid":  meta.UID,
		"bid":  meta.BID,
		"data": base64.StdEncoding.EncodeToString(encrypted),
	}

	out, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, 0644)
}

// wrapperFileKey derives the per-file seed key from the filename.
func wrapperFileKey(filename string) byte {
	name := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	t := uint32(0)
	for _, ch := range name {
		t = ((t << 5) - t + uint32(ch)) & 0xFFFFFFFF
	}
	return byte((247 ^ (t & 255)) & 0xFF)
}

// wrapperDecryptBytes decrypts the raw payload bytes.
func wrapperDecryptBytes(data []byte, filename string) []byte {
	fk := wrapperFileKey(filename)
	ls := uint32(fk)

	b := make([]byte, len(data))
	copy(b, data)

	for i := len(b) - 1; i >= 0; i-- {
		_c := uint32(fk ^ 42)
		_m := uint32(i % 128)
		_p := ((ls << 3) ^ (ls >> 2)) & 0xFFFFFFFF
		X := (_c + _m + _p) & 0xFFFFFFFF
		_k := byte(((X ^ 186) + 34) & 0xFF)
		v := b[i] ^ _k
		b[i] = v
		ls = uint32(v)
	}

	return b
}

// wrapperEncryptBytes encrypts the raw payload bytes.
func wrapperEncryptBytes(data []byte, filename string) []byte {
	fk := wrapperFileKey(filename)
	ls := uint32(fk)

	b := make([]byte, len(data))
	copy(b, data)

	for i := len(b) - 1; i >= 0; i-- {
		_c := uint32(fk ^ 42)
		_m := uint32(i % 128)
		_p := ((ls << 3) ^ (ls >> 2)) & 0xFFFFFFFF
		X := (_c + _m + _p) & 0xFFFFFFFF
		_k := byte(((X ^ 186) + 34) & 0xFF)
		plainByte := b[i]
		cipherByte := plainByte ^ _k
		b[i] = cipherByte
		ls = uint32(plainByte)
	}

	return b
}

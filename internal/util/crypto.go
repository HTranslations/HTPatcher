package util

var pngHeader = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52}

// DecryptPng decrypts an RPG Maker encrypted PNG file
func DecryptPng(data []byte) ([]byte, []byte, error) {
	header := data[:32]
	data = append(pngHeader, data[32:]...)
	return header, data, nil
}

// EncryptPng encrypts a PNG file with RPG Maker encryption
func EncryptPng(data []byte, header []byte) []byte {
	return append(header, data[16:]...)
}





package emitter

import (
	"encoding/hex"
	"fmt"
	"os"
)

// WriteHexLines 把 uint32 指令写入文件，每行 8 位大写十六进制（无 0x 前缀）
func WriteHexLines(path string, words []uint32) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, w := range words {
		b := []byte{byte((w >> 24) & 0xff), byte((w >> 16) & 0xff), byte((w >> 8) & 0xff), byte(w & 0xff)}
		s := hex.EncodeToString(b)
		_, err := fmt.Fprintln(f, s)
		if err != nil {
			return err
		}
	}
	return nil
}

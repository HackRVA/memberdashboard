package dbstore

import (
	"fmt"
	"strconv"
	"strings"
)

func encodeRFID(id string) string {
	// This probably doesn't matter, since leading zeros don't change
	// the parsed integer value. -PMW
	if id[0] == '0' {
		// remove the leading zero
		id = id[1:]
	}

	// Parse to 32 bit integer
	i, _ := strconv.ParseInt(id, 10, 64)

	// convert to base16
	idStr := fmt.Sprintf("%08x", i)

	// for some reason the bytes are backwards in the rfid reader
	//  let's reverse the bytes
	// the rfid reader also trims out any zero chars
	return strings.Join(reverse(chunks(idStr, 2)), "")
}

func chunks(s string, chunkSize int) []string {
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string
	chunk := make([]rune, chunkSize)
	len := 0
	for _, r := range s {
		chunk[len] = r
		len++
		if len == chunkSize {
			chunks = append(chunks, string(chunk))
			len = 0
		}
	}
	if len > 0 {
		chunks = append(chunks, string(chunk[:len]))
	}
	return chunks
}

func reverse(arr []string) []string {
	newArr := make([]string, 0, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		newArr = append(newArr, arr[i])
	}
	return newArr
}

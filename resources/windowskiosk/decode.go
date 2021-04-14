package main

import (
	"fmt"
	"strconv"
)

// Each byte needs to be converted to a decimal format
//  Unfortunately, the member server saves the UIDs as a string.
//  What's worse, leading zeros in our byte string are trimmed.  So,
//  We have to reverse this
func hexToDec(s string) int64 {
	output, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		fmt.Println(err)
		return int64(output)
	}
	return int64(output)
}

func decodeID(s string) string {
	var result string

	if len(s) < 8 {
		println("did not get a proper number of digits")
		return ""
	}

	for i, v := range chunks(s, 2) {
		result += strconv.FormatInt(hexToDec(v), 10)
		if i >= (len(s)/2)-1 {
			// don't print the space on the last element
			break
		}
		result += " "
	}

	println()
	println(result)

	return result
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

package main

import (
	"bytes"
	"log"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		log.Fatalln("no arguments given")
	}

	file := args[1]

	contents, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	lines := trimContents(contents)

	translateLinesToByteCode(lines)

}

func translateLinesToByteCode(lines []string) []string {
	lines_wo_symbols, jump_symbols := createSymbolTable(lines)

	byte_code_instructions := make([]string, len(lines_wo_symbols))

	for idx, line := range lines_wo_symbols {
		// TODO variable symbol tables
		byte_code_instructions[idx] = translateLineToByteCode(lines_wo_symbols)
	}

	return byte_code_instructions
}

func translateLineToByteCode(hack_stm string) string {
	return ""
}

func trimContents(contents []byte) []string {
	line_bytes := bytes.Split(contents, []byte{'\n'})

	lines := make([]string, len(line_bytes))

	for idx, val := range line_bytes {
		raw_line := string(val)
		lines[idx] = trimLine(raw_line)
	}

	return lines
}

// Removes whitespace and comments
func trimLine(raw_line string) string {
	before, _, _ := strings.Cut(raw_line, "//")
	return strings.ReplaceAll(before, " ", "")
}

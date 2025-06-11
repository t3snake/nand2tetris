package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
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

	file_wo_ext, found := strings.CutSuffix(file, ".asm")

	if !found {
		log.Fatalln("Invalid extension, use .asm files")
	}

	lines := trimContents(contents)
	fmt.Printf("Trimmed:\n%s\n", lines)

	bytecode := translateLinesToByteCode(lines)
	bytecode_content := []byte(strings.Join(bytecode, string('\n')))

	file_out := fmt.Sprintf("%s.hack", file_wo_ext)
	err = os.WriteFile(file_out, bytecode_content, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("bytecode written to %s", file_out)
}

func translateLinesToByteCode(lines []string) []string {
	variable_symbol_table := make(map[string]string, 0)
	lines_wo_symbols, label_symbol_table := createSymbolTable(lines)

	fmt.Printf("lines without labels:\n%s\n", lines_wo_symbols)
	fmt.Printf("label symbols:\n%s\n", label_symbol_table)

	byte_code_instructions := make([]string, len(lines_wo_symbols))

	for idx, line := range lines_wo_symbols {
		byte_code_instructions[idx] = translateLineToByteCode(line, &variable_symbol_table, label_symbol_table)
	}

	return byte_code_instructions
}

func translateLineToByteCode(
	hack_stm string,
	variable_symbol_table *map[string]string,
	label_symbol_table map[string]string,
) string {
	return ""
}

// Create label symbol table with predefined and user defined symbols.
func createSymbolTable(lines []string) (lines_wo_labels []string, label_symbol_table map[string]string) {
	label_symbol_table = makeLabelSymbolTable()
	lines_wo_labels = make([]string, 0)

	current_address := 17

	for _, line := range lines {
		if len(line) == 0 {
			continue
		} else if line[0] == '(' {
			temp := strings.ReplaceAll(line, "(", "")
			label := strings.ReplaceAll(temp, ")", "")

			label_symbol_table[label] = strconv.Itoa(current_address)
			current_address++
		} else {
			lines_wo_labels = append(lines_wo_labels, line)
		}
	}

	return lines_wo_labels, label_symbol_table
}

// Make label symbol table and prepopulate with predefined symbols
func makeLabelSymbolTable() map[string]string {
	symbol_table := make(map[string]string, 23)

	// Predefine R0 through R15
	for idx := range 16 {
		symbol_table["R"+strconv.Itoa(idx)] = strconv.Itoa(idx)
	}

	symbol_table["SCREEN"] = "16384"
	symbol_table["KBD"] = "24576"

	symbol_table["SP"] = "0"
	symbol_table["LCL"] = "1"
	symbol_table["ARG"] = "2"
	symbol_table["THIS"] = "3"
	symbol_table["THAT"] = "4"

	return symbol_table
}

// Removes whitespace and comments from given file contents.
func trimContents(contents []byte) []string {
	line_bytes := bytes.Split(contents, []byte{'\n'})

	lines := make([]string, len(line_bytes))

	for idx, val := range line_bytes {
		raw_line := string(val)
		lines[idx] = trimLine(raw_line)
	}

	return lines
}

// Removes whitespace and comments from given line.
func trimLine(raw_line string) string {
	before, _, _ := strings.Cut(raw_line, "//")
	return strings.ReplaceAll(before, " ", "")
}

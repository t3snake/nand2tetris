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
	lines_wo_symbols, symbol_table := createSymbolTable(lines)

	fmt.Printf("lines without labels:\n%s\n", lines_wo_symbols)
	fmt.Printf("label symbols:\n%s\n", symbol_table)

	byte_code_instructions := make([]string, len(lines_wo_symbols))

	new_memory_address := 16

	for idx, line := range lines_wo_symbols {
		byte_code_instructions[idx] = translateLineToByteCode(line, symbol_table, &new_memory_address)
	}

	return byte_code_instructions
}

func translateLineToByteCode(
	hack_stm string,
	symbol_table map[string]string,
	cur_mem *int,
) string {
	after, isATypeAddr := strings.CutPrefix(hack_stm, "@")

	if isATypeAddr {
		address, ok := symbol_table[after]
		if ok {
			return get16DigitBinary(address)
		}
		// if new symbol add to table with current address
		symbol_table[after] = strconv.Itoa(*cur_mem)
		*cur_mem++
		return get16DigitBinary(symbol_table[after])
	}
	// c instruction

	return ""
}

func get16DigitBinary(address_str string) string {
	address, err := strconv.Atoi(address_str)
	if err != nil {
		log.Fatal("Expected string got: " + address_str)
	}
	return fmt.Sprintf("%016b", address)
}

// Create symbol table with predefined symbols and user defined labels.
func createSymbolTable(lines []string) (lines_wo_labels []string, symbol_table map[string]string) {
	symbol_table = makeLabelSymbolTable()
	lines_wo_labels = make([]string, 0)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		} else if line[0] == '(' {
			temp := strings.ReplaceAll(line, "(", "")
			label := strings.ReplaceAll(temp, ")", "")

			// if there were 4 statements previously before label, 0 1 2 3
			// and label would point to 4 thus len of lines_wo_labels
			symbol_table[label] = strconv.Itoa(len(lines_wo_labels))
		} else {
			lines_wo_labels = append(lines_wo_labels, line)
		}
	}

	return lines_wo_labels, symbol_table
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

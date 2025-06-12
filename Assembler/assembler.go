package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var destination_bit_map map[string]string = map[string]string{
	"null": "000",
	"M":    "001",
	"D":    "010",
	"MD":   "011",
	"A":    "100",
	"AM":   "101",
	"AD":   "110",
	"AMD":  "111",
}

var jump_bit_map map[string]string = map[string]string{
	"null": "000",
	"JGT":  "001",
	"JEQ":  "010",
	"JGE":  "011",
	"JLT":  "100",
	"JNE":  "101",
	"JLE":  "110",
	"JMP":  "111",
}

var computation_bit_map map[string]string = map[string]string{
	"0":   "101010",
	"1":   "111111",
	"-1":  "111010",
	"D":   "001100",
	"A":   "110000",
	"!D":  "001101",
	"!A":  "110001",
	"-D":  "001111",
	"-A":  "110011",
	"D+1": "011111",
	"A+1": "110111",
	"D-1": "001110",
	"A-1": "110010",
	"D+A": "000010",
	"D-A": "010011",
	"A-D": "000111",
	"D&A": "000000",
	"D|A": "010101",
}

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
		_, err := strconv.Atoi(after)
		if err != nil {
			// if not numeric, it must be new symbol
			// if new symbol add to table with current address
			symbol_table[after] = strconv.Itoa(*cur_mem)
			*cur_mem++
			return get16DigitBinary(symbol_table[after])
		}
		// after must be numeric if no err then directly convert
		return get16DigitBinary(after)
	}
	// C instruction
	// dest = comp ; jmp
	// 111 + a bit + 6 computation bits + 3 destination bits + 3 jump bits
	var a_bit string
	var comp_bits string
	var dest_bits string
	var jump_bits string

	var hack_stm_wo_dest_bits string

	before, after, found := strings.Cut(hack_stm, "=")
	if found {
		dest_bits = destination_bit_map[before]
		hack_stm_wo_dest_bits = after
	} else {
		dest_bits = destination_bit_map["null"]
		hack_stm_wo_dest_bits = before
	}

	before, after, found = strings.Cut(hack_stm_wo_dest_bits, ";")

	if strings.ContainsAny(before, "M") {
		before = strings.ReplaceAll(before, "M", "A")
		a_bit = "1"
	} else {
		a_bit = "0"
	}

	comp_bits = computation_bit_map[before]

	if found {
		jump_bits = jump_bit_map[after]
	} else {
		jump_bits = jump_bit_map["null"]
	}

	return fmt.Sprintf("111%s%s%s%s", a_bit, comp_bits, dest_bits, jump_bits)
}

func get16DigitBinary(address_str string) string {
	address, err := strconv.Atoi(address_str)
	if err != nil {
		log.Fatal("Expected string got: " + address_str)
	}
	return fmt.Sprintf("0%015b", address)
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

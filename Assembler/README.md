# Hack Assembler

This is a simple assembler for the Hack computer. It translates `.asm` files (Hack assembly language) into `.hack` files (binary machine code).

[Book Chapter 5 Computer Architecture](https://www.nand2tetris.org/_files/ugd/44046b_552ed0898d5d491aabafd8a768a87c6f.pdf)  

## Build

```sh
go build assembler.go
```

## Usage

Run the assembler with a `.asm` file as input:

```sh
./assembler Rect.asm
```

This will generate a `Rect.hack` file containing the corresponding bytecode.

- Input: `Rect.asm` (Hack assembly)
- Output: `Rect.hack` (Hack machine code)

The output file can be run on the Hack computer that was made in the previous weeks.  
This can be played in the web ide.   
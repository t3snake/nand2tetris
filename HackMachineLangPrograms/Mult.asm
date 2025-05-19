// Mult - multiplies ram[0] and ram[1] and produces output at ram[2]

// Pseudo
// num0 = ram[0]
// num1 = ram[1]
// loop 1 to num0
// add num1 each time to accumulator

// init num0
@0
D=M
@num0
M=D

// init num1
@1
D=M
@num1
M=D

// init sum
@2
M=0

@i
M=1

// loop while i - num0 < 0 else end
(LOOP)
@i
D=M
@num0
D=D-M
@END
D;JGT

// add num1 to sum on ram[2]
@num1
D=M
@2
M=D+M

// i++
@i
M=M+1

@LOOP
0;JMP


(END)
@END
0;JMP
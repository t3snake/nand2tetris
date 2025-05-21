// Runs an infinite loop that listens to the keyboard input. 
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel. When no key is pressed, 
// the screen should be cleared.

// Pseudo
// When KBD register is not 0 
//    fill all screen registers with -1 : All black pixels
// else
//    fill all screen registers with 0 : All white pixels
// All screen register -> @SCREEN loop until 


// store screen address end
@8192
D=A
@addrend
M=D

(LOOP)
@i
M=0

// Jump based on keyboard input
@KBD
D=M
@WHITE
D;JEQ
@BLACK
D;JNE

(BLACK)
// if i - end == 0 of screen address jump to LOOP
@i
D=M
@addrend
D=D-M
@LOOP
D;JEQ

// set ram[SCREEN + i] to -1 setting 16 pixels black
@i
D=M
@SCREEN
A=D+A
M=-1

// i++
@i
M=M+1

@BLACK
0;JMP

(WHITE)
// if i - end == 0 of screen address jump to LOOP
@i
D=M
@addrend
D=D-M
@LOOP
D;JEQ

// set ram[SCREEN + i] to -1 setting 16 pixels black
@i
D=M
@SCREEN
A=D+A
M=0

// i++
@i
M=M+1

@WHITE
0;JMP
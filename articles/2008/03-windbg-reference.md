Id: 1096
Title: Windbg reference
Tags: debugging,windbg,win32
Date: 2008-03-13T07:48:23-07:00
Format: Markdown
Deleted: yes
--------------
```
k, kb, kp, kP

bp - set breakpoint
~1bp - set breakpoint only on thread 1
bl - list breakpoints
bd $n - disable breakpoint
bc - clear breakpoint
ba w4 $addr - set write breakpoint for 4 bytes at $addr

dv - display all local variables
dv /i - show symbol type
dv /V - show location of variable
dv /V argv - display address of argv variable

dt - display type
dt MyClass $addr - display $addr as type MyClass
dc @ecx l4 - display 4 longs at address ecx, showing hex + ascii
dd $addr $addr-range - display as hex
du $addr - display as unicode string
da $addr - display as ascii string
db $addr l4 - display 4 bytes as bytes
dyb $addr l4 - display as binary
df $addr l4 - display as float
dw $addr l4 - display as words (16-bit)
dW $addr l4 - display as 4 words + ascii

dps $addr l8 - try to resolve 8 pointer values at $addr as symbols
dpu $addr L4 - try to resolve 4 pointer values at $addr as unicode
strings

ln $addr - list symbols near address

g - go
g $addr - go until $addr is reached
gu - go until exit from function
~0 gu - like gu but only current thread is executing
t - trace (execute one assembly instruction)
p - like t but steps over functions
pc - trace until next call instruction
wt - like p but also gathers stats about the function (how many
instructions, what has been called)

.logopen c:\log.txt
.logclose

(address [$addr] - what is the type of $addr. . is $pc, @esp
peb, teb - show process/thread environment block
gle - show last error (like GetLastError()) and NTSTATUS
```

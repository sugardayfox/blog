Id: 1289
Title: Gdb basics
Tags: gdb,unix,programming
Date: 2006-09-06T17:00:00-07:00
Format: Markdown
--------------
**Customization**: File `.gdbinit` contains instructions that gdb will execute
at startup. Helpful for adding commonly used commands. Gdb first reads `.gdbinit`
in $HOME directory, then in current directory (if different than $HOME).

Gdb looks for source files in `$cdir` (directory embedded in executable
recorded during compilation) and `$cwd` (current working directory). You can
add to this list via `dir` e.g. `dir ..\..\WebCore\platform`. See current
list via `show dir`. If you need this, it's a perfect candidate for putting
inside `.gdbinit` file.

Basics:

* `run [<args>]` - run the program with (optional) arguments
* `c` (`continue`) - continue execution
* `bt` - print a stack trace
* `frame <n>` - switch to frame <n>
* `info locals` - show local variables
* `p <foo>` - print the value of a variable
* `x /fmt addr` - show data under a given addr using a given format
* `info breaks` - display current breakpoints
* `break <line-no>` - set breakpoint at a line number <line-no> in current file e.g. `break 24`
* `break *<addr>` - set a breakpoint at a given addr e.g. `break *foo + 24` (set a breakpoint 24 bytes
   after the beginning of function foo
* `delete <n>`, `disable <n>`, `enable <n>` - delete/disable/enable breakpoint <n>

Debugging at assembly level:

* `display /i $eip` so that gdb prints the next assembly instruction
* `nexti` and `stepi` for stepping by one instruction
* `set disassembly-flavor intel` changes assembly syntax from awful AT&T to less awful intel
* `info regs` - show content of registers

Useful macros to define in `.gdbinit`:

`dpc [<count>]` disassembles next <count> (or 24 if not given) bytes starting from current location.
```
define dpc
  if $argc ==1
    disass $pc $pc + $arg0
  end
  if $argc == 0
    disass $pc $pc+24
  end
end
```

`pu <addr>` print Unicode string under address `addr`
```
def pu
  set $uni = $arg0
  set $i = 0
  while (*$uni && $i++<100)
    if (*$uni < 0x80)
      print *(char*)$uni++
    else
      print /x *(short*)$uni++
    end
  end
end
```

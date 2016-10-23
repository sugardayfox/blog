Id: 301
Title: Basics of writing DOS .bat batch files
Tags: reference
Date: 2006-01-13T16:00:00-08:00
Format: Markdown
--------------
**Set a variable**: `SET FOO="bar"`

**Refer to a variable**: `echo %FOO%`

Command line arguments are variables `%1`, `%2` etc.

**Check if a variable is defined**: `IF NOT DEFINED FOO SET FOO="bar"`. Useful
for testing command line arguments.

**Check a result code from the last executed command**: `IF ERRORLEVEL $n $CMD-TO-EXECUTE`. It
executes the command if result code is **equal or greater** to $n, so to
check for any failures use 1.

**Do things in a loop**: `FOR /L %A IN (1,1,10) DO @echo hello`

**Execute another script**: `call another.bat`

**Check if directory exists**: `IF EXIST E:\NUL GOTO USE_E`

Here's an example of a batch file that does a few common things:

```bat
@ECHO OFF
@rem "pushd $dir" puts $dir on directory stack
pushd .
SET COMMENT=%1
IF NOT DEFINED COMMENT SET COMMENT=""

SET FOO="bar"

IF %FOO%=="bar" goto IS_BAR
echo The above should always be true
goto END

:IS_BAR
nmake -f Makefile.vc
IF ERRORLEVEL 1 goto ERR_ONE_OR_HIGHER
echo Compilation successful!
goto END

:ERR_ONE_OR_HIGHER
echo Compilation failed!
goto END

:END
@rem popd pops the directory name from the stack and does cd there
popd
```

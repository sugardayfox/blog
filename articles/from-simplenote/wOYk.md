Id: wOYk
Title: Advanced command execution in Go with os/exec
Format: Markdown
Tags: for-blog, go
Date: 2017-06-12T07:23:02Z
--------------
Go has excellent support for executing external programs. Let's start at the beginning.

## Running a command and capturing the output

Here's the simplest way to run `ls -lah` and capture its combined stdout/stderr.

```go
func main() {
	cmd := exec.Command("ls", "-lah")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
```

Full example: [advanced-exec/01-simple-exec.go](https://github.com/kjk/go-cookbook/blob/master/advanced-exec/01-simple-exec.go).

## Capture stdout and stderr separately

What if you want to do the same but capture stdout and stderr separately?

```go
func main() {
	cmd := exec.Command("ls", "-lah")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}
```

Full example: [advanced-exec/02-capture-stdout-stderr.go](https://github.com/kjk/go-cookbook/blob/master/advanced-exec/02-capture-stdout-stderr.go).

## Capture output but also show progress

What if the command takes a long time to finish? It would be nice to both capture stdout/stderr but also show the output of the program as it being generated (as opposed to dumping it at the very end).

It's a little bit more involved, but not terribly so.

```go
func copyAndCapture(w io.Writer, r io.Reader) []byte {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			break
		}
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			os.Stdout.Write(d)
		}
	}
	return out
}

func main() {
	cmd := exec.Command("ls", "-lah")
	var stdout, stderr []byte
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	cmd.Start()

	go func() {
		stdout = copyAndCapture(os.Stdout, stdoutIn)
	}()

	go func() {
		stderr = copyAndCapture(os.Stderr, stderrIn)
	}()

	err := cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout), string(stderr)
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
```

Full example: [advanced-exec/03-live-progress-and-capture.go](https://github.com/kjk/go-cookbook/blob/master/advanced-exec/03-live-progress-and-capture.go).

## Capture output but also show progress #2

Previous solution works but `copyAndCapture` looks like we're re-implementing `io.Copy` . Thanks to Go interfaces we can re-use `io.Copy` . We'll write `CapturingPassThroughWriter` struct implementing `io.Writer` interface. It'll capture everything that's written to it and also write it to underlying `io.Writer` . Possibly `CapturingPassThroughWriter` can be used in other contexts.

```go
/ CapturingPassThroughWriter is a writer that remembers
// data written to it and passes it to w
type CapturingPassThroughWriter struct {
	buf bytes.Buffer
	w io.Writer
}

// NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
func NewCapturingPassThroughWriter(w io.Writer) *CapturingPassThroughWriter {
	return &CapturingPassThroughWriter{
		w: w,
	}
}

func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
	w.buf.Write(d)
	return w.w.Write(d)
}

// Bytes returns bytes written to the writer
func (w *CapturingPassThroughWriter) Bytes() []byte {
	return w.buf.Bytes()
}

func main() {
	cmd := exec.Command("ls", "-lah")
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := NewCapturingPassThroughWriter(os.Stdout)
	stderr := NewCapturingPassThroughWriter(os.Stderr)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		io.Copy(stdout, stdoutIn)
	}()

	go func() {
		io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}
```

Full example: [advanced-exec/04-live-progress-and-capture-v2.go](https://github.com/kjk/go-cookbook/blob/master/advanced-exec/04-live-progress-and-capture-v2.go).

## Changing environment of executed program

Things you need to know about using of environment variables in Go:
* `os.Environ()` returns `[]string` where each string is in form of `FOO=bar`, where `FOO` is the name of environment variable and `bar` is the value
* `os.Getenv("FOO")` returns the value of environment variable.

Sometimes you need to modify the environment of the executed program.

You do it by setting `Env` member of `exec.Cmd` in the same format as `os.Environ()`. Usually you don't want to construct a completely new environment but pass your own environment augmented with more variables:

```go
	cmd := exec.Command("programToExecute")
	additionalEnv := "FOO=bar"
	newEnv := append(os.Environ(), additionalEnv))
	cmd.Env = newEnv
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("%s", out)
```

Full example: [advanced-exec/06-change-environment.go](https://github.com/kjk/go-cookbook/blob/master/advanced-exec/06-change-environment.go).


## Check early that a program is installed

Imagine you wrote a program that takes a long time to run. You call `foo` executable at the end to perform some essential task.

If `foo` executable is not present on your computer, the call will fail.

It's a good idea to detect that at the beginning of the program, not after doing a lot of work.


You can do it using `exec.LookPath`.

```go
func checkLsExists() {
	path, err := exec.LookPath("ls")
	if err != nil {
		fmt.Printf("didn't find 'ls' executable\n")
		return
	}
	fmt.Printf("'ls' executable is in '%s'\n", path)
}
```

Full example: [advanced-exec/05-check-exe-exists.go](https://github.com/kjk/go-cookbook/blob/master/advanced-exec/05-check-exe-exists.go).

In a real program you would call it at the beginning. If the program couldn't be found you would inform the user with descriptive error message and exit

Another way to check if program exists is to try to execute in a no-op mode (e.g. many programs support `--help` option).

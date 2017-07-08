Id: vEja
Title: Embedding build number in Go executable
Format: Markdown
Tags: for-blog, go, draft
CreatedAt: 2017-07-07T09:01:20Z
UpdatedAt: 2017-07-07T09:11:18Z
--------------
@header-image gfx/headers/header-10.jpg
@collection go-cookbook
@description How to embed build number in Go executable.
@status draft

So you've deployed your web application to production and it's running on a server far away.

When debugging problems it's good to know what version of the code is running.

If you're using git, that would be sha1 of the revision used to build the program.

We can embed that version in the executable during build thanks Go linker's `-X` option. It allows to externally set a value of a variable inside the executable.

In our Go program we would have:
```go
package main

var (
	sha1ver string // sha1 revision used to build the program
	buildTime string // when the executable was built
)
```

We can set that variable to sha1 of the git revision in our build script:
```sh
#!/bin/bash

now=$(date +'%Y-%m-%d_%T')
go build -ldflags "-X main.sha1ver=`git rev-parse HEAD` -X main.buildTime=$now"
```

// TODO: windows example

Let's deconstruct:
* `git rev-parse HEAD` returns sha1 of the current revision, e.g. `e5ce06c1f604efb1de91d515d5de865e7e164d59`
* `-X main.sha1ver=${foo}` tells Go linker to set variable `sha1ver` in package `main` to `${foo}`
* `-ldflags "${flags}"` tells Go build tool to pass `${flags}` to go linker

We also need an easy way to see that version. We can add `-version` cmd-line flag to print it out:
```go
var (
	flgVersion bool
)

func parseCmdLineFlags() {
	flags.BoolVar(&flgVersion, "version", false, "if true, print version and exit")
  flags.Parse()
  if flgVersion {
  	fmt.Printf("Build from sha1 %s on %s\n", sha1ver, buildTime)
    os.Exit(0)
  }
}
```

If this is a web application, we can additionally add a debug page that would show the version. I often do it like that:
```go

func servePlainText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

// /app/debug
func handleDebug(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("url: %s %s", r.Method, r.RequestURI)
	a := []string{s}

	s = "https: no"
	if r.TLS != nil {
		s = "https: yes"
	}
	a = append(a, s)

	s = fmt.Sprintf("RemoteAddr: %s", r.RemoteAddr)
	a = append(a, s)

	a = append(a, "Headers:")
	for k, v := range r.Header {
		if len(v) == 0 {
			a = append(a, k)
		} else if len(v) == 1 {
			s = fmt.Sprintf("  %s: %v", k, v[0])
			a = append(a, s)
		} else {
			a = append(a, "  "+k+":")
			for _, v2 := range v {
				a = append(a, "    "+v2)
			}
		}
	}

	a = append(a, "")
	a = append(a, fmt.Sprintf("ver: https://github.com/kjk/blog/commit/%s", sha1ver))
	a = append(a, fmt.Sprintf("built on: %s", buildTime)

	s = strings.Join(a, "\n")
	servePlainText(w, s)
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
  
	mux.HandleFunc("/app/debug", handleDebug)
	...

	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}
```

As you can see, in addition to printing code version I also show HTTP headers. During debugging, the more information, the better.

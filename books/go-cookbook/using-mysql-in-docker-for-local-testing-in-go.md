---
Id: w4re
Title: Using MySQL in Docker for local testing In Go
Format: Markdown
Tags: go
CreatedAt: 2017-06-12T06:16:54Z
UpdatedAt: 2017-07-09T01:50:59Z
PublishedOn: 2017-07-23
HeaderImage: gfx/headers/header-11.jpg
Collection: go-cookbook
Description: How and why to run MySQL in Docker when developing Go web apps locally.
---

Imagine you’re writing a web application that uses MySQL. You [deploy on Linux](/article/5/blueprint-for-deploying-web-apps-on-coreos.html) but code and test on Mac.

What is a good way to setup MySQL database for local developement and testing on Mac?

You can install MySQL on Mac using [official MySQL installer](https://dev.mysql.com/downloads/mysql/) or via [Homebrew](https://brew.sh/) but my favorite way is to use [docker](https://store.docker.com/editions/community/docker-ce-desktop-mac).

Docker better isolates MySQL from the rest of the system, which has  a couple of advantages:
* it's easier to install the exact version of MySQL that is running in production as there is a docker image for every version
* you don't have to worry that `brew upgrade` will upgrade MySQL. Auto-upgrade is desired for most software but not a database. You need to remember to use `brew pin` to disable that
* you can run several different versions of MySQL for different projects
* since it's running on Linux, it's closer to the code running in production

There is a downside to using Docker:

* you have to make sure that the database container is running
* you need to know the ip address of docker vm running the container

Doing this manually would be annoying and I like to automate so I wrote re-usable bit of Go code to do that.

It’s just a matter of running docker commands and parsing their outputs but it’s difficult enough to worth sharing the complete solution.

Conceptually, what we do is:

* run `docker ps -a` and parse the output
* if the container is not running at all, start it with `docker run`
* if the container is stopped, re-start it with `docker start`
* if the container is already running, extract ip address/port from the output

MySQL database is stored in a local directory mounted by the container. That way data persists even if the container is stopped.

The code is re-usable. You can customize it by changing:

* the base MySQL container. In my case it's [`mysql:5.6`](https://hub.docker.com/_/mysql/)
* where MySQL data is stored
* name of the container, which should by unique to the project
* port on which the database is exposed locally. In the container MySQL listens on standard port 3306. It must be exposed locally on a unique port

It can be adapted for other databases, like PostgreSQL.

We should only start docker when running locally. In my software I use cmd-line flag `-production` to distinguish between running in production and locally.

In production I would use the hard-coded host/ip of MySQL server. Locally I would call `startLocalDockerDbMust()` go get them.

```go
const (
	dockerStatusExited  = "exited"
	dockerStatusRunning = "running"
)

var (
	// using https://hub.docker.com/_/mysql/
	// to use the latest mysql, use mysql:8
	dockerImageName = "mysql:5.6"
	// name must be unique across containers runing on this computer
	dockerContainerName = "mysql-db-multi"
	// where mysql stores databases. Must be on local disk so that
	// database outlives the container
	dockerDbDir = "~/data/db-multi"
	// 3306 is standard MySQL port, I use a unique port to be able
	// to run multiple mysql instances for different projects
	dockerDbLocalPort = "7200"
)

type containerInfo struct {
	id       string
	name     string
	mappings string
	status   string
}

func quoteIfNeeded(s string) string {
	if strings.Contains(s, " ") || strings.Contains(s, "\"") {
		s = strings.Replace(s, `"`, `\"`, -1)
		return `"` + s + `"`
	}
	return s
}

func cmdString(cmd *exec.Cmd) string {
	n := len(cmd.Args)
	a := make([]string, n, n)
	for i := 0; i < n; i++ {
		a[i] = quoteIfNeeded(cmd.Args[i])
	}
	return strings.Join(a, " ")
}

func runCmdWithLogging(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func decodeContainerStaus(status string) string {
	// convert "Exited (0) 2 days ago" into statusExited
	if strings.HasPrefix(status, "Exited") {
		return dockerStatusExited
	}
	// convert "Up <time>" into statusRunning
	if strings.HasPrefix(status, "Up ") {
		return dockerStatusRunning
	}
	return strings.ToLower(status)
}

// given:
// 0.0.0.0:7200->3306/tcp
// return (0.0.0.0, 7200) or None if doesn't match
func decodeIPPortMust(mappings string) (string, string) {
	parts := strings.Split(mappings, "->")
	panicIf(len(parts) != 2, "invalid mappings string: '%s'", mappings)
	parts = strings.Split(parts[0], ":")
	panicIf(len(parts) != 2, "invalid mappints string: '%s'", mappings)
	return parts[0], parts[1]
}

func dockerContainerInfoMust(containerName string) *containerInfo {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Status}}|{{.Ports}}|{{.Names}}")
	outBytes, err := cmd.CombinedOutput()
	panicIfErr(err, "cmd.CombinedOutput() for '%s' failed with %s", cmdString(cmd), err)
	s := string(outBytes)
	// this returns a line like:
	// 6c5a934e00fb|Exited (0) 3 months ago|0.0.0.0:7200->3306/tcp|mysql-db-multi
	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		parts := strings.Split(line, "|")
		panicIf(len(parts) != 4, "Unexpected output from docker ps:\n%s\n. Expected 4 parts, got %d (%v)\n", line, len(parts), parts)
		id, status, mappings, name := parts[0], parts[1], parts[2], parts[3]
		if containerName == name {
			return &containerInfo{
				id:       id,
				status:   decodeContainerStaus(status),
				mappings: mappings,
				name:     name,
			}
		}
	}
	return nil
}

// returns host and port on which database accepts connection
func startLocalDockerDbMust() (string, string) {
	// docker must be running
	cmd := exec.Command("docker", "ps")
	err := cmd.Run()
	panicIfErr(err, "docker must be running! Error: %s", err)
	// ensure directory for database files exists
	dbDir := expandTildeInPath(dockerDbDir)
	err = os.MkdirAll(dbDir, 0755)
	panicIfErr(err, "failed to create dir '%s'. Error: %s", err)
	info := dockerContainerInfoMust(dockerContainerName)
	if info != nil && info.status == dockerStatusRunning {
		return decodeIPPortMust(info.mappings)
	}
	// start or resume container
	if info == nil {
		// start new container
		volumeMapping := dockerDbDir + "s:/var/lib/mysql"
		dockerPortMapping := dockerDbLocalPort + ":3306"
		cmd = exec.Command("docker", "run", "-d", "--name"+dockerContainerName, "-p", dockerPortMapping, "-v", volumeMapping, "-e", "MYSQL_ALLOW_EMPTY_PASSWORD=yes", "-e", "MYSQL_INITDB_SKIP_TZINFO=yes", dockerImageName)
	} else {
		// start stopped container
		cmd = exec.Command("docker", "start", info.id)
	}
	runCmdWithLogging(cmd)

	// wait max 8 seconds for the container to start
	for i := 0; i < 8; i++ {
		info := dockerContainerInfoMust(dockerContainerName)
		if info != nil && info.status == dockerStatusRunning {
			return decodeIPPortMust(info.mappings)
		}
		time.Sleep(time.Second)
	}

	panicIf(true, "docker container '%s' didn't start in time", dockerContainerName)
	return "", ""
}
```

Code for this chapter: https://github.com/kjk/go-cookbook/blob/master/start-mysql-in-docker-go

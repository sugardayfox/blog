Id: 5
Title: Blueprint for deploying web apps on CoreOS
Date: 2017-02-11T19:18:21-08:00
Format: Markdown
Tags: devops
--------------
I used to deploy my web apps on Ubuntu running on Digital Ocean but recently I switched to using [CoreOS](https://coreos.com/) instead of Ubuntu.

For a while I didn't understand CoreOS; a linux distro without package manager? How do I install more software on this thing?

Now I am a convert. CoreOS is not a Linux distro for end users. It's a distro for deploying applications packaged as docker containers.

The benefit of using CoreOS is less configuration needed compared to e.g. Ubuntu.

I used to deploy multiple apps per server but for operational simplicity I moved to using one server per app. At $5 per server (my apps are written in Go, so they run comfortably on the smallest servers) it's a reasonable cost.

Here's my playbook for deploying an app on CoreOS. This example is how I deploy my [blog](https://github.com/kjk/web-blog).

**1\. Create a unique ssh key for the machine**

For security it's good to have a unique ssh key for each machine:

* `ssh-keygen -t rsa -b 4096 -C "a comment"` : creates a new ssh key, save it as `id_rsa`
* content of `id_rsa.pub` is what you give DigitalOcean as ssh key when creating a server

After server is created, verify you can login: `ssh -i ./id_rsa core@<ip_address>`.

**2\. Initial server setup**

To make the scripts more re-usable, create `ipaddr.sh`:
```bash
# e.g. IPADDR=137.63.26.193
IPADDR=<ip address of the server>
```

For convenience I also write `login.sh`:
```bash
#!/bin/bash
. ./ipaddr.sh
ssh -i ./id_rsa core@${IPADDR}
```

Usually a kernel benefits from one or more tweaks and those tweaks have to be
applied at startup. We'll use `systemd` to run a script at startup that makes
adjustements to the kernel parameters. We need the following files on the
server:

`/etc/sysctl.d/startup.conf`:
```
# http://security.stackexchange.com/questions/43205/nf-conntrack-table-full-dropping-packet
# https://coreos.com/os/docs/latest/other-settings.html
net.netfilter.nf_conntrack_max=131072
# in seconds, default value was 600 == 5 mins
net.netfilter.nf_conntrack_generic_timeout = 60
# in seconds, default value was 432000 i.e. 5 days
net.netfilter.nf_conntrack_tcp_timeout_established = 54000
# was 46992
fs.file-max = 131072
```

This file contains the settings. Default values for connection tracking are so
low that it's easy for a malicious person to DoS your server just by opening
a small number of connections to your http server.

`/etc/systemd/system/startup_script.sh`:
```bash
#!/bin/bash

# put in /etc/systemd/system/startup_script.sh
# reload values from /etc/sysctl.d/startup.conf
sysctl --system
# make hashsize match net.netfilter.nf_conntrack_max (should be nf_conntrack_max / 4)
# https://security.stackexchange.com/questions/43205/nf-conntrack-table-full-dropping-packet
echo 32768 > /sys/module/nf_conntrack/parameters/hashsize
```

This script makes the changes. Now we need to make sure it gets called at
startup.

`/etc/systemd/system/startup.service`
```
# put in /etc/systemd/system/startup.service
# http://unix.stackexchange.com/questions/47695/how-to-write-startup-script-for-systemd
[Unit]
Description=things to do after each reboot
# just to be safe, run it after docker is booted
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
# per https://www.digitalocean.com/community/tutorials/how-to-create-and-run-a-service-on-a-coreos-cluster
EnvironmentFile=/etc/environment
# '=-' means it can fail
ExecStart=-/etc/systemd/system/startup_script.sh

[Install]
WantedBy=multi-user.target
```

I also automate putting those files on the server and configuring `systemd`
to notice them with a `initial-server-setup.sh` script.

```bash
#!/bin/bash
set -u -e -o pipefail

. ./ipaddr.sh

scp -i ./id_rsa ./startup.conf ./startup.service ./startup_script.sh core@${IPADDR}:/home/core/

ssh -i ./id_rsa core@${IPADDR} <<'ENDSSH'
sudo cp startup.conf /etc/sysctl.d/startup.conf
sudo chmod 0644 /etc/sysctl.d/startup.conf
sudo cp startup.service /etc/systemd/system
sudo cp startup_script.sh /etc/systemd/system
sudo chmod 0750 /etc/systemd/system/startup_script.sh
sudo systemctl daemon-reload
sudo systemctl enable /etc/systemd/system/startup.service

sudo bash -c "echo 32768 > /sys/module/nf_conntrack/parameters/hashsize"
sudo sysctl --system
rm startup.conf startup.service startup_script.sh
ENDSSH
```

After running `initial-server-setup.sh` you should test that the changes are
applied:

* login to the server with `./login.sh`
* reboot with `sudo shutdown -r now`
* login again
* `sudo sysctl -a | grep nf_conntrack_max`

**3\. Use systemctld to automatically restart the app**

When the server reboots we want the app to start automatically. We also want the app to automatically restart if it crashes.

CoreOS comes with `systemd` so we'll use that.

Create `blog.service` file which instructs systemd how to run a docker container named `blog`:
```bash
# put in /etc/systemd/system/blog.service
[Unit]
Description=blog
# this unit will only start after docker.service
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
# per https://www.digitalocean.com/community/tutorials/how-to-create-and-run-a-service-on-a-coreos-cluster
EnvironmentFile=/etc/environment
# before starting make sure it doesn't exist
# '=-' means it can fail
ExecStartPre=-/usr/bin/docker rm blog
ExecStart=/usr/bin/docker run --rm -p 80:80 -v /data-blog:/data --name blog blog:latest
ExecStop=/usr/bin/docker stop blog
# restart if the fails or is killed e.g. by oom
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

It's a one-time operation but I still like to have it as a script named `install-service.sh`:
```bash
#!/bin/bash
set -u -e -o pipefail

. ./ipaddr.sh

scp -i ./id_rsa ./blog.service core@${IPADDR}:/home/core/blog.service

ssh -i ./id_rsa core@${IPADDR} <<'ENDSSH'
cd /home/core
sudo cp blog.service /etc/systemd/system
sudo systemctl enable /etc/systemd/system/blog.service
rm blog.service
ENDSSH
```

You need to re-run it after updating `blog.service` file.


**4\. Package the app as a Docker image and upload to the server**

A script to build the app, package as Docker image and upload latest version to the server.

The most common advice for uploading/downloading docker images is to use docker registry. For simplicity I just use `docker save`/`docker load` and `scp`.

Here's a `docker_build_and_upload.sh` script:

```bash
#!/bin/bash

# build latest version of the app, upload to the server packaged as
# blog:latest docker image

set -u -e -o pipefail

. ./ipaddr.sh

dir=`pwd`
blog_dir=${GOPATH}/src/github.com/kjk/blog

echo "building"
cp config.json "${blog_dir}"
cd "${blog_dir}"
GOOS=linux GOARCH=amd64 go build -o blog_linux
docker build --no-cache --tag blog:latest .
rm blog_linux
cd "${dir}"

echo "docker save"
docker save blog:latest | bzip2 > blog-latest.tar.bz2
ls -lah blog-latest.tar.bz2

echo "uploading to the server"
scp -i ./id_rsa blog-latest.tar.bz2 core@${IPADDR}:/home/core/blog-latest.tar.bz2

echo "extracting on the server"
ssh -i ./id_rsa core@${IPADDR} <<'ENDSSH'
cd /home/core
bunzip2 --stdout blog-latest.tar.bz2 | docker load
rm blog-latest.tar.bz2
sudo systemctl restart blog
ENDSSH

rm -rf blog-latest.tar.bz2
```

You can see the content of Dockerfile [here](https://github.com/kjk/web-blog).

After setting things up and deploying the app, you should restart the OS (`shutdown -r`) to verify that the app will start up after reboot.

and `tail-logs.sh`:
```bash
#!/bin/bash

. ./ipaddr.sh

ssh -i ./id_rsa core@${IPADDR} <<'ENDSSH'
cd /home/core
docker logs -f blog
ENDSSH
```

---
Id: 9
Title: Notes on Ansible
Date: 2014-11-27T22:59:50-08:00
Tags: programming
Status: deleted
Format: Markdown
---
Ansible is a system for provisioning servers and automatic deployments.

I recently used it for deploying a website written in Go and it worked well for
both provisioning and deployment.

Provisioning is about setting up pre-requisites that only have to be done once.

Deployment is about updating the project to the latest version. It happens
as often as you update the project.

In case of a simple website, provisioning Ubuntu 14.10 server involved the
following steps:

* making sure nginx is installed
* creating a user for the project. This is how I like to organize projects
  if more than one runs on a single server
* add nginx configuration file so that it proxies a given server name to
  the port on which web server is running
* setup initd to watch and restart web server process

Deployment involves:

* cross-compiling the server
* packaging all files that need to be copied to the server in a .zip file
* copy the .zip file to the server
* shutting down currently running web server instance
* unpacking new version
* starting up new version

## Ansible process

Things are simpler if you you create `ansible.cfg`. For this project, mine was:

```
[defaults]
hostfile = inventory
remote_user = root
```

The file `inventory` contains list of servers we'll operate on:

```
sumatrawebsite-provision ansible_ssh_host=sumatrapdfreader.org ansible_ssh_user=root
sumatrawebsite ansible_ssh_host=sumatrapdfreader.org ansible_ssh_user=sumatrawebsite
```

Notice that we only operate against a single server but we have 2 different
configuration:

* `sumatrawebsite-provision`, among other things, creates `sumatrawebsite` user.
  Provisioning actions must be executed by a user with root priviledges.
* `sumatrawebsite` is a user that owns the server project, so it's easiest to use
  that user for deployment steps

Let's create `provisioning.sh` helper script:

```bash
#!/bin/bash
ansible-playbook provisioning.yml
```

This is how we run Ansible playbook. Playbook `provisioning.yml` is:

```yml
---
- name: initial server setup
  hosts: sumatrawebsite-initial
  sudo: True
  tasks:
    - name: create a user
      user: name=sumatrawebsite group=sumatrawebsite groups="sudo" shell=/bin/bash
    - name: make user a sudoer
      lineinfile: dest=/etc/sudoers state=present regexp='^%sumatrawebsite' line='%sumatrawebsite ALL=(ALL) NOPASSWD:ALL'
    - name: create user's .ssh directory
      file: path=/home/sumatrawebsite/.ssh state=directory owner=sumatrawebsite group=sumatrawebsite mode=0755
    - name: copy existing ssh key
      command: cp /root/.ssh/authorized_keys /home/sumatrawebsite/.ssh/authorized_keys
    - name: configure authorized_keys
      file: path=/home/sumatrawebsite/.ssh/authorized_keys mode=0644 owner=sumatrawebsite group=sumatrawebsite
    - name: create directory for nginx logs
      file: >
        path=/var/log/nginx/sumatrawebsite/ state=directory mode=0755
    - name: copy nginx config file
      copy: src=nginx.conf dest=/etc/nginx/sites-available/sumatrawebsite
    - name: enable website
      file: >
        dest=/etc/nginx/sites-enabled/sumatrawebsite
        src=/etc/nginx/sites-available/sumatrawebsite
        state=link
    - name: restart nginx
      service: name=nginx state=restarted
```

If you know Unix, you can figure out what those commands do. Ansible playbook
consists mostly of commands run on the server.

The good thing about Ansible is that it avoids re-doing commands. For example,
`user: name=sumatrawebsite group=sumatrawebsite groups="sudo" shell=/bin/bash`
creates a new user with certain properties.

If a user already exists, Ansible is smart enough to not execute any commands.



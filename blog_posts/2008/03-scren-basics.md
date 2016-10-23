Id: 1020
Title: screen basics
Tags: unix
Date: 2008-03-13T17:12:05-07:00
Format: Markdown
--------------
Basics of using screen:

Start: **screen**\
Deattach: `Ctrl-a + d` within screen session\
Reattach existing session: **screen -r [session-name]**\
List sessions: **screen -ls**

Commands: Ctrl-a and additional key:

-   `d` - detach session
-   `n`, `p` - move to next, previous screen session
-   `c` - create new window
-   `"` - list all windows
-   `Ctrl-a` - switch to other screen
-   `Shift-a` - rename window (shows better in window list)

Links: [LinuxJournal
article](http://www.linuxjournal.com/article/6340?page=0,0)

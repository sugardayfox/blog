Id: 4004
Title: Server monitoring idea
Tags: idea
Date: 2009-02-18T00:44:21-08:00
Format: Markdown
--------------
Server monitor is a service for monitoring servers, providing web-based
UI and iPhone UI.

Architecture overview:

-   a monitoring client (written in C for efficiency but can be
    prototyped in Python or C\#) runs on every server being monitored.
    It periodically (e.g. every second) records important stats about
    the server (current memory/cpu/disk usage, basic info about each
    process running etc.), logs to disk for persistency, and
    periodically (e.g. every 30 mins or after reaching some threshold of
    size of data written to disk) sends this data to the collection
    server. Monitoring client maintains persistent TCP connection to
    collection server (alternative: use UDP, probably could make it more
    efficient but at the expense of making things more complicated (has
    to ensure all data was sent))
-   collection server just saves the data on disk and backs them up on
    s3
-   a web ui front end gets the data about all servers being monitored
    and shows them in useful ways

Misc notes:

-   we need to minimize the amount of data being sent by designing
    compact way of representing the data and compressing it on the wire
    (zlib? gzip? bzip2? lzo?)
-   collection server can talk back to C monitoring client e.g. to ask
    for more frequent updates (e.g. when a user is logged in to web ui
    we want to show near-realtime stats for the server, so we need to
    get the data more often than every 30 mins) so the protocol would
    actually be two-way. It should also be extensible so that we can add
    new commands in the future

Todo:

-   design the format of data we send from monitoring client to
    collection server
-   design the protocol between monitoring client and collection server
-   design a way to load-balance for scalability. We can’t have all
    monitoring clients talk to just one collection server since one
    server will be able to handle a limited number of clients. Also, we
    need this for reliability (one collection server should be able to
    go down and other servers should be able to pick up). We need to be
    able to shutdown/bring up collection servers. We need a way for the
    monitoring client to pick up an available collection server.
-   design how web ui front-end works. It also needs to talk to one of
    many load-balanced data servers who serve relevant data to web ui
    client
-   design web ui, what do we show, how do we show it etc.
-   security of monitoring client - it probably shouldn’t run as root
-   implement monitoring client
-   implement collection server
-   implement basic web ui
-   implement basic iPhone ui
-   continuously improve web ui and iPhone ui

Web ui ideas:

-   starts with overview screen showing all servers at once, one or more
    graphs per server showing basic stats like cpu load, number of
    processes, disk usage, network load, memory usage, updated in real
    time
-   a server detail screen which shows more detailed information e.g.
    the equivalent of ‘top’ output, for just one server. possibly also a
    miniature graphs for other servers, for easy switching between them
-   a death/birth graph which shows a timeline of when processes die and
    get started. Not sure how useful that would be
-   long-term graphs for an overview of how one variable does over time
    e.g. how a CPU load behaves over time, how network traffic behaves
    over time etc. Maybe a way to compare the same value on the same
    graph between two different periods of time (e.g. cpu load today vs.
    yesterday)
-   process inspector - looking about a data for just one process e.g.
    how its memory usage evolved over time

Ideas for other capabilities:

-   e-mail notifications for events e.g. “notify when free space on
    volume /mnt drops below 1 GB”, “notify when apache dies”

Business ideas

-   pricing is monthly, with different plans based on the number of
    servers monitored e.g. low-priced \$7/month for one server (personal
    use) and drastically going up for more servers (small/big company
    usage). Also at least a month free for one server
-   a demo account for one server to let people play with functionality
    even without creating an account


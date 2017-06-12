Id: 59001
Title: VisualAck 0.3.3 released
Tags: software,mac,VisualAck
Date: 2010-01-21T18:51:19-08:00
Format: Markdown
--------------
This idea of "release early, release often" might actually work. A day
after releasing first public version of VisualAck, my visual grep tool
for mac, I'm releasing another version which includes an improvement
[suggested in the
forum](http://visualack.uservoice.com/forums/34392-general/suggestions/453734-remove-escape-sequences-from-output-when-redirecte?ref=title):
to not highlight matches with control sequences when not printing output
to a tty (e.g. when output is redirected to a file).

To complement that, I've added ack-compatible options `--[no]color` and
`--[no]colour` to force or disable highlighting matches.

This, of course, only applies to command-line `vack` utility and not the
UI.

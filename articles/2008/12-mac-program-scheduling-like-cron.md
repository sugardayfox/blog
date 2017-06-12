Id: 243
Title: Mac program scheduling (like crontab)
Tags: mac
Date: 2008-12-13T06:49:35-08:00
Format: Markdown
--------------
User scripts in `~/Library/LaunchAgents`

After adding a script, run: `launchctl load .` to make sure the script
gets loaded by `launchd`. Use `launchctl list` to verify it was loaded.

After modifying a script, do `launchctl unload ${label}` followed by
`launchctl load ${label}`, to make sure the script was updated.

Do `launchctl list` to see list of loaded script. This is undocumented
by appears to be true: 1 in status column means there’s a problem.

Example script `local.kjk.sumatrastatsdaily.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>local.kjk.sumatrastatsdaily</string>

    <key>ProgramArguments</key>
    <array>
      <string>/Users/kkowalczyk/src/kjk-priv/scripts/do-sumatra-stats-daily.sh</string>
    </array>

    <key>LowPriorityIO</key>
    <true/>

    <key>Nice</key>
    <integer>1</integer>

    <key>StartCalendarInterval</key>
    <dict>
      <key>Hour</key>
      <integer>4</integer>

      <key>Minute</key>
      <integer>15</integer>
    </dict>
  </dict>
</plist>
```

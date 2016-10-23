Id: 1307
Title: Getting user-specific application data directory for .NET WinForms apps
Tags: .net,c#,win32,winforms
Date: 2005-12-31T16:00:00-08:00
Format: Markdown
--------------
Windows defines a number of directories for specific things.

One of those directories is a directory that a well-behaved Windows application should use to store user-specific data (e.g. a list of recently open documents by this user, customization information etc.).

Usually this is a hidden directory `C:\Documents and Settings\$UserName\Application Data`. Most applications choose to create a sub-directory named after application. Bigger companies, like Macromedia or Microsoft choose, to use 2-level hierarchy: `$CompanyName\$ApplicationName`.

In .NET framework, `Application` object (from `System.Windows.Forms` namespace) has a static property `UserAppDataPath` that returns a string naming directory for storing user-specific application data. However, designers went overboard and used 3-level hierarchy: `$CompanyName\$ApplicationName\$VersionNumber`.

Adding version number is a mistake. When user upgrades the application, he doesn't want to loose his settings and customizations.

Another way to get this is to use `Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData)` and append your application name, like this:

```c#
class Util
{
    static public string GetUserDataPath()
    {
        string dir = Environment.GetFolderPath(Environment.SpecialFolder.ApplicationData);
        dir = System.IO.Path.Combine(dir, "MySoftware");
        if (!Directory.Exists(dir))
            Directory.CreateDirectory(dir);
        return dir;
    }
}
```

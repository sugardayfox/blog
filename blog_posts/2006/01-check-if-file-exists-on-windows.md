Id: 305
Title: Check if file exists on Windows
Tags: win32,c,programming
Date: 2006-01-01T16:00:00-08:00
Format: Markdown
--------------
Windows doesn't have a built-in function that checks if a file with a given name
exists. It can be trivially written using `GetFileAttributes` or `FindFirstFile`
APIs. Version below uses `GetFileAttributes`.

```c
/* Return TRUE if file 'fileName' exists */
bool FileExists(const TCHAR *fileName)
{
    DWORD       fileAttr;

    fileAttr = GetFileAttributes(fileName);
    if (0xFFFFFFFF == fileAttr)
        return false;
    return true;
}
```

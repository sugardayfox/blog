Id: 303
Title: Get file size under windows
Tags: win32,c,programming
Date: 2006-01-06T16:00:00-08:00
Format: Markdown
--------------
Windows doesn't have an API to get a file size based on file name. This small function does that.

It returns -1 if a file doesn't exist.

It doesn't handle files > 2 GB (max positive number for 32 bit signed value). It's quite easy to extend it to 64-bits if you know what is the 64 bit integer type in your compiler (unfortunately there's no standard).

A better design might be `BOOL GetFileSize(const TCHAR *fileName, unsigned long *fileSizeOut)` i.e. returning false if file doesn't exist and putting the file size into `fileSizeOut`.

```c
long GetFileSize(const TCHAR *fileName)
{
    BOOL                        fOk;
    WIN32_FILE_ATTRIBUTE_DATA   fileInfo;

    if (NULL == fileName)
        return -1;

    fOk = GetFileAttributesEx(fileName, GetFileExInfoStandard, (void*)&fileInfo);
    if (!fOk)
        return -1;
    assert(0 == fileInfo.nFileSizeHigh);
    return (long)fileInfo.nFileSizeLow;
}
```

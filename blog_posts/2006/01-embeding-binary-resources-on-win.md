Id: 1294
Title: Embedding binary resources on Windows
Tags: win32,c,programming
Date: 2006-01-29T16:00:00-08:00
Format: Markdown
--------------
Windows binary can have resources embedded in them. Most resources are of
predetermined type (e.g. a menu, an icon or a bitmap) but you can also
embed arbitrary binary data (e.g. a text file). The proper syntax is hard to figure out
just from reading msdn docs.

This snippet shows how to embed a binary resource from a file.

First you need to define a resource identifier in a header file (e.g. `resource.h`)
that will be used by both C compiler and resource compiler:
```c
#define MY_RESOURCE 300
```

Then you need to add to your resource file (e.g. `resource.rc`):
```
MY_RESOURCE RCDATA "file-with-data.txt"
```

And finally, this is how you can get to this data:
```c
void WorkOnResource(void)
{
    HGLOBAL     res_handle = NULL;
    HRSRC       res;
    char *      res_data;
    DWORD       res_size;

    // NOTE: providing g_hInstance is important, NULL might not work
    res = FindResource(g_hInstance, MAKEINTRESOURCE(MY_RESOURCE), RT_RCDATA);
    if (!res)
        return;
    res_handle = LoadResource(NULL, res);
    if (!res_handle)
        return;
    res_data = (char*)LockResource(res_handle);
    res_size = SizeofResource(NULL, res);
    /* you can now use the resource data */
}
```

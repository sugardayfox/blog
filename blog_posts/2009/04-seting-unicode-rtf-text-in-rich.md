Id: 19006
Title: Setting unicode rtf text in rich edit control
Tags: win32,programming
Date: 2009-04-14T23:36:41-07:00
Format: Markdown
--------------
Despite being around for thousands of years, win32 API is poorly
documented in many areas.

Today I spent an hour battling a baffling problem with rich edit
control.

I was setting rtf-formatted text using `WM_SETTEXT` message.

In non-unicode build, where text was ansi, it worked and the text showed
up nicely formatted.

In unicode build, however, the text wasn’t parsed as rtf and showed
verbatim.

An hour of trying various ways to make it work ended with me finding a
working work-around: using `EM_SETTEXTEX` and utf-8 encoded text with
`CP_UTF8` code-page. I’m using WTL so the code ended up being:

```c++
#ifdef UNICODE
 // Don’t know why I have to do this, but SetWindowText() with unicode
 // doesn’t work (rtf codes are not being recognized)
 const char *sUtf = WstrToUtf8(s);
 m_statusMsgEdit.SetTextEx((LPCTSTR)sUtf, ST_DEFAULT, CP_UTF8);
#else
 m_statusMsgEdit.SetWindowText(s);
#endif
```

MSDN [claims](http://msdn.microsoft.com/en-us/library/bb774284.aspx) I
could have used code-page 1200 for unicode text, but that didn’t work
either.

It’s quite possible I’ve botched something somewhere in my code and this
work-around isn’t needed but I can hardly think of anything.


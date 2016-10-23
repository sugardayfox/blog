Id: 928
Title: Debugging adventure
Tags: programming
Date: 2006-01-13T15:42:02-08:00
Format: Markdown
--------------
This is a debugging story and the lessons learned.\
\
 I was writing an application using a
[MailSwitchToAccount](http://msdn.microsoft.com/library/default.asp?url=/library/en-us/mobilesdk5/html/wce51lrfmailswitchtoaccount.asp)
API new in WIndows Mobile 5.0. According to docs it's defined in
cemapi.h and an app needs to link with cemapi.lib.\
\
 That's how I called it:\

> hr = MailSwitchToAccount(\_T("SMS"));

Visual Studio 2005 refused to link the app, claiming that:\

> TestDevice.obj : error LNK2019: unresolved external symbol "long
> \_\_cdecl MailSwitchToAccount(wchar\_t const \*,unsigned long)"
> (?MailSwitchToAccount@@YAJPB\_WK@Z) referenced in function "void
> \_\_cdecl OnMailSwitchToAccount(void)" (?OnMailSwitchToAccount@@YAXXZ)

My first thought was that this symbol is simply missing from cemapi.lib.
So I used dumpbin.exe (part of Visual Studio tools) to dump all symbols
exported from cemapi.lib:

> dumpbin -exports cemapi.lib

I found that it does have MailSwitchToAccount:

> ?MailSwitchToAccount@@YAJPBGK@Z (long \_\_cdecl
> MailSwitchToAccount(unsigned short const \*,unsigned long))

Howerver, if you look closely, you'll see that signatures don't match:
expected type for the first parameter is "unsigned short \*" while I'm
calling it as "wchar\_t const \*". Now, "unsigned short" is the WCHAR
Windows UNICODE type. C++ also defines wchar\_t which is also UNICODE
char.\
\
 I vaguely remembered that some C++ compilers have an option to treat
wchar\_t as a native language type (as opposed to just a typedef for
existing type, "unsigned short" in Windows' case). At indeed, there it
was, in project properties, C/C++/Language there's an option "Treat
wchar\_t as Built-in Type", set by default to Yes. You can set it to
"No", which corresponds to passing "/Zc:wchar\_t-" to cl.exe.\
\
 It seems wrong that you have to do that. It seems like
cemapi.dll/cemapi.lib were compiled with "/Zc:wchar\_t-" which forces
everyone who links to them also be compiled like that.\
 Lessons learned:\

-   C++ is evil
-   dumpbin.exe is your friend
-   Visual Studio's code browsing is your friend (it's easy to find out
    what a given typedef or \#define really is)
-   be aware of "wchar\_t as built-in type or not" issue

\
 <span style="font-weight: bold;">Update:</span> turns out it's a known
problem and has been blogged about on [official Visual Studio
blog](http://blogs.msdn.com/vsdteam/archive/2005/11/16/linker_error_lnk2019_lnk2001.aspx).
Good news is that it will probably be fixed in future versions of Visual
Studio.

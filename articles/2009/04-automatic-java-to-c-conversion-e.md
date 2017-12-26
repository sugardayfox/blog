Id: 21002
Title: Automatic Java to C# conversion - experience using Java Language Conversion Assistant
Tags: c#,.net,tenderbase
Date: 2009-04-22T00:10:19-07:00
Format: Markdown
Status: hidden
--------------
Recently I’ve been working on porting perst (a small, fast,
object-oriented database) from Java to C\#.

The thought of manually converting 10k+ lines of code was scary so I
decided to give automatic translation a try.

First I tried
[sharpen](http://evain.net/blog/articles/2008/05/20/sharpen-an-open-source-java-to-c-converter),
which is a part of db4o project. That was a fail. What I expected was a
command-line tool that I can run on **.java files that will produce**.cs
files. What I got was an eclipse plugin. I can barely create a working
project in eclipse so I failed to install it properly (even after trying
to follow supposedly step-by-step instructions).

Then I tried Microsoft’s JLCA (Java Language Conversion Assistant). It’s
part of Visual Studio 2005. Sadly, it’s been discontinued (apparently
Microsoft no longer considers Java to be threatening enough to continue
investing in automatic java-to-c\# conversion).

Here are the exact commands I used to convert perst sources:

    mkdir c:\kjk\src\nuperst-csharp
    mkdir c:\kjk\src\nuperst15-csharp
    cd C:\Program Files\Microsoft Visual Studio 8\JavaLanguageConversionAssistant
    jconvert c:\kjk\src\nuperst\java\src /JDK J2EE /OutDir c:\kjk\src\nuperst-csharp
    "C:\Program Files\Microsoft Visual Studio 8\JavaLanguageConversionAssistant\jconvert" c:\kjk\src\nuperst\java\src /JDK J2EE /Out c:\kjk\src\nuperst-csharp
    "C:\Program Files\Microsoft Visual Studio 8\JavaLanguageConversionAssistant\jconvert" c:\kjk\src\nuperst\java\src15 /JDK J2EE /Out c:\kjk\src\nuperst15-csharp /ProjectType Library /ProjectName Nuperst /Verbose
    "C:\Program Files\Microsoft Visual Studio 8\JavaLanguageConversionAssistant\jconvert" c:\kjk\src\nuperst\java\src15 /JDK J2EE /Out c:\kjk\src\nuperst15-csharp

The results are far from perfect but I’m certainly happy I tried it
first. I’ve had to fix parts that were converted incorrectly, parts that
were not converted at all, fix style issues (make the code more
idiomatic C\#) etc. It took me several days and there’s still plenty to
do. It doesn’t change the fact that I probably saved lots of time
compared to an alternative i.e. converting everything manually.

Java and C\# are remarkably similar but have enough dissimilarities to
trip automatic conversion. There are many big and small, obvious and
subtle differences. Parts where JLCA had problems:

-   perst uses reflection heavily. JLCA couldn’t translate class loading
    by name Java idioms (that use ClassLoader) to .NET equivalent (which
    must fish out things from assemblies)
-   Iterator in Java has remove(), IEnumerator in .NET doesn’t. JLCA
    doesn’t quite know what to do about it
-   classes implementing IO.Stream interfaces were not fully converted
    due to differences in interfaces
-   some final values were translated as readonly, instead of const, and
    initialized at runtime (vs. compile time) which caused other values
    that dependent on it be incorrect
-   the code is littered with ‘internal’ keyword - it’s the most
    faithful translation but I don’t think internal is very idiomatic
    C\#
-   no attempt was made to fix translation naming conventions (e.g. Java
    has method() and C\# should have Method()) other than when
    translating between corresponding, already existing classes

Automatic Java to C\# translation is promising. It is harder than it
looks like at first but if an implementation abandoned several years ago
can do a decent job, a better implementation should be possible.

GWT translates Java to JavaScript so it seems like it should be possible
to do even better job converting to a language that is much similar.

Sharpen might be doing a better job, but I would need extremely
detailed, step-by-step-with-screenshots guide on how to use it. While I
understand why it’s done this way (eclipse has built-in Java parser that
sharpen uses), it’s a shame that mechanics of using a tool are
prohibitively complex (at least for those who don’t know ins-and-outs of
eclipse)

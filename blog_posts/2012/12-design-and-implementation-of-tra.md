Id: 1518013
Title: Design and implementation of translation system for desktop software
Tags: programming,sumatra
Date: 2012-12-16T18:50:28-08:00
Format: Markdown
--------------
A web-based, crowd-sourced translation system.
----------------------------------------------

I just finished building [AppTranslator](http://www.apptranslator.org/),
which is a third iteration of a system that allows people to contribute
translations for
[SumatraPDF](http://blog.kowalczyk.info/software/sumatrapdf/).

This post describes the big picture of how it’s designed and
implemented. Everything I discuss is
[open-source](https://github.com/kjk/apptranslator), so you can dig into
code for more details.

The goals of translation system were to make it easy for me to add new
strings for translations and easy for other people to contribute
translations.

I call it a “system” because it consists of several parts that are
designed to act as a sensible whole.

How translations are stored in C++ code
---------------------------------------

I try to design for simplicity of implementation to minimize the amount
of code to write. Often it means that I reject popular solutions if I
can design a simpler one.

SumatraPDF is a C++ desktop software for Windows. There already is a
popular system for maintaining translation for such programs: use
resource dlls, one dll per supported language. The translations are kept
in .rc files as strings and referred to by numeric ids from code via
[LoadString()](http://msdn.microsoft.com/en-us/library/windows/desktop/ms647486.aspx)
API.

That would be very clumsy both for me (when adding new translations) and
for potential translators (they have to edit files in a cryptic format).
It also requires distributing multiple dlls (I prefer to ship my apps as
single executables).

Managing translation strings in C++ code
----------------------------------------

The first part of the system was therefore a new design for managing
translation strings in the app.

All translations are marked with [\_TR()
macro](https://code.google.com/p/sumatrapdf/source/browse/trunk/src/Translations.h)
in the code.

A [python
script](https://code.google.com/p/sumatrapdf/source/browse/trunk/scripts/update_translations.py)
extracts those strings from the source. It then reads translations
provided by contributors and generates a C++ file with strings. \_TR()
macro evaluates to calling a simple supporting code in
[Translations.cpp](https://code.google.com/p/sumatrapdf/source/browse/trunk/src/Translations.cpp)
that looks up a translation of a string given currently selected
language.

How translations are contributed by users
-----------------------------------------

SumatraPDF is an open-source project so I rely on contributions from
users to keep it translated.

The system for collecting contributions had 3 iterations.

In first iteration, I had a single text file with all translations.

It looked like this:

```
String to translate
de:German translation
fr:French translation

Another string to translate
de:another German translation
```

People would download the latest version of the file from svn, add
missing translations and e-mail it to me. I would check-in that to svn
and re-run the script that rebuilds C++ file with strings.

At some point the file became very big so I split it into multiple
files, one per language.

It was working ok but the process of submitting translation was time
consuming for translators and the process of updating the code was time
consuming for me.

For that reason I built a [web-based
service](http://www.apptranslator.org) which makes it much easier to
contribute a translation.

I also added the necessary API endpoints in the server to allow writing
scripts for automating uploading strings to translate and downloading
latest translations.

The design of web-based UI
--------------------------

A web-based UI for editing translation is not a novel idea. However, my
brief research shows that few do it well.

I try to be pragmatic about things. If a decent 3rd party system with
enough flexibility to meet my needs already did exist, I would rather
use it than develop my own (writing code takes time).

I’m tempted to say that I designed for simplicity but, while being true,
it’s also rather vague. More accurately: I designed for simplicity of
the translator workflow.

You can judge the UI [yourself](http://www.apptranslator.org) but let me
point out few specific points:

-   it takes 2 clicks from main page to a point where you can submit
    translation. The first click is to select the project (as
    AppTranslator is designed to support multiple projects) and the
    second click to select the language
-   on the page for a given language, the untranslated strings are at
    the top. The rest is sorted
-   I don’t require creating an account for AppTranslator but
    authenticate with Twitter

Those might seem like obvious points but I found that other systems I’ve
surveyed were baroque enough to inspire Kafka.

Ubuntu’s [translation system](https://wiki.ubuntu.com/Translations) gets
a special recognition for amount of bureaucracy, complicated workflows
(joining a translation team to submit a translation?) and bad copy (they
feel it’s important to inform you that in order to contribute a
translation you need an internet connection, among many other useless
bits of information).

Mozilla [is marginally better](https://wiki.mozilla.org/L10n:Home_Page)

Technical specs
---------------

AppTranslator is written in Go and
[open-source](https://github.com/kjk/apptranslator) (BSD license).

I run it on Ubuntu server. It’s possible, but [not
easy](https://github.com/kjk/apptranslator/blob/master/docs/deploy_your_own.txt)
to run your own instance.

SumatraPDF is also [open-source](https://code.google.com/p/sumatrapdf/).
The code is in C++ with a bunch of python helper scripts to automate
interaction with AppTranslator server.

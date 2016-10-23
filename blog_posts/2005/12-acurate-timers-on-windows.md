Id: 1309
Title: Accurate timers on Windows
Tags: win32,programming
Date: 2005-12-30T16:00:00-08:00
Format: Markdown
--------------
From [this blog post](http://ooeygui.typepad.com/ooey_gui/2005/10/animations.html) on using accurate timers on Windows.

[CraeteWaitableTimer](http://msdn.microsoft.com/library/default.asp?url=/library/en-us/dllproc/base/createwaitabletimer.asp) API is similar to [SetTimer](http://msdn.microsoft.com/library/default.asp?url=/library/en-us/winui/winui/windowsuserinterface/windowing/timers/timerreference/timerfunctions/settimer.asp). You can specify timeout and a callback function. It has the following advantages:

* you can specify a delay
* it returns a kernel handle you can wait on
* it allows pasing context data pointer to a callback function
* waitable timers are much more accurate as they are not calculated during idle and messages aren't collapsed
* you can adjust the parameters of the timer with [SetWaitableTimer](http://msdn.microsoft.com/library/default.asp?url=/library/en-us/dllproc/base/setwaitabletimer.asp) and cancel it with [CancelWaitableTimer](http://msdn.microsoft.com/library/default.asp?url=/library/en-us/dllproc/base/cancelwaitabletimer.asp)..

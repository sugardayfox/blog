Id: 346001
Title: Which technology for writing desktop software?
Tags: programming
Date: 2010-12-05T23:20:58-08:00
Format: Markdown
--------------
This post is inspired by [a
question](http://discuss.joelonsoftware.com/default.asp?biz.5.829264.15)
posted on JoelOnSoftware software business forum.

To paraphrase: “I’m a one person coding machine. I want to write small,
consumer oriented desktop software. Which technology should I use?
Should I write cross-platform code or stick with one platform?”

It is a good question because the amount of available technologies is
vast: C/C++, .NET, Java, Adobe Air, Cocoa with Objective-C, QT,
wxWidgets and many more lesser known technologies.

I have an answer.

Don’t write cross-platform code
-------------------------------

The logic behind writing cross-platform code is deceptively attractive:
if I target both Mac and Windows, I’ll sell twice as many copies of my
software.

However, at the beginning your biggest problem is that your idea hasn’t
been tested with the market. Writing for a single platform is faster so
you’ll find out sooner if your software is a hit or a flop.

If it’s a flop, the effort of supporting two platforms was lot.

If it’s a success: congratulation, you now have the money to invest in
supporting the other platform and a much higher chance that the result
will be successful as well.

Don’t delude yourself that a cross-platform toolkit will allow you to
ship a cross-platform as fast as shipping a single-platform code. Life
is not as simple as that.

Today Qt is the best cross-platform toolkit, a truly impressive piece of
work. It is, however, limited by its legacy. Qt requires you to use C++.
You’ll be much more productive using C\# with WPF on Windows or
Objective C with Coca on Mac.

Writing code takes only a fraction of time and effort needed to ship an
application to end users. Other activities that you’ll need to do that
also require doubling the time and effort:

-   testing the code on both platforms
-   writing documentation and marketing materials for 2 platforms
-   learning how to market to both Windows and Mac users
-   learning to become an expert user of both OSes to help your users
    when inevitable support requests come in

What you need from a platform
-----------------------------

Sometimes platform are chosen based on vague understanding of what they
provide. I see no end of people who e.g. will blindly suggest Adobe Air
just because it’s marketed as allowing cross-platform development and
happens to get lots of buzz recently, without considering if the
technical limitations of Air would even make a given application
possible to write.

Here are things to consider when choosing a platform for your software.

### Is it capable enough?

Adobe Air might be a great technology for some applications. If you’re
building a Twitter client, it’ll do just fine. If you’re building a
novel compression engine, it’s not appropriate.

By their nature, cross platform toolkits tend to target the lowest
common denominator and don’t support important platform-specific
features. Those fancy Core Animation transitions in your Mac app? You
probably won’t get that from Java or Air or Qt. Support for Windows 7
task bar integration? Probably not there either.

### Is it documented well?

This is not only about documentation that comes in the box but also
about the availability of greater ecosystem. You’ll run into problems,
bugs, things you don’t understand. It’s important to resolve those
issues quickly.

Mainstream technologies like Java or .NET are documented in books from
respected publishers, numerous blog posts and have communities where you
can show up, ask a question and get an answer quickly. There’s much less
of that for e.g Adobe Air or Clojure.

### Are there third-party libraries available?

You don’t want to write a zip compression or http client from scratch.
Availability of 3rd party libraries that you can drop in your code and
save tons of time is important. In that respect Java and .NET are
unmatched. There’s lots of C/C++ libraries out there, but integrating
them in your project will be harder.

### Are other programs like the one you’re thinking about implemented with that toolkit?

Can you write a great text editor in Air? I don’t know, I haven’t seen
one.

The fact that Microsoft used .NET/WPF for writing editor in Visual
Studio or that Eclipse is written in Java is a proof that a complex text
editor can be written in those technologies.

### Do you know it well?

Learning a new technology is hard. Even switching to a similar platforms
(e.g. from Java to C\#) requires a large amount of effort and time. All
else being equal, use the technology you already know.

Now for a final verdict.

Cocoa is for Mac
----------------

There really is no contest: for writing Mac software use Cocoa with
Objective-C. While Java or C\# are more productive than Objective-C, you
can’t beat the combination of tool support (XCode, especially in version
4, is doing a lot to help you program in Objective-C), documentation,
community. It’s also the most capable, being the native language of the
platform.

.NET and WPF is for Windows
---------------------------

While Java is close to C\# in productivity, documentation and
availability of 3rd party libraries, .NET wins by having WPF, better
integration with Windows and being better optimized for desktop apps.

WPF is what seals the deal: there’s nothing even remotely as good in
Java land for writing UI code.

Unsurprisingly, .NET is more capable than Java when it comes to all the
little things you need to do to integrate with the OS really well.
Things like support for Win 7 taskbar, reading/writing registry etc. I’m
sure all that can eventually be achieved with Java, but .NET already has
everything you can possibly need to make your app a good citizen in
Windows.

Finally, Java platform is mostly used on servers and optimized for that.
In that context, it’s much more important that an app runs fast in a
stead state than to have short startup time.

On desktop, things like startup time matters and Microsoft cares about
minimizing that, since their own apps are increasingly using .NET. Your
apps get the benefit of Microsoft’s work on optimizing startup time.

What if you absolutely must target both platforms?
--------------------------------------------------

What if your startup idea requires multi-platform support and you just
got few millions in funding?

It’s still better to hire two competent programmers, each an expert in
writing software for one platform.

I’ve seen very few successful products that target both Windows and Mac.
I exclude titans like Microsoft or Adobe, who have extremely complex
application and infinite amounts of cash to pour into development.

Within the confines of small company, at best you might succeed despite
sucking, like Libox (their Mac application doesn’t even respect Cmd-Q).
All the runaway successes that come to my mind, like Evernote or
Dropbox, did invest in a native clients for each platform they support,
because their apps are their core differentiator (there are tens of
note-taking applications out there or apps that allow to backup or share
files in some way).

Id: 98001
Title: Things I've learned this week
Date: 2010-04-24T01:26:24-07:00
Format: Markdown
--------------
Things I’ve learned this week:

-   using pt instead of px for CSS font-size property is a good way to
    have text that looks good on monitors with different DPIs. 8pt is a
    good value for small text
-   typing takes a lot of time. I’m working on [JavaScript quick
    reference](/articles/jsref.html) and it turns out that just typing
    all that text takes hours (and yes, I do touch type). No wonder
    books take years to write
-   therefore shorthand notations are good. I started writing it html
    and it was painful. I improved things by writing small [python
    script](http://github.com/kjk/blog/blob/master/www/articles/jsrefgen.py)
    that generates html from a more compact text format
-   [etherpad](http://etherpad.com) was open-sourced but pretty much
    abandoned by (now) google folks. There is an [active
    project](http://github.com/ether/pad) that continues the work
-   CSS: `overflow: auto;` in a div makes it scrollable (instead, for
    example, overflowing it)
-   HTML5: canvas element has `toDataURL()` which returns PNG of itself.
    Can be used to add a way to easily export canvas data for external
    consumption e.g. `window.location = canvasEl.toDataURL()` will
    navigate to a page that is the PNG file (which can then be save to
    disk or printed using browser’s context menu). More sophisticated
    solution could dynamically add PNG to current page
-   [Tornado](http://www.tornadoweb.org/) is a decent python web server
    framework


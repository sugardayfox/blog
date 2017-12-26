Id: 954
Title: A simple catchpa scheme
Date: 2006-08-17T12:37:09-07:00
Format: Markdown
--------------
Captcha are those blurry and transformed images that recently became so
popular on many websites that accept user-contributed content.

They are an unevitable consequences of spammers becoming more sophisticated.

My forum for Sumatra PDF recently received an annoyingly high amount of spam, so I decided to put an end to it. Or at least die trying.

I extended the [FruitShow](http://sourceforge.net/projects/fruitshow)
forum software with a captcha scheme I've stolen from
[CVSTrac](http://www.cvstrac.org/) software.

It's very simple: instead of showing blurry images, it asks people to enter a result of a very simple arithmetic expression, like 1+3.

It seems to work for CVSTrac so maybe it'll stop spam on my forum as
well (so far it's been a couple of days without a single spam).

Technically, it's not hard to defeat - all the data needed for correct
response are in the html.

I'm not even trying to do anything fancy like hide the numbers in JavaScript (so that the spam bot needs full JavaScript evaluation engine).

I'm counting more on obscurity of the method.

While it would be easy to manually modify the spam bot to defeat this particular captcha, I'm hoping that no-one will bother to put the effort just so that they can spam one little website.

FruitShow, by the way, rocks.

It took me an hour and just a few lines of PHP to add this. Too bad it doesn't seem to be developed anymore (3 months of checkin silence).

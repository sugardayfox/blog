Id: 998
Title: Google App Engine tip
Tags: appengine
Date: 2008-07-05T14:21:52-07:00
Format: Markdown
--------------
**Update**: the bug below has been fixed and SDK behavior now matches production.

Don't put template files in directories marked as static dirs in app.yaml
file. For whatever reason those files are not available to template.render()
function.

What's worse: this works locally, in dev environment and only breaks when you
upload the app. I've learned that the hard way.

I've opened a [bug for this issue][1] (so go there and star it) but I wonder
if Google even looks at them anymore, after all the requests to add Pascal
support.

   [1]: http://code.google.com/p/googleappengine/issues/detail?id=550



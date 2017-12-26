Id: 3081
Title: NSCopying, NSMutableCopying or NSCoding
Date: 2009-02-18T00:57:00-08:00
Format: Markdown
Status: hidden
Tags: objective c
--------------
If you get weird crashes that look as if some other Cocoa
class is releasing your object prematurely and an access to
instance variables crashes because object has already been freed,
check what protocols your class’s superclass conforms to.

Chances are, it conforms to `NSCopying`, `NSMutableCopying` or `NSCoding` and you forgot to override `copyWithZone:` or `mutableCopyWithZone:`.

Sometimes this mistake also looks as if there were somehow two
copies of your object, one valid, and one for which constructor was never called
and which has invalid instance variables and behaves zombie-like, but neither `NSZombie` nor any of your other memory debug tricks trigger for it.

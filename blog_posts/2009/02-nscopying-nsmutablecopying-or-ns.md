Id: 3081
Title: NSCopying, NSMutableCopying or NSCoding
Tags: objective c
Date: 2009-02-18T00:57:00-08:00
Format: Markdown
--------------
If you get weird crashes that look suspiciously like some other Cocoa
class (that should really be holding on to an object of your class) is
releasing your object prematurely, and thus an access to one of your
instance variables causes a crash because that’s already been freed,
check what protocols your class’s superclass conforms to. Chances are,
it conforms to `NSCopying`, `NSMutableCopying` or `NSCoding` and you
silly sleepwalker forgot to override `copyWithZone:` or
`mutableCopyWithZone:`.

Sometimes this mindless mistake also looks as if there were somehow two
copies of your object, one valid, and one for which there was never a
constructor called, and which thus has some invalid instance variables
and behaves zombie-like, but neither `NSZombie` nor any of your other
memory debug tricks really trigger for it.

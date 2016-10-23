Id: 1292
Title: What makes a CD bootable
Date: 2006-03-02T16:00:00-08:00
Format: Markdown
--------------
This is a copy of [this article](http://technopedia.info/tech/2006/03/03/what-makes-a-cd-bootable.html)

In order for a CD to be bootable, it must contain two files: BOOTCAT. BIN and
BOOTIMG.BIN. BOOTCAT.BIN is a catalog file, and BOOTIMG.BIN is an image file an
image of a bootable floppy disk. That's why you need an existing bootable 
floppy in order to make a bootable CD.

When you browse the contents of a bootable CD in Windows or at a command prompt, 
you won't see any of the files that you would find on a bootable floppy; that's 
because they are all stored within BOOTIMG.BIN.

When you boot from the bootable CD, everything in BOOTIMG.BIN and BOOTCAT.BIN 
shows up as being on the A: drive, while everything else on the CD shows up on 
the regular CD drive letter. It tricks the PC into thinking that there is 
actually a floppy in the A: drive.

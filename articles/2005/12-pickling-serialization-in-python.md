Id: 312
Title: Pickling (serialization) in Python
Tags: python
Date: 2005-12-30T16:00:00-08:00
Format: Markdown
--------------
Pickling is an easy way to serialize data in Python. One possible use for that is preserving the state across script executions (like saving preferences).

There are few things worth knowing:

* python has `pickle` and `cPickle` modules. They are almost the same (`pickle` handles more cases but `cPickle` is faster)
* you can specify protocol parameter to `dump()` function. Use `cPickle.HIGHEST_PROTOCOL` - it's the most efficient one
* the simplest thing to do is to stuff everything you want to serialize in a hash and serialize the hash

The code snippet below shows how to save and load some data to a file. It removes the file if unpickling fails (which can happen if e.g. file is corrupted or not in the right format). The retry logic comes from experience - I found that `os.remove()` right after `close()` might fail.

```python
import sys, os, string, time, cPickle

DATA_FILE_NAME = "settings.dat"

def saveData():
    fo = open(DATA_FILE_NAME, "wb")
    version = 1.0
    aString = "some data"
    cPickle.dump(version, fo, protocol = cPickle.HIGHEST_PROTOCOL)
    cPickle.dump(aString, fo, protocol = cPickle.HIGHEST_PROTOCOL)
    fo.close()

def loadData():
    try:
        fo = open(DATA_FILE_NAME, "rb")
    except IOError:
        # it's ok to not have the file
        print "didn't find file %s with data" % DATA_FILE_NAME
        return
    try:
        version = cPickle.load(fo)
        aString = cPickle.load(fo)
    except:
        fo.close()
        removeRetryCount = 0
        while removeRetryCount < 3:
            try:
                os.remove(filePath)
                break
            except:
                time.sleep(1) # try to sleep to make the time for the file not be used anymore
                print "exception: n  %s, n  %s, n  %s n  when trying to remove file %s" % (sys.exc_info()[0], sys.exc_info()[1], sys.exc_info()[2], filePath)
            removeRetryCount += 1
        return
    fo.close()
```

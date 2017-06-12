Id: 311
Title: Serialization in C#
Tags: c#,.net
Date: 2005-12-27T16:00:00-08:00
Format: Markdown
--------------
Serialization (known as pickling in python) is an easy way to convert an object
to a binary representation that can then be e.g. written to disk or sent over
a wire.

It's useful e.g. for easy saving of settings to a file.

You can serialize your own classes if you mark them with `[Serializable]` attribute.
This serializes all members of a class, except those marked as `[NonSerialized]`.

.NET offers 2 serializers: binary, SOAP, XML. The difference between binary and SOAP is:

* binary is more efficient (time and memory used)
* binary is not human-readable. SOAP isn't much better.

XML is slightly different:

* it lives in `System.Xml.Serialization`
* it uses `[XmlIgnore]` instead of `[NonSerialized]` and ignores `[Serializable]`
* it doesn't serialize private class members

An example of serialization/deserialization to a file:

```c#
using System.IO;
using System.Diagnostics;
using System.Runtime.Serialization;
using System.Runtime.Serialization.Formatters;
using System.Runtime.Serialization.Formatters.Binary;

[Serializable]
public class MySettings {
    public int screenDx;
    public ArrayList recentlyOpenedFiles;
    [NonSerialized]public string dummy;
}

public class Settings {
    const int VERSION = 1;
    static void Save(MySettings settings, string fileName) {
            Stream stream = null;
            try {
                IFormatter formatter = new BinaryFormatter();
                stream = new FileStream(fileName, FileMode.Create, FileAccess.Write, FileShare.None);
                formatter.Serialize(stream, VERSION);
                formatter.Serialize(stream, settings);
            } catch {
                // do nothing, just ignore any possible errors
            } finally {
                if (null != stream)
                    stream.Close();
            }
    }

    static MySettings Load(string fileName) {
        Stream stream = null;
        MySettings settings = null;
        try {
            IFormatter formatter = new BinaryFormatter();
            stream = new FileStream(fileName, FileMode.Open, FileAccess.Read, FileShare.None);
            int version = (int)formatter.Deserialize(stream);
            Debug.Assert(version == VERSION);
            settings = (MySettings)formatter.Deserialize(stream);
        } catch {
            // do nothing, just ignore any possible errors
        } finally {
            if (null != stream)
                stream.Close();
        }
        return settings;
    }
}
```

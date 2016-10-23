Id: 1987
Title: Too much oo
Date: 2008-01-10T21:41:59-08:00
Format: Markdown
--------------
Just found in an otherwise good code:\
\
<div style="margin-left: 40px; font-family: Courier New;">

class DocDecompressor {\
\
public:\
 DocDecompressor() {}\
 \~DocDecompressor() {}\
\
 size\_t decompress(ZLInputStream &stream, char \*buffer, size\_t
compressedSize, size\_t maxUncompressedSize);\
};\

</div>

\
 Why is it a class and not just a function?\
\
Rhetorical question, I know. But it illustrates what I've seen more than
once: people get so enamored with some technique that they blindly use
it even when it doesn't make any sense.\
\


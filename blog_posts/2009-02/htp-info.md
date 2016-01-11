Id: 7040
Title: HTTP info
Tags: http,reference
Date: 2009-02-26T00:13:33-08:00
Format: Markdown
--------------
### Response codes

```
200 - OK
206 - Partial Content (if successfully returned part of the file)

301 - Moved permanently
302 - Found (temporary redirect)
303 - See other
304 - Not Modified

400 - Bad Request
401 - Not Authorized
403 - Forbidden
404 - Not Found
405 - Method Not Allowed
406 - Not Acceptable
```

### Basic and Digest authentication

Covered by [rfc 2617](http://www.ietf.org/rfc/rfc2617.txt)

Basic adds Authorization: header, e.g.:
`Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==` where the value
after Basic is base64 encoding of string \${userid} “:” \${password}

Digest is more complicated. Server responds with e.g.:
```
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Digest
 realm=“testrealm@host.com”,
 qop=“auth,auth-int”,
 nonce=“dcd98b7102dd2f0e8b11d0f600bfb0c093”,
 opaque=“5ccc069c403ebaf9f0171e9517f40e41”
```

And client has to reply with:
```
Authorization: Digest username=“Mufasa”,
 realm=“testrealm@host.com”,
 nonce=“dcd98b7102dd2f0e8b11d0f600bfb0c093”,
 uri=“/dir/index.html”,
 qop=auth,
 nc=00000001,
 cnonce=“0a4f113b”,
 response=“6629fae49393a05397450978507c4ef1”,
 opaque=“5ccc069c403ebaf9f0171e9517f40e41”
```

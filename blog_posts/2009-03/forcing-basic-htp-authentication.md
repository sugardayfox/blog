Id: 14007
Title: Forcing basic http authentication for HttpWebRequest (in .NET/C#)
Tags: .net,c#
Date: 2009-03-14T16:12:55-07:00
Format: Markdown
--------------
HttpWebRequest is a handy .NET class for doing HTTP requests. It has
built-in support for HTTP basic authentication via credentials. However,
it doesn’t work the way I expected: supplying credentials doesn’t send
Authorization HTTP header with the request but only in response to
server’s challenge. It often breaks in real world, where servers might
not issue a challenge and simply not authenticate a request.

Fortunately fixing it by manually adding Authorization HTTP header to
the request is simple and this code snippet shows how to do it:

```c#
public void SetBasicAuthHeader(WebRequest req, String userName, String
userPassword)
{
 string authInfo = userName + “:” + userPassword;
 authInfo = Convert.ToBase64String(Encoding.Default.GetBytes(authInfo));
 req.Headers[“Authorization”] = “Basic ” + authInfo;
}
```

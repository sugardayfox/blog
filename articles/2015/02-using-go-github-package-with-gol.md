Id: 15
Title: Accessing GitHub API from Go
Date: 2015-02-12T23:18:58-08:00
Format: Markdown
tags: go, programming
--------------
@header-image gfx/headers/header-06.jpg
@collection go-cookbook

To use GitHub API from Go you can use [github.com/google/go-github](http://github.com/google/go-github) package.
GitHub, like many other sites, uses OAuth 2.0 protocol for authentication. There are few packages for Go that implement OAuth. This article describes how to use [golang.org/x/oauth2](http://golang.org/x/oauth2) package.
If you prefer to read code, read [this sample](https://github.com/kjk/kjkpub/blob/master/go/github_sample/sample_1.go).

## Short intro to OAuth for GitHub

To access private info about the user via the API in your web app you need to authenticate your app against GitHub website.

First, you need to [register your application with GitHub](https://github.com/settings/applications) (under Developer applications).
You provide description of the app, including a callback url (we'll get back to that).
GitHub creates client id and client secret that you'll need to provide to OAuth libraries (and keep them secret).

The flow of OAuth is:

* the user is on your website and clicks “login with GitHub” link
* you redirect the user to GitHub's authorization page. In that url you specify desired access level and a random secret
* the user  authorizes your app by clicking on a link
* GitHub redirects to a callback url on your website (which you provided when registering the app with GitHub)
* in the url handler, extract “secret” and “code” args
* you have to check that the secret is the same as the one you sent to GitHub (security measure that prevents forgery)
* you call another GitHub url to exchange code for access token

Access token is what you use to authenticate your API calls.

## Notes on callback url

You define callback url when you register the app. This is a problem for testing, because your app might use different hosts for development or production.
One option is to create different apps for testing and production, with different callback urls.
OAuth protocol allows providing callback url as an argument when calling authorization page (and oauth2 library allows setting it as `oauth2.Config.RedirectURL`) but GitHub [has restrictions](https://developer.github.com/v3/oauth/#redirect-urls) and doesn't allow changing host.

## Implementing OAuth using [golang.org/x/oauth2](http://golang.org/x/oauth2)

Create OAuth config object:

```go
import (
    "fmt"
    "net/http"

    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
    githuboauth "golang.org/x/oauth2/github"
)

var (
    // You must register the app at https://github.com/settings/applications
    // Set callback to http://127.0.0.1:7000/github_oauth_cb
    // Set ClientId and ClientSecret to
    oauthConf = &oauth2.Config{
        ClientID:     "",
        ClientSecret: "",
        // select level of access you want https://developer.github.com/v3/oauth/#scopes
        Scopes:       []string{"user:email", "repo"},
        Endpoint:     githuboauth.Endpoint,
    }
    // random string for oauth2 API calls to protect against CSRF
    oauthStateString = "thisshouldberandom"
)
```

This is your main html page:

```go
const htmlIndex = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`

// /
func handleMain(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(htmlIndex))
}
```

In the handler for /login url, redirect to GitHub's authorization page:

```go
// /login
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
    url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
```

GitHub will show authorization page to your user. If the user authorizes your app, GitHub will re-direct to OAuth callback. Here's how you can turn it into a token, token into http client and use that client to list GitHub information about user:

```go
// /github_oauth_cb. Called by github after authorization is granted
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
    if state != oauthStateString {
        fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    code := r.FormValue("code")
    token, err := oauthConf.Exchange(oauth2.NoContext, code)
    if err != nil {
        fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    oauthClient := oauthConf.Client(oauth2.NoContext, token)
    client := github.NewClient(oauthClient)
    user, _, err := client.Users.Get("")
    if err != nil {
        fmt.Printf("client.Users.Get() faled with '%s'\n", err)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }
    fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
```

And finally this is how you tie everything together:

```go
func main() {
    http.HandleFunc("/", handleMain)
    http.HandleFunc("/login", handleGitHubLogin)
    http.HandleFunc("/github_oauth_cb", handleGitHubCallback)
    fmt.Print("Started running on http://127.0.0.1:7000\n")
    fmt.Println(http.ListenAndServe(":7000", nil))
}
```

## Remembering OAuth access token

Once you retrieve the token, you need to remember it e.g. in a database. Here's how to serialize/deserialize token to/from JSON:

```go
func tokenToJSON(token *oauth2.Token) (string, error) {
    if d, err := json.Marshal(token); err != nil {
        return "", err
    } else {
        return string(d), nil
    }
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
    var token oauth2.Token
    if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
        return nil, err
    }
    return &token, nil
}
```

## Using personal access token

If you want to access someone else's GitHub info, you need to authorize your app with GitHub as described above.
If you want to access your own account, you can generate a [personal access token](https://github.com/settings/applications#personal-access-tokens). Here's how to do it:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
)

var (
    // you need to generate personal access token at
    // https://github.com/settings/applications#personal-access-tokens
    personalAccessToken = ""
)

type TokenSource struct {
    AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
    token := &oauth2.Token{
        AccessToken: t.AccessToken,
    }
    return token, nil
}

func main() {
    tokenSource := &TokenSource{
        AccessToken: personalAccessToken,
    }
    oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
    client := github.NewClient(oauthClient)
    user, _, err := client.Users.Get("")
    if err != nil {
        fmt.Printf("client.Users.Get() faled with '%s'\n", err)
        return
    }
    d, err := json.MarshalIndent(user, "", "  ")
    if err != nil {
        fmt.Printf("json.MarshlIndent() failed with %s\n", err)
        return
    }
    fmt.Printf("User:\n%s\n", string(d))
}
```

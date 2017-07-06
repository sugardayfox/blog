PublishedOn: 2017-07-02T00:00:00Z
Id: Jl3G
Title: HTTPS for free in Go, with little help of Let's Encrypt
Format: Markdown
Tags: for-blog, published, go
CreatedAt: 2017-06-23T09:56:26Z
UpdatedAt: 2017-07-06T06:09:00Z
--------------
@header-image gfx/headers/header-04.jpg
@collection go-cookbook

## Why HTTPS?

Having HTTPS for your website is important:

* HTTPS encrypts the traffic between browser and server. Passwords of your users are protected from traffic sniffing on naughty intermediary servers or miscreants sniffing wifi packets roaming through the air in a cafe.
* new HTTP/2 protocol is faster than HTTP/1.1 but only works over HTTPS.
* if you care about SEO, Google ranks HTTPS websites higher than HTTP ones.

## Using someone else to provide HTTPS

Before we learn how to support HTTPS directly in your Go web server, let's talk about simpler options.

You can use a third-party service like [Cloudflare](https://www.cloudflare.com/).

Their free plan offers acting as HTTPS proxy to your HTTP-only website.

To use Cloudflare:
* configure your domain to use CloudFlare's DNS servers
* in CloudFlare's DNS settings point the domain to your server and set Status to "DNS and HTTP proxy (CDN)"
* in CloudFlare's Crypto settings, set SSL to "Flexible" (browser talks to CloudFlare via HTTPS, CloudFlare talks to 
* configure CloudFlare's HTTPS proxy in their web interface by providing IP address of your server.
* for good measure, also enable "Always use HTTPS"

Browser talks to CloudFlare, which takes care of provisioning SSL certificate and proxies the traffic to your server. This might be slower due to additional traffic or faster due to CloudFlare servers being faster than yours (being faster is their business).

AWS, Google Cloud and some other hosting providers also provide free HTTPS for servers hosted on their infrastructures.

Another option is to run your server behind reverse-proxy server supporting HTTPS, like [Caddy](https://caddyserver.com/).

## Directly supporting HTTPS

Not long ago, if you wanted a SSL certificate, you had to pay many dollars a year for each domain.

[Let's Encrypt](https://letsencrypt.org/) changed that. It's a non-profit organization that provides certificates for free and offers HTTP API for obtaining certificates. API allows automating the process.

Before Let's Encrypt you would buy a certificate, which is just a bunch of bytes. You would save certificate to a file and configure your web server to use it.

With Let's Encrypt you can use their API to obtain the certificate for free, automatically, when your server starts. 

Thankfully all the hard work of talking to the API has already bee done by others. We just need to plug it in.

There are a couple of Go libraries that implement Let's Encrypt support.

I've been using [golang.org/x/crypto/acme/autocert](https://godoc.org/golang.org/x/crypto/acme/autocert), which is developed by Go core developers.

It's been several months now and it works flawlessly.

Here's how to start an HTTPS web server that uses free SSL certificates from Let's Encrypt.

Full example: [free-ssl-certificates/main.go](https://github.com/kjk/go-cookbook/blob/master/free-ssl-certificates/main.go).


```go
const (
	htmlIndex    = `<html><body>Welcome!</body></html>`
	inProduction = true
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, htmlIndex)
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleIndex)

	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func main() {
	var httpsSrv *http.Server
  
	// when testing locally it doesn't make sense to start
	// HTTPS server, so only do it in production.
	// In real code, I control this with -production cmd-line flag
	if inProduction {
		// Note: use a sensible value for data directory
		// this is where cached certificates are stored
		dataDir := "."
		hostPolicy := func(ctx context.Context, host string) error {
			// Note: change to your real domain
			allowedHost := "www.mydomain.com"
			if host == allowedHost {
				return nil
			}
			return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
		}

		httpsSrv = makeHTTPServer()
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: hostPolicy,
			Cache:      autocert.DirCache(dataDir),
		}
		httpsSrv.Addr = ":443"
		httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		go func() {
			err := httpsSrv.ListenAndServeTLS("", "")
			if err != nil {
				log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
			}
		}()
	}

	httpSrv := makeHTTPServer()
	httpSrv.Addr = ":80"
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}
}
```

There are some important things to note.

**1\.** The standard port for HTTPS is 443

**2\.** You can run only HTTP, only HTTPS or both.

**3\.** If the server doesn't have a certificate, it'll use HTTP API to ask Let's Encrypt servers for it.

Those requests are throttled to 20 per week to avoid over-loading Let's Encrypt servers.

It's therefore important to cache the certificate somewhere. In our example we cache them on disk, using [`autocert.DirCache`](https://godoc.org/golang.org/x/crypto/acme/autocert#DirCache) cache.

Cache is an interface so you could implement your own storage e.g. in a SQL database or Redis.

**4\.** You must set up DNS correctly. 

To verify that you're the owner of domain for which you want a certificate, Let's Encrypt server calls back your server.

For that to work, DNS name must resolve to the IP address of your server.

This means that local testing of HTTPS code-path is hard. I usually don't bother.

If you really want to, you can use [ngrok](https://ngrok.com/) to expose your local port to the internet, set DNS to resolve your domain to public DNS name that ngrok creates, wait a bit to make sure that DNS information propagates to Let's Encrypt computers.

**5\.** You might be wondering: what is this HostPolicy business?

As I mentioned, Let's Certificate throttles certificate provisioning so you need to ensure the server won't ask for certificates for domains you don't care about. Autocert docs [explain this well](https://github.com/golang/crypto/blob/5a033cc77e57eca05bdb50522851d29e03569cbe/acme/autocert/autocert.go#L104)

Our example assumes most common case: a server that only responds to a single domain. You can easily change the logic.

**6\.** We're not running HTTPS when testing locally.

When testing locally on your laptop, there's no point in running HTTPS version. Your computer most likely doesn't have publicly visible IP address so Let's Encrypt servers can't reach you, so you won't get the certificate.

We also won't be able to bind to HTTPS port 443 (only root processes can bind to ports lower than 1024).

In the example I use `inProduction` flag to decide if I should start HTTPS server.

In real code I add code that checks for `-production` cmd-line flag and use that.

## Redirecting from HTTP to HTTPS

If you can do HTTPS there's no point in providing plain HTTP. 

We can re-direct all HTTP request to HTTPS equivalent, for better security and SEO (Google doesn't like duplicate content so your SEO rank will be better with a single version of the website).

```go
func makeServerFromMux(mux *http.ServeMux) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func makeHTTPToHTTPSRedirectServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, req *http.Request) {
		newURI := "https://" + req.Host + req.URL.String()
		http.Redirect(w, req, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return makeServerFromMux(mux)
}

func main() {
	httpSrv := makeHTTPToHTTPSRedirectServer()
	httpSrv.Addr = ":80"
	fmt.Printf("Starting HTTP server on %s\n", httpSrv.Addr)
	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
	}
}
```

This concludes the technical part.

## How free certificates came to be

Arguably due to a design mistake, SSL protocol not only encrypts but also proves site's identity to the browser.

It provides accountability so that we can trace the ownership of google.com and see that it is indeed owned by Google, Inc in US, and not Ivan The Hacker in Moscow.

We implement that accountability by trusting a very small number of companies (Certificate Authorities) to issue certificates that prove the identity of the website owner.

When you apply for a certificate, Certificate Authority has to verify your identity. They do it by checking your papers.

Verifying the papers requires labor. Keeping certificates safe requires labor. It's reasonable that Certificate Authorities charge for the service of issuing certifcates.

The trust doesn't scale. Browser and OS vendors can trust 10 companies to not issue invalid certificates, but they can't trust a thousand.

We don't want any random company to become a rogue certificate authority and start issuing certificates for google.com domains to Ivan The Hacker.

It would be too much effort to continuosly audit thousands of Certificate Authority companies so as a result we ended up with just a few. 

A market controlled by small number of companies tends to become a cartel that keeps prices high due to lack of competition.

That's exactly what happened in SSL certificates market. You can have a low-end server for $60/year and a certificate alone would cost more than that.

That was a problem because the cost of SSL certificates was a significant barrier to adopting encryption by all websites.

A few companies decided to poll their resources and solve that problem for the greater good of the web.

They funded Let's Encrypt which became a Certificate Authority, wrote necessary software and is running the servers that do the work of issuing certificates. It's been a [raging success](//letsencrypt.org/2017/06/28/hundred-million-certs.html).

And that's how free certificates came to be.

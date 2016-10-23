Id: 1002039
Title: Buying a certificate for signing windows applications
Tags: programming
Date: 2012-04-15T22:28:01-07:00
Format: Markdown
--------------
Recently I’ve bought a code signing certificate so that I can sign my
Windows application
[SumatraPDF](http://blog.kowalczyk.info/software/sumatrapdf/free-pdf-reader.html)
(I’m hoping it’ll reduce number of false positive from various
anti-virus programs that claim that SumatraPDF is suspect).

There were some things that I wish I knew before starting the process,
so I’m documenting the process here for posterity.

There are many places to buy a code signing certificate. I bought it
from [K Software](http://ksoftware.net/) as they have low prices and
there’s enough info on the internet to assure me that it’s the right
kind of certificate for signing Windows apps. The certificate is
actually issued by Comodo, K Software is only a reseller (but they have
lower prices than Comodo - go figure).

I bought a certificate valid for 3 years (\$245). The validity only
affects whether you can use the certificate to sign. After you sign the
app, its signature is valid forever.

After certificate expires, you can renew it. The shortest (and cheapest)
validity period is 1 year. I opted for 3 years because it minimizes the
hassle of renewing every year.

Important note before you start: you also need to have some internet
domain registered in your name (or in your organization’s name) and to
minimize the troubles make sure there is a valid e-mail address with
that domain that you can receive (I use Google Apps for that domain and
use it to forward e-mails for that domain to my personal gmail account).

The domain is necessary to complete verification process of your
identity that Comodo does. It’s a strange requirement for a certificate
for signing applications but I’m guessing it’s because certificate have
roots in SSL/internet and in that case domain name is required.

Signing an app with a certificate basically serves as a stamp that says
“this application has been signed by company/individual X”. For the
system to work, X must be a legitimate company/person and not, say, a
hacker.

For that reason the organization that issues certificates (in this case
Comodo) needs to verify the identity of the person buying the
certificate so that e.g. I can’t order a certificate that says my name
is “Microsoft” and start signing my apps as coming from Microsoft.

The verification process starts after your purchase the certificate. The
details depend on whether you’re a company or an individual. I ordered
as an individual and the verification process was:

-   They asked for a copy of valid id. I e-mailed them a photo of my
    driver’s license (taken with iPhone)
-   Then they asked for a copy of a phone bill. I e-mailed them the PDF
    bill I downloaded from AT&T’s website
-   Then they called me on my phone to verify the phone number on the
    bill is my phone number

All in all, the back and forth took the whole day.

After the certificate is issued you need to download it to a file. To do
that you need to visit a web page that Comodo created for you in a
supported browser (FireFox or IE, Chrome is not supported, I used
FireFox) on the same computer that was used to order the certificate.

That creates a certificate and adds it to browser’s certificate store.
Finally, you export the certificate to a file. The steps are detailed at
<http://blog.ksoftware.net/>.

As to actual signing, I use K Software’s ksign tool (the command-line
version ksigncmd that I call from my build script).

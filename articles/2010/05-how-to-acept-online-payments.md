Id: 148001
Title: How to accept online payments
Tags: business
Date: 2010-05-05T00:33:33-07:00
Format: Markdown
--------------
Summary of [this
article](http://blog.meatinthesky.com/introduction-to-online-payments-tldr-its-a-to)
on how to accept online payments.

To accept credit cards you need an internet merchant account. Despite
the name, it’s not like a bank account. It only allows you to process
credit card transactions, but the money must be deposited to your real
checking account.

Use [TransFS](http://transfs.com/) to find a provider for internet
merchant account. Don’t use regular banks for it (like Bank Of America
etc.)

TransFS merchant accounts are “interchange plus” i.e. to process your
transaction they charge you interchange fee (fee that credit card
companies like Visa charge) plus a little bit on top of that. It takes
3-4 weeks to get merchant account.

Then you need to find a payment gateway. Gateway takes credit card
information that users input on your site and talk to the processors
(First Data, Paymentech, Global Payments etc.) to get the money
transfered to your merchant account (which then gets transfered to your
regular bank account).

You need to write the code to talk to a gateway yourself.

[Authorize.net](http://www.authorize.net/solutions/merchantsolutions/onlinemerchantaccount/)
is a popular gateway.
[Braintree](http://www.braintreepaymentsolutions.com/services/payment-gateway)
has better API but costs more and your merchant account has to support
First Data Nashville processor.

For micro transactions don use merchant account (fees too high). Use
[PayPal](http://www.paypal.com/IntegrationCenter/ic_micropayments.html)
or [Amazon
Payments](https://payments.amazon.com/sdui/sdui/about?nodeId=73479#suse_micro)

Or you can use PayPal, which has three offerings:

-   Website Payments standard:
    -   combined merchant account and gateway
    -   can only accept PayPal payments
-   Website Payments pro:
    -   combined merchant account and gateway
    -   can only accept credit card payments
-   PayFlow Pro - a standard gateway that works with any internet
    merchant account

PayPal is simpler and can even be cheaper. Downside: they will
oftentimes hold back up to 25% of your proceeds for three months as a
fraud prevention/risk management effort

Other options:

-   [Recurly](http://recurly.com/) helps to outsource subscription
    billings. Doesn’t require merchant account or payment gateway.
-   PayPal, Google Checkout, Amazon Payments offer similar service (see
    [this
    comparison](http://www.merchantaccountblog.com/1027/a-comparison-of-3-alternative-payment-methods))


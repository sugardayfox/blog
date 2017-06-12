Id: 14005
Title: setting up s3 logging
Tags: aws
Date: 2009-03-13T22:19:22-07:00
Format: Markdown
--------------
An example:

Getting acl on kjklogs bucket: `s3curl.pl —id kjk — -s -v http://s3.amazonaws.com/kjklogs?acl`

Ssetting acl on kjklogs bucket: `s3curl.pl —id kjk —put mylogs.acl — -s -v ‘http://s3.amazonaws.com/kjklogs?acl’`

Getting logging status for source bucket: `s3curl.pl —id kjk — -s -v http://s3.amazonaws.com/kjkpub?logging >kjkpub.logging.xml`

```xml
<Grant>
  <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
    <URI>http://acs.amazonaws.com/groups/s3/LogDelivery</URI>
  </Grantee>
  <Permission>WRITE</Permission>
</Grant>

<Grant>
  <Grantee xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="Group">
    <URI>http://acs.amazonaws.com/groups/s3/LogDelivery</URI>
  </Grantee>
  <Permission>READ_ACP</Permission>
</Grant>
```

Note: I mistakenly reset logging on a bucket by blindly using
[boto’s](http://boto.s3.amazonaws.com) `create_bucket()` function in my
scripts on an existing bucket. It didn’t destroy the existing data in
the bucket but it would reset logging permissions. A fix was to use
`get_bucket()` instead.

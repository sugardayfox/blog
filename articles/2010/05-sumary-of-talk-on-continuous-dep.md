Id: 134001
Title: Summary of talk on continuous deployment
Date: 2010-05-03T12:28:26-07:00
Format: Markdown
--------------
Summary of [this
talk](http://www.justin.tv/startuplessonslearned/b/262657677) from
Startup Lessons Learned conference about transitioning to continuous
deployment.

WiredReach: simple sharing software (BoxCloud, CloudFire), 7 years in
business. Hybrid web/desktop model.

Before continuous deployment releases were on a 2 week cycle: 1 week
dev, 1 week QA.

After: release multiple times a day.

Before: common staging area. After: standalone sandboxes for each dev
and for QA.

Before: release was all day event (code freeze, testing, packaging
etc.). After: release is a non-event.

Release is triggered by a checkin. It runs a battery of tests and only
goes into production if tests pass. The whole process takes \<20min.

Before: release changed hundreds lines of code. After: \< 25 lines of
code.

Switching to continuous deployment was scary, feeling of lack of safety
net.

**How they transitioned to continuous deployment**

Code in small batch sizes (2 hrs worth of coding).

Gradual automation: first deployed manually, then automated more and
more.

Wrote functional tests, starting with test for user activation
(registration etc.)

Functional tests take time, they wanted release cycle to be \<30 min.
Solution: parallelize tests on multiple machines.

Problem: with time tests get out of date and start failing. Solution: a
rule that says can only deploy if all tests pass. If test became
invalid, had to fix it.

Incrementally build cluster immune system, which can detected bad
changes. Started by monitoring with off-the-shelf tools like nagios,
ganglia, over time added custom monitoring, business level metrics.\
Used 5 whys based on production problems to figure out what to monitor.

For client software developed background update process, ability to
selectively push updates.

**How they develop new features now**

Remaining problem: how to know they are implementing the right features
and not just implementing random features faster?

They spend more time measuring and optimizing existing features than
adding new features.

Constrain features pipeline. Don’t start working on new features until
deployed features have been validated.

**How to validate a feature?**

Start with qualitative analysis: contact customers who asked for the
feature and get their feedback, focusing not on coolness but on whether
it solved their problem and made the difference in ability to make or
keep the sale.

On quantitative side they use mixpanel, kissmetrics, google analytics
and focus on macro-level changes (have 5 metrics they track: revenue,
retention etc.)

**My commentary**

I was disappointed that the talk doesn’t talk about “cluster immune
system” in more technical detail. The idea is easy to express but the
actual implementation seems to me frighteningly complex, especially
given the need for testing new code on realistic data (e.g. production
accounts) and at the same time in a sandboxed environment that is not
production.

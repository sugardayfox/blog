Id: 8047
Title: Summary of David Ditzel talk on binary translation
Tags: talk,summary
Date: 2009-02-22T18:12:55-08:00
Format: Markdown
--------------
Summary of David Ditzel
[talk](http://norfolk.cs.washington.edu/htbin-post/unrestricted/colloq/details.cgi?id=759)
on binary translation.

David Ditzel: worked at Transmeta, now at Intel

A 25 year perspective on binary translation: what worked, what didn’t
work.

Examples of binary translation:

-   pentium pro and later: translates x86 into internal UOPS via
    hardware
-   intel ia32el runs user level x86 programs on Itanium with \~60% perf
    of native machine
-   Java JITs, .NET (MSIL), VmWare
-   Apple’s Rosetta runs PowerPC user programs on x86
-   Transmeta Crusoe and Efficeon processors - transparent, full system
    level translation of x86

SYMBOL computer - implemented os, editor etc. in logic gates. Lessons:
don’t do that. Use the right combination of software, hardware and
micro-ops.

AT&T Crisp:

-   1987, first CMOS superscalar chip (superscalar: multiple
    instructions per clock)
-   Reduced Instruction Set Processor targeted at C
-   hw translated instructions from external, compact version into 180
    bit wide UOP cache
-   optimization tricks:\
     \* branch folding - make branches disappear from pipeline\
     \* “Stack Cache” as registers - to reduce memory references

Lessons from AT&T Crisp:

-   translation from external instruction set to internal instruction
    set works well

Binary Instrumentation

MIPS had tools pixie and pixstats (\~1987) to statically modify binaries
to count instructions.

Sun followed (\~1988) with spix, spixstats etc. Also were able to run
MIPS on Sparc (at 1/3rd speed).

Sun tried to extend Sparc with instructions to help x86 emulation but
decided that hw mismatch was too big - Sparc was not the right
architecture for this.

Lessons realized in 1995:

-   dynamic translation was reaching 1/3 speed of native, static 1/2 of
    native
-   processor designed from scratch for binary translation might improve
    efficiency of dynamic translation by 2-3x
-   full system level binary translation might soon become practical and
    even exceed perf of standard microprocessor

That led to Transmeta in 1995. Transmeta:

-   spent \$600M over 12 years in R&D
-   5 generations of processors although only 3 announced

Key challenges for hybrid processors:

-   must be 100% compatible. When doing binary translation between
    commodity processors (e.g. x86\<-\>PowerPC), there are corner cases
    where emulating things exactly is tricky due to hardware mismatch,
    which causes inefficiencies in translation =\> must design for
    translation
-   precise control over user visible state, including precise exception
    and interrupt semantics
-   delivered performance, including overhead of translation

Hardware support for hybrid processors:

-   private, non-volatile storage (FLASH ROM), for storing translation
    software
-   private memory for storing translated code (stole 5% of DRAM during
    boot)
-   software controlled state commit/rollback/abort
-   more registers than x86
-   alias detection under software control
-   fine grain detection of self-modifying code
-   auto-typing of pure memory vs I/O (because I/O can be memory-mapped
    which prohibits some optimizations)
-   fast traps supported by underlying runtime system
-   instruction primitives for fast interpretation

Software controlled atomic execution - execute in temporary space and
ability to rollback to previous commit point. Used to perform
not-always-safe optimization which are considered ok as long as we hit
next commit point without problems. If not, rollback and re-execute
without optimizations. Needed to be able to undo stores to memory.

Transmeta’s code morphing:

-   first level - interpreter
-   second compiler
-   translation - it can cost 10000 instructions to translate 1

Efficeon improvements used 4 levels (gears):

-   first level - interpreter 15 instruction per 1, gathers
-   after executing basic block 50 uses quick translation (cost 500
    instructions per 1 native), also gathers more information
-   after executing more than few hundred times - more optimized
    translation, more costly, classic optimization like common
    sub-expression elimination, memory re-ordering, critical path
    scheduling
-   for hot loops, optimize multiple code blocks and use more aggressive
    optimization

Lessons: optimization pay off. The bigger the blocks, the bigger
optimization payoff.

Binary translation myths:

-   myth: translation is slow. It’s only \~20% overhead
-   myth: saving translation to disk is a good idea. In efficeon they
    improved translation so that it was faster to translate than read
    from disk
-   myth: static translation is faster. They compiled Linux kernel to
    run natively but it was slower than dynamic translation because
    dynamic translation could use runtime information to optimize.
-   myth: software isn’t reliable. Doesn’t match transmeta’s nor
    Transiitive’s experience: they didn’t have any x86 compatibility
    bugs

Why binary translation now: power usage since increasing cores requires
more power so we might not have enough power to light up all processors
at full speed.

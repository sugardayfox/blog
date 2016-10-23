Id: 306
Title: High-resolution timer for timing code fragments
Tags: win32,c,programming
Date: 2005-12-31T16:00:00-08:00
Format: Markdown
--------------
Timing of code fragments is simple:

* record the start time
* execute the code
* record the end time
* calucate the difference between end and start

For timing exuction time of pieces of code you need a high-resolution timer. On Windows, Windows CE and Pocket PC/Smartphone (which are Windows CE variations) you can use `QueryPerformanceCounter` and `QueryPerformanceFrequence` API calls.

Here's a simple implementation in C:

```c
typedef struct prof_timer_t {
    LARGE_INTEGER time_start;
    LARGE_INTEGER time_stop;
} prof_timer_t;

void prof_timer_start(prof_timer_t *timer) {
    QueryPerformanceCounter(&timer->time_start);
}

void prof_timer_stop(prof_timer_t *timer) {
    QueryPerformanceCounter(&timer->time_stop);
}

double prof_timer_get_duration_in_secs(prof_timer_t *timer) {
    LARGE_INTEGER freq;
    double duration;
    QueryPerformanceFrequency(&freq);
    duration = (double)(timer->time_stop.QuadPart-timer->time_start.QuadPart)/(double)freq.QuadPart;
    return duration;
}
```

And in C++:

```c++
// very simple, high-precision (at least in theory) timer for timing API calls
struct ProfTimer {
    void Start(void) {
        QueryPerformanceCounter(&mTimeStart);
    };
    void Stop(void) {
        QueryPerformanceCounter(&mTimeStop);
    };
    double GetDurationInSecs(void)
    {
        LARGE_INTEGER freq;
        QueryPerformanceFrequency(&freq);
        double duration = (double)(mTimeStop.QuadPart-mTimeStart.QuadPart)/(double)freq.QuadPart;
        return duration;
    }

    LARGE_INTEGER mTimeStart;
    LARGE_INTEGER mTimeStop;
};
```

And here's an example of using the C++ version:

```c++
    ProfTimer t;
    t.Start();
    foo();
    t.Stop();
    double dur = t.GetDurationInSecs();
    printf("executing foo() took %f seconds\n" , dur);
</code>

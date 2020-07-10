# gologmine

Log clustering library heavily based on the LogMine algorithm ( https://www.cs.unm.edu/~mueen/Papers/LogMine.pdf )

## Details

When viewing thousands/tens of thousands or even millions of lines of logs, sometimes it is more useful to just see the patterns that are appearing rather than the exact log entries that are emerging. 

For example, if you have an application starting to log the message "disconnected socket against client IP <ip>" where the IP is dynamic, then unless you've already set up alerts against the beginning of the string, you might miss these warnings signs.

By being able to generate patterns based on your logs (without any additional human configuration/tuning) and alert on a new pattern being discovered can definitely be useful.

This was intiailly intended to be a complete implementation of the LogMine algorithm, but based off the public papers, presentations and videos I was honestly unable to conclude some of the details required. So this is an implementation of what I could determine plus some "fill in the blanks".

An example from the LogMine presentations/papers has the following example log:

```
2017/02/24 09:01:00 login 127.0.0.1 user=bear12
2017/02/24 09:02:00 DB Connect 127.0.0.1 user=bear12
2017/02/24 09:02:00 DB Disconnect 127.0.0.1 user=bear12
2017/02/24 09:04:00 logout 127.0.0.1 user=bear12
2017/02/24 09:05:00 login 127.0.0.1 user=bear34
2017/02/24 09:06:00 DB Connect 127.0.0.1 user=bear34
2017/02/24 09:07:00 DB Disconnect 127.0.0.1 user=bear34
2017/02/24 09:08:00 logout 127.0.0.1 user=bear34
2017/02/24 09:09:00 login 127.0.0.1 user=bear#1
2017/02/24 09:10:00 DB Connect 127.0.0.1 user=bear#1
2017/02/24 09:11:00 DB Disconnect 127.0.0.1 user=bear#1
2017/02/24 09:12:00 logout 127.0.0.1 user=bear#1
```


In this case you can see that if you ignore timestamps, IPs and even user names, then you can see there are basically 4 patterns of log messages that are appearing. "login", "logout", "DB Connect" and "DB Disconnect". You might be in a situation where you're interested in these patterns.

When using this tool (or as a library) you can mine the logs and cluster the messages into groups of similar messages. Then you can keep expanding the boundries of the clusters depending on how much detail is required).

In this specific case, we can perform the most general form of clustering with the results:

```
count 1 : pattern DATE TIME DB Disconnect IPV4 user = NOTSPACE
count 1 : pattern DATE TIME logout IPV4 user = NOTSPACE
count 1 : pattern DATE TIME login IPV4 user = NOTSPACE
count 1 : pattern DATE TIME DB Connect IPV4 user = NOTSPACE
count 2 : pattern DATE TIME DB Disconnect IPV4 user = WORD
count 2 : pattern DATE TIME logout IPV4 user = WORD
count 2 : pattern DATE TIME login IPV4 user = WORD
count 2 : pattern DATE TIME DB Connect IPV4 user = WORD
```

In this case you can see some elements of the messages have been generalised out (time, date, unique words etc)

Next, we can generalise the clustering to the next level and get:

```
count 3 : pattern DATE TIME login IPV4 user = NOTSPACE
count 3 : pattern DATE TIME DB Connect IPV4 user = NOTSPACE
count 3 : pattern DATE TIME DB Disconnect IPV4 user = NOTSPACE
count 3 : pattern DATE TIME logout IPV4 user = NOTSPACE
```

Again, we can seee we have 12 entries (totalling the counts) and we can also see that in general we have 4 types of messages. login, logout, DB Connect and DB Disconnect.

If we take it one step further we start getting into the situation of very generalised messages:

```
count 6 : pattern DATE TIME WORD IPV4 user = NOTSPACE
count 6 : pattern DATE TIME DB WORD IPV4 user = NOTSPACE
```

Here we're simply breaking it down to messages that happened to have "DB" in it and those that didn't.

We can take it one final level, but in this situation (although producing correct results) aren't overly useful:

```
count 12 : pattern DATE TIME * WORD IPV4 user = NOTSPACE
```







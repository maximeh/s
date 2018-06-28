s
=

A pet project to have some fun in Go.
It's called s because you use it to *s*erve file or directory.

This is a rip off of the idea of [Simon Budig's woof](http://www.home.unix-ag.org/simon/woof.html).

You can serve either:
   - serve a file (1 or N times; at the end of the count; the app closes)
   - serve a directory until you close the app. (useful for dev purposes)

The default port is 4242 (so you don't need special ACL to run it; you can
change it with '--port'.

It will listen on all your interfaces; so if the IP found is not the one you
want to use, you can simply change the URL.

#Install

Still here ?  Wow.

Well, if your *REALLY* want to use it, you can do:

```sh
go install github.com/maximeh/s
make
```

And that should be it, it should download and compile everything.

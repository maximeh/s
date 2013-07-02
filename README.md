s
=

A pet project to have some fun in Go.
It's called s because you use it to *s*erve file or directory.

This is a rip off of the idea of [Simon Budig's
woof](http://www.home.unix-ag.org/simon/woof.html) but with less capability (I
know...).

You can serve a file once (and only once currently) or a directory until you
close the app.
The only parameter let you change the default port (which is 4242).

You can't setup the IP you want, it will try to find one itself (using some
stolen method from Woof) and it may not work with IPv6.

#Usage

Still here ?  Wow.

Well, if your *REALLY* want to use it, you can do:

```sh
go install github.com/maximeh/s
```

And that should be it, it should download everything, compile it and install it.



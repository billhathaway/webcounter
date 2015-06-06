Webcounter challenge
==
This project is based on https://github.com/SDGophers/2015-04-challenge and currently does the initial challenge phase and the extra-credit of supporting multiple image types based on the request suffix (.txt/.jpg/.gif/.png).

The main logic is exposed as a library github.com/billhathaway/webcounter  

There is an executable available under the 'server' directory, which will run on port 8080 by default.  

Startup is a little slow since per the challenge guidelines, it fetches the source image from the challenge's GH repo.  

Running
--
    go get github.com/billhathaway/webcounter/wcserver

    wcserver // listening on port 8080 by default

Microbenchmarks
--
Made using wrk against server running on localhost, iMac Mid-2011 (3.1Ghz Core i5), GOMAXPROCS=4  

Request Type| Rate/sec
------------|-------:
.txt|74555
.png|267
.jpg|480
.gif|55

Testing
--
    go test -v

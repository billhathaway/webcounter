Webcounter challenge
==
[![Build Status](https://travis-ci.org/billhathaway/webcounter.svg?branch=master)](https://travis-ci.org/billhathaway/webcounter)  

This project is based on https://github.com/SDGophers/2015-04-challenge and currently handles:  
* initial challenge phase
* extra-credit of supporting multiple image types based on the request suffix (.txt/.jpg/.gif/.png)  
* 1,000 visitor level of per-referer counts  

The main logic is exposed as a library github.com/billhathaway/webcounter  

There is an executable available under the 'server' directory, which will run on port 8080 by default.  

Startup is a little slow (3-5 secs) since per the challenge guidelines, it fetches the source image for the glyphs from the challenge's GH repo.  

Running
--
    go get github.com/billhathaway/webcounter/wcserver

    wcserver // listening on port 8080 by default


hit server at http://localhost:8080/someid.png and refresh your browser a bunch and watch the count go up  

![Screen Shot of app](/images/screenshot.png)  

Microbenchmarks
--
Made using wrk against server running on localhost, iMac Mid-2011 (3.1Ghz Core i5), GOMAXPROCS=4  

Request Type| Rate/sec
------------|-------:
.txt|74555
.jpg|480
.png|267
.gif|55

Testing
--
    go test -v

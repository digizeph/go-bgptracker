# go-bgptracker 

[![Build Status](https://travis-ci.org/digizeph/go-bgptracker.svg?branch=master)](https://travis-ci.org/digizeph/go-bgptracker)

A simple web application that shows the oldest and most recent BGP dump files from all the collectors in RIPE RRC and RouteViews.

An demo is running at http://ht3.mwzhang.com:9999/

![](https://screenshots.firefoxusercontent.com/images/0c442d1c-cd3a-4a67-a38c-55db414716b7.png)

## Build and Run

```
go install github.com/digizeph/go-bgptracker/
go-bgptracker
```

On default, the service runs on port `9999`. You can change the part in `main.go` where it says `":9999"` to any other port number as you wish.

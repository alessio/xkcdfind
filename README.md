# xkcdfind [![Travis-CI Status](https://api.travis-ci.org/alessio/xkcdfind.png?branch=master)](http://travis-ci.org/#!/alessio/xkcdfind)

__xkcdfind__ is a small command line tool that downloads and maintains an offline index of [Xkcd](https://xkcd.com/)'s comics descriptors. It features a simple search engine which can be used to find comics that match the search terms provided on the command line. Regular expressions are supported too.

# Installation

```bash
$ go get github.com/alessio/xkcdfind
```

# Usage

```console
alessio@bizet:~$ xkcdfind -h
Usage of xkcdfind:
  -index string
    	Index file (default: 'index.json')
  -update
    	Force the update of the index
```

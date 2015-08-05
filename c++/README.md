StaticServer
==================

[![Build Status](https://travis-ci.org/breezechen/static_server.svg?branch=master)](https://travis-ci.org/breezechen/static_server)

A cross-platform http static server based on boost asio example http server2. Now it's supporting large file(> 4G) download.

How to build?
============

**Boost(>= 1.50) should be installed.**

Windows
-------
Open .sln in vs2008, change include dir to yours, then compile.

Linux
-----
Install boost if needed.
``` shell
cd $BOOST_ROOT; 
./bootstrap.sh --with-libraries=system,filesystem,thread
sudo ./b2 threading=multi link=static --prefix=/usr/local -d0 install
```
```shell
g++ *.cpp -o server -pthread -lboost_system -lboost_filesystem -lboost_thread -lrt -O2 -DNDEBUG
```


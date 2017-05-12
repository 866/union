# Union HTTP Server #

This is a web server written in Go for communication with clients/*Trigger* server.
*Trigger* server is not for public use so it is on the private repository.
If you want to collaborate please contact me <mailto:comrazvictor@example.com>

### What is this server for? ###

* Playing on stock markets
* Communication with Trigger TCP/IP server
* Providing the service for communication between traders(chat, forum)
* Providing web service for trading for the large number of clients


### How do I get set up? ###

Do the following steps:

* Clone the repo:

```
#!bash

git clone  https://github.com/866/union
```

* Go to the source directory:

```
#!bash

cd union
```

* Install all necessary packages:

```
#!bash

glide install
```

* If you are using Windows:
  - Install [mingw-w64](https://sourceforge.net/projects/mingw-w64/)
  - Setup the next environment variables(make sure that mingw/bin in your PATH):
  ```
  CXX=x86_64-w64-mingw32-g++
  CC=x86_64-w64-mingw32-gcc
  ```

* Build the project:

```
#!bash
go build
```

### Contribution guidelines ###

* Write clean and commented code
* Do Gometalinter checking
* Use glide for development

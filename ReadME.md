# S-Org #

In `/handlers` there is the code to use for the api endpoints.

In `/resources` are stored files such as dynamodb  and api gateway configuration.

Makefile has some script for making development easier.

* Makefile
  * clean - cleans all the binaries in the `/bin` dir
  * build -  creates the binaries and puts them in the `/bin` dir
  * deploy - delpoys the stack to aws
  * db - starts the database

# S-Org #

In `/handlers` there is the code to use for the api endpoints.

In `/resources` are stored files such as dynamodb  and api gateway configuration.

Makefile has some script for making development easier.

* Makefile
  * clean - cleans all the binaries in the `/bin` dir
  * build -  creates the binaries and puts them in the `/bin` dir
  * deploy - delpoys the stack to aws
  * db - starts the database



### Status Codes ###

|Code|Message|
|----|-------|
|1|Invalid body|
|2|Bcrypt error|
|3|DynamoDBError| <!--- Not sure but meybe sub errors can be 30,31 , 32, 33 .... --->
|4|Could not find user|
|5|Password incorect|
|6|Error by getting key|
|7|Error by encrypting token|
|8|Error by getting key but previous action succesfull|
|9|Error by encrypting token but previous action successfull|
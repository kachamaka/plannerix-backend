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

|   |   |   |   |   |   |   |   |
|---|---|---|---|---|---|---|---|
|100|Invalid Body|200|MarshalMapError|300|PutItemError - Users|42|NO ERROR|
|101|Invalid Username|201|MarshalJSONToMap|301|PutItemError - Schedules|400|Invalid|
|102|Invalid Pass|202|MarshalListError|302|PutItemError - Grades|401|DoesNotExist|
|103|Invalid Subjects|203|UnmarshalMapError|303|PutItemError - Subjects|402|ErrorWith|
|104|Invalid Schedule|204|UnmarshalListOfMapsError|304|PutItemError - Events|403|Internal Server Error|
|105|Invalid Email|205|OutputError|305|DeleteItemError - Grades|404|NotFound|
|106|Email does not exist|206|ExpressionBuilderError|306|DeleteItemError - Events|---|---|
|107|Hash Error|---|---|307|UpdateItemError - Events|---|---|
|108|Username taken|---|---|308|UpdateItemError - Schedules|---|---|
<!-- 
|Code|Message|
|----|-------|
|100|Invalid Body|
|101|Invalid Username|
|102|Invalid Pass|
|103|Invalid Subjects|
|104|Invalid Schedule|
|105|Invalid Email|
|106|Email does not exist|
|107|Hash Error|
|108|Username taken|
|200|MarshalMapError|
|201|MarshalJSONToMap|
|202|MarshalListError|
|203|UnmarshalMapError|
|204|UnmarshalListOfMapsError|
|205|OutputError|
|206|ExpressionBuilderError|
|300|PutItemError - Users|
|301|PutItemError - Schedules|
|302|PutItemError - Grades|
|303|PutItemError - Subjects|
|304|PutItemError - Events|
|305|DeleteItemError - Grades|
|306|DeleteItemError - Events|
|307|UpdateItemError - Events|
|308|UpdateItemError - Schedules|
|400|Invalid|
|401|DoesNotExist|
|402|ErrorWith|
|403|Internal Server Error|
|404|NotFound|
|42|NO ERRORS| -->

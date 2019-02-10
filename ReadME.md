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

|Code|Message||Code|Message||Code|Message||Code|Message||Code|Message|
|----|-------||----|-------||----|-------||----|-------||----|-------|
|100|Invalid Body||200|MarshalMapError||300|PutItemError - Users||400|Invalid|
|101|Invalid Username||201|MarshalJSONToMap||301|PutItemError - Schedules||401|DoesNotExist|
|102|Invalid Pass||202|MarshalListError||302|PutItemError - Grades||402|ErrorWith|
|103|Invalid Subjects||203|UnmarshalMapError||303|PutItemError - Subjects||403|Internal Server Error|
|104|Invalid Schedule||204|UnmarshalListOfMapsError||304|PutItemError - Events||404|NotFound|
|105|Invalid Email||205|OutputError||305|DeleteItemError - Grades|
|106|Email does not exist||206|ExpressionBuilderError||306|DeleteItemError - Events|
|107|Hash Error||||||||||307|UpdateItemError - Events|
|108|Username taken||||||||||308|UpdateItemError - Schedules|

|42|NO ERRORS|

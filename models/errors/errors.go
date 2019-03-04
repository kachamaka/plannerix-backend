package errors

import "errors"

var (
	LambdaError              = errors.New("Lambda Error")                 //-1
	InvalidBody              = errors.New("Invalid Body")                 //100
	KeyError                 = errors.New("Getting Key Error")            //109
	TokenError               = errors.New("Token Error")                  //110
	MarshalMapError          = errors.New("Marshal Map Error")            //200
	MarshalJsonToMapError    = errors.New("Marshal JSON To Map Error")    //201
	MarshalListError         = errors.New("Marshal List Error")           //202
	UnmarshalMapError        = errors.New("Unmarshal Map Error")          //203
	UnmarshalListOfMapsError = errors.New("Unmarshal List Of Maps Error") //204
	OutputError              = errors.New("Output Error")                 //205
	ExpressionBuilderError   = errors.New("Expression Builder Error")     //206
	PutItemError             = errors.New("Error with PutItem")
	//300 - users
	//301 - schedules
	//302 - grades
	//303 - subject
	//304 - events
	DeleteItemError = errors.New("Error with DeleteItem")
	//305 - grades
	//306 - events
	UpdateItemError = errors.New("Error with UpdateItem")
	//307 - events
	//308 - schedule
)

func NotFound(something string) error {
	return errors.New(something + " not found")
} //404
func Invalid(something string) error {
	return errors.New("Грешка с " + something)
} //400
func DoesNotExist(something string) error {
	return errors.New(something + " не съществува")
} //401
func ErrorWith(something string) error {
	return errors.New("Грешка с " + something)
} //402

// Invalid Username 101
// Invalid Pass 102
// Invalid Subjects 103
// Invalid Schedule 104
// Invalid Email 105
// Email does not exist 106
// Hash Error 107
// Username taken 108
// Internal server error 403
//
//
//
//
//
//
//
//
// NO ERRORS 42

package errno


var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK."}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error."}
	BadRequest          = &Errno{Code: 10002,Message: "Bad request."}

	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}
	ErrParseURLArgs =&Errno{Code:2004,Message:"Parse URL args error."}
	ErrDecodeToken =&Errno{Code:2005,Message:"Error occirred while Decode jwt str"}


	// user errors
	ErrEncrypt           = &Errno{Code: 20101, Message: "Error occurred while encrypting the user password."}
	ErrUserNotFound      = &Errno{Code: 20102, Message: "The user was not found."}
	ErrTokenInvalid      = &Errno{Code: 20103, Message: "The token was invalid."}
	ErrPasswordIncorrect = &Errno{Code: 20104, Message: "The password was incorrect."}
	ErrNeedToken         = &Errno{Code: 20105,Message:"Need token to access"}
	ErrTokenExpired      = &Errno{Code:20106,Message: "Token expired"}
	ErrTokenRequired     = &Errno{Code:20107,Message:"Token required"}
	ErrUsernamePasswordRequired =&Errno{Code:20108,Message:"Username and Password required"}
	ErrUnauth           =&Errno{Code:20109,Message:"Unauthorized"}

	ErrGetNamespace        =&Errno{Code:20201,Message:"Error occurred while get namespace"}
	ErrNameSpaceExists    =&Errno{Code:20202,Message:"Namespace already exists"}
	ErrCreateNamespace    =&Errno{Code:20203,Message:"Error occurred while create namespace"}
	ErrDeleteNamespace    =&Errno{Code:20204,Message:"Error occurred while delete namespace"}
)
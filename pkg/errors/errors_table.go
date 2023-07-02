package errors

type ErrorCode int

const (
	errGenericErrorCode ErrorCode = iota
	ErrInvalidUserMail
	ErrInvalidUserName
	ErrInvalidPassword
	ErrUserAlreadyExists
	ErrUserCreationFailure
	ErrNoSuchUser

	ErrNoSuchToken
	ErrTokenAlreadyExists
	ErrAuthenticationFailure
	ErrAuthenticationExpired
	ErrNotLoggedIn

	ErrFailedToGetBody
	ErrBodyParsingFailed

	ErrNoSuchHeader
	ErrNonUniqueHeader

	ErrNoResponse
	ErrResponseIsError

	ErrPostInvalidData
	ErrPostRequestFailed
	ErrGetRequestFailed

	ErrDbConnectionInvalid
	ErrInvalidQuery
	ErrInvalidSqlTable
	ErrInvalidSqlProp
	ErrDuplicatedSqlProp
	ErrInvalidSqlScript
	ErrInvalidSqlScriptArg
	ErrSqlTranslationFailed

	ErrNotImplemented

	lastErrorCode
)

var errorsCodeToMessage = map[ErrorCode]string{
	ErrInvalidUserMail:     "user mail is invalid",
	ErrInvalidUserName:     "user name is invalid",
	ErrInvalidPassword:     "password is invalid",
	ErrUserAlreadyExists:   "user already exists",
	ErrUserCreationFailure: "internal error while creating user",
	ErrNoSuchUser:          "no such user",

	ErrNoSuchToken:           "no such token",
	ErrTokenAlreadyExists:    "token already exists",
	ErrAuthenticationFailure: "authentication failure",
	ErrAuthenticationExpired: "authentication expired",
	ErrNotLoggedIn:           "not logged in",

	ErrFailedToGetBody:   "failed to get request body",
	ErrBodyParsingFailed: "failed to parse request body",

	ErrNoSuchHeader:    "no such header in request",
	ErrNonUniqueHeader: "header is defined multiple times in request",

	ErrNoResponse:      "no response",
	ErrResponseIsError: "response returned error code",

	ErrPostInvalidData:   "invalid post request data",
	ErrPostRequestFailed: "post request failed",
	ErrGetRequestFailed:  "get request failed",

	ErrDbConnectionInvalid:  "db connection is invalid",
	ErrInvalidQuery:         "invalid sql query",
	ErrInvalidSqlTable:      "invalid table for sql query",
	ErrInvalidSqlProp:       "invalid property for sql query",
	ErrDuplicatedSqlProp:    "duplicated property for sql query",
	ErrInvalidSqlScript:     "invalid script for sql query",
	ErrInvalidSqlScriptArg:  "invalid script argument for sql query",
	ErrSqlTranslationFailed: "failed to generate sql query",

	ErrNotImplemented: "not implemented",
}

var defaultErrorMessage = "unexpected error occurred"

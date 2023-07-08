package errors

type ErrorCode int

const (
	errGenericErrorCode ErrorCode = iota
	ErrInvalidUserMail
	ErrInvalidUserName
	ErrInvalidPassword
	ErrUserAlreadyExists
	ErrUserCreationFailure
	ErrUserGetFailure
	ErrUserDeletionFailure
	ErrNoSuchUser

	ErrFailedToGetBody
	ErrBodyParsingFailed

	ErrNoSuchHeader
	ErrNonUniqueHeader

	ErrNoResponse
	ErrResponseIsError

	ErrPostInvalidData
	ErrPostRequestFailed
	ErrGetRequestFailed

	ErrDbConnectionFailed
	ErrDbConnectionInvalid
	ErrInvalidQuery
	ErrInvalidSqlTable
	ErrInvalidSqlProp
	ErrDuplicatedSqlProp
	ErrInvalidSqlFilter
	ErrInvalidSqlScript
	ErrInvalidSqlScriptArg
	ErrSqlTranslationFailed
	ErrNoPropInSqlSelectQuery
	ErrInvalidSqlComparisonKey
	ErrInvalidSqlComparisonValue
	ErrNoValuesInSqlComparison
	ErrInvalidSqlColumn
	ErrDuplicatedSqlColumn
	ErrNoColumnInSqlInsertQuery
	ErrNoColumnInSqlUpdateQuery

	ErrDbCorruptedData
	ErrDbRequestCreationFailed
	ErrDbRequestFailed
	ErrMultiValuedDbElement
	ErrInvalidSqlQueryReceiverType
	ErrNoRowsReturnedForSqlQuery
	ErrSqlRowParsingFailed

	ErrNotImplemented

	lastErrorCode
)

var errorsCodeToMessage = map[ErrorCode]string{
	ErrInvalidUserMail:     "user mail is invalid",
	ErrInvalidUserName:     "user name is invalid",
	ErrInvalidPassword:     "password is invalid",
	ErrUserAlreadyExists:   "user already exists",
	ErrUserCreationFailure: "error while creating user",
	ErrUserGetFailure:      "error while getting user",
	ErrUserDeletionFailure: "error while deleting user",
	ErrNoSuchUser:          "no such user",

	ErrFailedToGetBody:   "failed to get request body",
	ErrBodyParsingFailed: "failed to parse request body",

	ErrNoSuchHeader:    "no such header in request",
	ErrNonUniqueHeader: "header is defined multiple times in request",

	ErrNoResponse:      "no response",
	ErrResponseIsError: "response returned error code",

	ErrPostInvalidData:   "invalid post request data",
	ErrPostRequestFailed: "post request failed",
	ErrGetRequestFailed:  "get request failed",

	ErrDbConnectionFailed:        "db connection failed",
	ErrDbConnectionInvalid:       "db connection is invalid",
	ErrInvalidQuery:              "invalid sql query",
	ErrInvalidSqlTable:           "invalid table for sql query",
	ErrInvalidSqlProp:            "invalid property for sql query",
	ErrDuplicatedSqlProp:         "duplicated property for sql query",
	ErrInvalidSqlFilter:          "invalid filter for sql query",
	ErrInvalidSqlScript:          "invalid script for sql query",
	ErrInvalidSqlScriptArg:       "invalid script argument for sql query",
	ErrSqlTranslationFailed:      "failed to generate sql query",
	ErrNoPropInSqlSelectQuery:    "no property set for sql query",
	ErrInvalidSqlComparisonKey:   "invalid comparison key for sql query",
	ErrInvalidSqlComparisonValue: "invalid comparison value for sql query",
	ErrNoValuesInSqlComparison:   "no comparison values set for sql query",
	ErrInvalidSqlColumn:          "invalid column for sql query",
	ErrDuplicatedSqlColumn:       "duplicated column for sql query",
	ErrNoColumnInSqlInsertQuery:  "no column set for sql query",
	ErrNoColumnInSqlUpdateQuery:  "no column set for sql query",

	ErrDbCorruptedData:             "failed to interpret data from database",
	ErrDbRequestCreationFailed:     "failed to create database request",
	ErrDbRequestFailed:             "failed to query data from database",
	ErrMultiValuedDbElement:        "multiple values for expected unique database entry",
	ErrInvalidSqlQueryReceiverType: "invalid receiver of a sql query",
	ErrNoRowsReturnedForSqlQuery:   "sql query returned no rows",
	ErrSqlRowParsingFailed:         "parsing of sql row failed",

	ErrNotImplemented: "not implemented",
}

var defaultErrorMessage = "unexpected error occurred"

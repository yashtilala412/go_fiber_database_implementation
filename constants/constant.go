package constants

// Query Parameters
const (
	ParamAppID       = "appID"
	Limit            = "limit"
	Offset           = "page"
	ParamFilterPrice = "price"
)

// Error Messages
const (
	ErrorInvalidLimit  = "Invalid limit value"
	ErrorInvalidOffset = "Invalid page value"
	ErrorInvalidAppID  = "Invalid App ID"
	ErrorAppNotFound   = "App Not Found"
	ErrorLoadingCache  = "Error loading app data into cache"
)

// Defaults
const (
	DefaultLimit  = "30"
	DefaultPage   = "1"
	DefaultOffset = "1"
)
const (
	// Query Parameter Names
	ParamAppName     = "appname"
	ParamSentiment   = "sentiment"
	ParamPolarityMin = "polarity_min"
	ParamPolarityMax = "polarity_max"

	// Default Query Values
	DefaultAppName     = "10 Best Foods for You"
	DefaultSentiment   = "Positive"
	DefaultPolarityMin = "0.5"
	DefaultPolarityMax = "1.0"

	// Error Messages
	ErrorInvalidPolarityMin = "Invalid minimum polarity value"
	ErrorInvalidPolarityMax = "Invalid maximum polarity value"

	//review controller

	//review controller
	LogDeletingReviews         = "Deleting reviews for app with name"
	ErrDeletingReviews         = "Error deleting reviews"
	ErrDeleteReviews           = "Failed to delete reviews"
	ReviewsDeletedSuccessfully = "Reviews deleted successfully"

	//review model

)
const (
	//error constants

	ErrInvalidAppNameFormat = "Invalid app name format"
	ErrAppNotFound          = "App not found"
	ErrDeleteApp            = "Failed to delete app"
	AppDeletedSuccessfully  = "App deleted successfully"
	ErrDecodingAppName      = "Error decoding app name"
	ErrDeletingApp          = "Error deleting app"
	LogDeletingApp          = "Deleting app with name"
	AppNotFoundErrorMessage = "App not found"

	// ... other constants ...
)

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
	ErrorInvalidLimit     = "Invalid limit value"
	ErrorInvalidOffset    = "Invalid page value or offset value"
	ErrorInvalidAppID     = "Invalid App ID"
	ErrorAppNotFound      = "App Not Found"
	ErrorLoadingCache     = "Error loading app data into cache"
	FailedToGetApp        = "Failed To get app"
	ErrorFiledToUpdateApp = "Failed to update app data: "
)

const (
	// ParamUid is the key for user ID route parameter
	ParamUid = "userId"

	// ParamAppId is the key for app ID route parameter
	ParamAppId = "id"
	// DefaultLimit is the default page size
	DefaultLimit = 30

	// DefaultOffset is the default page offset
	DefaultOffset = 1
)
const (
	ParamReviewID = "id"
)
const (
	// Error codes for reviews
	ErrorInvalidReviewID  = "Invalid Review ID"
	ErrorReviewNotFound   = "Review not found"
	FailedToGetReview     = "Failed to get review"
	FailedToGetReviews    = "Failed to get reviews"
	FailedToUpdateReviews = "Failed to update review"
)
const (
	ErrorInvalidRequestBody     = "Invalid request body"
	ErrorFiledToCreateApp       = "Failed to create app data: "
	ErrorFiledToCreateReviewApp = "Failed to create review data: "
	ErrorlimitAccess            = "Limit exceeded: max 500 apps per request allowed for this PC"
)
const (
	ErrorFaiedToDeleteReview = "Failed to delete review"
	ErrorFaiedToDeleteApp    = "Failed to delete app"
)
const (
	ReviewsDeletedSuccessfully = "Review deleted successfully"
	AppsDeletedSuccessfully    = "Apps deleted successfully"
)

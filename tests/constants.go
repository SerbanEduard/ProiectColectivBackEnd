package tests

const (
	// Test User Data
	TestFirstName     = "John"
	TestLastName      = "Doe"
	TestUsername      = "johndoe"
	TestEmail         = "john@example.com"
	TestPassword      = "password123"
	ExistingUsername  = "existinguser"
	ExistingEmail     = "existing@example.com"

	// Test IDs
	TestUserID   = "123"
	TestTeamID   = "team123"
	TestTeamID2  = "team456"
	TestUserID1  = "user1"
	TestUserID2  = "user2"

	// Error Messages
	ErrUserNotFound        = "user not found"
	ErrUsernameExists      = "username already exists"
	ErrEmailExists         = "email already exists"
	ErrInvalidDuration     = "Invalid timeSpentOnApp format"

	// Success Messages
	MsgStatisticsUpdated = "Statistics updated successfully"

	// HTTP Methods and Paths
	HTTPMethodPOST = "POST"
	HTTPMethodPUT  = "PUT"
	PathUsersSignup     = "/users/signup"
	PathUserStatistics  = "/users/123/statistics"

	// Content Types
	ContentTypeJSON = "application/json"

	// Test Duration Strings
	TestDurationApp   = "2h30m"
	TestDurationTeam  = "1h15m"
	TestDurationInvalid = "invalid"

	// Gin Param Keys
	ParamKeyID = "id"

	// JSON Keys
	JSONKeyError   = "error"
	JSONKeyMessage = "message"
)
package variables

// Server Errors
const (
	JsonPackFailedError     = "Failed to marshal JSON object"
	ResponseSendFailedError = "Failed to send response to client"
	ListenAndServeError     = "Failed to listen and serve"
)

// Authorization Errors
const (
	UserNotAuthorized = "User not authorized"
)

// API Messages
const (
	StatusMethodNotAllowedError = "Method not allowed"
	StatusBadRequestError       = "Bad request"
	StatusInternalServerError   = "Internal server error"
	StatusUnauthorizedError     = "Unauthorized"
	SessionCreateError          = "Session create failed"
	StatusOkMessage             = "Succesful response"
	SessionKilledError          = "Session killed failed"
	SessionNotFoundError        = "Session not found"
	UserAlreadyExistsError      = "User already exists"
)

// Middleware types
type (
	contextKey string
	sessionKey string
)

// Middleware keys constants
const (
	UserIDKey    contextKey = "userId"
	SessionIDKey sessionKey = "sessionId"
)

// Configs types
type (
	AuthorizationAppConfig struct {
		Port string `yaml:"port"`
	}

	CacheDataBaseConfig struct {
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		DbNumber int    `yaml:"db"`
		Timer    int    `yaml:"timer"`
	}

	RelationalDataBaseConfig struct {
		User         string `yaml:"user"`
		DbName       string `yaml:"dbname"`
		Password     string `yaml:"password"`
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Sslmode      string `yaml:"sslmode"`
		MaxOpenConns int    `yaml:"max_open_conns"`
		Timer        uint32 `yaml:"timer"`
	}
)

// Cookies data
const (
	SessionCookieName = "session_id"
	HttpOnly          = true
)

// Repository messages
const (
	AuthorizationCachePingRetryError      = "Authorization cache: ping failed"
	AuthorizationCachePingMaxRetriesError = "Authorization cache: ping error. Maximum number of retries reached"
	SessionRemoveError                    = "Delete session request could not be completed:"
)

// Repository constants
const (
	MaxRetries = 5
)

// Core Messages
const (
	InvalidEmailOrPasswordError     = "Invalid email or password"
	SessionRepositoryNotActiveError = "Session repository not active"
	ProfileRepositoryNotActiveError = "Profile repository not active"
	CreateProfileError              = "Create profile failed"
	ProfileNotFoundError            = "Profile not found"
	GetProfileError                 = "Get profile failed"
	GetProfileRoleError             = "Get profile role failed"
)

// Core variables
var (
	LetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// Logger constants
const (
	ModuleLogger     = "Module"
	CoreModuleLogger = "CoreModuleLogger"
)

// Regexp
const (
	LoginRegexp = `^[a-zA-Z0-9]+$`
)

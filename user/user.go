package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	gotruevault "github.com/FirstVisit/go-truevault"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// User contains the base access fields required for a TrueVault user
type User struct {
	AccessToken string `json:"access_token"`
	AccountID   string `json:"account_id"`
	APIKey      string `json:"api_key"`
	ID          string `json:"id"`
	Status      string `json:"status"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	MFAEnrolled bool   `json:"mfa_enrolled"`
}

// Status indicates the state of the user
type Status struct {
	status string
}

var (
	// Activated the user is active in TV
	Activated = Status{status: "ACTIVATED"}
	// Pending the user is pending in TV
	Pending = Status{status: "PENDING"}
	// Locked the user is locked in TV
	Locked = Status{status: "LOCKED"}
	// Deactivated the user is deactivated in TV
	Deactivated = Status{status: "DEACTIVATED"}
)

func (u *Status) String() string {
	return u.status
}

type crudResponse struct {
	Result        string            `json:"result"`
	TransactionID string            `json:"transaction_id"`
	User          User              `json:"user"`
	Error         gotruevault.Error `json:"error"`
}

type getUserResponse struct {
	Result        string `json:"result"`
	TransactionID string `json:"transaction_id"`
	Users         []User `json:"users"`
}

type createAPIKeyResponse struct {
	ApiKey        string `json:"api_key"`
	Result        string `json:"result"`
	TransactionID string `json:"transaction_id"`
}

//go:generate mockery --name Client
type Client interface {
	Get(ctx context.Context, userId []string, full bool) ([]User, error)
	Create(ctx context.Context, username, password, attributes string, groupIds []string, status Status, accessTokenNotValueAfter time.Time) (User, error)
	List(ctx context.Context, status Status, full bool) ([]User, error)
	Update(ctx context.Context, userId, username, password, accessToken string, accessTokenNotValueAfter time.Time, attributes string, status Status) (User, error)
	UpdatePassword(ctx context.Context, userId, password string) error
	Delete(ctx context.Context, userID string) error
	CreateAccessToken(ctx context.Context, userId string, notValidAfter time.Time) error
	CreateAPIKey(ctx context.Context, userID string) error
}

// Service implements the Client interface
type Service struct {
	*gotruevault.Client
}

// New creates a new Service service
func New(client gotruevault.Client) Service {
	return Service{&client}
}

// Get returns information about one or more users. If any IDs aren't valid UUIDs, returns a 400. If any can’t be
// found or the user doesn't have permission to read them, returns a 404. Otherwise, returns 200.
//
// Note: When full=true, this endpoint consumes an Operation for every user returned, so a request with 50 ids will
//       count as 50 Operations. When full=false, it consumes 1 operation regardless of how many users are returned
//
// userIds - string(req’d) - comma separated list of user IDs to retrieve. At most 100 ids can be fetched at a time.
// full – boolean(optional, default: ‘false’) - return Service attributes and Group IDs. Note: If true, then this
//				 endpoint consumes an Operation for every user returned. If false, only a single Operation is used.
func (u *Service) Get(ctx context.Context, userId []string, full bool) ([]User, error) {
	if userId == nil || len(userId) == 0 {
		return nil, errors.New("user id required")
	}

	q := make(url.Values)
	q.Set("full", strconv.FormatBool(full))

	req, err := u.NewRequest(ctx, http.MethodGet, u.URLBuilder.GetUserURL(userId), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, err
	}

	var msg getUserResponse
	return msg.Users, u.Do(req, &msg)
}

// Create creates a new TrueVault Service. The username given must be unique to ACTIVATED and LOCKED Users for an
// Account. Upon creation, both an API_KEY and an ACCESS_TOKEN will be automatically vended to the user. For security
// reasons, the API_KEY will only be shown upon creation or via the TrueVault Management Console for the account’s
// administrators. If group_ids is provided, the newly created user will be added to all given groups. The user making
// the request must have the C Group::GROUPID::GroupMembership::.* or U Group::GROUPID permission for all given groups.
// Please see authorization for more information regarding recommendations for API_KEY and ACCESS_TOKEN usage.
//
// username – string(req’d) - username for the Service being created
// password – string(optional) - password for the Service being created. If created without a password, the user
//					 can’t authenticate using the login endpoint, but it can still have an API key. This allows creating
//		 			 service accounts for backups or other server-to-TrueVault communication.
// attributes – b64 string(optional) - base64 encoded JSON document describing the Service attributes
// groupIds   – (optional) - list of group IDs where the new user will be placed
// status     – (optional) - the user’s status, one of ACTIVATED (default), PENDING, or LOCKED
// accessTokenNotValueAfter – (optional) - expiration time of generated access token
func (u *Service) Create(ctx context.Context, username, password, attributes string, groupIds []string, status Status, accessTokenNotValueAfter time.Time) (User, error) {
	if username == "" {
		return User{}, errors.New("username required to create user")
	}

	data := url.Values{}
	data.Set("username", username)

	if password != "" {
		data.Set("password", password)
	}

	if attributes != "" {
		data.Set("attributes", attributes)
	}

	if groupIds != nil {
		data.Set("group_ids", strings.Join(groupIds, ","))
	}

	if status.String() != "" {
		data.Set("status", status.String())
	}

	if !accessTokenNotValueAfter.IsZero() {
		data.Set("access_token_not_value_after", accessTokenNotValueAfter.String())
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return User{}, err
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.CreateUserURL(), gotruevault.ContentTypeApplicationJSON, buf)
	if err != nil {
		return User{}, err
	}

	var msg crudResponse
	return msg.User, u.Do(req, &msg)
}

// List returns all Users belonging to an Account.
// status – string(optional, default: ‘ACTIVATED’) - comma separated list of statuses (inclusive). Accepts any
//			combination of ACTIVATED, DEACTIVATED, or LOCKED.
// full – boolean(optional, default: ‘false’) - return Service attributes and Group IDs. Note: If true, then this endpoint
//		  consumes an Operation for every user returned. If false, only a single Operation is used.
func (u *Service) List(ctx context.Context, status Status, full bool) ([]User, error) {
	q := make(url.Values)
	q.Set("status", status.String())
	q.Set("full", strconv.FormatBool(full))

	req, err := u.NewRequest(ctx, http.MethodGet, u.URLBuilder.ListUserURL(q), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, err
	}

	var msg getUserResponse
	return msg.Users, u.Do(req, &msg)
}

// Update a given Service’s properties. Strictly overwrites existing values.
// userId – string(required)
// full – boolean(optional, default: ‘false’) - return Service attributes and Group IDs. Note: If true, then this
//				 endpoint consumes an Operation for every user returned. If false, only a single Operation is used.
func (u *Service) Update(ctx context.Context, userId, username, password, accessToken string, accessTokenNotValueAfter time.Time, attributes string, status Status) (User, error) {
	if userId == "" {
		return User{}, errors.New("user id required to update user")
	}

	data := url.Values{}
	if username != "" {
		data.Set("username", username)
	}

	if password != "" {
		data.Set("password", password)
	}

	if accessToken != "" {
		data.Set("access_token", accessToken)
	}

	if !accessTokenNotValueAfter.IsZero() {
		data.Set("access_token_not_value_after", accessTokenNotValueAfter.String())
	}

	if attributes != "" {
		data.Set("attributes", attributes)
	}

	if status.String() != "" {
		data.Set("status", status.String())
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return User{}, err
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.UpdateUserURL(userId), gotruevault.ContentTypeApplicationJSON, buf)
	if err != nil {
		return User{}, err
	}

	var msg crudResponse
	return msg.User, u.Do(req, &msg)
}

// UpdatePassword Updates a given Service’s password. Requires the `U` activity on the `Service::USERID::Password` or
// `Service::USERID resource`.
//
// userId – string(required)
// returns - nil on success otherwise, ErrorNotFound when user does not exist
func (u *Service) UpdatePassword(ctx context.Context, userId, password string) error {
	if userId == "" {
		return errors.New("user id required")
	}

	if password == "" {
		return errors.New("password is required")
	}

	data := url.Values{}
	data.Set("password", password)

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.UpdateUserPasswordURL(userId), gotruevault.ContentTypeApplicationJSON, buf)
	if err != nil {
		return err
	}

	var msg crudResponse
	if err := u.Do(req, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

// Delete deactivates a user: frees the associated username, all ACCESS_TOKENs, and removes user_id from all Groups.
// warning - This endpoint does not delete any data permanently, unlike the Document and BLOB delete endpoints. If you
//           need to completely purge a user’s data for policy or compliance reasons, first update the user’s attributes
//           to be {}, then update their username to be a unique random string, then call this endpoint.
// warning - Once the user has been deactivated, it cannot be reactivated via a status update.
func (u *Service) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user id required")
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.DeleteUserURL(userID), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return err
	}

	var msg crudResponse
	if err := u.Do(req, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

// CreateAccessToken Vends a new `ACCESS_TOKEN` for user_id.
func (u *Service) CreateAccessToken(ctx context.Context, userId string, notValidAfter time.Time) error {
	if userId == "" {
		return errors.New("user id required")
	}

	data := url.Values{}
	data.Set("not_valid_after", notValidAfter.String())

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return err
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.CreateAccessTokenURL(userId), gotruevault.ContentTypeApplicationJSON, buf)
	if err != nil {
		return err
	}

	var msg crudResponse
	if err := u.Do(req, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

// CreateAPIKey replaces the current `API_KEY` for user_id. Companion to `ACCESS_TOKEN` method. Must have `U` group
// permissions for the user.
func (u *Service) CreateApiKey(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user id required")
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.CreateApiKeyURL(userID), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return "", err
	}

	var msg createAPIKeyResponse
	return msg.ApiKey, u.Do(req, &msg)
}

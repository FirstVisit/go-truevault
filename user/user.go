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

// TvUser contains the base access fields required for a TrueVault user
type TvUser struct {
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
type Status string

const (
	Activated   Status = "ACTIVATED"
	Pending     Status = "PENDING"
	Locked      Status = "LOCKED"
	Deactivated Status = "DEACTIVATED"
)

// CRUDResponse contains the response from creating a new TrueVault User
type CRUDResponse struct {
	Result        string            `json:"result"`
	TransactionID string            `json:"transaction_id"`
	User          TvUser            `json:"user"`
	Error         gotruevault.Error `json:"error"`
}

// GetUserResponse contains the response from creating a new TrueVault User
type GetUserResponse struct {
	Result        string   `json:"result"`
	TransactionID string   `json:"transaction_id"`
	Users         []TvUser `json:"users"`
}

type CreateAPIKeyResponse struct {
	ApiKey        string `json:"api_key"`
	Result        string `json:"result"`
	TransactionID string `json:"transaction_id"`
}

//go:generate mockery --name Client
type Client interface {
	Get(ctx context.Context, userId []string, full bool) ([]TvUser, error)
	Create(ctx context.Context, username, password, attributes string, groupIds []string, status Status, accessTokenNotValueAfter time.Time) (TvUser, error)
	List(ctx context.Context, status Status, full bool) ([]TvUser, error)
	Update(ctx context.Context, userId, username, password, accessToken string, accessTokenNotValueAfter time.Time, attributes string, status Status) (TvUser, error)
	UpdatePassword(ctx context.Context, userId, password string) error
	Delete(ctx context.Context, userID string) error
	CreateAccessToken(ctx context.Context, userId string, notValidAfter time.Time) error
	CreateAPIKey(ctx context.Context, userID string) error
}

// User implements the Client interface
type User struct {
	*gotruevault.Client
}

// New creates a new User service
func New(client gotruevault.Client) User {
	return User{&client}
}

// Get returns information about one or more users. If any IDs aren't valid UUIDs, returns a 400. If any can’t be
// found or the user doesn't have permission to read them, returns a 404. Otherwise, returns 200.
//
// Note: When full=true, this endpoint consumes an Operation for every user returned, so a request with 50 ids will
//       count as 50 Operations. When full=false, it consumes 1 operation regardless of how many users are returned
//
// @param userIds - string(req’d) - comma separated list of user IDs to retrieve. At most 100 ids can be fetched at a time.
// @param full – boolean(optional, default: ‘false’) - return User attributes and Group IDs. Note: If true, then this
//				 endpoint consumes an Operation for every user returned. If false, only a single Operation is used.
func (u *User) Get(ctx context.Context, userId []string, full bool) ([]TvUser, error) {
	if userId == nil || len(userId) == 0 {
		return nil, errors.New("user id required")
	}

	q := make(url.Values)
	q.Set("full", strconv.FormatBool(full))

	req, err := u.NewRequest(ctx, http.MethodGet, u.URLBuilder.GetUserURL(userId), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, err
	}

	var msg GetUserResponse
	return msg.Users, u.Do(req, &msg)
}

// Create creates a new TrueVault User. The username given must be unique to ACTIVATED and LOCKED Users for an
// Account. Upon creation, both an API_KEY and an ACCESS_TOKEN will be automatically vended to the user. For security
// reasons, the API_KEY will only be shown upon creation or via the TrueVault Management Console for the account’s
// administrators. If group_ids is provided, the newly created user will be added to all given groups. The user making
// the request must have the C Group::GROUPID::GroupMembership::.* or U Group::GROUPID permission for all given groups.
// Please see authorization for more information regarding recommendations for API_KEY and ACCESS_TOKEN usage.
//
// @param username – string(req’d) - username for the User being created
// @param password – string(optional) - password for the User being created. If created without a password, the user
//					 can’t authenticate using the login endpoint, but it can still have an API key. This allows creating
//		 			 service accounts for backups or other server-to-TrueVault communication.
// @param attributes – b64 string(optional) - base64 encoded JSON document describing the User attributes
// @param groupIds   – (optional) - list of group IDs where the new user will be placed
// @param status     – (optional) - the user’s status, one of ACTIVATED (default), PENDING, or LOCKED
// @param accessTokenNotValueAfter – (optional) - expiration time of generated access token
func (u *User) Create(ctx context.Context, username, password, attributes string, groupIds []string, status Status, accessTokenNotValueAfter time.Time) (TvUser, error) {
	if username == "" {
		return TvUser{}, errors.New("username required to create user")
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

	// TODO: Force empty string to be invalid. Should default (zero) to ACTIVATED
	if status != "" {
		data.Set("status", string(status))
	}

	if !accessTokenNotValueAfter.IsZero() {
		data.Set("access_token_not_value_after", accessTokenNotValueAfter.String())
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return TvUser{}, err
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.CreateUserURL(), gotruevault.ContentTypeApplicationJSON, buf)
	if err != nil {
		return TvUser{}, err
	}

	var msg CRUDResponse
	return msg.User, u.Do(req, &msg)
}

// List returns all Users belonging to an Account.
//
// status – string(optional, default: ‘ACTIVATED’) - comma separated list of statuses (inclusive). Accepts any
//			combination of ACTIVATED, DEACTIVATED, or LOCKED.
// full – boolean(optional, default: ‘false’) - return User attributes and Group IDs. Note: If true, then this endpoint
//		  consumes an Operation for every user returned. If false, only a single Operation is used.
func (u *User) List(ctx context.Context, status Status, full bool) ([]TvUser, error) {
	q := make(url.Values)
	q.Set("status", string(status))
	q.Set("full", strconv.FormatBool(full))

	req, err := u.NewRequest(ctx, http.MethodGet, u.URLBuilder.ListUserURL(q), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return nil, err
	}

	var msg GetUserResponse
	return msg.Users, u.Do(req, &msg)
}

// Update a given User’s properties. Strictly overwrites existing values.
//
// @Param userId – string(required)
// @Param full – boolean(optional, default: ‘false’) - return User attributes and Group IDs. Note: If true, then this
//				 endpoint consumes an Operation for every user returned. If false, only a single Operation is used.
func (u *User) Update(ctx context.Context, userId, username, password, accessToken string, accessTokenNotValueAfter time.Time, attributes string, status Status) (TvUser, error) {
	if userId == "" {
		return TvUser{}, errors.New("user id required to update user")
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

	// TODO: Force empty string to be invalid. Should default (zero) to ACTIVATED
	if status != "" {
		data.Set("status", string(status))
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return TvUser{}, err
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.UpdateUserURL(userId), gotruevault.ContentTypeApplicationJSON, buf)
	if err != nil {
		return TvUser{}, err
	}

	var msg CRUDResponse
	return msg.User, u.Do(req, &msg)
}

// UpdatePassword Updates a given User’s password. Requires the `U` activity on the `User::USERID::Password` or
// `User::USERID resource`.
//
// @Param userId – string(required)
// @Returns - nil on success
//			- ErrorNotFound when user does not exist
func (u *User) UpdatePassword(ctx context.Context, userId, password string) error {
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

	var msg CRUDResponse
	if err := u.Do(req, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

// Delete deactivates a user: frees the associated username, all ACCESS_TOKENs, and removes user_id from all Groups.
// @Warning: This endpoint does not delete any data permanently, unlike the Document and BLOB delete endpoints. If you
//           need to completely purge a user’s data for policy or compliance reasons, first update the user’s attributes
//           to be {}, then update their username to be a unique random string, then call this endpoint.
// @Warning: Once the user has been deactivated, it cannot be reactivated via a status update.
func (u *User) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("user id required")
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.DeleteUserURL(userID), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return err
	}

	var msg CRUDResponse
	if err := u.Do(req, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

// CreateAccessToken Vends a new `ACCESS_TOKEN` for user_id.
func (u *User) CreateAccessToken(ctx context.Context, userId string, notValidAfter time.Time) error {
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

	var msg CRUDResponse
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
func (u *User) CreateApiKey(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user id required")
	}

	req, err := u.NewRequest(ctx, http.MethodPost, u.URLBuilder.CreateApiKeyURL(userID), gotruevault.ContentTypeApplicationJSON, nil)
	if err != nil {
		return "", err
	}

	var msg CreateAPIKeyResponse
	return msg.ApiKey, u.Do(req, &msg)
}

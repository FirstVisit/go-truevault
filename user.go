package go_truevault

import (
	"encoding/json"
	"errors"
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

// UserStatus indicates the state of the user
type UserStatus string

const (
	Activated   UserStatus = "ACTIVATED"
	Pending                = "PENDING"
	Locked                 = "LOCKED"
	Deactivated            = "DEACTIVATED"
)

// CreateOrUpdateUserResponse contains the response from creating a new TrueVault User
type CreateOrUpdateUserResponse struct {
	Result        string `json:"result"`
	TransactionID string `json:"transaction_id"`
	User          User   `json:"user"`
	Error         Error  `json:"error"`
}

// GetUserResponse contains the response from creating a new TrueVault User
type GetUserResponse struct {
	Result        string `json:"result"`
	TransactionID string `json:"transaction_id"`
	Users         []User `json:"users"`
}

type CreateAPIKeyResponse struct {
	ApiKey        string `json:"api_key"`
	Result        string `json:"result"`
	TransactionID string `json:"transaction_id"`
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
func (c *Client) Create(username string, password string, attributes string, groupIds []string, status UserStatus, accessTokenNotValueAfter time.Time) (User, error) {
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

	// TODO: Force empty string to be invalid. Should default (zero) to ACTIVATED
	if status != "" {
		data.Set("status", string(status))
	}

	if !accessTokenNotValueAfter.IsZero() {
		data.Set("access_token_not_value_after", accessTokenNotValueAfter.String())
	}

	resp, err := c.post("https://api.truevault.com/v1/users", data)
	if err != nil {
		return User{}, err
	}

	var msg CreateOrUpdateUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
		return User{}, errors.New("failed to unmarshal json data to create user")
	}

	if msg.Error.Message != "" {
		return User{}, errors.New(msg.Error.Message)
	}

	return msg.User, err
}

// Get returns information about one or more users. If any IDs aren’t valid UUIDs, returns a 400. If any can’t be
// found or the user doesn't have permission to read them, returns a 404. Otherwise, returns 200.
//
// Note: When full=true, this endpoint consumes an Operation for every user returned, so a request with 50 ids will
//       count as 50 Operations. When full=false, it consumes 1 operation regardless of how many users are returned
//
// @param userIds - string(req’d) - comma separated list of user IDs to retrieve. At most 100 ids can be fetched at a time.
// @param full – boolean(optional, default: ‘false’) - return User attributes and Group IDs. Note: If true, then this
//				 endpoint consumes an Operation for every user returned. If false, only a single Operation is used.
func (c *Client) Get(userId []string, full bool) ([]User, error) {
	if userId == nil || len(userId) == 0 {
		return nil, errors.New("user id required")
	}

	q := make(url.Values)
	q.Set("full", strconv.FormatBool(full))

	resp, err := c.get("https://api.truevault.com/v2/users/"+strings.Join(userId, ","), q)
	if err != nil {
		return nil, err
	}

	var msg GetUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
		return nil, err
	}

	// TODO: Does TV actually not return an error here?

	return msg.Users, nil
}

// List returns all Users belonging to an Account.
//
// status – string(optional, default: ‘ACTIVATED’) - comma separated list of statuses (inclusive). Accepts any
//			combination of ACTIVATED, DEACTIVATED, or LOCKED.
// full – boolean(optional, default: ‘false’) - return User attributes and Group IDs. Note: If true, then this endpoint
//		  consumes an Operation for every user returned. If false, only a single Operation is used.
func (c *Client) List(status UserStatus, full bool) ([]User, error) {
	q := make(url.Values)
	q.Set("status", string(status))
	q.Set("full", strconv.FormatBool(full))

	resp, err := c.get("https://api.truevault.com/v2/users/", q)
	if err != nil {
		return nil, err
	}

	var msg GetUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
		return nil, err
	}

	// TODO: Does TV actually not return an error here?

	return msg.Users, nil
}

// Update a given User’s properties. Strictly overwrites existing values.
//
// @Param userId – string(required)
// @Param full – boolean(optional, default: ‘false’) - return User attributes and Group IDs. Note: If true, then this
//				 endpoint consumes an Operation for every user returned. If false, only a single Operation is used.
func (c *Client) Update(userId, username, password, accessToken string, accessTokenNotValueAfter time.Time, attributes string, status UserStatus) (User, error) {
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

	// TODO: Force empty string to be invalid. Should default (zero) to ACTIVATED
	if status != "" {
		data.Set("status", string(status))
	}

	resp, err := c.post("https://api.truevault.com/v1/users/"+userId, data)
	if err != nil {
		return User{}, err
	}

	var msg CreateOrUpdateUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
		return User{}, errors.New("failed to unmarshal json data to create user")
	}

	if msg.Error.Message != "" {
		return User{}, errors.New(msg.Error.Message)
	}

	return msg.User, err
}

// UpdatePassword Updates a given User’s password. Requires the `U` activity on the `User::USERID::Password` or
// `User::USERID resource`.
//
// @Param userId – string(required)
// @Returns - nil on success
//			- ErrorNotFound when user does not exist
func (c *Client) UpdatePassword(userId, password string) error {
	if userId == "" {
		return errors.New("user id required")
	}

	if password == "" {
		return errors.New("password is required")
	}

	data := url.Values{}
	data.Set("password", password)

	resp, err := c.post("https://api.truevault.com/v1/users/"+userId, data)
	if err != nil {
		return err
	}

	var msg CreateOrUpdateUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
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
func (c *Client) Delete(userID string) error {
	if userID == "" {
		return errors.New("user id required")
	}

	resp, err := c.post("https://api.truevault.com/v1/users/"+userID, nil)
	if err != nil {
		return err
	}

	var msg CreateOrUpdateUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

// CreateAccessToken Vends a new `ACCESS_TOKEN` for user_id.
func (c *Client) CreateAccessToken(userId string, notValidAfter time.Time) error {
	if userId == "" {
		return errors.New("user id required")
	}

	data := url.Values{}
	data.Set("not_valid_after", notValidAfter.String())

	resp, err := c.post("https://api.truevault.com/v1/users", data)
	if err != nil {
		return err
	}

	var msg CreateOrUpdateUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

// CreateAPIKey Rreplaces the current `API_KEY` for user_id. Companion to `ACCESS_TOKEN` method. Must have `U` group
// permissions for the user.
func (c *Client) CreateAPIKey(userID string) error {
	if userID == "" {
		return errors.New("user id required")
	}

	URL := "https://api.truevault.com/v1/users/" + userID + "/api_key"
	resp, err := c.post(URL, nil)
	if err != nil {
		return err
	}

	var msg CreateOrUpdateUserResponse
	if err := json.Unmarshal(resp, &msg); err != nil {
		return err
	}

	if msg.Error.Message != "" {
		return errors.New(msg.Error.Message)
	}

	return nil
}

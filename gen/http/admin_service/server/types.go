// Code generated by goa v3.11.3, DO NOT EDIT.
//
// admin_service HTTP server types
//
// Command:
// $ goa gen github.com/inst-api/poster/design

package server

import (
	"unicode/utf8"

	adminservice "github.com/inst-api/poster/gen/admin_service"
	goa "goa.design/goa/v3/pkg"
)

// AddManagerRequestBody is the type of the "admin_service" service
// "add_manager" endpoint HTTP request body.
type AddManagerRequestBody struct {
	Login    *string `form:"login,omitempty" json:"login,omitempty" xml:"login,omitempty"`
	Password *string `form:"password,omitempty" json:"password,omitempty" xml:"password,omitempty"`
}

// PushBotsOKResponseBody is the type of the "admin_service" service
// "push_bots" endpoint HTTP response body.
type PushBotsOKResponseBody struct {
	// количество ботов, которых мы отправили
	SentBots int `json:"sent_bots"`
	// количество ботов, которых сохранили в проксе
	SavedBots int32 `json:"saved_bots"`
	// имена ботов, которые мы сохранили
	Usernames []string `form:"usernames" json:"usernames" xml:"usernames"`
}

// NewPushBotsOKResponseBody builds the HTTP response body from the result of
// the "push_bots" endpoint of the "admin_service" service.
func NewPushBotsOKResponseBody(res *adminservice.PushBotsResult) *PushBotsOKResponseBody {
	body := &PushBotsOKResponseBody{
		SentBots:  res.SentBots,
		SavedBots: res.SavedBots,
	}
	if res.Usernames != nil {
		body.Usernames = make([]string, len(res.Usernames))
		for i, val := range res.Usernames {
			body.Usernames[i] = val
		}
	}
	return body
}

// NewAddManagerPayload builds a admin_service service add_manager endpoint
// payload.
func NewAddManagerPayload(body *AddManagerRequestBody, token *string) *adminservice.AddManagerPayload {
	v := &adminservice.AddManagerPayload{
		Login:    *body.Login,
		Password: *body.Password,
	}
	v.Token = token

	return v
}

// NewPushBotsPayload builds a admin_service service push_bots endpoint payload.
func NewPushBotsPayload(token string) *adminservice.PushBotsPayload {
	v := &adminservice.PushBotsPayload{}
	v.Token = token

	return v
}

// ValidateAddManagerRequestBody runs the validations defined on
// add_manager_request_body
func ValidateAddManagerRequestBody(body *AddManagerRequestBody) (err error) {
	if body.Login == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("login", "body"))
	}
	if body.Password == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("password", "body"))
	}
	if body.Password != nil {
		if utf8.RuneCountInString(*body.Password) < 4 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.password", *body.Password, utf8.RuneCountInString(*body.Password), 4, true))
		}
	}
	return
}

package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/stretchr/testify/assert"
)

func TestValidateMixinToken(t *testing.T) {
	const token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkxOTcxNzEsImlhdCI6MTczOTE5NzExMSwianRpIjoiMjZlYWY2NDAtMDg3ZS00OTY3LThiYmYtYjlhMjUzM2NjNzA0Iiwic2NwIjoiRlVMTCIsInNpZCI6ImRiMmYzMmJiLWYyYTUtNDJiMS1iOTQ2LTYzYTRlMTI5YjAyYyIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiI1YzRmMzBhNi0xZjQ5LTQzYzMtYjM3Yi1jMDFhYWU1MTkxYWYifQ.JnOPghKeIP7IgSqeZ_d9mPTjFXJhZxm3HXPg_3ixVEtmRzxz_46AylOm01UiwhsK1gq7RvXDn5BbBZ2cr5EDbT96vjXASz6Im8GHNjj0i9xBnTZdGvrPVjUuJD58rGxTgiRUXHfbS4bWRm6K2VeeDq-92uJOgn6vfreYdj80NfQ"
	err := validateMixinToken(token, nil)
	assert.True(t, IsError(err, ErrCodeBadMixinToken))
}

func TestMixinClientUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Log(r.Header.Get("User-Agent"))
		assert.Equal(t, getUserAgent(), r.Header.Get("User-Agent"))
		w.Header().Set("X-Request-Id", r.Header.Get("X-Request-Id"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":{"identity_number":"1234567890","user_id":"5c4f30a6-1f49-43c3-b37b-c01aae5191af"}}`))
	}))
	defer server.Close()

	client := mixin.GetRestyClient()
	client.SetBaseURL(server.URL)

	const token = "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjIzNjk5MTc1NDMsImlhdCI6MTczOTE5NzU0MywianRpIjoiYzlkMzIyNzktMWQ4Zi00ZmQ2LTk4ZGQtNzNmZjlmNDYzYTlmIiwic2NwIjoiRlVMTCIsInNpZCI6ImRiMmYzMmJiLWYyYTUtNDJiMS1iOTQ2LTYzYTRlMTI5YjAyYyIsInNpZyI6IjVlNmI1OGZmYTEwYjNiYzUxNzI0ZmYwYmJkMmFmYjkxYzQ3NzFlZTM0MGY1ZDY4NTM0MGRmYTRjODU0YmFmYmEiLCJ1aWQiOiI1YzRmMzBhNi0xZjQ5LTQzYzMtYjM3Yi1jMDFhYWU1MTkxYWYifQ.YP3cPeBQK0fFHLClsX4ro61qM18_1svfEM1cZrWwBMQaKdIIm7GnY6QadK7UTwl0ZIPJvmqHV3BXF0LOivWTRZ4onXTezpRYUX6JveYXU1H6EP_YAvn06j8SxmTGeiYUCjA0_PqVMovnNjmu3O6Wi6r0osgeFBlKchgOxdkWNG8"

	auth := Authorizer{}
	user, err := auth.AuthorizeMixinToken(context.Background(), token)
	assert.NoError(t, err)
	assert.Equal(t, "1234567890", user.IdentityNumber)
	assert.Equal(t, "5c4f30a6-1f49-43c3-b37b-c01aae5191af", user.UserID)
}

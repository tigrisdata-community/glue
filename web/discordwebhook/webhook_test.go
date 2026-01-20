package discordwebhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSend(t *testing.T) {
	tests := []struct {
		name     string
		whurl    string
		webhook  Webhook
		wantBody string
	}{
		{
			name:  "simple webhook with content",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content: "Hello, Discord!",
			},
			wantBody: `{"content":"Hello, Discord!","avatar_url":"","allowed_mentions":null}`,
		},
		{
			name:  "webhook with username",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content:  "Test message",
				Username: "TestBot",
			},
			wantBody: `{"content":"Test message","username":"TestBot","avatar_url":"","allowed_mentions":null}`,
		},
		{
			name:  "username longer than 32 characters is truncated",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content:  "Test",
				Username: "ThisUsernameIsWayTooLongAndExceeds32",
			},
			wantBody: `{"content":"Test","username":"ThisUsernameIsWayTooLongAndExcee","avatar_url":"","allowed_mentions":null}`,
		},
		{
			name:  "webhook with avatar URL",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content:   "Test",
				AvatarURL: "https://example.com/avatar.png",
			},
			wantBody: `{"content":"Test","avatar_url":"https://example.com/avatar.png","allowed_mentions":null}`,
		},
		{
			name:  "webhook with embeds",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content: "Test",
				Embeds: []Embeds{
					{
						Title:       "Test Title",
						Description: "Test Description",
						URL:         "https://example.com",
						Fields: []EmbedField{
							{Name: "Field1", Value: "Value1", Inline: true},
							{Name: "Field2", Value: "Value2", Inline: false},
						},
					},
				},
			},
			wantBody: `{"content":"Test","avatar_url":"","embeds":[{"title":"Test Title","description":"Test Description","url":"https://example.com","fields":[{"name":"Field1","value":"Value1","inline":true},{"name":"Field2","value":"Value2","inline":false}]}],"allowed_mentions":null}`,
		},
		{
			name:  "webhook with embed author",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content: "Test",
				Embeds: []Embeds{
					{
						Title:  "Test",
						Author: &EmbedAuthor{Name: "Author", URL: "https://example.com", IconURL: "https://example.com/icon.png"},
					},
				},
			},
			wantBody: `{"content":"Test","avatar_url":"","embeds":[{"title":"Test","author":{"name":"Author","url":"https://example.com","icon_url":"https://example.com/icon.png"}}],"allowed_mentions":null}`,
		},
		{
			name:  "webhook with embed footer",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content: "Test",
				Embeds: []Embeds{
					{
						Title:  "Test",
						Footer: &EmbedFooter{Text: "Footer text", IconURL: "https://example.com/footer.png"},
					},
				},
			},
			wantBody: `{"content":"Test","avatar_url":"","embeds":[{"title":"Test","footer":{"text":"Footer text","icon_url":"https://example.com/footer.png"}}],"allowed_mentions":null}`,
		},
		{
			name:  "webhook with allowed mentions",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content:         "Test",
				AllowedMentions: map[string][]string{"parse": {"users", "roles"}},
			},
			wantBody: `{"content":"Test","avatar_url":"","allowed_mentions":{"parse":["users","roles"]}}`,
		},
		{
			name:  "empty webhook",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				AllowedMentions: map[string][]string{},
			},
			wantBody: `{"avatar_url":"","allowed_mentions":{}}`,
		},
		{
			name:  "complete webhook with all fields",
			whurl: "https://discord.com/api/webhooks/test",
			webhook: Webhook{
				Content:         "Full test",
				Username:        "FullBot",
				AvatarURL:       "https://example.com/avatar.png",
				AllowedMentions: map[string][]string{"parse": {"everyone"}},
				Embeds: []Embeds{
					{
						Title:       "Full Embed",
						Description: "Full Description",
						URL:         "https://example.com/full",
						Author:      &EmbedAuthor{Name: "Full Author", URL: "https://example.com", IconURL: "https://example.com/icon.png"},
						Fields: []EmbedField{
							{Name: "FullField", Value: "FullValue", Inline: true},
						},
						Footer: &EmbedFooter{Text: "Full Footer", IconURL: "https://example.com/footer.png"},
					},
				},
			},
			wantBody: `{"content":"Full test","username":"FullBot","avatar_url":"https://example.com/avatar.png","embeds":[{"title":"Full Embed","description":"Full Description","url":"https://example.com/full","author":{"name":"Full Author","url":"https://example.com","icon_url":"https://example.com/icon.png"},"fields":[{"name":"FullField","value":"FullValue","inline":true}],"footer":{"text":"Full Footer","icon_url":"https://example.com/footer.png"}}],"allowed_mentions":{"parse":["everyone"]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Send(tt.whurl, tt.webhook)

			if req.Method != http.MethodPost {
				t.Errorf("Send() method = %v, want %v", req.Method, http.MethodPost)
			}

			if req.URL.String() != tt.whurl {
				t.Errorf("Send() URL = %v, want %v", req.URL.String(), tt.whurl)
			}

			if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
				t.Errorf("Send() Content-Type = %v, want application/json", contentType)
			}

			var gotBody map[string]interface{}
			if err := json.NewDecoder(req.Body).Decode(&gotBody); err != nil {
				t.Fatalf("Failed to decode request body: %v", err)
			}

			var wantBody map[string]interface{}
			if err := json.Unmarshal([]byte(tt.wantBody), &wantBody); err != nil {
				t.Fatalf("Failed to decode want body: %v", err)
			}

			if !jsonEqual(gotBody, wantBody) {
				gotJSON, _ := json.Marshal(gotBody)
				wantJSON, _ := json.Marshal(wantBody)
				t.Errorf("Send() body = %s, want %s", gotJSON, wantJSON)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful validation with 204 No Content",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "error on 200 OK",
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name:       "error on 400 Bad Request",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name:       "error on 401 Unauthorized",
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name:       "error on 404 Not Found",
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
		{
			name:       "error on 500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("Failed to create test response: %v", err)
			}
			defer resp.Body.Close()

			err = Validate(resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil {
				var webErr interface{ Error() string }
				if _, ok := err.(interface{ Error() string }); !ok {
					t.Errorf("Validate() error does not implement Error() method")
				} else {
					_ = webErr
				}
			}
		})
	}
}

func TestUsernameTruncation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantLen  int
	}{
		{
			name:     "username at exactly 32 characters is not truncated",
			username: "12345678901234567890123456789012",
			wantLen:  32,
		},
		{
			name:     "username at 33 characters is truncated to 32",
			username: "123456789012345678901234567890123",
			wantLen:  32,
		},
		{
			name:     "very long username is truncated to 32",
			username: "ThisIsAVeryLongUsernameThatExceedsTheLimit",
			wantLen:  32,
		},
		{
			name:     "empty username remains empty",
			username: "",
			wantLen:  0,
		},
		{
			name:     "short username is unchanged",
			username: "Short",
			wantLen:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Send("https://discord.com/api/webhooks/test", Webhook{
				Username: tt.username,
			})

			var webhook Webhook
			if err := json.NewDecoder(req.Body).Decode(&webhook); err != nil {
				t.Fatalf("Failed to decode webhook: %v", err)
			}

			if len(webhook.Username) != tt.wantLen {
				t.Errorf("Username length = %d, want %d", len(webhook.Username), tt.wantLen)
			}

			if tt.wantLen > 0 && len(webhook.Username) > 0 && len(webhook.Username) != len(tt.username) {
				// If truncated, verify it's a prefix
				if tt.username[:tt.wantLen] != webhook.Username {
					t.Errorf("Username = %s, want prefix %s", webhook.Username, tt.username[:tt.wantLen])
				}
			}
		})
	}
}

// jsonEqual performs a deep comparison of two JSON objects.
func jsonEqual(a, b map[string]interface{}) bool {
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

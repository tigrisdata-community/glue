package discordwebhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tigrisdata-community/glue/web"
)

type ThreadCreationResponse struct {
	Type            int       `json:"type"`
	Content         string    `json:"content"`
	Mentions        []any     `json:"mentions"`
	MentionRoles    []any     `json:"mention_roles"`
	Attachments     []any     `json:"attachments"`
	Embeds          []any     `json:"embeds"`
	Timestamp       time.Time `json:"timestamp"`
	EditedTimestamp any       `json:"edited_timestamp"`
	Flags           int       `json:"flags"`
	Components      []any     `json:"components"`
	ID              string    `json:"id"`
	ChannelID       string    `json:"channel_id"`
	Author          Author    `json:"author"`
	Pinned          bool      `json:"pinned"`
	MentionEveryone bool      `json:"mention_everyone"`
	Tts             bool      `json:"tts"`
	WebhookID       string    `json:"webhook_id"`
	Position        int       `json:"position"`
}

type Author struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Avatar        any    `json:"avatar"`
	Discriminator string `json:"discriminator"`
	PublicFlags   int    `json:"public_flags"`
	Flags         int    `json:"flags"`
	Bot           bool   `json:"bot"`
	GlobalName    any    `json:"global_name"`
	Clan          any    `json:"clan"`
	PrimaryGuild  any    `json:"primary_guild"`
}

func ParseThreadCreation(resp *http.Response) (*ThreadCreationResponse, error) {
	if resp.StatusCode != http.StatusOK {
		return nil, web.NewError(http.StatusOK, resp)
	}

	var result ThreadCreationResponse
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("can't parse ThreadCreationResponse: %w", err)
	}

	return &result, nil
}

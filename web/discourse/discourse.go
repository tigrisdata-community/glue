package discourse

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/tigrisdata-community/glue/web/useragent"
)

func GetCategoryAndTag(ctx context.Context, u string) (*CategoryAndTagResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", useragent.Generate("tigris-gtm-glue", "https://tigrisdata.com"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var result CategoryAndTagResult
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type CategoryAndTagResult struct {
	Users         []Users         `json:"users"`
	PrimaryGroups []PrimaryGroups `json:"primary_groups"`
	FlairGroups   []FlairGroups   `json:"flair_groups"`
	TopicList     TopicList       `json:"topic_list"`
}

type Users struct {
	ID               int    `json:"id"`
	Username         string `json:"username"`
	Name             string `json:"name"`
	AvatarTemplate   string `json:"avatar_template"`
	TrustLevel       int    `json:"trust_level"`
	Admin            bool   `json:"admin,omitempty"`
	Moderator        bool   `json:"moderator,omitempty"`
	PrimaryGroupName string `json:"primary_group_name,omitempty"`
	FlairName        string `json:"flair_name,omitempty"`
	FlairURL         string `json:"flair_url,omitempty"`
	FlairBgColor     string `json:"flair_bg_color,omitempty"`
	FlairColor       string `json:"flair_color,omitempty"`
	FlairGroupID     int    `json:"flair_group_id,omitempty"`
}

type PrimaryGroups struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type FlairGroups struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	FlairURL     string `json:"flair_url"`
	FlairBgColor string `json:"flair_bg_color"`
	FlairColor   string `json:"flair_color"`
}

type Tags struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TopicCount  int    `json:"topic_count"`
	Staff       bool   `json:"staff"`
	Description any    `json:"description"`
}

type Posters struct {
	Extras         string `json:"extras"`
	Description    string `json:"description"`
	UserID         int    `json:"user_id"`
	PrimaryGroupID any    `json:"primary_group_id"`
	FlairGroupID   any    `json:"flair_group_id"`
}

type Topics struct {
	FancyTitle         string    `json:"fancy_title"`
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Slug               string    `json:"slug"`
	PostsCount         int       `json:"posts_count"`
	ReplyCount         int       `json:"reply_count"`
	HighestPostNumber  int       `json:"highest_post_number"`
	ImageURL           any       `json:"image_url"`
	CreatedAt          time.Time `json:"created_at"`
	LastPostedAt       time.Time `json:"last_posted_at"`
	Bumped             bool      `json:"bumped"`
	BumpedAt           time.Time `json:"bumped_at"`
	Archetype          string    `json:"archetype"`
	Unseen             bool      `json:"unseen"`
	Pinned             bool      `json:"pinned"`
	Unpinned           any       `json:"unpinned"`
	Visible            bool      `json:"visible"`
	Closed             bool      `json:"closed"`
	Archived           bool      `json:"archived"`
	Bookmarked         any       `json:"bookmarked"`
	Liked              any       `json:"liked"`
	Tags               []string  `json:"tags"`
	TagsDescriptions   any       `json:"tags_descriptions"`
	Views              int       `json:"views"`
	LikeCount          int       `json:"like_count"`
	HasSummary         bool      `json:"has_summary"`
	LastPosterUsername string    `json:"last_poster_username"`
	CategoryID         int       `json:"category_id"`
	OpLikeCount        int       `json:"op_like_count"`
	PinnedGlobally     bool      `json:"pinned_globally"`
	FeaturedLink       any       `json:"featured_link"`
	HasAcceptedAnswer  bool      `json:"has_accepted_answer"`
	CanVote            bool      `json:"can_vote"`
	Posters            []Posters `json:"posters"`
}

// JSONURL returns the relative path to the JSON API for this topic.
// The result is like "/t/slug/123.json" which can be appended to a base URL.
func (t Topics) JSONURL() string {
	return "/t/" + t.Slug + "/" + strconv.Itoa(t.ID) + ".json"
}

type TopicList struct {
	CanCreateTopic bool     `json:"can_create_topic"`
	MoreTopicsURL  string   `json:"more_topics_url"`
	PerPage        int      `json:"per_page"`
	TopTags        []string `json:"top_tags"`
	Tags           []Tags   `json:"tags"`
	Topics         []Topics `json:"topics"`
}

// GetTopic fetches a Discourse topic by URL (e.g., https://community.fly.io/t/slug/123.json)
func GetTopic(ctx context.Context, u string) (*TopicResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", useragent.Generate("tigris-gtm-glue", "https://tigrisdata.com"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TopicResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type TopicResult struct {
	PostStream       PostStream       `json:"post_stream"`
	TimelineLookup   [][]int          `json:"timeline_lookup"`
	SuggestedTopics  []Topics         `json:"suggested_topics"`
	Tags             []string         `json:"tags"`
	TagsDescriptions any              `json:"tags_descriptions"`
	FancyTitle       string           `json:"fancy_title"`
	ID               int              `json:"id"`
	Title            string           `json:"title"`
	PostsCount       int              `json:"posts_count"`
	CreatedAt        time.Time        `json:"created_at"`
	Views            int              `json:"views"`
	ReplyCount       int              `json:"reply_count"`
	LikeCount        int              `json:"like_count"`
	LastPostedAt     time.Time        `json:"last_posted_at"`
	Visible          bool             `json:"visible"`
	Closed           bool             `json:"closed"`
	Archived         bool             `json:"archived"`
	HasSummary       bool             `json:"has_summary"`
	Archetype        string           `json:"archetype"`
	Slug             string           `json:"slug"`
	CategoryID       int              `json:"category_id"`
	WordCount        int              `json:"word_count"`
	DeletedAt        any              `json:"deleted_at"`
	UserID           int              `json:"user_id"`
	FeaturedLink     any              `json:"featured_link"`
	PinnedGlobally   bool             `json:"pinned_globally"`
	PinnedAt         any              `json:"pinned_at"`
	PinnedUntil      any              `json:"pinned_until"`
	ImageURL         any              `json:"image_url"`
	SlowModeSeconds  int              `json:"slow_mode_seconds"`
	Draft            any              `json:"draft"`
	DraftKey         string           `json:"draft_key,omitempty"`
	DraftSequence    any              `json:"draft_sequence,omitempty"`
	Unpinned         any              `json:"unpinned"`
	Pinned           bool             `json:"pinned"`
	CurrentPostNumber int             `json:"current_post_number"`
	HighestPostNumber int             `json:"highest_post_number"`
	DeletedBy        any              `json:"deleted_by"`
	ActionsSummary   []ActionsSummary `json:"actions_summary"`
	ChunkSize        int              `json:"chunk_size"`
	Bookmarked       bool             `json:"bookmarked"`
	TopicTimer       *TopicTimer      `json:"topic_timer,omitempty"`
	MessageBusLastID int              `json:"message_bus_last_id"`
	ParticipantCount int              `json:"participant_count"`
	ShowReadIndicator bool            `json:"show_read_indicator"`
	Thumbnails       any              `json:"thumbnails"`
	SlowModeEnabledUntil any           `json:"slow_mode_enabled_until"`
	AcceptedAnswer   *AcceptedAnswer  `json:"accepted_answer,omitempty"`
	CanVote          bool             `json:"can_vote"`
	VoteCount        int              `json:"vote_count"`
	UserVoted        bool             `json:"user_voted"`
	Details          TopicDetails     `json:"details"`
	Bookmarks        any              `json:"bookmarks"`
}

// JSONURL returns the relative path to the JSON API for this topic.
// The result is like "/t/slug/123.json" which can be appended to a base URL.
func (t TopicResult) JSONURL() string {
	return "/t/" + t.Slug + "/" + strconv.Itoa(t.ID) + ".json"
}

type PostStream struct {
	Posts  []Post `json:"posts"`
	Stream []int  `json:"stream"`
}

type Post struct {
	ID                  int              `json:"id"`
	Name                string           `json:"name"`
	Username            string           `json:"username"`
	AvatarTemplate      string           `json:"avatar_template"`
	CreatedAt           time.Time        `json:"created_at"`
	Cooked              string           `json:"cooked"`
	PostNumber          int              `json:"post_number"`
	PostType            int              `json:"post_type"`
	PostsCount          int              `json:"posts_count"`
	UpdatedAt           time.Time        `json:"updated_at"`
	ReplyCount          int              `json:"reply_count"`
	ReplyToPostNumber   any              `json:"reply_to_post_number"`
	QuoteCount          int              `json:"quote_count"`
	IncomingLinkCount   int              `json:"incoming_link_count"`
	Reads               int              `json:"reads"`
	ReadersCount        int              `json:"readers_count"`
	Score               float64          `json:"score"`
	Yours               bool             `json:"yours"`
	TopicID             int              `json:"topic_id"`
	TopicSlug           string           `json:"topic_slug"`
	DisplayUsername     string           `json:"display_username"`
	PrimaryGroupName    any              `json:"primary_group_name"`
	FlairName           any              `json:"flair_name"`
	FlairURL            any              `json:"flair_url"`
	FlairBgColor        any              `json:"flair_bg_color"`
	FlairColor          any              `json:"flair_color"`
	FlairGroupID        any              `json:"flair_group_id"`
	BadgesGranted       []any            `json:"badges_granted"`
	Version             int              `json:"version"`
	CanEdit             bool             `json:"can_edit"`
	CanDelete           bool             `json:"can_delete"`
	CanRecover          bool             `json:"can_recover"`
	CanSeeHiddenPost    bool             `json:"can_see_hidden_post"`
	CanWiki             bool             `json:"can_wiki"`
	LinkCounts          []LinkCount      `json:"link_counts"`
	Read                bool             `json:"read"`
	UserTitle           any              `json:"user_title"`
	Bookmarked          bool             `json:"bookmarked"`
	ActionsSummary      []ActionsSummary `json:"actions_summary"`
	Moderator           bool             `json:"moderator"`
	Admin               bool             `json:"admin"`
	Staff               bool             `json:"staff"`
	UserID              int              `json:"user_id"`
	Hidden              bool             `json:"hidden"`
	TrustLevel          int              `json:"trust_level"`
	DeletedAt           any              `json:"deleted_at"`
	UserDeleted         bool             `json:"user_deleted"`
	EditReason          any              `json:"edit_reason"`
	CanViewEditHistory  bool             `json:"can_view_edit_history"`
	Wiki                bool             `json:"wiki"`
	PostURL             string           `json:"post_url"`
	CanAcceptAnswer     bool             `json:"can_accept_answer"`
	CanUnacceptAnswer   bool             `json:"can_unaccept_answer"`
	AcceptedAnswer      bool             `json:"accepted_answer"`
	TopicAcceptedAnswer bool             `json:"topic_accepted_answer"`
	CanVote             bool             `json:"can_vote"`
	ReplyToUser         *User            `json:"reply_to_user,omitempty"`
}

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	AvatarTemplate string `json:"avatar_template"`
}

type LinkCount struct {
	URL         string `json:"url"`
	Internal    bool   `json:"internal"`
	Reflection  bool   `json:"reflection"`
	Title       string `json:"title"`
	Clicks      int    `json:"clicks"`
}

type ActionsSummary struct {
	ID     int `json:"id"`
	Count  int `json:"count"`
	Hidden any `json:"hidden,omitempty"`
}

type TopicTimer struct {
	ID               int       `json:"id"`
	ExecuteAt        time.Time `json:"execute_at"`
	DurationMinutes  int       `json:"duration_minutes"`
	BasedOnLastPost  bool      `json:"based_on_last_post"`
	StatusType       string    `json:"status_type"`
	CategoryID       any       `json:"category_id"`
}

type AcceptedAnswer struct {
	PostNumber  int    `json:"post_number"`
	Username    string `json:"username"`
	Name        any    `json:"name"`
	Excerpt     string `json:"excerpt"`
	AccepterName any    `json:"accepter_name"`
}

type TopicDetails struct {
	CanEdit          bool          `json:"can_edit"`
	NotificationLevel int          `json:"notification_level"`
	Participants     []User        `json:"participants"`
	CreatedBy        User          `json:"created_by"`
	LastPoster       User          `json:"last_poster"`
	Links            []TopicLink   `json:"links"`
}

type TopicLink struct {
	URL        string `json:"url"`
	Title      string `json:"title"`
	Internal   bool   `json:"internal"`
	Attachment bool   `json:"attachment"`
	Reflection bool   `json:"reflection"`
	Clicks     int    `json:"clicks"`
	UserID     int    `json:"user_id"`
	Domain     string `json:"domain"`
	RootDomain string `json:"root_domain"`
}

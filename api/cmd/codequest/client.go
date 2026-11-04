package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// apiClient talks to the CodeQuest API. It is the only place the token is
// sent; lesson HTTP steps use a separate client that never sees it.
type apiClient struct {
	base  string
	token string
	http  *http.Client
	ua    string
}

func newAPIClient(base, token string) *apiClient {
	return &apiClient{
		base: base, token: token, ua: "codequest-cli/" + version,
		http: &http.Client{
			Timeout: 30 * time.Second,
			// The API never redirects; refusing keeps the token on one host.
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
}

// apiError is the API's `{"error":{"code","message"}}` body.
type apiError struct {
	Status     int
	Code       string
	Message    string
	RetryAfter string
}

func (e *apiError) Error() string {
	switch {
	case e.Status == http.StatusUnauthorized:
		return "you are not signed in (or your token expired); run `codequest login`"
	case e.Status == http.StatusPaymentRequired:
		return "this lesson is part of a paid plan: " + e.Message
	case e.Code == "chapter_locked":
		return "this lesson is locked: " + e.Message
	case e.Code == "rate_limited":
		if e.RetryAfter != "" {
			return fmt.Sprintf("too many requests; try again in %ss", e.RetryAfter)
		}
		return "too many requests; try again shortly"
	case e.Message != "":
		return e.Message
	default:
		return fmt.Sprintf("the API answered %d", e.Status)
	}
}

func (c *apiClient) do(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+"/api/v1"+path, body)
	if err != nil {
		return err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.ua)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("cannot reach CodeQuest at %s: %w", c.base, err)
	}
	defer func() { _ = resp.Body.Close() }()
	b, err := io.ReadAll(io.LimitReader(resp.Body, 4<<20))
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 300 {
		var eb struct {
			Error struct{ Code, Message string } `json:"error"`
		}
		_ = json.Unmarshal(b, &eb)
		return &apiError{Status: resp.StatusCode, Code: eb.Error.Code, Message: eb.Error.Message, RetryAfter: resp.Header.Get("Retry-After")}
	}
	if out == nil || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("unexpected response from %s%s: %w", c.base, path, err)
	}
	return nil
}

func isAPIError(err error, code string) bool {
	var ae *apiError
	return errors.As(err, &ae) && ae.Code == code
}

// API shapes the CLI uses.

type tokenResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type meResponse struct {
	User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
	Stats struct {
		TotalXP      int    `json:"total_xp"`
		Level        int    `json:"level"`
		NextLevelXP  int    `json:"next_level_xp"`
		StreakDays   int    `json:"streak_days"`
		LeagueTier   string `json:"league_tier"`
		DailyGoal    int    `json:"daily_goal"`
		TodayLessons int    `json:"today_lessons"`
		Bytes        int    `json:"bytes"`
	} `json:"stats"`
}

type nextResponse struct {
	ExerciseID  string `json:"exercise_id"`
	Title       string `json:"title"`
	CourseTitle string `json:"course_title"`
	ChapterNo   int    `json:"chapter_no"`
	LessonNo    int    `json:"lesson_no"`
}

type exerciseResponse struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Version    int                    `json:"version"`
	Title      string                 `json:"title"`
	Chapter    struct{ Title string } `json:"chapter"`
	Position   struct{ Index int }    `json:"position"`
	PublicView json.RawMessage        `json:"public_view"`
}

type checkRow struct {
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	Feedback string `json:"feedback"`
	Got      string `json:"got"`
	Expected string `json:"expected"`
}

type gradeResult struct {
	Status string     `json:"status"`
	Score  float64    `json:"score"`
	Checks []checkRow `json:"checks"`
}

type outcome struct {
	XPAwarded      int  `json:"xp_awarded"`
	TotalXP        int  `json:"total_xp"`
	Level          int  `json:"level"`
	LevelUp        bool `json:"level_up"`
	StreakDays     int  `json:"streak_days"`
	StreakExtended bool `json:"streak_extended"`
	Rewards        *struct {
		Bytes []struct {
			Label  string `json:"label"`
			Amount int    `json:"amount"`
		} `json:"bytes"`
		Medals []struct {
			Name string `json:"name"`
			Tier string `json:"tier"`
		} `json:"medals"`
		QuestsDone      []string `json:"quests_done"`
		ChapterComplete bool     `json:"chapter_complete"`
	} `json:"rewards"`
}

// submitResponse covers both the 200 and the 202 (pending) bodies, and the
// GET /submissions/{id} poll.
type submitResponse struct {
	SubmissionID string          `json:"submission_id"`
	Status       string          `json:"status"`
	Result       json.RawMessage `json:"result"`
	XPAwarded    int             `json:"xp_awarded"`
	Outcome      *outcome        `json:"outcome"`
}

func (s submitResponse) pending() bool {
	return s.Status == "pending" || (len(s.Result) == 0 || strings.TrimSpace(string(s.Result)) == "null")
}

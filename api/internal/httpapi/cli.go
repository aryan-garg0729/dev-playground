package httpapi

import (
	"errors"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/aryangarg/codequest/api/internal/auth"
	"github.com/aryangarg/codequest/api/internal/store"
)

type cliCodeResponse struct {
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

// createCLICode shows a signed-in browser user a one-time code to paste into
// `codequest login`. A CLI token cannot mint codes.
func (a *api) createCLICode(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.Valid || viaCLI(r) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in on the website first")
		return
	}
	var req struct{}
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	code, exp, err := a.CLI.CreateLoginCode(r.Context(), u.UUID)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, cliCodeResponse{Code: code, ExpiresAt: exp})
}

type cliTokenRequest struct {
	Code  string `json:"code"`
	Label string `json:"label"` // e.g. the machine's hostname
}

type cliTokenResponse struct {
	Token string `json:"token"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// cliToken trades a login code for a bearer token.
func (a *api) cliToken(w http.ResponseWriter, r *http.Request) {
	var req cliTokenRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	tok, user, err := a.CLI.Redeem(r.Context(), req.Code, req.Label)
	if errors.Is(err, auth.ErrInvalidCode) {
		writeError(w, http.StatusBadRequest, "invalid_code", "that code is unknown, used or expired; get a new one in your browser")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	row, err := store.New(a.Pool).GetUser(r.Context(), user)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, cliTokenResponse{Token: tok, Name: row.Name, Email: row.Email})
}

// cliLogout revokes the bearer token it was called with.
func (a *api) cliLogout(w http.ResponseWriter, r *http.Request) {
	tok, ok := bearerToken(r)
	if !ok || !viaCLI(r) {
		writeError(w, http.StatusUnauthorized, "unauthorized", "send the CLI token to revoke")
		return
	}
	if err := a.CLI.Revoke(r.Context(), tok); err != nil {
		a.internal(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type cliNextResponse struct {
	ExerciseID  string `json:"exercise_id"`
	Title       string `json:"title"`
	CourseID    string `json:"course_id"`
	CourseTitle string `json:"course_title"`
	ChapterNo   int    `json:"chapter_no"`
	LessonNo    int    `json:"lesson_no"`
}

// cliNext is the learner's next unpassed project_cli lesson in course order.
// Locks are enforced when the CLI then fetches it with GET /exercises/{id}.
func (a *api) cliNext(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if !u.Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return
	}
	row, err := store.New(a.Pool).NextCLILesson(r.Context(), u.UUID)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "not_found", "no command-line lessons left; pass a lesson id to redo one")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, cliNextResponse{
		ExerciseID: row.ID, Title: row.Title, CourseID: row.CourseID, CourseTitle: row.CourseTitle,
		ChapterNo: int(row.ChapterIndex) + 1, LessonNo: int(row.LessonIndex) + 1,
	})
}

package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/aryangarg/codequest/api/internal/economy"
	"github.com/aryangarg/codequest/api/internal/store"
)

// signedIn writes 401 and returns false for anonymous requests.
func signedIn(w http.ResponseWriter, r *http.Request) bool {
	if !currentUser(r).Valid {
		writeError(w, http.StatusUnauthorized, "unauthorized", "sign in first")
		return false
	}
	return true
}

// econError maps economy errors; every one of them left state unchanged.
func (a *api) econError(w http.ResponseWriter, r *http.Request, err error) {
	var limit economy.LimitError
	switch {
	case errors.Is(err, economy.ErrInsufficientBytes):
		writeError(w, http.StatusConflict, "insufficient_bytes", "not enough Bytes")
	case errors.As(err, &limit):
		writeError(w, http.StatusConflict, "limit_reached", limit.Msg)
	case errors.Is(err, economy.ErrNotOwned):
		writeError(w, http.StatusConflict, "not_owned", err.Error())
	case errors.Is(err, economy.ErrNothingToPeek):
		writeError(w, http.StatusConflict, "nothing_to_peek", "no failed hidden check with a declared input to peek at; nothing was used")
	case errors.Is(err, economy.ErrRequestReused):
		writeError(w, http.StatusConflict, "request_reused", err.Error())
	case errors.Is(err, economy.ErrUnavailable), errors.Is(err, economy.ErrBadRequest):
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
	default:
		a.internal(w, r, err)
	}
}

func writeRaw(w http.ResponseWriter, status int, b json.RawMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

func (a *api) shop(w http.ResponseWriter, r *http.Request) {
	if !signedIn(w, r) {
		return
	}
	v, err := a.Economy.Shop(r.Context(), currentUser(r).UUID)
	if err != nil {
		a.econError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

type requestBody struct {
	RequestID string `json:"request_id"`
}

func (a *api) buy(w http.ResponseWriter, r *http.Request) {
	if !signedIn(w, r) {
		return
	}
	var req requestBody
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	out, err := a.Economy.Buy(r.Context(), currentUser(r).UUID, chi.URLParam(r, "powerup"), req.RequestID)
	if err != nil {
		a.econError(w, r, err)
		return
	}
	writeRaw(w, http.StatusOK, out)
}

func (a *api) usePowerup(w http.ResponseWriter, r *http.Request) {
	if !signedIn(w, r) {
		return
	}
	var req economy.UseRequest
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	user := currentUser(r)
	if req.ExerciseID != "" {
		if locked, err := a.exerciseLocked(r.Context(), r, req.ExerciseID); err != nil || locked {
			a.lockedOr(w, r, err)
			return
		}
	}
	out, err := a.Economy.Use(r.Context(), user.UUID, chi.URLParam(r, "powerup"), req)
	if err != nil {
		a.econError(w, r, err)
		return
	}
	writeRaw(w, http.StatusOK, out)
}

func (a *api) quests(w http.ResponseWriter, r *http.Request) {
	if !signedIn(w, r) {
		return
	}
	v, err := a.Economy.Quests(r.Context(), currentUser(r).UUID)
	if err != nil {
		a.econError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (a *api) openChest(w http.ResponseWriter, r *http.Request) {
	if !signedIn(w, r) {
		return
	}
	var req requestBody
	if err := decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	out, err := a.Economy.OpenDailyChest(r.Context(), currentUser(r).UUID, req.RequestID)
	if err != nil {
		a.econError(w, r, err)
		return
	}
	writeRaw(w, http.StatusOK, out)
}

func (a *api) medals(w http.ResponseWriter, r *http.Request) {
	if !signedIn(w, r) {
		return
	}
	v, err := a.Economy.Medals(r.Context(), currentUser(r).UUID)
	if err != nil {
		a.econError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

// exerciseLocked reports whether the exercise's chapter sits behind an
// undefeated boss for this user.
func (a *api) exerciseLocked(ctx context.Context, r *http.Request, exerciseID string) (bool, error) {
	meta, err := store.New(a.Pool).ExerciseChapter(ctx, exerciseID)
	if err != nil {
		return false, nil // unknown exercises 404 later
	}
	locked, err := a.Economy.LockedChapters(ctx, currentUser(r), meta.CourseID)
	if err != nil {
		return false, err
	}
	return locked[meta.ChapterID], nil
}

func (a *api) lockedOr(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		a.internal(w, r, err)
		return
	}
	writeError(w, http.StatusForbidden, "chapter_locked", "defeat the previous chapter's boss to unlock this chapter")
}

// requireUnlocked refuses lessons in a locked chapter (GET, run and submit).
func (a *api) requireUnlocked(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if locked, err := a.exerciseLocked(r.Context(), r, chi.URLParam(r, "exerciseID")); err != nil || locked {
			a.lockedOr(w, r, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

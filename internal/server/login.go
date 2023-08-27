package server

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func (s *server) handleLogin(w http.ResponseWriter, r *http.Request) {

	var ctx = r.Context()

	session, err := s.sessions.Get(r, "poker-session")
	if err != nil {
		// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
		s.logger.WithError(err).Error("failed to load session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	query := r.URL.Query()

	state := query.Get("state")
	code := query.Get("code")

	if state == "" || code == "" {
		state, err := generateRandomState()
		if err != nil {
			// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
			s.logger.WithError(err).Error("failed to generate state for authentication request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		session.Values["state"] = state

		err = session.Save(r, w)
		if err != nil {
			s.logger.WithError(err).Error("failed to save state for authentication request")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Location", s.authenticator.AuthCodeURL(state))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	spew.Dump(session.Values)

	sessionState, ok := session.Values["state"]
	if !ok {
		s.logger.Error("invalid session, no state stored in session")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if sessionState.(string) != state {
		s.logger.Error("session state does equal query state, discarding")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.authenticator.Exchange(ctx, code)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	idToken, err := s.authenticator.VerifyIDToken(ctx, token)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var profile = make(map[string]any)
	err = idToken.Claims(&profile)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile

	err = session.Save(r, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", "/profile")
	w.WriteHeader(http.StatusTemporaryRedirect)
	return
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

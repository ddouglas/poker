package server

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"poker"

	"github.com/google/uuid"
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
		s.logger.WithError(err).Error("failed to verify id token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var profile = make(map[string]any)
	err = idToken.Claims(&profile)
	if err != nil {
		s.logger.WithError(err).Error("failed to provision claims token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	emailInf, ok := profile["name"]
	if !ok {
		s.logger.WithError(err).Error("profile is missing informaiton necessary to identify user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := s.userRepo.UserByEmail(ctx, emailInf.(string))
	if err != nil {
		s.logger.WithError(err).Error("failed to look up user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user == nil {
		user = &poker.User{
			ID:    uuid.New().String(),
			Name:  fmt.Sprintf("%s %s", profile["given_name"], profile["family_name"]),
			Email: emailInf.(string),
		}

		// Support reaching out to the profile api to retrive emaployee id and profile uri

		err := s.userRepo.SaveUser(ctx, user)
		if err != nil {
			s.logger.WithError(err).Error("failed to save user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	session.Values["userID"] = user.ID

	err = session.Save(r, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(s.authenticator.Provider.Endpoint().AuthURL)

	s.writeRedirectRouteName(w, "dashboard")

}

func (s *server) handleLogout(w http.ResponseWriter, r *http.Request) {

	session, err := s.sessions.Get(r, "poker-session")
	if err != nil {
		// Create an error page and redirect to that. Use session flashing to flash an internal error message of sorts
		s.logger.WithError(err).Error("failed to load session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1

	err = session.Save(r, w)
	if err != nil {
		s.logger.WithError(err).Error("failed to exchange code for token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%sv2/logout?client_id=%s&returnTo=%s", s.authenticator.IssuerURL, s.authenticator.ClientID, url.QueryEscape(s.appURL)))
	w.WriteHeader(http.StatusTemporaryRedirect)

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

package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/cookie"
	"github.com/deeploy-sh/deeploy/internal/server/forms"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/deeploy-sh/deeploy/internal/server/ui/pages"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) LandingView(w http.ResponseWriter, r *http.Request) {
	pages.Landing().Render(r.Context(), w)
}

func (h *UserHandler) AuthView(w http.ResponseWriter, r *http.Request) {
	isCLI := r.URL.Query().Get("cli") == "true"
	session := r.URL.Query().Get("session")

	hasUser, err := h.service.HasUser()
	if err != nil {
		slog.Error("failed to check if user exists", "error", err)
	}

	if hasUser {
		pages.Login(forms.LoginErrors{}, forms.LoginForm{}, isCLI, session).Render(r.Context(), w)
		return
	}

	pages.Register(forms.RegisterErrors{}, forms.RegisterForm{}, isCLI, session).Render(r.Context(), w)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	isCLI := r.URL.Query().Get("cli") == "true"
	session := r.URL.Query().Get("session")

	form := forms.LoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	formErrs := form.Validate()
	if formErrs.HasErrors() {
		pages.Login(formErrs, form, isCLI, session).Render(r.Context(), w)
		return
	}

	token, err := h.service.Login(form.Email, form.Password)
	if err != nil {
		slog.Warn("login failed", "error", err)
		formErrs.Email = "Email or password incorrect"
		formErrs.Password = "Email or password incorrect"
		pages.Login(formErrs, form, isCLI, session).Render(r.Context(), w)
		return
	}

	if isCLI && session != "" {
		auth.SetSessionToken(session, token)
		pages.CliAuthSuccess().Render(r.Context(), w)
		return
	}

	cookie.SetCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	isCLI := r.URL.Query().Get("cli") == "true"
	session := r.URL.Query().Get("session")

	// Single-tenant mode: Only allow registration if no user exists
	// TODO: Remove this check for multi-user support
	hasUser, _ := h.service.HasUser()
	if hasUser {
		http.Redirect(w, r, "/auth", http.StatusSeeOther)
		return
	}

	form := forms.RegisterForm{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("passwordConfirm"),
	}
	formErrs := form.Validate()
	if formErrs.HasErrors() {
		pages.Register(formErrs, form, isCLI, session).Render(r.Context(), w)
		return
	}

	token, err := h.service.Register(form)
	if err == errs.ErrDuplicateEmail {
		formErrs.Email = "Email address is already in use"
		pages.Register(formErrs, form, isCLI, session).Render(r.Context(), w)
		return
	}
	if err != nil {
		slog.Error("user creation failed", "error", err)
		formErrs.General = "Something went wrong. Please try again."
		pages.Register(formErrs, form, isCLI, session).Render(r.Context(), w)
		return
	}

	if isCLI && session != "" {
		auth.SetSessionToken(session, token)
		pages.CliAuthSuccess().Render(r.Context(), w)
		return
	}

	cookie.SetCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie.ClearCookie(w)
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

// PollCLISession is called by TUI to check if auth is complete
func (h *UserHandler) PollCLISession(w http.ResponseWriter, r *http.Request) {
	session := r.URL.Query().Get("session")
	if session == "" {
		http.Error(w, "session required", http.StatusBadRequest)
		return
	}

	token, ready := auth.GetSessionToken(session)

	w.Header().Set("Content-Type", "application/json")

	if !ready {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{"status": "pending"})
		return
	}

	// Token is ready - delete session and return token
	auth.DeleteSession(session)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

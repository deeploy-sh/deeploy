package handlers

import (
	"log/slog"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/cookie"
	"github.com/deeploy-sh/deeploy/internal/server/forms"
	"github.com/deeploy-sh/deeploy/internal/server/ui/pages"
	"github.com/deeploy-sh/deeploy/internal/shared/errs"
	"github.com/deeploy-sh/deeploy/internal/server/service"
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
	port := r.URL.Query().Get("port")

	hasUser, err := h.service.HasUser()
	if err != nil {
		slog.Error("failed to check if user exists", "error", err)
	}

	if hasUser {
		pages.Login(forms.LoginErrors{}, forms.LoginForm{}, isCLI, port).Render(r.Context(), w)
		return
	}

	pages.Register(forms.RegisterErrors{}, forms.RegisterForm{}, isCLI, port).Render(r.Context(), w)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	isCLI := r.URL.Query().Get("cli") == "true"
	port := r.URL.Query().Get("port")

	form := forms.LoginForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	errs := form.Validate()
	if errs.HasErrors() {
		pages.Login(errs, form, isCLI, port).Render(r.Context(), w)
		return
	}

	token, err := h.service.Login(form.Email, form.Password)
	if err != nil {
		slog.Warn("login failed", "error", err)
		errs.Email = "Email or password incorrect"
		errs.Password = "Email or password incorrect"
		pages.Login(errs, form, isCLI, port).Render(r.Context(), w)
		return
	}

	if isCLI && port != "" {
		pages.CliAuthSuccess(port, token).Render(r.Context(), w)
		return
	}

	cookie.SetCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	isCLI := r.URL.Query().Get("cli") == "true"
	port := r.URL.Query().Get("port")

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
	formeErrs := form.Validate()
	if formeErrs.HasErrors() {
		pages.Register(formeErrs, form, isCLI, port).Render(r.Context(), w)
		return
	}

	token, err := h.service.Register(form)
	if err == errs.ErrDuplicateEmail {
		formeErrs.Email = "Email address is already in use"
		pages.Register(formeErrs, form, isCLI, port).Render(r.Context(), w)
		return
	}
	if err != nil {
		slog.Error("user creation failed", "error", err)
		formeErrs.General = "Something went wrong. Please try again."
		pages.Register(formeErrs, form, isCLI, port).Render(r.Context(), w)
		return
	}

	if isCLI && port != "" {
		pages.CliAuthSuccess(port, token).Render(r.Context(), w)
		return
	}

	cookie.SetCookie(w, token)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie.ClearCookie(w)
	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

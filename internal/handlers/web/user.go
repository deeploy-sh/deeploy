package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/axadrn/deeploy/internal/cookie"
	"github.com/axadrn/deeploy/internal/errs"
	"github.com/axadrn/deeploy/internal/forms"
	"github.com/axadrn/deeploy/internal/services"
	"github.com/axadrn/deeploy/internal/ui/pages"
)

type UserHandler struct {
	service services.UserServiceInterface
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) AuthView(w http.ResponseWriter, r *http.Request) {
	isCLI := r.URL.Query().Get("cli") == "true"
	port := r.URL.Query().Get("port")

	hasUser, err := h.service.HasUser()
	if err != nil {
		// TODO: handle error correctly
		fmt.Println(err)
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
		log.Printf("Login failed: %v", err)
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
		log.Printf("User creation failed: %v", err)
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
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

package forms

type ProjectForm struct {
	Title string `json:"title"`
}

type ProjectFormErrors struct {
	Title string `json:"title"`
}

func (f *ProjectForm) Validate() ProjectFormErrors {
	var errors ProjectFormErrors
	if f.Title == "" {
		errors.Title = "Title is required"
	}
	return errors
}

func (e *ProjectFormErrors) HasErrors() bool {
	return e.Title != ""
}

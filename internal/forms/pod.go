package forms

type PodForm struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ProjectID   string `json:"project_id"`
}

type PodFormErrors struct {
	Title     string `json:"title"`
	ProjectID string `json:"project_id"`
}

func (f *PodForm) Validate() PodFormErrors {
	var errors PodFormErrors
	if f.Title == "" {
		errors.Title = "Title is required"
	}
	if f.ProjectID == "" {
		errors.ProjectID = "Project ID is required"
	}
	return errors
}

func (e *PodFormErrors) HasErrors() bool {
	return e.Title != "" || e.ProjectID != ""
}

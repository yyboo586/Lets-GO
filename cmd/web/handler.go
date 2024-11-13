package main

import (
	"errors"
	"fmt"
	"net/http"
	"snippetbox/internal/models"
	"snippetbox/internal/validator"
	"strconv"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-`
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data := &models.TemplateData{
			CurrentYear: time.Now().Year(),
			Form: snippetCreateForm{
				Title:   "",
				Content: "",
				Expires: 365,
			},
		}
		app.render(w, http.StatusOK, "create.html", data)

		return
	}

	var formData snippetCreateForm
	if err := app.decodePostForm(r, &formData); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	formData.CheckField(validator.NotBlank(formData.Title), "title", "This field cannot be blank")
	formData.CheckField(validator.MaxChars(formData.Title, 100), "title", "This field cannot be more than 100 characters long")
	formData.CheckField(validator.NotBlank(formData.Content), "content", "This field cannot be blank")
	formData.CheckField(validator.PermittedInt(formData.Expires, 1, 7, 365), "expires", "This field must be equal to 1, 7 or 365")

	if !formData.Valid() {
		data := &models.TemplateData{
			CurrentYear: time.Now().Year(),
			Form:        formData,
		}
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(formData.Title, formData.Content, formData.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/views/%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := &models.TemplateData{
		CurrentYear: time.Now().Year(),
		Snippet:     snippet,
	}

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *application) snippetList(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &models.TemplateData{
		CurrentYear: time.Now().Year(),
		Snippets:    snippets,
	}
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) decodePostForm(r *http.Request, dst any) (err error) {
	if err = r.ParseForm(); err != nil {
		return err
	}

	if err = app.formDecoder.Decode(dst, r.PostForm); err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		// For all other errors, we return them as normal.
		return err
	}

	return nil
}

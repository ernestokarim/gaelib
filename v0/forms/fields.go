package forms

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

type Control struct {
	Id, Name, Value, Error string
	Help                   string
	Validations            []*Validator

	// If there's an error for this field, reset the value.
	// Useful for passwords, for example.
	ResetValue bool
}

func (c *Control) Build() string {
	err := ""
	if c.Error != "" {
		err = " error"
	}

	// Fields without name are built without label
	if c.Name == "" {
		return fmt.Sprintf(`
			<div class="control-group%s">
				%%s
				<p class="help-block">%s</p>
				<p class="help-block">%s</p>
			</div>
		`, err, c.Error, c.Help)
	}

	return fmt.Sprintf(`
		<div class="control-group%s">
			<label class="control-label" for="%s">%s</label>
			<div class="controls">
				%%s
				<p class="help-block">%s</p>
				<p class="help-block">%s</p>
			</div>
		</div>
	`, err, c.Id, c.Name, c.Error, c.Help)
}

// --------------------------------------------------------

type InputField struct {
	Control            *Control
	Class              []string
	Disabled, ReadOnly bool
	Type, PlaceHolder  string
}

func (f *InputField) Build() string {
	// Tag attributes
	attrs := map[string]string{
		"type": f.Type,
		"id":   f.Control.Id,
		"name": f.Control.Id,
	}

	// Add the value if we don't need to reset it every time
	if !f.Control.ResetValue {
		attrs["value"] = template.HTMLEscapeString(f.Control.Value)
	}

	// Add the disabled flag
	if f.Disabled {
		attrs["disabled"] = "disabled"
	}

	// Add the read-only flag
	if f.ReadOnly {
		attrs["readonly"] = "readonly"
	}

	// The place holder
	if f.PlaceHolder != "" {
		attrs["placeholder"] = f.PlaceHolder
	}

	// The CSS classes
	if f.Class != nil {
		attrs["class"] = strings.Join(f.Class, " ")
	}

	// Build the control HTML
	ctrl := "<input"
	for k, v := range attrs {
		ctrl += fmt.Sprintf(" %s=\"%s\"", k, v)
	}
	ctrl += ">"

	return fmt.Sprintf(f.Control.Build(), ctrl)
}

// --------------------------------------------------------

type SubmitField struct {
	Label                  string
	CancelUrl, CancelLabel string
}

func (f *SubmitField) Build() string {
	// Build the cancel button if present
	cancel := ""
	if f.CancelLabel != "" && f.CancelUrl != "" {
		cancel = fmt.Sprintf(`&nbsp;&nbsp;&nbsp;<a href="%s" class="btn">%s</a>`,
			f.CancelUrl, f.CancelLabel)
	}

	// Build the control
	return fmt.Sprintf(`
		<div class="form-actions">
			<button type="submit" class="btn btn-primary">%s</button>
			%s
		</div>
	`, f.Label, cancel)
}

// --------------------------------------------------------

type SelectField struct {
	Control        *Control
	Class          []string
	Labels, Values []string
}

func (f *SelectField) Build() string {
	// The select tag attributes
	attrs := map[string]string{
		"id":   f.Control.Id,
		"name": f.Control.Id,
	}

	// The CSS classes
	if f.Class != nil {
		attrs["class"] = strings.Join(f.Class, " ")
	}

	ctrl := "<select"
	for k, v := range attrs {
		ctrl += fmt.Sprintf(" %s=\"%s\"", k, v)
	}
	ctrl += ">"

	// Assert the same length precondition, because the error is not
	// very descriptive
	if len(f.Labels) != len(f.Values) {
		panic("labels and values should have the same size")
	}

	for i, label := range f.Labels {
		// Option tag attributes
		attrs := map[string]string{}

		if f.Values[i] == "" {
			// Hide the option if it's the default blank one
			attrs["style"] = "display: none;"
		} else {
			// If it's the currently select one, select it again
			if f.Control.Value == f.Values[i] {
				attrs["selected"] = "selected"
			}

			// Set the value
			attrs["value"] = f.Values[i]
		}

		// Build the HTML of the option tag
		ctrl += "<option"
		for k, v := range attrs {
			ctrl += fmt.Sprintf(" %s=\"%s\"", k, v)
		}
		ctrl += ">" + label + "</option>"
	}

	// Finish the control build
	ctrl += "</select>"

	return fmt.Sprintf(f.Control.Build(), ctrl)
}

// --------------------------------------------------------

type TextAreaField struct {
	Control     *Control
	Class       []string
	Rows        int
	PlaceHolder string
}

func (f *TextAreaField) Build() string {
	// Tag attributes
	attrs := map[string]string{
		"rows":        strconv.FormatInt(int64(f.Rows), 10),
		"id":          f.Control.Id,
		"name":        f.Control.Id,
		"placeholder": f.PlaceHolder,
	}

	// The CSS classes
	if f.Class != nil {
		attrs["class"] = strings.Join(f.Class, " ")
	}

	// Build the control HTML
	ctrl := "<textarea"
	for k, v := range attrs {
		ctrl += fmt.Sprintf(" %s=\"%s\"", k, v)
	}
	ctrl += ">" + template.HTMLEscapeString(f.Control.Value) + "</textarea>"

	return fmt.Sprintf(f.Control.Build(), ctrl)
}

// --------------------------------------------------------

type HiddenField struct {
	Name, Value string
}

func (f *HiddenField) Build() string {
	return fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`, f.Name,
		template.HTMLEscapeString(f.Value))
}

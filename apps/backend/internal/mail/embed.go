package mail

import (
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed template.html
var verifyTemplate string

func loadVerifyEmailTemplate() (*template.Template, error) {
	tmpl, err := template.New("verify-email").Parse(verifyTemplate)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse verify email template: %w",
			err,
		)
	}

	return tmpl, nil
}

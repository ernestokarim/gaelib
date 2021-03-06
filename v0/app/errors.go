package app

import (
	"bytes"

	"conf"
	"github.com/ernestokarim/gaelib/v0/errors"
	"github.com/ernestokarim/gaelib/v0/mail"

	"appengine"
)

func LogError(c appengine.Context, err error) {
	e := errors.New(err).(*errors.Error)
	c.Errorf("%s", e.Error())
	sendErrorByEmail(c, e.Error())
}

func NotFound() error {
	return errors.Code(404)
}

func Forbidden() error {
	return errors.Code(403)
}

func NotAllowed() error {
	return errors.Code(405)
}

func sendErrorByEmail(c appengine.Context, errorStr string) {
	appid := appengine.AppID(c)

	// Try to send an email to the admin if the app is in production
	if !appengine.IsDevAppServer() {
		for _, admin := range conf.ADMIN_EMAILS {
			// Build the template data
			data := map[string]interface{}{
				"Error":    errorStr,
				"UserMail": admin,
				"AppId":    appid,
			}

			// Execute the template
			html := bytes.NewBuffer(nil)
			if err := Template(html, []string{"mails/error"}, data); err != nil {
				c.Errorf("cannot prepare an error email to the admin %s: %s", admin, err)
				continue
			}

			// Send the email to the admin
			m := &mail.Mail{
				To:       admin,
				ToName:   "Administrador",
				From:     "errors@" + appid + ".appspotmail.com",
				FromName: "Aviso de Errores",
				Subject:  "Se ha producido un error en la aplicación",
				Html:     string(html.Bytes()),
			}
			if err := mail.SendMail(c, m); err != nil {
				c.Errorf("cannot send an error email to the admin %s: %s", admin, err)
				continue
			}
		}
	}
}

package emaildelivery

import (
	"errors"
	"net"
	"net/smtp"
	"strconv"
)

type Ingredient struct {
	IngredientInterfaceImpl EmailDeliveryInterface
}

func MakeIngredient(config TypeInputWithService) Ingredient {

	result := Ingredient{
		IngredientInterfaceImpl: config.Service,
	}

	if config.Override != nil {
		result.IngredientInterfaceImpl = config.Override(result.IngredientInterfaceImpl)
	}

	return result
}

func SendSMTPEmail(config SMTPServiceConfig, content SMTPGetContentResult) error {

	fromHeader := "From: " + config.From.Name + " <" + config.From.Email + ">\r\n"
	subject := "Subject: " + content.Subject + "\r\n"
	body := content.Body + "\r\n"
	msg := []byte(fromHeader + subject + body)
	if content.IsHtml {
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
		msg = []byte(fromHeader + subject + mime + body)
	}

	servername := config.Host + ":" + strconv.Itoa(config.Port)

	host, _, err := net.SplitHostPort(servername)
	if err != nil {
		return err
	}

	smtpAuth := smtp.PlainAuth("", config.From.Email, config.Password, host)
	secure := false
	if config.Secure != nil {
		secure = *config.Secure
	}

	if secure {
		// TODO: figure out how to send secure emails
		// tlsconfig := &tls.Config{
		// 	InsecureSkipVerify: true,
		// 	ServerName:         host,
		// }

		// conn, err := tls.Dial("tcp", servername, tlsconfig)
		// if err != nil {
		// 	return err
		// }

		// c, err := smtp.NewClient(conn, host)
		// if err != nil {
		// 	return err
		// }

		// err = c.Auth(smtpAuth)
		// if err != nil {
		// 	return err
		// }

		// if err = c.Mail("From: " + config.From.Name + " <" + config.From.Email + ">"); err != nil {
		// 	return err
		// }

		// if err = c.Rcpt(content.ToEmail); err != nil {
		// 	return err
		// }

		// // Data
		// w, err := c.Data()
		// if err != nil {
		// 	return err
		// }

		// _, err = w.Write([]byte(msg))
		// if err != nil {
		// 	return err
		// }

		// err = w.Close()
		// if err != nil {
		// 	return err
		// }

		// c.Quit()

		return errors.New("secure connection is not supported")
	} else {
		return smtp.SendMail(host+":"+strconv.Itoa(config.Port), smtpAuth, config.From.Email, []string{content.ToEmail}, msg)
	}
}

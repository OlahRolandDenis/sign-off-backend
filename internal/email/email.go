package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

type data_failed struct {
	Token    string
	Decision string
	Message  string
}

type data struct {
	Title       string
	Content     string
	ApprovalURL string
}

type dataagency struct {
	Title       string
	Content     string
	ApprovalURL string
	ClientEmail string
	Token       string
	Dest        string
}

func SendEmail(approvalURL, title, content, to string) error {
	key := os.Getenv("RESEND_API_KEY")
	if key == "" {
		log.Print("Error no ReSend_KEY")
		return fmt.Errorf("No key")
	}

	tmp, err := template.ParseFiles("internal/templates/email.html")
	if err != nil {
		log.Print("Error parsing html")
		return err
	}

	var de data
	de.ApprovalURL = approvalURL
	de.Content = content
	de.Title = title

	var body bytes.Buffer
	err = tmp.Execute(&body, de)
	if err != nil {
		log.Print("Error saving html")
		return err
	}

	var destinatari []string
	destinatari = append(destinatari, to)

	htmlString := body.String()
	client := resend.NewClient(key)

	params := &resend.SendEmailRequest{
		From:    "NodGo <approvals@nodgo.app>",
		To:      destinatari,
		Subject: "Action Required: " + title,
		Html:    htmlString,
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent to %s", to)
	return nil
}

func SendEmailAgency(title, content, clientemail, token, to string) error {
	key := os.Getenv("RESEND_API_KEY")
	if key == "" {
		log.Print("Error no ReSend_KEY")
		return fmt.Errorf("No key")
	}

	tmp, err := template.ParseFiles("internal/templates/emailagency.html")
	if err != nil {
		log.Print("Error parsing html")
		return err
	}

	var de dataagency
	de.ClientEmail = clientemail
	de.Content = content
	de.Title = title
	de.Dest = to
	de.Token = token

	var body bytes.Buffer
	err = tmp.Execute(&body, de)
	if err != nil {
		log.Print("Error saving html")
		return err
	}

	var destinatari []string
	destinatari = append(destinatari, to)

	htmlString := body.String()
	client := resend.NewClient(key)

	params := &resend.SendEmailRequest{
		From:    "NodGo <approvals@nodgo.app>",
		To:      destinatari,
		Subject: "Client hasn't approved: " + title,
		Html:    htmlString,
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent to %s", to)
	return nil
}

func SendCallbackFailedEmail(to, token, decision, errMsg string) error {

	key := os.Getenv("RESEND_API_KEY")
	if key == "" {
		log.Print("Error no ReSend_KEY")
		return fmt.Errorf("No key")
	}

	tmp, err := template.ParseFiles("internal/templates/callback_failed.html")
	if err != nil {
		log.Print("Error parsing html")
		return err
	}

	var de data_failed
	de.Token = token
	de.Decision = decision
	de.Message = errMsg

	var body bytes.Buffer
	err = tmp.Execute(&body, de)
	if err != nil {
		log.Print("Error saving html")
		return err
	}

	var destinatari []string
	destinatari = append(destinatari, to)

	htmlString := body.String()
	client := resend.NewClient(key)

	params := &resend.SendEmailRequest{
		From:    "NodGo <approvals@nodgo.app>",
		To:      destinatari,
		Subject: "Approval callback failed",
		Html:    htmlString,
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email sent to %s", to)
	return nil
}

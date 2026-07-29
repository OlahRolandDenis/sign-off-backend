package jobs

import (
	"Desktop/signoff/internal/callback"
	"Desktop/signoff/internal/db"
	"Desktop/signoff/internal/email"
	"context"
	"log"
	"os"

	"github.com/robfig/cron/v3"
)

func StartReminderJob() {
	c := cron.New()
	c.AddFunc("@hourly", checkPendingApprovals)
	c.Start()
	log.Print("Reminder job started")
}

func checkPendingApprovals() {
	var ctx = context.Background()
	var q = `SELECT a.title, a.content, a.client_email, a.token FROM approvals a WHERE a.status='pending' AND a.nudge_sent=false AND a.created_at < NOW() - INTERVAL '4 hours' LIMIT 100`
	rows, err := db.Pool.Query(ctx, q)
	if err != nil {
		log.Printf("Error finding approvals: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var title string
		var content string
		var to string
		var token string

		err = rows.Scan(&title, &content, &to, &token)
		if err != nil {
			log.Printf("Error extracting values: %v", err)
		}

		Base_URL := os.Getenv("BASE_URL")
		var APV_URL string
		APV_URL = Base_URL + "/approve/" + token

		err = email.SendEmail(APV_URL, title, content, to)
		if err != nil {
			log.Printf("Error sending reminder email: %v", err)
		}

		var q = `UPDATE approvals SET nudge_sent=true WHERE token=$1`
		_, err := db.Pool.Exec(ctx, q, token)
		if err != nil {
			log.Printf("Error setting send true: %v", err)
		} else {
			log.Printf("Nudge Send OK")
		}
	}

	q = `
    SELECT a.title, a.content, a.token, a.client_email, a.callback_url, ag.email 
    FROM approvals a 
    JOIN agencies ag ON a.agency_id = ag.id 
    WHERE a.status = 'pending' 
    AND a.hang_alert_sent = false 
    AND a.created_at < NOW() - INTERVAL '47 hours' 
    LIMIT 100
`
	rows2, err := db.Pool.Query(ctx, q)
	if err != nil {
		log.Printf("Error finding approvals: %v", err)
	}
	defer rows2.Close()

	for rows2.Next() {
		var title string
		var content string
		var to string
		var token string
		var cemail string
		var callback_url string

		err = rows2.Scan(&title, &content, &token, &cemail, &callback_url, &to)
		if err != nil {
			log.Printf("Error extracting values: %v", err)
		}

		err = email.SendEmailAgency(title, content, cemail, token, to)
		if err != nil {
			log.Printf("Error sending reminder email: %v", err)
		}

		if callback_url != "" {
			go func() {
				err := callback.SendCallback(callback_url, token, "expired", "No response within 48 hours")
				if err != nil {
					log.Printf("Callback failed: %v", err)
					email.SendCallbackFailedEmail(to, token, "expired", err.Error())
				}
			}()
		}

		var q = `UPDATE approvals SET hang_alert_sent=true WHERE token=$1`
		_, err := db.Pool.Exec(ctx, q, token)
		if err != nil {
			log.Printf("Error setting hang_alert_sent true: %v", err)
		} else {
			log.Printf("Hang alert send OK")
		}
	}

}

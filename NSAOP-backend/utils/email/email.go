package email

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/jordan-wright/email"

	"nsaop/config"
	"nsaop/model"
	"nsaop/utils"
	"nsaop/utils/constant"
)

type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	command := string(fromServer)
	command = strings.TrimSpace(command)
	command = strings.TrimSuffix(command, ":")
	command = strings.ToLower(command)
	if more {
		if command == "username" {
			return []byte(fmt.Sprintf("%s", a.username)), nil
		} else if command == "password" {
			return []byte(fmt.Sprintf("%s", a.password)), nil
		} else {
			// We've already sent everything.
			return nil, fmt.Errorf("unexpected server challenge: %s", command)
		}
	}
	return nil, nil
}

type Email struct {
	Sender       string
	SenderName   string
	Password     string
	Reciever     string
	RecieverName string
	Subject      string
	Body         []byte
	Server       string
}

func SendEmail(info Email) error {
	from := mail.Address{Name: info.SenderName, Address: info.Sender}
	// to := mail.Address{Name: email.RecieverName, Address: email.Reciever}
	to := mail.Address{Name: info.RecieverName, Address: info.Reciever}

	e := email.NewEmail()
	e.From = from.String()
	e.To = []string{to.String()}
	e.Subject = info.Subject
	e.Cc = []string{from.String()}
	e.HTML = info.Body

	var timeout time.Duration = 30
	errCh := make(chan error, 1)
	go func() {
		err := e.Send(info.Server, &loginAuth{info.Sender, info.Password})
		if err == nil {
			err = errors.Errorf("nil")
		}
		errCh <- err
	}()
	for {
		select {
		case err := <- errCh:
			if err.Error() != "nil" {
				log.Fatalf("fail to send: %v", err.Error())
				return err
			} else {
				return nil
			}
		case <- time.After(timeout * time.Second):
			return errors.Errorf("send email timeout")
		}
	}
}

func SendEmailToUser(user model.User, subject string, body string) error {
	return SendEmail(Email{
		Server:       config.Email.GetString("server"),
		Sender:       config.Email.GetString("sender"),
		SenderName:   config.Email.GetString("sender_name"),
		Password:     config.Email.GetString("password"),
		Reciever:     user.Email,
		RecieverName: user.Username,
		Subject:      subject,
		Body:         []byte(body),
	})
}

type billItem struct {
	Comment string
	UUID    string
	Traffic string
	RMB     string
}

func SendBillEmail(user model.User) {
	now := time.Now().AddDate(0, -1, 0)
	html, subject := GetBillHtml(user, now.Year(), int(now.Month()))
	if html == "" {
		return
	}
	SendEmailToUser(user, subject, html)
}

func GetBillHtml(user model.User, year, month int) (string, string) {
	byt, errEmail := ioutil.ReadFile("config/body.html")
	if errEmail != nil {
		log.Fatalf("%v", errEmail)
	}
	tmpl, _ := template.New("email").Parse(string(byt))
	data := struct {
		Title        string
		Date         string
		Username     string
		Homepage     string
		Year         string
		Items        []billItem
		TotalTraffic string
		TotalRMB     string
	}{
		Username: user.Username,
		Title:    "NSAOP",
		Homepage: "https://nsaop.enjoycolin.top",
		Date:     fmt.Sprintf("%d年%d月", year, month),
		Year:     fmt.Sprintf("%d", year),
	}
	var services []model.Service
	var err error
	whereClause := model.DB
	if whereClause, err = utils.GenerateStatusClause(whereClause, []string{
		constant.StatusOn,
		constant.StatusSuspend,
		constant.StatusPause,
		constant.StatusRetrieve,
	}); err != nil {
		log.Fatalf("internal server error: send email whereClause bug")
	}
	err = model.DB.Model(&user).Where(whereClause).Association(constant.TableService).Find(&services)
	if err != nil {
		log.Fatalf(err.Error())
	}
	if len(services) == 0 {
		return "", ""
	}
	data.Items = make([]billItem, len(services))
	totalTraffic := 0
	totalRMB := 0
	costPerGB := float64(1)
	for i, service := range services {
		traffic := int(rand.Float64() * 10000)
		num := int(math.Ceil(float64(traffic) / 1024))
		data.Items[i] = billItem{
			Comment: service.Comment,
			UUID:    service.ID.String(),
			Traffic: fmt.Sprintf("%d", traffic),
			RMB:     fmt.Sprintf("%.2f", float64(num)*costPerGB),
		}
		totalTraffic += traffic
		totalRMB += num
	}
	data.TotalTraffic = fmt.Sprintf("%d", totalTraffic)
	data.TotalRMB = fmt.Sprintf("%.2f", float64(totalRMB)*costPerGB)

	var htmlbuf bytes.Buffer
	tmpl.Execute(&htmlbuf, data)
	html := htmlbuf.String()
	return html, data.Date + data.Title
}

func SendBillAll() {
	go func() {
		rows, _ := model.DB.Model(&model.User{}).Rows()
		defer rows.Close()

		for rows.Next() {
			var user model.User
			model.DB.ScanRows(rows, &user)
			if user.Role == constant.RoleCustomer {
				SendBillEmail(user)
			}
		}
	}()
}

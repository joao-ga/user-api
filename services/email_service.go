package services

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
	"time"
	"user-api/models"
)

func SendBirthdayEmails(userCollection *mongo.Collection) error {
	today := time.Now().Format("02/01")

	birthdayUsers := bson.M{
		"$expr": bson.M{
			"$eq": bson.A{
				bson.M{"$substr": bson.A{"birthday", 0, 5}},
				today,
			},
		},
	}

	cursor, err := userCollection.Find(context.Background(), birthdayUsers)
	if err != nil {
		return fmt.Errorf("erro ao decodificar usuários: %w", err)
	}

	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return fmt.Errorf("erro ao decodificar usuários: %w", err)
	}

	for _, user := range users {
		if err := sendEmail(user.Email, user.Name); err != nil {
			fmt.Printf("Erro ao enviar e-mail para %s: %v\n", user.Email, err)
		} else {
			fmt.Printf("E-mail enviado com sucesso para %s\n", user.Email)
		}
	}

	return nil
}

func sendEmail(toEmail string, toName string) error {
	from := "joaogbiazon@gmail.com"
	password := "123456"

	subject := "Feliz Aniversário!"
	body := fmt.Sprintf("Olá %s,\n\nParabéns pelo seu aniversário! Tenha um excelente dia!", toName)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", toName)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	return d.DialAndSend(m)
}

package services

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/gomail.v2"
	"os"
	"time"
	"user-api/models"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Erro ao carregar arquivo .env:", err)
	}
}

func SendBirthdayEmails(userCollection *mongo.Collection) error {
	today := time.Now().Format("01-02")

	birthdayUsers := bson.M{
		"$expr": bson.M{
			"$eq": bson.A{
				bson.M{"$dateToString": bson.M{
					"format": "%m-%d",
					"date":   "$birthday",
				}},
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
	from := os.Getenv("EMAIL_USER")
	password := os.Getenv("EMAIL_PASSWORD")

	subject := "Feliz Aniversário!"
	body := fmt.Sprintf("Olá %s,\n\nParabéns pelo seu aniversário! Tenha um excelente dia!", toName)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	return d.DialAndSend(m)
}

package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
	"user-api/models"
)

var userCollection *mongo.Collection

func InitCollection(collection *mongo.Collection) {
	userCollection = collection
}

func CreateUser(c *gin.Context) {
	var userInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	}

	// Bind os dados de entrada
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Definir o formato esperado para a data
	layout := `02/01/2006` // dd/MM/yyyy

	// Tentar fazer o parsing da data
	birthday, err := time.Parse(layout, userInput.Birthday)
	if err != nil {
		// Logar o erro de parsing para mais informações
		fmt.Println("Erro ao parsear a data:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use dd/MM/yyyy"})
		return
	}

	// Criar um usuário com o campo Birthday convertido
	user := models.User{
		ID:       primitive.NewObjectID(),
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: userInput.Password,
		Birthday: birthday,
	}

	// Inserir no banco
	result, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created", "id": result.InsertedID})
}

func GetAllUsers(c *gin.Context) {
	cursor, err := userCollection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(context.Background())

	var users []models.User
	if err := cursor.All(context.Background(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

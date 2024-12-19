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
	"user-api/services"
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

func GetUserById(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var newUser struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	}

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	layout := "02/01/2006"
	birthday, err := time.Parse(layout, newUser.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use dd/MM/yyyy"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"name":     newUser.Name,
			"email":    newUser.Email,
			"password": newUser.Password,
			"birthday": birthday,
		},
	}

	_, err = userCollection.UpdateOne(context.Background(), bson.M{"_id": objectId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated", "id": objectId})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = userCollection.DeleteOne(context.Background(), bson.M{"_id": objectId})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func TestSendBirthdayEmails(c *gin.Context) {
	if err := services.SendBirthdayEmails(userCollection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully sent birthday emails to user"})
}

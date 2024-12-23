package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"user-api/controllers"
	"user-api/routes"
)

var userCollection *mongo.Collection

func init() {
	// Conectar ao MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://jbiazonferreira:qwerty123456@cluster0.82ixr.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"))
	if err != nil {
		log.Fatalf("Erro ao criar cliente MongoDB: %v", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	// Definir a coleção que você vai usar
	userCollection = client.Database("userdb").Collection("users")

	// Inicializar o controlador de usuários
	controllers.InitCollection(userCollection)
}

func main() {
	// Criar o roteador Gin
	router := gin.Default()

	// Rota de teste na "/"
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API está funcionando!"})
	})

	// Rotas de usuários
	routes.UserRoutes(router)

	// Rodar o servidor
	router.Run(":8080")
}

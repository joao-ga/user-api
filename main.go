package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"user-api/routes"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-api/controllers"
)

var userCollection *mongo.Collection

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Conectar ao MongoDB
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI não encontrado no arquivo .env")
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Erro ao criar cliente MongoDB: %v", err)
	}

	err = client.Connect(context.Background())
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	userCollection = client.Database("userdb").Collection("users")

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

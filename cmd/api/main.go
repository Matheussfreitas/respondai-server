package main

import (
	"fmt"
	"goserver/docs"
	"goserver/internal/config"
	"goserver/internal/database"
	"goserver/internal/handler"
	"net/http"
)

// @title RespondAI Server API
// @version 1.0
// @description API para o servidor RespondAI, que gerencia usuários, quizzes e interações com o modelo de linguagem.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	dbConn, err := database.NewPostgres(cfg.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	docs.SwaggerInfo.Host = "localhost:" + cfg.Port
	docs.SwaggerInfo.Schemes = []string{"http"}

	router := handler.NewRouter(dbConn)
	router.RegisterRoutes()

	fmt.Printf("Conectado ao banco de dados: %s\n", cfg.DatabaseUrl)
	fmt.Printf("Swagger disponível em: http://localhost:%s/swagger/index.html\n", cfg.Port)
	fmt.Printf("Scalar disponível em: http://localhost:%s/reference\n", cfg.Port)
	fmt.Printf("Servidor rodando na porta %s\n", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, router.GetHandler()); err != nil {
		panic(err)
	}
}

// @title API de Validação de E-mails
// @version 1.0
// @description API para validação de e-mails
// @contact.name Marcos Oleniuk (Autor)
// @contact.email marcos@moleniuk.com
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
package main

import (
	_ "api-verifica-email-golang/docs"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerfiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/truemail-rb/truemail-go"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	_ "strings"
	"sync"
)

const secretKey = "lkWVR82OF5Ffo3VjVlL8GRkXPn5dt261gq7wOfxDUIokVE0nnGcJ5EN1NeNbUliQ"

// BulkEmailRequest representa a estrutura para validação de múltiplos e-mails
// @description Lista de e-mails para validação
type BulkEmailRequest struct {
	Emails []string `json:"emails" binding:"required"`
}

// ValidationResponse representa a estrutura da resposta de validação de e-mail
// @description Resposta detalhada da validação de e-mails
type ValidationResponse struct {
	Email      string `json:"email"`
	IsValid    bool   `json:"is_valid"`
	Validation string `json:"validation_type"`
	Domain     string `json:"domain,omitempty"`
	Error      string `json:"error,omitempty"`
}

func jsonResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"status": statusCode, "data": data})
}

func jsonErrorResponse(c *gin.Context, statusCode int, errorMessage string) {
	c.JSON(statusCode, gin.H{"status": statusCode, "error": errorMessage})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != secretKey {
			jsonErrorResponse(c, http.StatusUnauthorized, "Chave de autenticação inválida")
			c.Abort()
			return
		}
		c.Next()
	}
}

// @Summary Valida um único e-mail
// @BasePath /api
// @Security ApiKeyAuth
// @Summary Valida um único e-mail
// @Description Valida um e-mail passado via query string
// @Param email query string true "E-mail para validação"
// @Param X-API-KEY header string true "Chave de autenticação"
// @Success 200 {object} ValidationResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /valida-email [get]
func validateEmail(c *gin.Context, configuration *truemail.Configuration) {
	email := c.Query("email")
	if email == "" {
		jsonErrorResponse(c, http.StatusBadRequest, "Email é obrigatório")
		return
	}

	result, err := truemail.Validate(email, configuration)
	if err != nil {
		jsonErrorResponse(c, http.StatusInternalServerError, "Erro na validação do email")
		return
	}

	response := ValidationResponse{
		Email:      email,
		IsValid:    len(result.Errors) == 0,
		Validation: result.ValidationType,
		Domain:     result.Domain,
	}

	if len(result.Errors) > 0 {
		response.Error = result.Errors[strconv.Itoa(0)]
	}

	jsonResponse(c, http.StatusOK, response)
}

// @Summary Valida múltiplos e-mails
// @BasePath /api
// @Security ApiKeyAuth
// @Summary Valida múltiplos e-mails
// @Description Valida uma lista de e-mails passados no corpo da requisição
// @Param data body BulkEmailRequest true "Lista de e-mails para validação"
// @Param X-API-KEY header string true "Chave de autenticação"
// @Success 200 {array} ValidationResponse
// @Failure 400 {object} map[string]interface{}
// @Router /emails [post]
func validateEmails(c *gin.Context, configuration *truemail.Configuration) {
	var request BulkEmailRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonErrorResponse(c, http.StatusBadRequest, "Formato inválido ou lista de e-mails ausente")
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []ValidationResponse

	for _, email := range request.Emails {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()

			result, err := truemail.Validate(email, configuration)
			response := ValidationResponse{
				Email:      email,
				IsValid:    len(result.Errors) == 0,
				Validation: result.ValidationType,
				Domain:     result.Domain,
			}
			if err != nil || len(result.Errors) > 0 {
				response.Error = "Erro na validação "
			}

			mu.Lock()
			results = append(results, response)
			mu.Unlock()
		}(email)
	}

	wg.Wait()
	jsonResponse(c, http.StatusOK, results)
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	configuration, err := truemail.NewConfiguration(
		truemail.ConfigurationAttr{
			VerifierEmail:         "marcos@ajrorato.ind.br",
			ConnectionTimeout:     3,
			ResponseTimeout:       3,
			ConnectionAttempts:    2,
			ValidationTypeDefault: "smtp",
			SmtpFailFast:          true,
			SmtpErrorBodyPattern:  `.*(550|user not found|account).*`,
			EmailPattern:          `\A[a-zA-Z0-9._%+-]{1,32}@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}\z`,
		},
	)
	if err != nil {
		log.Fatalf("Erro ao criar configuração da API: %v", err)
	}

	router := gin.Default()
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erro ao obter o diretório de trabalho: %v", err)
	}
	templatePath := filepath.Join(wd, "templates", "*.tmpl")
	log.Printf("Carregando templates de: %s", templatePath)
	router.LoadHTMLGlob(templatePath)

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "API Validação de Emails",
			"body":  "Link para a documentação da API",
			"url":   "/docs",
		})
	})

	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	protected := router.Group("/api")
	protected.Use(authMiddleware())
	{
		protected.GET("/valida-email", func(c *gin.Context) {
			validateEmail(c, configuration)
		})

		protected.POST("/emails", func(c *gin.Context) {
			validateEmails(c, configuration)
		})
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

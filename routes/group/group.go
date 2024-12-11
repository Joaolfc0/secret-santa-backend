package group

import (
	groupHandler "service-secret-santa/handlers/group"

	"github.com/gin-gonic/gin"
)

// Routes sets up the routes for the group resource
func Routes(defaultGroup *gin.RouterGroup, handler groupHandler.Handler) {
	groupsGroup := defaultGroup.Group("/group")
	{
		//deafault group = http://localhost:8080/secret-santa/
		// Rota para criar um grupo
		groupsGroup.POST("", handler.CreateGroup)

		// Rota para obter um grupo pelo ID
		groupsGroup.GET("/:id", handler.GetGroup)

		// Rota para atualizar um grupo pelo ID
		groupsGroup.PUT("/:id", handler.UpdateGroup)

		// Rota para deletar um grupo pelo ID
		groupsGroup.DELETE("/:id", handler.DeleteGroup)

		// Rota para adicionar um participante ao grupo
		groupsGroup.POST("/:id/add-participant", handler.AddParticipant)

		// Rota para gerar os matches dos participantes do grupo
		groupsGroup.POST("/:id/match-participants", handler.MatchParticipants)

		// Rota para obter o match de um participante
		groupsGroup.GET("/:id/my-match", handler.GetMyMatch)

		//Rota para obter todos grupos disponíveis.
		groupsGroup.GET("", handler.GetAllGroups)

	}
}

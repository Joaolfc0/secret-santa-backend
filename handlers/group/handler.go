package group

import (
	"net/http"
	"service-secret-santa/customError"
	"service-secret-santa/models"
	"service-secret-santa/services/group"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateGroup(c *gin.Context)
	GetGroup(c *gin.Context)
	UpdateGroup(c *gin.Context)
	DeleteGroup(c *gin.Context)
	GetMyMatch(c *gin.Context)
	GetAllGroups(c *gin.Context)
	MatchParticipants(c *gin.Context)
	AddParticipant(c *gin.Context)
}

type resource struct {
	svc group.Service
}

// CreateGroup godoc
//
// @Summary 	Create a new group
// @Description Create a new group with a name
// @Tags 		group
// @Accept  	json
// @Produce  	json
// @Param 		body 		body 		models.Group 	true 	"Group object"
// @Success 	201 		{object} 	models.Group
// @Failure		400 		"{"error": "Bad Request."}"
// @Failure 	500 		"{"error": "Internal Server Error."}"
// @Router 		/group [post]
func (r *resource) CreateGroup(c *gin.Context) {
	var group models.Group

	if err := c.ShouldBindJSON(&group); err != nil {
		customErr := customError.NewCustomError(customError.WithBadRequest(err.Error(), "Invalid request body"))
		c.JSON(customErr.Status, customErr)
		return
	}

	err := group.Validate()
	if err != nil {
		customErr := customError.NewCustomError(customError.WithBadRequest(err.Error(), "Validation error"))
		c.JSON(customErr.Status, customErr)
		return
	}

	result, createErr := r.svc.CreateGroup(&group)
	if createErr != nil {
		c.JSON(createErr.Status, createErr)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetGroup godoc
//
// @Summary 	Get a group by ID
// @Description Retrieve details of a specific group by its ID
// @Tags 		group
// @Produce  	json
// @Param 		id 			path 		string 		true 	"Group ID"
// @Success 	200 		{object} 	models.Group
// @Failure		404 		"{"error": "Not Found."}"
// @Failure 	500 		"{"error": "Internal Server Error."}"
// @Router 		/group/{id} [get]
func (r *resource) GetGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		customErr := customError.NewCustomError(customError.WithBadRequest("Group id is empty", "Invalid request params"))
		c.JSON(customErr.Status, customErr)
	}

	group, err := r.svc.GetGroupByID(id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, group)
}

// UpdateGroup godoc
//
// @Summary 	Update a group
// @Description Update details of a specific group
// @Tags 		group
// @Accept  	json
// @Produce  	json
// @Param 		id 			path 		string 		true 	"Group ID"
// @Param 		body 		body 		models.Group 	true 	"Updated group object"
// @Success 	200 		{object} 	models.Group
// @Failure		400 		"{"error": "Bad Request."}"
// @Failure		404 		"{"error": "Not Found."}"
// @Failure 	500 		"{"error": "Internal Server Error."}"
// @Router 		/group/{id} [put]
func (r *resource) UpdateGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		customErr := customError.NewCustomError(customError.WithBadRequest("Group id is empty", "Invalid request params"))
		c.JSON(customErr.Status, customErr)
	}

	var group models.Group

	if err := c.ShouldBindJSON(&group); err != nil {
		customErr := customError.NewCustomError(customError.WithBadRequest(err.Error(), "Invalid request body"))
		c.JSON(customErr.Status, customErr)
		return
	}

	result, err := r.svc.UpdateGroup(id, &group)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// DeleteGroup godoc
//
// @Summary 	Delete a group
// @Description Delete a specific group by its ID
// @Tags 		group
// @Produce  	json
// @Param 		id 			path 		string 		true 	"Group ID"
// @Success 	204 		"{}"
// @Failure		404 		"{"error": "Not Found."}"
// @Failure 	500 		"{"error": "Internal Server Error."}"
// @Router 		/group/{id} [delete]
func (r *resource) DeleteGroup(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		customErr := customError.NewCustomError(customError.WithBadRequest("Group id is empty", "Invalid request params"))
		c.JSON(customErr.Status, customErr)
	}

	err := r.svc.DeleteGroup(id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddParticipant godoc
//
// @Summary 	Add a participant to a group
// @Description Add a new participant to an existing group
// @Tags 		group
// @Accept  	json
// @Produce  	json
// @Param 		id 			path 		string 		true 	"Group ID"
// @Param 		body 		body 		models.AddParticipantRequest true "Participant to add"
// @Success 	200 		{object} 	models.Group
// @Failure		400 		"{"error": "Bad Request."}"
// @Failure		404 		"{"error": "Not Found."}"
// @Failure 	500 		"{"error": "Internal Server Error."}"
// @Router 		/group/{id}/add-participant [post]
func (r *resource) AddParticipant(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		customErr := customError.NewCustomError(customError.WithBadRequest("Group id is empty", "Invalid request params"))
		c.JSON(customErr.Status, customErr)
	}
	var body models.Participant

	if err := c.ShouldBindJSON(&body); err != nil {
		customErr := customError.NewCustomError(customError.WithBadRequest(err.Error(), "Invalid request body"))
		c.JSON(customErr.Status, customErr)
		return
	}

	result, err := r.svc.AddParticipant(id, &body)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// MatchParticipants godoc
//
// @Summary 	Match participants in a group
// @Description Generate secret matches for participants in a group
// @Tags 		group
// @Produce  	json
// @Param 		id 			path 		string 		true 	"Group ID"
// @Success 	200 		{object} 	models.Group
// @Failure		400 		"{"error": "Bad Request."}"
// @Failure		404 		"{"error": "Not Found."}"
// @Failure 	500 		"{"error": "Internal Server Error."}"
// @Router 		/group/{id}/match-participants [post]
func (r *resource) MatchParticipants(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		customErr := customError.NewCustomError(customError.WithBadRequest("Group id is empty", "Invalid request params"))
		c.JSON(customErr.Status, customErr)
	}

	result, err := r.svc.MatchParticipants(id)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetMyMatch godoc
//
// @Summary 	Get the match for a participant
// @Description Retrieve the participant you are matched to gift in a group
// @Tags 		group
// @Produce  	json
// @Param 		id 			path 		string 		true 	"Group ID"
// @Param 		username	query 		string 		true 	"Participant username"
// @Success 	200 		"string"
// @Failure		400 		"{"error": "Bad Request."}"
// @Failure		404 		"{"error": "Not Found."}"
// @Failure		500 		"{"error": "Internal Server Error."}"
// @Router 		/group/{id}/my-match [get]
func (r *resource) GetMyMatch(c *gin.Context) {
	id := c.Param("id")
	username := c.Query("username")

	if username == "" {
		customErr := customError.NewCustomError(customError.WithBadRequest("Username is required", "Validation error"))
		c.JSON(customErr.Status, customErr)
		return
	}

	if id == "" {
		customErr := customError.NewCustomError(customError.WithBadRequest("Group id is empty", "Invalid request params"))
		c.JSON(customErr.Status, customErr)
	}

	match, err := r.svc.GetMyMatch(id, username)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"match": match})
}

// GetAllGroups godoc
//
// @Summary 	Get all groups
// @Description Retrieve a list of all groups
// @Tags 		group
// @Produce  	json
// @Success 	200 		{array} 	models.Group
// @Failure		500 		"{"error": "Internal Server Error."}"
// @Router 		/group [get]
func (r *resource) GetAllGroups(c *gin.Context) {
	groups, err := r.svc.GetAllGroups()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, groups)
}

func NewGroupHandler(svc group.Service) Handler {
	return &resource{svc: svc}
}

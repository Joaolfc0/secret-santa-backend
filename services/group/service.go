package group

import (
	"math/rand"
	"service-secret-santa/customError"
	"service-secret-santa/models"
	"service-secret-santa/repositories/group"
	"time"
)

type Service interface {
	CreateGroup(group *models.Group) (*models.Group, *customError.CustomError)
	GetGroupByID(id string) (*models.Group, *customError.CustomError)
	UpdateGroup(id string, group *models.Group) (*models.Group, *customError.CustomError)
	DeleteGroup(id string) *customError.CustomError
	AddParticipant(id string, participant *models.Participant) (*models.Group, *customError.CustomError)
	MatchParticipants(id string) (*models.Group, *customError.CustomError)
	GetMyMatch(id string, username string) (string, *customError.CustomError)
	GetAllGroups() ([]*models.Group, *customError.CustomError)
}

type resource struct {
	repo group.Repository
}

func (r *resource) CreateGroup(group *models.Group) (*models.Group, *customError.CustomError) {
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	return r.repo.CreateGroup(group)
}

func (r *resource) GetGroupByID(id string) (*models.Group, *customError.CustomError) {
	return r.repo.GetGroupByID(id)
}

func (r *resource) UpdateGroup(id string, group *models.Group) (*models.Group, *customError.CustomError) {
	group.UpdatedAt = time.Now()
	return r.repo.UpdateGroup(id, group)
}

func (r *resource) DeleteGroup(id string) *customError.CustomError {
	return r.repo.DeleteGroup(id)
}

func (r *resource) AddParticipant(id string, participant *models.Participant) (*models.Group, *customError.CustomError) {
	return r.repo.AddParticipant(id, participant)
}

func (r *resource) MatchParticipants(id string) (*models.Group, *customError.CustomError) {
	group, err := r.repo.GetGroupByID(id)
	if err != nil {
		return nil, err
	}

	if len(group.Participants) < 2 {
		return nil, customError.NewCustomError(customError.WithBadRequest("Not enough participants", "At least two participants are required for matching"))
	}

	// Mapeia os nomes dos participantes
	remainingParticipants := make([]string, len(group.Participants))
	for i, participant := range group.Participants {
		remainingParticipants[i] = participant.Name
	}

	var matches []models.Match

	rand.Seed(time.Now().UnixNano())

	for _, participant := range group.Participants {
		// Filtra a lista de participantes restantes para excluir o próprio participante
		filtered := make([]string, 0)
		for _, name := range remainingParticipants {
			if name != participant.Name {
				filtered = append(filtered, name)
			}
		}

		// Verifica se há participantes disponíveis para match
		if len(filtered) == 0 {
			return nil, customError.NewCustomError(customError.WithInternalServerError("Matching failed", "Could not generate unique matches"))
		}

		// Seleciona um participante aleatório
		randomIndex := rand.Intn(len(filtered))
		match := filtered[randomIndex]

		// Remove o participante selecionado da lista de participantes restantes
		for i, name := range remainingParticipants {
			if name == match {
				remainingParticipants = append(remainingParticipants[:i], remainingParticipants[i+1:]...)
				break
			}
		}

		// Cria um novo match
		matches = append(matches, models.Match{
			First:  participant.Name,
			Second: match,
		})
	}

	// Atualiza os matches no repositório
	updateErr := r.repo.UpdateMatches(id, matches)
	if updateErr != nil {
		return nil, updateErr
	}

	// Atualiza os matches no grupo
	group.Matches = matches

	return group, nil
}

func (r *resource) GetMyMatch(id string, username string) (string, *customError.CustomError) {
	return r.repo.GetMyMatch(id, username)
}

func (r *resource) GetAllGroups() ([]*models.Group, *customError.CustomError) {
	return r.repo.GetAllGroups()
}

func NewGroupService(repo group.Repository) Service {
	return &resource{repo: repo}
}

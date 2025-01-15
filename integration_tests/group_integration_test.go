package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"service-secret-santa/config"
	handlers "service-secret-santa/handlers/group"
	"service-secret-santa/models"
	repos "service-secret-santa/repositories/group"
	routes "service-secret-santa/routes/group"
	services "service-secret-santa/services/group"

	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbClient *mongo.Client
	handler  handlers.Handler
	router   *gin.Engine
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	apiGroup := router.Group("/secret-santa")
	routes.Routes(apiGroup, handler)
	return router
}

func executeRequest(method, url string, body interface{}) *httptest.ResponseRecorder {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestMain(m *testing.M) {
	config.LoadConfig()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "4.4.10",
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017/tcp": {{HostPort: "27017"}},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		dbClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
		return err
	}); err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err)
	}

	groupRepo := repos.NewGroupRepository(dbClient)
	groupSvc := services.NewGroupService(groupRepo)
	handler = handlers.NewGroupHandler(groupSvc)
	router = setupRouter()

	collection := dbClient.Database(config.Cfg.MongoDB).Collection("groups")
	collection.InsertOne(context.TODO(), models.CreateMockGroup())
	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	if err := dbClient.Disconnect(context.Background()); err != nil {
		log.Fatalf("Could not disconnect MongoDB client: %s", err)
	}

	os.Exit(code)
}
func TestCreateGroupSuccess(t *testing.T) {
	newGroup := map[string]interface{}{
		"Name": "Amigos do Trabalho",
	}
	w := executeRequest("POST", "/secret-santa/group", newGroup)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var createdGroup models.Group
	if err := json.Unmarshal(w.Body.Bytes(), &createdGroup); err != nil {
		t.Fatalf("Could not parse response body: %s", err)
	}

	if createdGroup.Name != newGroup["Name"].(string) {
		t.Errorf("Expected group name %s, got %s", newGroup["Name"].(string), createdGroup.Name)
	}

	if createdGroup.Id.IsZero() {
		t.Error("Expected a valid group ID, got zero value")
	}
}

func TestGetGroupByIDSuccess(t *testing.T) {
	newGroup := models.CreateMockGroup()
	w := executeRequest("GET", "/secret-santa/group/"+newGroup.Id.Hex(), nil)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var fetchedGroup models.Group
	if err := json.Unmarshal(w.Body.Bytes(), &fetchedGroup); err != nil {
		t.Fatalf("Could not parse response body: %s", err)
	}

	if fetchedGroup.Name != newGroup.Name {
		t.Errorf("Expected group name %s, got %s", newGroup.Name, fetchedGroup.Name)
	}
}

func TestUpdateGroupSuccess(t *testing.T) {
	updatedGroup := models.CreateMockGroup()
	w := executeRequest("PUT", "/secret-santa/group/"+updatedGroup.Id.Hex(), updatedGroup)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var updatedGroupResponse models.Group
	if err := json.Unmarshal(w.Body.Bytes(), &updatedGroupResponse); err != nil {
		t.Fatalf("Could not parse response body: %s", err)
	}

	if updatedGroupResponse.Name != updatedGroup.Name {
		t.Errorf("Expected group name %s, got %s", updatedGroup.Name, updatedGroupResponse.Name)
	}
}

func TestAddParticipantSuccess(t *testing.T) {
	createdGroup := models.CreateMockGroup()
	newParticipant := map[string]interface{}{
		"Name": "Carlos",
	}
	w := executeRequest("POST", "/secret-santa/group/"+createdGroup.Id.Hex()+"/add-participant", newParticipant)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var updatedGroup models.Group
	if err := json.Unmarshal(w.Body.Bytes(), &updatedGroup); err != nil {
		t.Fatalf("Could not parse response body: %s", err)
	}

	if len(updatedGroup.Participants) != 2 {
		t.Errorf("Expected 1 participant, got %d", len(updatedGroup.Participants))
	}

	if updatedGroup.Participants[1].Name != newParticipant["Name"].(string) {
		t.Errorf("Expected participant name %s, got %s", newParticipant["Name"].(string), updatedGroup.Participants[0].Name)
	}
}

func TestDeleteGroupSuccess(t *testing.T) {
	createdGroup := models.CreateMockGroup()
	w := executeRequest("DELETE", "/secret-santa/group/"+createdGroup.Id.Hex(), nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, w.Code)
	}

	w = executeRequest("GET", "/secret-santa/group/"+createdGroup.Id.Hex(), nil)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

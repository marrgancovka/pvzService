package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/marrgancovka/pvzService/internal/models"
	"github.com/marrgancovka/pvzService/internal/pkg/db"
	"github.com/marrgancovka/pvzService/internal/pkg/grpcconn"
	"github.com/marrgancovka/pvzService/internal/pkg/jwter"
	"github.com/marrgancovka/pvzService/internal/pkg/metrics"
	"github.com/marrgancovka/pvzService/internal/pkg/middleware"
	"github.com/marrgancovka/pvzService/internal/pkg/servers/mainServer"
	"github.com/marrgancovka/pvzService/internal/services/auth"
	authHandler "github.com/marrgancovka/pvzService/internal/services/auth/delivery/http"
	authRepository "github.com/marrgancovka/pvzService/internal/services/auth/repo"
	authUsecase "github.com/marrgancovka/pvzService/internal/services/auth/usecase"
	"github.com/marrgancovka/pvzService/internal/services/pvz"
	pvzHandler "github.com/marrgancovka/pvzService/internal/services/pvz/delivery/http"
	pvzRepository "github.com/marrgancovka/pvzService/internal/services/pvz/repo"
	pvzUsecase "github.com/marrgancovka/pvzService/internal/services/pvz/usecase"
	"github.com/marrgancovka/pvzService/migrations"
	"github.com/marrgancovka/pvzService/pkg/builder"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PVZIntegrationSuite struct {
	suite.Suite
	app       *fxtest.App
	serverURL string
	client    *http.Client
}

func (s *PVZIntegrationSuite) SetupSuite() {
	s.app = fxtest.New(
		s.T(),
		fx.Provide(
			func() *slog.Logger {
				return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelDebug,
				}))
			},
			Config,
			jwter.New,
			builder.SetupBuilder,
			db.NewPostgresPool,
			db.NewPostgresConnect,

			grpcconn.Provide,

			fx.Annotate(metrics.New, fx.As(new(metrics.Metrics))),
			middleware.NewAuthMiddleware,
			middleware.NewMetricsMiddleware,

			fx.Annotate(jwter.New, fx.As(new(auth.JWTer))),
			authHandler.NewHandler,
			fx.Annotate(authUsecase.NewUsecase, fx.As(new(auth.Usecase))),
			fx.Annotate(authRepository.NewRepository, fx.As(new(auth.Repository))),

			pvzHandler.NewHandler,
			fx.Annotate(pvzUsecase.NewUsecase, fx.As(new(pvz.Usecase))),
			fx.Annotate(pvzRepository.NewRepository, fx.As(new(pvz.Repository))),

			mainServer.NewRouter,
		),
		fx.Invoke(
			migrations.RunMigrations,
			mainServer.RunServer,
		),
	)

	s.serverURL = "http://" + testURL
	s.client = http.DefaultClient
	s.app.RequireStart()
}

func (s *PVZIntegrationSuite) authenticate(role *models.DummyLogin) string {
	regBody, err := json.Marshal(role)
	assert.NoError(s.T(), err)

	resp, err := s.makeRequest(http.MethodPost, "/dummyLogin", regBody, "")
	require.NoError(s.T(), err)

	assert.Equal(s.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	require.NoError(s.T(), err)

	token := string(data)
	token = strings.Trim(token, "\"")
	return token
}

func (s *PVZIntegrationSuite) makeRequest(method string, path string, body []byte, token string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(s.T().Context(), method, s.serverURL+"/api/v1"+path, bytes.NewBuffer(body))
	require.NoError(s.T(), err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return s.client.Do(req)
}

func (s *PVZIntegrationSuite) TestFullPVZWorkflow() {
	t := s.T()

	moderatorToken := s.authenticate(&models.DummyLogin{Role: models.RoleModerator})
	employeeToken := s.authenticate(&models.DummyLogin{Role: models.RoleEmployee})

	pvzBody, err := json.Marshal(models.Pvz{
		ID:               uuid.New(),
		RegistrationDate: time.Now(),
		City:             models.CityMoscow,
	})
	require.NoError(t, err)

	resp, err := s.makeRequest(http.MethodPost, "/pvz", pvzBody, moderatorToken)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	createdPvz := &models.Pvz{}
	err = json.NewDecoder(resp.Body).Decode(&createdPvz)
	require.NoError(t, err)
	pvzID := createdPvz.ID

	receptionBody, err := json.Marshal(&models.ReceptionRequest{PvzID: pvzID})
	require.NoError(t, err)

	resp, err = s.makeRequest(http.MethodPost, "/receptions", receptionBody, employeeToken)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	createdReception := &models.Reception{}
	err = json.NewDecoder(resp.Body).Decode(&createdReception)
	require.NoError(t, err)

	for i := 1; i <= 50; i++ {
		var typeProduct models.ProductType
		switch rand.Intn(3) {
		case 0:
			typeProduct = models.TypeElectronics
		case 1:
			typeProduct = models.TypeClothes
		case 2:
			typeProduct = models.TypeShoes
		}

		productData := &models.ProductRequest{
			Type:  typeProduct,
			PvzID: pvzID,
		}
		productBody, err := json.Marshal(productData)
		require.NoError(t, err)

		resp, err = s.makeRequest(http.MethodPost, "/products", productBody, employeeToken)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	}

	resp, err = s.makeRequest(http.MethodPost, fmt.Sprintf("/pvz/%s/close_last_reception", pvzID), []byte{}, employeeToken)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPVZWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	suite.Run(t, new(PVZIntegrationSuite))
}

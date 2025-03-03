package auth

import (
	"context"
	"crypto/subtle"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"

	appconfig "github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/admin"
)

type Service struct {
	userRepo       UserRepo
	apiKeyRepo     APIKeyRepo
	config         Config
	sessionManager *scs.SessionManager
}

type UserRepo interface {
	FindByEmailAndPassword(ctx context.Context, email, password string) (User, error)
}

type APIKeyRepo interface {
	Access(ctx context.Context, keyID uuid.UUID) (APIKey, error)
}

type Config struct {
	SessionStore      scs.Store
	SecretKey         []byte
	SuperUserLogin    []byte
	SuperUserPassword []byte
}

func NewAuthService(userRepo UserRepo, apiKeyRepo APIKeyRepo, config Config) *Service {
	sm := scs.New()

	sm.Lifetime = 72 * time.Hour
	sm.Cookie.Secure = appconfig.GetEnv() == appconfig.ProdEnv
	if config.SessionStore != nil {
		sm.Store = config.SessionStore
	}

	return &Service{
		userRepo:       userRepo,
		apiKeyRepo:     apiKeyRepo,
		config:         config,
		sessionManager: sm,
	}
}

func (s *Service) GetSessionManager() *scs.SessionManager {
	return s.sessionManager
}

func (s *Service) NewSessionAuthContext(ctx context.Context) admin.AuthContext {
	if s.sessionManager.Token(ctx) == "" {
		return nil
	}

	return &sessionAuthContext{
		sm:  s.sessionManager,
		ctx: ctx,
	}
}

func (s *Service) LogInWithAccessToken(ctx context.Context, r LogInRequest) (*LogInResponse, error) {
	user, err := s.userRepo.FindByEmailAndPassword(ctx, r.Email, r.Password)
	if err != nil {
		return nil, err
	}

	claims := newJWTClaims(user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.GetSecretKey())
	if err != nil {
		return nil, err
	}

	return &LogInResponse{
		User:        user,
		AccessToken: accessToken,
	}, nil
}

func (s *Service) LogInWithSession(ctx context.Context, r LogInRequest) error {
	user, err := s.userRepo.FindByEmailAndPassword(ctx, r.Email, r.Password)
	if err != nil {
		return err
	}

	err = s.sessionManager.RenewToken(ctx)
	if err != nil {
		return err
	}

	s.sessionManager.Put(ctx, "user_id", user.ID)
	s.sessionManager.Put(ctx, "is_admin", user.IsAdmin)

	return nil
}

func (s *Service) DestroySession(ctx context.Context) error {
	return s.sessionManager.Destroy(ctx)
}

func (s *Service) GetSecretKey() []byte {
	return s.config.SecretKey
}

func (s *Service) IsSuperUser(username, password string) bool {
	return subtle.ConstantTimeCompare([]byte(username), s.config.SuperUserLogin) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), s.config.SuperUserPassword) == 1
}

func (s *Service) ResolveAPIKey(ctx context.Context, key string) (admin.AuthContext, error) {
	keyID, err := ParseAPIKey(key)
	if err != nil {
		return nil, err
	}

	apiKey, err := s.apiKeyRepo.Access(ctx, keyID)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

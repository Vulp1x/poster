package service

import (
	"context"

	adminservice "github.com/inst-api/poster/gen/admin_service"
	authservice "github.com/inst-api/poster/gen/auth_service"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/pb/instaproxy"
	"github.com/inst-api/poster/pkg/logger"
	"goa.design/goa/v3/security"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type botsStore interface {
	FindReadyBots(ctx context.Context) (domain.BotAccounts, error)
}

// admin_service service example implementation.
// The example methods log the requests and return zero values.
type adminServicesrvc struct {
	auth         authservice.Auther
	userStore    userStore
	instProxyCli instaproxy.InstaProxyClient
	botsStore    botsStore
}

// NewAdminService returns the admin_service service implementation.
func NewAdminService(auth authservice.Auther, userStore userStore, botsStore botsStore, conn *grpc.ClientConn) adminservice.Service {
	return &adminServicesrvc{auth: auth, userStore: userStore, botsStore: botsStore, instProxyCli: instaproxy.NewInstaProxyClient(conn)}
}

// JWTAuth implements the authorization logic for service "admin_service" for
// the "jwt" security scheme.
func (s *adminServicesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {
	//
	// TBD: add authorization logic.
	//
	// In case of authorization failure this function should return
	// one of the generated error structs, e.g.:
	//
	//    return ctx, myservice.MakeUnauthorizedError("invalid token")
	//
	// Alternatively this function may return an instance of
	// goa.ServiceError with a Name field value that matches one of
	// the design error names, e.g:
	//
	//    return ctx, goa.PermanentError("unauthorized", "invalid token")
	//

	return s.auth.JWTAuth(ctx, token, scheme)
}

// AddManager Add Manager by unique_id. Only admins could add Managers
func (s *adminServicesrvc) AddManager(ctx context.Context, p *adminservice.AddManagerPayload) error {
	logger.Debug(ctx, "admin_service.AddManager")

	hashPass, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf(ctx, "failed to generate hash for pass: %v", err)
		return adminservice.InternalError("")
	}

	err = s.userStore.Create(ctx, p.Login, string(hashPass))
	if err != nil {
		logger.Errorf(ctx, "failed to create user: %v", err)
		return adminservice.InternalError("")
	}

	return nil
}

func (s *adminServicesrvc) PushBots(ctx context.Context, p *adminservice.PushBotsPayload) (*adminservice.PushBotsResult, error) {
	bots, err := s.botsStore.FindReadyBots(ctx)
	if err != nil {
		return nil, adminservice.InternalError(err.Error())
	}

	protoBots := bots.ToGRPCProto(ctx)

	req := instaproxy.SaveBotsRequest{Bots: protoBots}

	resp, err := s.instProxyCli.SaveBots(ctx, &req)
	if err != nil {
		return nil, adminservice.InternalError(err.Error())
	}

	logger.Infof(ctx, "saved %d bots, sent %d", resp.BotsSaved, len(protoBots))
	return &adminservice.PushBotsResult{
		SentBots:  len(bots),
		SavedBots: resp.BotsSaved,
		Usernames: resp.Usernames,
	}, nil
}

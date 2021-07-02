package http

import (
	"net/http"

	"github.com/asim/go-micro/v3"
	"github.com/cernbox/ocis-eosprojects/pkg/assets"
	"github.com/cernbox/ocis-eosprojects/pkg/proto/v0"
	"github.com/cernbox/ocis-eosprojects/pkg/version"
	revauser "github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/ocis-pkg/account"
	"github.com/owncloud/ocis/ocis-pkg/middleware"
	ohttp "github.com/owncloud/ocis/ocis-pkg/service/http"
)

// Server initializes the http service and server.
func Server(opts ...Option) ohttp.Service {
	options := newOptions(opts...)
	handler := options.Handler

	svc := ohttp.NewService(
		ohttp.Logger(options.Logger),
		ohttp.Name(options.Name),
		ohttp.Version(options.Config.Server.Version),
		ohttp.Address(options.Config.HTTP.Addr),
		ohttp.Namespace(options.Config.HTTP.Namespace),
		ohttp.Context(options.Context),
		ohttp.Flags(options.Flags...),
	)

	mux := chi.NewMux()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	mux.Use(middleware.NoCache)
	mux.Use(middleware.Cors)
	mux.Use(middleware.Secure)
	mux.Use(middleware.ExtractAccountUUID(
		account.Logger(options.Logger),
		account.JWTSecret(options.Config.TokenManager.JWTSecret)),
	)

	mux.Use(middleware.Version(
		options.Name,
		version.String,
	))

	mux.Use(middleware.Logger(
		options.Logger,
	))

	mux.Use(middleware.Static(
		options.Config.HTTP.Root,
		assets.New(
			assets.Logger(options.Logger),
			assets.Config(options.Config),
		),
		options.Config.HTTP.CacheTTL,
	))

	mux.Route(options.Config.HTTP.Root, func(r chi.Router) {
		r.Get("/api/v0/projects", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := revauser.ContextGetUser(ctx)
			if !ok {
				return
			}

			projects := handler.GetProjects(user)

			rsp := &proto.Response{
				Projects: projects,
			}

			render.Status(r, http.StatusCreated)
			render.JSON(w, r, rsp)
		})
	})

	err := micro.RegisterHandler(svc.Server(), mux)
	if err != nil {
		options.Logger.Fatal().Err(err).Msg("failed to register the handler")
	}

	svc.Init()
	return svc
}

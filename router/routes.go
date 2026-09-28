package router

import (
	"jobsite/handler"
	"jobsite/middleware"
	"jobsite/service"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
)

func NewRouter(logger log.Logger, authSvc service.UserService, jobSvc service.JobService, lookupSvc service.LookupService, jwtSecret string) http.Handler {
	endpoints := handler.NewEndpoints(authSvc, jobSvc, lookupSvc)
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(handler.EncodeErrorResponse),
	}

	r.Handle("/api/register", kithttp.NewServer(
		endpoints.RegisterUser(),
		handler.DecodeRegisterUserRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("POST")

	r.Handle("/api/login", kithttp.NewServer(
		endpoints.Login(),
		handler.DecodeLoginRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("POST")

	r.Handle("/api/forgot-password", kithttp.NewServer(
		endpoints.ForgetPassword(),
		handler.DecodeForgetPasswordRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("POST")

	r.Handle("/api/reset-password", kithttp.NewServer(
		endpoints.ResetPassword(),
		handler.DecodeResetPasswordRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("POST")

	r.Handle("/api/jobs", kithttp.NewServer(
		endpoints.ListJobs(),
		handler.DecodeListJobsRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("GET")

	r.Handle("/api/countries", kithttp.NewServer(
		endpoints.ListCountries(),
		handler.DecodeEmptyRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("GET")

	r.Handle("/api/job-types", kithttp.NewServer(
		endpoints.ListJobTypes(),
		handler.DecodeEmptyRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("GET")

	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.JWTAuthMiddleware(logger, authSvc, jwtSecret, "admin"))

	admin.Handle("/jobs", kithttp.NewServer(
		endpoints.PostJob(),
		handler.DecodePostJobRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("POST")

	admin.Handle("/jobs", kithttp.NewServer(
		endpoints.UpdateJob(),
		handler.DecodeUpdateJobRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("PUT")

	admin.Handle("/jobs/{id}", kithttp.NewServer(
		endpoints.DeleteJob(),
		handler.DecodeDeleteJobRequest,
		handler.EncodeResponse,
		options...,
	)).Methods("DELETE")

	return r
}

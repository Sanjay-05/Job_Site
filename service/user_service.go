package service

import (
	"context"
	"errors"
	"jobsite/helper"
	"jobsite/models"
	"jobsite/repository"
	"jobsite/request"
	"jobsite/utility"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"golang.org/x/crypto/bcrypt"
)

var (
	OtpStore = make(map[string]utility.OTPData)
	OtpMutex sync.Mutex
)

type UserService interface {
	RegisterUser(ctx context.Context, payload request.RegisterUserInput) error
	Login(ctx context.Context, userName, password string) (string, error)
	ForgetPassword(ctx context.Context, userName string) error
	ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error
}

type userService struct {
	userRepo  repository.UserRepository
	logger    log.Logger
	jwtSecret string
	tokenTime time.Duration
}

func NewUserService(userRepo repository.UserRepository, logger log.Logger, jwtSecret string) UserService {
	return &userService{
		userRepo:  userRepo,
		logger:    logger,
		jwtSecret: jwtSecret,
		tokenTime: 24 * time.Hour,
	}
}

func (s *userService) RegisterUser(ctx context.Context, payload request.RegisterUserInput) error {

	if err := helper.ValidateJson(payload); err != nil {
		return err
	}

	if existing, err := s.userRepo.FindUserByName(ctx, payload.UserName); err == nil && existing != nil {
		return errors.New("username already taken")
	}

	if err := helper.ValidatePassword(payload.Password); err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		UserName: payload.UserName,
		Password: string(hashed),
		Email:    payload.Email,
		PhoneNo:  payload.PhoneNo,
		Role:     utility.Role(payload.Role),
	}
	if err := s.userRepo.CreateUser(ctx, &user); err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(ctx context.Context, userName, password string) (string, error) {
	payload := request.LoginInput{
		UserName: userName,
		Password: password,
	}
	if err := helper.ValidateJson(payload); err != nil {
		return "", errors.New("validation error")
	}

	user, err := s.userRepo.FindUserByName(ctx, userName)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := helper.GenerateToken(user, s.tokenTime, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) ForgetPassword(ctx context.Context, userName string) error {
	user, err := s.userRepo.FindUserByName(ctx, userName)
	if err != nil {
		return nil
	}
	otp, err := helper.GenerateOTP()
	if err != nil {
		return err
	}
	hashedOTP, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	OtpMutex.Lock()
	OtpStore[user.UserName] = utility.OTPData{
		OTP:     string(hashedOTP),
		Expires: time.Now().Add(5 * time.Minute),
	}
	OtpMutex.Unlock()
	level.Info(s.logger).Log("message", "otp send to mail", "otp", otp)
	return nil
}

func (s *userService) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {

	if req.NewPassword != req.ConfirmNewPassword {
		return errors.New("confirm password and new password should be same")
	}

	if err := helper.ValidatePassword(req.NewPassword); err != nil {
		return err
	}

	user, err := s.userRepo.FindUserByName(ctx, req.Username)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword))
	if err == nil {
		return errors.New("new password cannot be the same as the old password")
	}
	OtpMutex.Lock()
	data, ok := OtpStore[req.Username]
	OtpMutex.Unlock()
	if !ok || time.Now().After(data.Expires) {
		return errors.New("otp expired or invalid")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(data.OTP), []byte(req.OTP)); err != nil {
		return errors.New("invalid otp")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err := s.userRepo.ResetPassword(ctx, req.Username, string(hashedPassword)); err != nil {
		return err
	}
	helper.StartOTPCleanup(&OtpMutex, OtpStore)
	return nil
}

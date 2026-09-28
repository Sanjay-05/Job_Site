package utility

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kit/log"
)

type OTPData struct {
	OTP     string
	Expires time.Time
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func InitLogger(folderName string, fileName string) (*os.File, log.Logger, error) {
	if err := os.MkdirAll(folderName, 0755); err != nil {
		return nil, nil, err
	}

	path := filepath.Join(folderName, fileName)

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		return nil, nil, err
	}

	multi := io.MultiWriter(os.Stdout, file)

	logger := log.NewLogfmtLogger(multi)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	return file, logger, nil
}

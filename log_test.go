package log_test

import (
	"testing"

	"github.com/nzqpeace/log"
	"os"
)

func TestLogCommon(t *testing.T) {
	f, err := os.OpenFile("t.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		t.Error("open file failed")
	}

	logger := log.NewLogger(f)
	logger.Debug("test debug")
}

func TestMaxFileSize(t *testing.T) {
	logFileName := "./log/max_file_size.log"
	log.SetLogFile(logFileName)
	log.SetMaxFileSize(10 * 1024 * 1024)

	for i := 0; i < 10*1024*1024; i++ {
		log.Info("Test max file size, line:%d", i)
	}
}

func BenchmarkDebug(t *testing.B) {
	f, err := os.OpenFile("t.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		t.Error("open file failed")
	}

	logger := log.NewLogger(f)
	logger.SetFlags(log.LstdFlags | log.Lmicroseconds)

	for i := 0; i < t.N; i++ {
		logger.Debug("test debug")
	}
}

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

	i := 0
	d := 1.456
	s := "this is a test string"

	logger := log.NewLogger(f)
	logger.Debug("test debug")
	logger.Debug("test debug, %d/%f/%s", i, d, s)
	logger.Debugln("test debug")
	logger.Debugln("test debug, ", i, d, s)

	logger.Info("test info, %d/%f/%s", i, d, s)
	logger.Infoln("test info, ", i, d, s)

	logger.Error("test error, %d/%f/%s", i, d, s)
	logger.Errorln("test error, ", i, d, s)
}

func TestMaxFileSize(t *testing.T) {
	logFileName := "./log/max_file_size.log"
	log.SetLogFile(logFileName)
	log.SetMaxFileSize(1 * 1024)

	for i := 0; i < 1*1024; i++ {
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

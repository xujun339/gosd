package start

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type LogConfig struct {
	Level string `json:"level"`
	FileLog bool `json:"fileLog"` // 是否日志记录
	Filename string `json:"filename"`
	MaxSize int `json:"maxsize"` // 单位 M
	MaxAge int `json:"max_age"` // 单位 天
	MaxBackups int `json:"max_backups"`
}

var defaultConfig *LogConfig = &LogConfig{
	"debug",
	true,
	"runtime/logs/sys.log",
	100,
	30,
	100,
}

var Logger *zap.Logger

func InitLogger(cfg *LogConfig)  {
	if cfg == nil {
		cfg = defaultConfig
	}
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.FileLog)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic("initLogger fail, err:" + err.Error())
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	Logger = zap.New(core, zap.AddCaller())
	Logger.Debug("init logger success")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}


func getLogWriter(filename string, maxSize, maxBackup, maxAge int, fileLog bool) zapcore.WriteSyncer {
	writeSyncers := make([]zapcore.WriteSyncer, 2)
	writeSyncers[0] = zapcore.AddSync(os.Stdout)
	if fileLog {
		lumberJackLogger := &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    maxSize,
			MaxBackups: maxBackup,
			MaxAge:     maxAge,
		}
		writeSyncers[1] = zapcore.AddSync(lumberJackLogger)
	}
	return zapcore.NewMultiWriteSyncer(writeSyncers...)
}

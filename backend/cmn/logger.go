package cmn

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var z *zap.Logger

//zapcore.Core一般需要三个配置 --Encoder,writeSyncer,LogLevel

func GetLogger() *zap.Logger {
	if z == nil {
		InitLogger()
	}
	return z
}

func InitLogger() {
	//创建一个logger实例
	//logger, _ = zap.NewProduction()
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel)

	//第一个配置为调用函数信息记录到日志中
	z = zap.New(core, zap.AddCaller())
	//替换logger实例
	zap.ReplaceGlobals(z)
	return
}

// 编写编码器（如何写入日志），根据自身需求修改一些配置
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	//修改时间标签
	encoderConfig.TimeKey = "time"
	//使用可读方式展示时间--默认是float类型的时间
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//使用大写字母记录日志级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 指定日志写入位置
func getLogWriter() zapcore.WriteSyncer {
	//file, _ := os.Create("./log")
	//return zapcore.AddSync(file)
	//lumberjack实现日志切割归档
	//MaxSize:在进行切割之前，日志文件的最大大小(MB)
	//MaxBackups:保留旧文件的最大个数
	//MaxAges :保留旧文件的最大天数
	//Compress:是否压缩/归档旧文件
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "logging",
		MaxSize:    10,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberjackLogger)
}

//
//import (
//	"go.uber.org/zap"
//	"go.uber.org/zap/zapcore"
//)
//
//var isConsoleEnabled = false
//var isPostgresqlEnabled = false
//
//const consoleLoggerTimeLayout = "15:04:05.000"
//
//var rLogLevel = zapcore.InfoLevel
//
//var z *zap.Logger
//
//// InitLogger initialize zap logger
////func InitLogger() {
////	if z != nil {
////		return
////	}
////
////	if viper.IsSet("logger.console.enable") {
////		isConsoleEnabled = viper.GetBool("logger.console.enable")
////	}
////
////	if viper.IsSet("logger.postgresql.enable") {
////		isPostgresqlEnabled = viper.GetBool("logger.postgresql.enable")
////	}
////
////	consoleTimeLayout := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
////		enc.AppendString(t.Format(consoleLoggerTimeLayout))
////	}
////
////	logLevel := zap.LevelEnablerFunc(func(v zapcore.Level) bool {
////		return v >= rLogLevel
////	})
////
////	// 自定义日志配置
////	consoleLoggerCfg := zapcore.EncoderConfig{
////		TimeKey:   "Time",
////		LevelKey:  "Level",
////		NameKey:   "Name",
////		CallerKey: "Caller",
////
////		MessageKey:    "Msg",
////		StacktraceKey: "Stack",
////
////		LineEnding:  zapcore.DefaultLineEnding,
////		EncodeLevel: zapcore.CapitalColorLevelEncoder,
////		EncodeTime:  consoleTimeLayout,
////
////		EncodeDuration: zapcore.StringDurationEncoder,
////		EncodeCaller:   zapcore.ShortCallerEncoder,
////
////		ConsoleSeparator: " ",
////	}
////
////	consoleEncoder := zapcore.NewConsoleEncoder(consoleLoggerCfg)
////
////	consoleSink := zapcore.Lock(os.Stdout)
////
////	// 将某个相对路径转为绝对路径
////	basePath, err := filepath.Abs("./logger")
////
////}

package countlog

type Config struct {
	LogLevel int
	LogFile string
}

func Setup(config Config) {
	if config.LogLevel == 0 {
		config.LogLevel = LevelTrace
	}
	if config.LogFile == "" {
		config.LogFile = "STDOUT"
	}
	logWriter := NewAsyncLogWriter(
		config.LogLevel,
		NewFileLogOutput(config.LogFile))
	logWriter.Start()
	LogWriters = append(LogWriters, logWriter)
}
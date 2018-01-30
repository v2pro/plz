package countlog

//
//type Executor interface {
//	Go(handler func(ctx *Context))
//}
//
//type defaultExecutor struct {
//}
//
//func (executor *defaultExecutor) Go(handler func(ctx *Context)) {
//	go func() {
//		handler(Ctx(context.Background()))
//	}()
//}
//
//var AsyncLogExecutor Executor = &defaultExecutor{}
//
//type AsyncLogWriter struct {
//	MinLevel       int
//	EventWhitelist map[string]bool
//	eventChan      chan Event
//	isClosed       chan bool
//	LogFormat      LogFormatter
//	LogOutput      LogOutput
//}
//
//func (logWriter *AsyncLogWriter) ShouldLog(level int, event string, properties []interface{}) bool {
//	if logWriter.EventWhitelist[event] {
//		return true
//	}
//	return level >= logWriter.MinLevel
//}
//
//func (logWriter *AsyncLogWriter) WriteLog(level int, event string, properties []interface{}) {
//	select {
//	case logWriter.eventChan <- Event{Level: level, Event: event, Properties: properties}:
//	default:
//		if ShouldLog(LevelTrace) {
//			logWriter.eventChan <- Event{Level: level, Event: event, Properties: properties}
//		} else {
//			// drop on the floor
//		}
//	}
//}
//
//func (logWriter *AsyncLogWriter) Close() {
//	close(logWriter.isClosed)
//}
//
//func (logWriter *AsyncLogWriter) Start() {
//	AsyncLogExecutor.Go(func(ctx *Context) {
//		defer func() {
//			recovered := recover()
//			if recovered != nil {
//				os.Stderr.WriteString(fmt.Sprintf("countlog panic: %v\n", recovered))
//				buf := make([]byte, 1<<16)
//				runtime.Stack(buf, true)
//				os.Stderr.Write(buf)
//				os.Stderr.Sync()
//			}
//			if logWriter.LogOutput != nil {
//				logWriter.LogOutput.Close()
//			}
//		}()
//		for {
//			select {
//			case event := <-logWriter.eventChan:
//				logWriter.handleEvent(event)
//			case <-ctx.Done():
//				for {
//					select {
//					case event := <-logWriter.eventChan:
//						logWriter.handleEvent(event)
//					default:
//						return
//					}
//				}
//			case <-logWriter.isClosed:
//				for {
//					select {
//					case event := <-logWriter.eventChan:
//						logWriter.handleEvent(event)
//					default:
//						return
//					}
//				}
//			}
//		}
//	})
//}
//
//func (logWriter *AsyncLogWriter) handleEvent(event Event) {
//	formattedEvent := logWriter.LogFormat.FormatLog(event)
//	if formattedEvent == nil {
//		return
//	}
//	logWriter.LogOutput.OutputLog(event.Level, event.Properties[1].(int64), formattedEvent)
//}
//
//func NewAsyncLogWriter(minLevel int, output LogOutput) *AsyncLogWriter {
//	writer := &AsyncLogWriter{
//		MinLevel:       minLevel,
//		eventChan:      make(chan Event, 1024),
//		isClosed:       make(chan bool),
//		LogFormat:      &HumanReadableFormat{},
//		LogOutput:      output,
//		EventWhitelist: map[string]bool{},
//	}
//	return writer
//}

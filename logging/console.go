package logging

import (
	"fmt"
	"golang.org/x/net/context"
	"runtime"
)

type ConsoleErrorLogger struct {
}

func (l *ConsoleErrorLogger) LogTracefCtx(ctx context.Context, s string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogDebugfCtx(ctx context.Context, s string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogInfofCtx(ctx context.Context, s string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogWarnfCtx(ctx context.Context, s string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogErrorfCtx(ctx context.Context, format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func (l *ConsoleErrorLogger) LogErrorfCtxWithTrace(ctx context.Context, format string, a ...interface{}) {
	l.LogErrorfWithTrace(format, a...)
}

func (l *ConsoleErrorLogger) LogFatalfCtx(ctx context.Context, format string, a ...interface{}) {
	l.LogErrorf(format, a...)
}

func (l *ConsoleErrorLogger) LogAtLevelfCtx(ctx context.Context, level LogLevel, levelLabel string, format string, a ...interface{}) {
	if l.IsLevelEnabled(level) {
		l.LogErrorf(format, a...)
	}
}

func (l *ConsoleErrorLogger) LogTracef(format string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogDebugf(format string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogInfof(format string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogWarnf(format string, a ...interface{}) {
	return
}

func (l *ConsoleErrorLogger) LogErrorf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func (l *ConsoleErrorLogger) LogErrorfWithTrace(format string, a ...interface{}) {
	trace := make([]byte, 2048)
	runtime.Stack(trace, false)

	format = format + "\n%s"
	a = append(a, trace)

	l.LogErrorf(format, a...)
}

func (l *ConsoleErrorLogger) LogFatalf(format string, a ...interface{}) {
	l.LogErrorf(format, a...)
}

func (l *ConsoleErrorLogger) LogAtLevelf(level LogLevel, levelLabel string, format string, a ...interface{}) {
	if l.IsLevelEnabled(level) {
		l.LogErrorf(format, a...)
	}
}

func (l *ConsoleErrorLogger) IsLevelEnabled(level LogLevel) bool {
	return level >= Error
}

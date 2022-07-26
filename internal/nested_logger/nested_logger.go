package nested_logger

import (
	"io"
	"log"
	"strings"
	"time"
)

type trace struct {
	level       string // name of the traced group
	levelStart  time.Time
	levelBuffer *strings.Builder
}

func NewTrace(level string) trace {
	return trace{
		level:       level,
		levelStart:  time.Now(),
		levelBuffer: &strings.Builder{},
	}
}

type NestedTraceLogger struct {
	logger log.Logger

	traceStack []trace
	mainOutput io.Writer // root log output
	mainPrefix string    // root log prefix
}

func (l *NestedTraceLogger) SetPrefix(prefix string) {
	l.logger.SetPrefix(prefix)
	l.mainPrefix = prefix
}

// StartTrace creates a child logging group with time tracking
func (l *NestedTraceLogger) StartTrace(traceLevel string) {
	l.traceStack = append([]trace{NewTrace(traceLevel)}, l.traceStack...)
	// write all logs at this level to that level's buffer instead of root buffer
	l.logger.SetOutput(l.traceStack[0].levelBuffer)
	l.logger.SetPrefix(l.getPadding() + "⎢ [" + traceLevel + "] ")
}

// EndTrace ends the current child logging group, saving the time it took from the start
func (l *NestedTraceLogger) EndTrace() {
	if len(l.traceStack) == 0 {
		return
	}

	level := l.traceStack[0]
	// Print current group end and duration
	l.logger.SetPrefix(l.getPadding() + "⎣ [" + level.level + "] ")
	l.logger.Println("duration: " + time.Since(level.levelStart).String())

	// Get the accumulated log for the current log group
	accumulatedLog := level.levelBuffer.String()
	nextPadding := l.getPadding()

	// Remove the first element of stack
	if len(l.traceStack) > 1 {
		l.traceStack = l.traceStack[1:len(l.traceStack)]
	} else {
		l.traceStack = make([]trace, 0)
	}

	if len(l.traceStack) == 0 {
		// exited all logging groups, time to write to the root
		l.logger.SetPrefix(l.mainPrefix)
		l.logger.SetOutput(l.mainOutput)
		l.mainOutput.Write([]byte(accumulatedLog))
	} else {
		nextLevel := l.traceStack[0]
		l.logger.SetPrefix(nextPadding + "⎢ [" + nextLevel.level + "] ")
		l.logger.SetOutput(nextLevel.levelBuffer)
		// Append accumulated logs to the parent log group
		nextLevel.levelBuffer.Write([]byte(accumulatedLog))
	}
}

// SetOutput is a wrapper aroung log.Logger's SetOutput method
func (l *NestedTraceLogger) SetOutput(output io.Writer) {
	l.mainOutput = output

	l.logger.SetOutput(output)
}

// SetFlags is a wrapper aroung log.Logger's SetFlags method
func (l *NestedTraceLogger) SetFlags(flag int) {
	l.logger.SetFlags(flag)
}

// Print is a wrapper aroung log.Logger's Print method
func (l *NestedTraceLogger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

// Printf is a wrapper aroung log.Logger's Printf method
func (l *NestedTraceLogger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Println is a wrapper aroung log.Logger's Println method
func (l *NestedTraceLogger) Println(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *NestedTraceLogger) getPadding() string {
	var b strings.Builder
	if len(l.traceStack) == 0 {
		return ""
	}
	for range l.traceStack[1:] {
		b.WriteString("⎢ ")
	}

	return b.String()
}

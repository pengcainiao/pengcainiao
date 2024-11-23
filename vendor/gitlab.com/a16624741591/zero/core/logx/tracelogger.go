package logx

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	defaultLogger = zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "2006-01-02 15:04:05"
	})).With().Timestamp().Logger()
	errorLogger = zerolog.New(
		zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = os.Stderr
			w.TimeFormat = "2006-01-02 15:04:05"
		}),
	).With().Timestamp().Caller().Stack().Logger()
)

type (
	TraceLogger struct {
		logEntry zerolog.Logger
	}
)

type TraceLoggerEvent struct {
	event *zerolog.Event
	level zerolog.Level
	ctx   context.Context
}

func (l *TraceLogger) Debug() *TraceLoggerEvent {
	l.logEntry = defaultLogger
	return &TraceLoggerEvent{
		level: zerolog.InfoLevel,
		event: l.logEntry.Debug()}
}

func (l *TraceLogger) Info() *TraceLoggerEvent {
	l.logEntry = defaultLogger
	return &TraceLoggerEvent{
		level: zerolog.InfoLevel,
		event: l.logEntry.Info()}
}

func (l *TraceLogger) Warn() *TraceLoggerEvent {
	l.logEntry = defaultLogger
	return &TraceLoggerEvent{
		level: zerolog.InfoLevel,
		event: l.logEntry.Warn(),
	}
}

func (l *TraceLogger) Error() *TraceLoggerEvent {
	l.logEntry = errorLogger
	return &TraceLoggerEvent{
		level: zerolog.ErrorLevel,
		event: l.logEntry.Error(),
	}
}

func (l *TraceLogger) Err(err error) *TraceLoggerEvent {
	if err != nil {
		l.logEntry = errorLogger
		event := &TraceLoggerEvent{
			event: l.logEntry.Error().Err(err),
			level: zerolog.ErrorLevel,
		}
		event.traceError(err)
		return event
	}
	return l.Info()
}

func (e *TraceLoggerEvent) traceError(err error) {
	if e.level == zerolog.ErrorLevel {
		span := trace.SpanFromContext(e.ctx)
		if span != nil {
			span.RecordError(err, trace.WithStackTrace(true))
		}
	}
}

func (l *TraceLogger) Fatal() *TraceLoggerEvent {
	l.logEntry = errorLogger
	return &TraceLoggerEvent{event: l.logEntry.Fatal()}
}

func (l *TraceLogger) Panic() *TraceLoggerEvent {
	l.logEntry = errorLogger
	return &TraceLoggerEvent{event: l.logEntry.Panic()}
}

// Enabled return false if the *TraceLoggerEvent is going to be filtered out by
// log level or sampling.
func (e *TraceLoggerEvent) Enabled() bool {
	return e.event.Enabled()
}

// Discard disables the event so Msg(f) won't print it.
func (e *TraceLoggerEvent) Discard() *TraceLoggerEvent {
	e.event.Discard()
	return e
}

// Msg sends the *TraceLoggerEvent with msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *TraceLoggerEvent should be disposed.
// Calling Msg twice can have unexpected result.
func (e *TraceLoggerEvent) Msg(msg string) {
	if e.level == zerolog.ErrorLevel {
		span := trace.SpanFromContext(e.ctx)
		if span != nil {
			span.AddEvent("message", trace.WithAttributes(attribute.String("message", msg)))
		}
	}
	e.event.Msg(msg)
}

// Send is equivalent to calling Msg("").
//
// NOTICE: once this method is called, the *TraceLoggerEvent should be disposed.
func (e *TraceLoggerEvent) Send() {
	e.event.Send()
}

// Msgf sends the event with formatted msg added as the message field if not empty.
//
// NOTICE: once this method is called, the *TraceLoggerEvent should be disposed.
// Calling Msgf twice can have unexpected result.
func (e *TraceLoggerEvent) Msgf(format string, v ...interface{}) {
	e.event.Msgf(fmt.Sprintf(format, v...))
}

// Fields is a helper function to use a map to set fields using type assertion.
func (e *TraceLoggerEvent) Fields(fields map[string]interface{}) *TraceLoggerEvent {
	e.event.Fields(fields)
	return e
}

// Dict adds the field key with a dict to the event context.
// Use zerolog.Dict() to create the dictionary.
func (e *TraceLoggerEvent) Dict(key string, dict *TraceLoggerEvent) *TraceLoggerEvent {
	e.event.Dict(key, dict.event)
	return e
}

// Array adds the field key with an array to the event context.
// Use zerolog.Arr() to create the array or pass a type that
// implement the LogArrayMarshaler interface.
func (e *TraceLoggerEvent) Array(key string, arr zerolog.LogArrayMarshaler) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Array(key, arr)
	return e
}

// Object marshals an object that implement the LogObjectMarshaler interface.
func (e *TraceLoggerEvent) Object(key string, obj zerolog.LogObjectMarshaler) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Object(key, obj)
	return e
}

// Func allows an anonymous func to run only if the event is enabled.
func (e *TraceLoggerEvent) Func(f func(e *TraceLoggerEvent)) *TraceLoggerEvent {
	if e != nil && e.Enabled() {
		f(e)
	}
	return e
}

// EmbedObject marshals an object that implement the LogObjectMarshaler interface.
func (e *TraceLoggerEvent) EmbedObject(obj zerolog.LogObjectMarshaler) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	obj.MarshalZerologObject(e.event)
	return e
}

// Str adds the field key with val as a string to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Str(key, val string) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Str(key, val)
	return e
}

// Strs adds the field key with vals as a []string to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Strs(key string, vals []string) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Strs(key, vals)
	return e
}

// Stringer adds the field key with val.String() (or null if val is nil) to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Stringer(key string, val fmt.Stringer) *TraceLoggerEvent {
	if e == nil {
		return e
	}

	e.event.Stringer(key, val)
	return e
}

// Bytes adds the field key with val as a string to the *TraceLoggerEvent context.
//
// Runes outside of normal ASCII ranges will be hex-encoded in the resulting
// JSON.
func (e *TraceLoggerEvent) Bytes(key string, val []byte) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Bytes(key, val)
	return e
}

// Hex adds the field key with val as a hex string to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Hex(key string, val []byte) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Hex(key, val)
	return e
}

// RawJSON adds already encoded JSON to the log line under key.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (e *TraceLoggerEvent) RawJSON(key string, b []byte) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.RawJSON(key, b)
	return e
}

// AnErr adds the field key with serialized err to the *TraceLoggerEvent context.
// If err is nil, no field is added.
func (e *TraceLoggerEvent) AnErr(key string, err error) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.AnErr(key, err)
	return e
}

// Errs adds the field key with errs as an array of serialized errors to the
// *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Errs(key string, errs []error) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Errs(key, errs)

	return e
}

// Err adds the field "error" with serialized err to the *TraceLoggerEvent context.
// If err is nil, no field is added.
//
// To customize the key name, change zerolog.ErrorFieldName.
//
// If Stack() has been called before and zerolog.ErrorStackMarshaler is defined,
// the err is passed to ErrorStackMarshaler and the result is appended to the
// zerolog.ErrorStackFieldName.
func (e *TraceLoggerEvent) Err(err error) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.traceError(err)
	e.event.Err(err)
	return e
}

// Stack enables stack trace printing for the error passed to Err().
//
// ErrorStackMarshaler must be set for this method to do something.
func (e *TraceLoggerEvent) Stack() *TraceLoggerEvent {
	e.event.Stack()
	return e
}

// Bool adds the field key with val as a bool to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Bool(key string, b bool) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Bool(key, b)
	return e
}

// Bools adds the field key with val as a []bool to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Bools(key string, b []bool) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Bools(key, b)
	return e
}

// Int adds the field key with i as a int to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Int(key string, i int) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Int(key, i)
	return e
}

// Ints adds the field key with i as a []int to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Ints(key string, i []int) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Ints(key, i)
	return e
}

// Int8 adds the field key with i as a int8 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Int8(key string, i int8) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Int8(key, i)
	return e
}

// Ints8 adds the field key with i as a []int8 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Ints8(key string, i []int8) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Ints8(key, i)
	return e
}

// Int16 adds the field key with i as a int16 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Int16(key string, i int16) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Int16(key, i)
	return e
}

// Ints16 adds the field key with i as a []int16 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Ints16(key string, i []int16) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Ints16(key, i)
	return e
}

// Int32 adds the field key with i as a int32 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Int32(key string, i int32) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Int32(key, i)
	return e
}

// Ints32 adds the field key with i as a []int32 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Ints32(key string, i []int32) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Ints32(key, i)
	return e
}

// Int64 adds the field key with i as a int64 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Int64(key string, i int64) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Int64(key, i)
	return e
}

// Ints64 adds the field key with i as a []int64 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Ints64(key string, i []int64) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Ints64(key, i)
	return e
}

// Uint adds the field key with i as a uint to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uint(key string, i uint) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uint(key, i)
	return e
}

// Uints adds the field key with i as a []int to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uints(key string, i []uint) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uints(key, i)
	return e
}

// Uint8 adds the field key with i as a uint8 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uint8(key string, i uint8) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uint8(key, i)
	return e
}

// Uints8 adds the field key with i as a []int8 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uints8(key string, i []uint8) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uints8(key, i)
	return e
}

// Uint16 adds the field key with i as a uint16 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uint16(key string, i uint16) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uint16(key, i)
	return e
}

// Uints16 adds the field key with i as a []int16 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uints16(key string, i []uint16) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uints16(key, i)
	return e
}

// Uint32 adds the field key with i as a uint32 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uint32(key string, i uint32) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uint32(key, i)
	return e
}

// Uints32 adds the field key with i as a []int32 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uints32(key string, i []uint32) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uints32(key, i)
	return e
}

// Uint64 adds the field key with i as a uint64 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uint64(key string, i uint64) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uint64(key, i)
	return e
}

// Uints64 adds the field key with i as a []int64 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Uints64(key string, i []uint64) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Uints64(key, i)
	return e
}

// Float32 adds the field key with f as a float32 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Float32(key string, f float32) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Float32(key, f)
	return e
}

// Floats32 adds the field key with f as a []float32 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Floats32(key string, f []float32) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Floats32(key, f)
	return e
}

// Float64 adds the field key with f as a float64 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Float64(key string, f float64) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Float64(key, f)
	return e
}

// Floats64 adds the field key with f as a []float64 to the *TraceLoggerEvent context.
func (e *TraceLoggerEvent) Floats64(key string, f []float64) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Floats64(key, f)
	return e
}

// Timestamp adds the current local time as UNIX timestamp to the *TraceLoggerEvent context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
//
// NOTE: It won't dedupe the "time" key if the *TraceLoggerEvent (or *Context) has one
// already.
func (e *TraceLoggerEvent) Timestamp() *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Timestamp()
	return e
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *TraceLoggerEvent) Time(key string, t time.Time) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Time(key, t)
	return e
}

// Times adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (e *TraceLoggerEvent) Times(key string, t []time.Time) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Times(key, t)
	return e
}

// Dur adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *TraceLoggerEvent) Dur(key string, d time.Duration) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Dur(key, d)
	return e
}

// Durs adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *TraceLoggerEvent) Durs(key string, d []time.Duration) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Durs(key, d)
	return e
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *TraceLoggerEvent) TimeDiff(key string, t time.Time, start time.Time) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.TimeDiff(key, t, start)
	return e
}

// Interface adds the field key with i marshaled using reflection.
func (e *TraceLoggerEvent) Interface(key string, i interface{}) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.Interface(key, i)
	return e
}

// CallerSkipFrame instructs any future Caller calls to skip the specified number of frames.
// This includes those added via hooks from the context.
func (e *TraceLoggerEvent) CallerSkipFrame(skip int) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.CallerSkipFrame(skip)
	return e
}

// Caller adds the file:line of the caller with the zerolog.CallerFieldName key.
// The argument skip is the number of stack frames to ascend
// Skip If not passed, use the global variable CallerSkipFrameCount
func (e *TraceLoggerEvent) Caller(skip ...int) *TraceLoggerEvent {
	e.event.Caller(skip...)
	return e
}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *TraceLoggerEvent) IPAddr(key string, ip net.IP) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.IPAddr(key, ip)
	return e
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *TraceLoggerEvent) IPPrefix(key string, pfx net.IPNet) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.IPPrefix(key, pfx)
	return e
}

// MACAddr adds MAC address to the event
func (e *TraceLoggerEvent) MACAddr(key string, ha net.HardwareAddr) *TraceLoggerEvent {
	if e == nil {
		return e
	}
	e.event.MACAddr(key, ha)
	return e
}

// NewTraceLogger sets ctx to log, for keeping tracing information.
func NewTraceLogger(ctx context.Context) *TraceLogger {
	return &TraceLogger{
		logEntry: defaultLogger,
	}
}

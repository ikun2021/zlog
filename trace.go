package zlog

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"math/rand"
	"sync"
)

type customTracer struct {
}

func (tr *customTracer) Start(ctx context.Context, name string, options ...trace.SpanStartOption) (context.Context, trace.Span) {
	span := NewSpan()
	ctx = trace.ContextWithSpan(ctx, span)
	return ctx, span
}

func InitTrace() context.Context {
	ct := &customTracer{}
	//开启一个span
	ctx, _ := ct.Start(context.Background(), initSpan)

	return ctx
}
func DebugfCtx(ctx context.Context, msg string, fields ...interface{}) {
	fs := extractTrace(ctx)
	for _, v := range fs {
		msg = msg + " %s"
		fields = append(fields, v)
	}
	DefaultSugarLog.Debugf(msg, fields...)
}
func InfofCtx(ctx context.Context, msg string, fields ...interface{}) {
	fs := extractTrace(ctx)
	for _, v := range fs {
		msg = msg + " %s"
		fields = append(fields, v)
	}
	DefaultSugarLog.Infof(msg, fields...)
}

func WarnfCtx(ctx context.Context, msg string, fields ...interface{}) {
	fs := extractTrace(ctx)
	for _, v := range fs {
		msg = msg + " %s"
		fields = append(fields, v)
	}
	DefaultSugarLog.Warnf(msg, fields...)
}
func ErrorfCtx(ctx context.Context, msg string, fields ...interface{}) {
	fs := extractTrace(ctx)
	for _, v := range fs {
		msg = msg + " %s"
		fields = append(fields, v)
	}
	DefaultSugarLog.Errorf(msg, fields...)
}
func PanicfCtx(ctx context.Context, msg string, fields ...interface{}) {
	fs := extractTrace(ctx)
	for _, v := range fs {
		msg = msg + " %s"
		fields = append(fields, v)
	}
	DefaultSugarLog.Panicf(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fs := extractTraceField(ctx)
	fields = append(fields, fs...)
	DefaultLogger.Debug(msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fs := extractTraceField(ctx)
	fields = append(fields, fs...)
	DefaultLogger.Info(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fs := extractTraceField(ctx)
	fields = append(fields, fs...)
	DefaultLogger.Warn(msg, fields...)
}
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fs := extractTraceField(ctx)
	fields = append(fields, fs...)
	DefaultLogger.Error(msg, fields...)
}
func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	fs := extractTraceField(ctx)
	fields = append(fields, fs...)
	DefaultLogger.Panic(msg, fields...)
}

func extractTraceField(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		fields = append(fields, zap.String("traceId", spanCtx.TraceID().String()))
	}
	if spanCtx.HasSpanID() {
		fields = append(fields, zap.String("spanId", spanCtx.SpanID().String()))
	}
	return fields
}

type idGenerator struct {
	randSource *rand.Rand
	lock       sync.Mutex
}

var (
	defaultIDGenerator = NewIDGenerator()
)

func NewIDGenerator() *idGenerator {
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	randSource := rand.New(rand.NewSource(rngSeed))
	return &idGenerator{
		randSource: randSource,
	}
}
func (g *idGenerator) NewID() trace.TraceID {
	g.lock.Lock()
	defer g.lock.Unlock()
	sid := trace.TraceID{}
	for {
		_, _ = g.randSource.Read(sid[:])
		if sid.IsValid() {
			break
		}
	}
	return sid
}

type Span struct {
	trace.Span
	TraceID trace.TraceID
}

func NewSpan() *Span {
	span := &Span{
		TraceID: defaultIDGenerator.NewID(),
	}
	return span
}
func (s *Span) SpanContext() trace.SpanContext {
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    s.TraceID,
		SpanID:     trace.SpanID{},
		TraceFlags: 0,
		TraceState: trace.TraceState{},
		Remote:     false,
	})

}

func extractTrace(ctx context.Context) []interface{} {
	fields := make([]interface{}, 0)
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		fields = append(fields, "traceId")
		fields = append(fields, spanCtx.TraceID().String())
	}
	if spanCtx.HasSpanID() {
		fields = append(fields, zap.String("spanId", spanCtx.SpanID().String()))
	}
	return fields
}

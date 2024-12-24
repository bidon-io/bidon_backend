package config

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc/status"

	"github.com/getsentry/sentry-go"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// NewGRPCServer creates a new gRPC server with all middleware configured
func NewGRPCServer(logger *zap.Logger) *grpc.Server {
	metrics := NewMetrics()

	// Enable gRPC system log redirection to Zap
	grpcLogger := logger.Named("grpc_server")
	grpczap.ReplaceGrpcLoggerV2(grpcLogger)

	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			metrics.UnaryServerInterceptor(),
			grpczap.UnaryServerInterceptor(grpcLogger),
			unaryErrorInterceptor(),
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(
					newPanicRecoveryHandler(metrics, grpcLogger),
				),
			),
		),
		grpc.ChainStreamInterceptor(
			metrics.StreamServerInterceptor(),
			grpczap.StreamServerInterceptor(grpcLogger),
			streamErrorInterceptor(),
			recovery.StreamServerInterceptor(
				recovery.WithRecoveryHandler(
					newPanicRecoveryHandler(metrics, grpcLogger),
				),
			),
		),
	)

	metrics.ServerMetrics.InitializeMetrics(srv)
	return srv
}

type Metrics struct {
	ServerMetrics *grpcprom.ServerMetrics
	PanicsTotal   prometheus.Counter
}

// NewMetrics creates a new metrics instance with configured collectors
func NewMetrics() *Metrics {
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	panicsTotal := promauto.With(reg).NewCounter(prometheus.CounterOpts{
		Name: "grpc_req_panics_recovered_total",
		Help: "Total number of gRPC requests recovered from internal panic.",
	})
	reg.MustRegister(srvMetrics)

	return &Metrics{
		ServerMetrics: srvMetrics,
		PanicsTotal:   panicsTotal,
	}
}

func (m *Metrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return m.ServerMetrics.UnaryServerInterceptor(
		grpcprom.WithExemplarFromContext(exemplarFromContext),
	)
}

func (m *Metrics) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return m.ServerMetrics.StreamServerInterceptor(
		grpcprom.WithExemplarFromContext(exemplarFromContext),
	)
}

// unaryErrorInterceptor logs errors and sends them to Sentry
func unaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			logger := ctxzap.Extract(ctx)
			logger.Error("gRPC Unary call error",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)

			if hub := sentry.CurrentHub(); hub != nil {
				hub.Scope().SetContext("gRPC", map[string]interface{}{
					"method": info.FullMethod,
					"error":  err.Error(),
				})
				hub.CaptureException(err)
			}
		}
		return resp, err
	}
}

// streamErrorInterceptor logs errors and sends them to Sentry
func streamErrorInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, stream)
		if err != nil {
			logger := ctxzap.Extract(stream.Context())
			logger.Error("gRPC Stream call error",
				zap.String("method", info.FullMethod),
				zap.Bool("isClientStream", info.IsClientStream),
				zap.Bool("isServerStream", info.IsServerStream),
				zap.Error(err),
			)

			if hub := sentry.CurrentHub(); hub != nil {
				hub.Scope().SetContext("gRPC", map[string]interface{}{
					"method":         info.FullMethod,
					"error":          err.Error(),
					"isClientStream": info.IsClientStream,
					"isServerStream": info.IsServerStream,
				})
				hub.CaptureException(err)
			}
		}
		return err
	}
}

func newPanicRecoveryHandler(metrics *Metrics, logger *zap.Logger) func(p interface{}) error {
	return func(p interface{}) error {
		metrics.PanicsTotal.Inc()

		logger.Error("Recovered from panic",
			zap.Any("panic", p),
			zap.String("stack", string(debug.Stack())),
		)

		sentry.CurrentHub().Recover(p)
		return status.Errorf(codes.Internal, "%s", p)
	}
}

// exemplarFromContext extracts trace ID from context for Prometheus exemplars
func exemplarFromContext(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{"traceID": span.TraceID().String()}
	}
	return nil
}

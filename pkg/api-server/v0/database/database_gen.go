// generated by 'threeport-sdk gen' for database init boilerplate - do not edit

package database

import (
	"context"
	"fmt"
	v0 "github.com/threeport/threeport/pkg/api/v0"
	v1 "github.com/threeport/threeport/pkg/api/v1"
	log "github.com/threeport/threeport/pkg/log/v0"
	zap "go.uber.org/zap"
	postgres "gorm.io/driver/postgres"
	gorm "gorm.io/gorm"
	logger "gorm.io/gorm/logger"
	"os"
	"reflect"
	"strings"
	"time"
)

// ZapLogger is a custom GORM logger that forwards log messages to a Zap logger.
type ZapLogger struct {
	Logger *zap.Logger
}

// Init initializes the API database.
func Init(autoMigrate bool, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &ZapLogger{Logger: logger},
		NowFunc: func() time.Time {
			utc, _ := time.LoadLocation("UTC")
			return time.Now().In(utc).Truncate(time.Microsecond)
		},
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// LogMode overrides the standard GORM logger's LogMode method to set the logger mode.
func (zl *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	return zl
}

// Info overrides the standard GORM logger's Info method to forward log messages
// to the zap logger.
func (zl *ZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	fields := make([]zap.Field, 0, len(data))
	for i := 0; i < len(data); i += 2 {
		fields = append(fields, zap.Any(data[i].(string), data[i+1]))
	}
	zl.Logger.Info(msg, fields...)
}

// Warn overrides the standard GORM logger's Warn method to forward log messages
// to the zap logger.
func (zl *ZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	fields := make([]zap.Field, 0, len(data))
	for i := 0; i < len(data); i += 2 {
		fields = append(fields, zap.Any(data[i].(string), data[i+1]))
	}
	zl.Logger.Warn(msg, fields...)
}

// Error overrides the standard GORM logger's Error method to forward log messages
// to the zap logger.
func (zl *ZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	fields := make([]zap.Field, 0, len(data))
	for i := 0; i < len(data); i += 2 {
		if reflect.TypeOf(data[i]).Kind() == reflect.Ptr {
			data[i] = fmt.Sprintf("%+v", data[i])
		}
		fields = append(fields, zap.Any(data[i].(string), data[i+1]))
	}
	zl.Logger.Error(msg, fields...)
}

// Trace overrides the standard GORM logger's Trace method to forward log messages
// to the zap logger.
func (zl *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// use the fc function to get the SQL statement and execution time
	sql, rows := fc()

	// create a new logger with some additional fields
	logger := zl.Logger.With(
		zap.String("type", "sql"),
		zap.String("sql", suppressSensitive(sql)),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", time.Since(begin)),
	)

	// if an error occurred, add it as a field to the logger
	if err != nil {
		logger = logger.With(zap.Error(err))
	}

	// log the message using the logger
	logger.Debug("gorm query")
}

// Return all database init object interfaces.
func GetDbInterfaces() []interface{} {
	return []interface{}{

		&v0.AttachedObjectReference{},
		&v0.AwsAccount{},
		&v0.AwsEksKubernetesRuntimeDefinition{},
		&v0.AwsEksKubernetesRuntimeInstance{},
		&v0.AwsObjectStorageBucketDefinition{},
		&v0.AwsObjectStorageBucketInstance{},
		&v0.AwsRelationalDatabaseDefinition{},
		&v0.AwsRelationalDatabaseInstance{},
		&v0.ControlPlaneComponent{},
		&v0.KubernetesRuntimeDefinition{},
		&v0.KubernetesRuntimeInstance{},
		&v0.Definition{},
		&v0.DomainNameDefinition{},
		&v0.DomainNameInstance{},
		&v0.GatewayDefinition{},
		&v0.GatewayHttpPort{},
		&v0.GatewayInstance{},
		&v0.GatewayTcpPort{},
		&v0.HelmWorkloadDefinition{},
		&v0.HelmWorkloadInstance{},
		&v0.Instance{},
		&v0.ControlPlaneDefinition{},
		&v0.ControlPlaneInstance{},
		&v0.LogBackend{},
		&v0.LogStorageDefinition{},
		&v0.LogStorageInstance{},
		&v0.LoggingDefinition{},
		&v0.LoggingInstance{},
		&v0.MetricsDefinition{},
		&v0.MetricsInstance{},
		&v0.ObservabilityDashboardDefinition{},
		&v0.ObservabilityDashboardInstance{},
		&v0.ObservabilityStackDefinition{},
		&v0.ObservabilityStackInstance{},
		&v0.Profile{},
		&v0.SecretDefinition{},
		&v0.SecretInstance{},
		&v0.TerraformDefinition{},
		&v0.TerraformInstance{},
		&v0.Tier{},
		&v0.WorkloadDefinition{},
		&v0.WorkloadEvent{},
		&v0.WorkloadInstance{},
		&v0.WorkloadResourceDefinition{},
		&v0.WorkloadResourceInstance{},
		&v1.AttachedObjectReference{},
		&v1.WorkloadInstance{},
	}
}

// suppressSensitive supresses messages containing sesitive strings.
func suppressSensitive(msg string) string {
	for _, str := range log.SensitiveStrings() {
		if strings.Contains(msg, str) {
			return fmt.Sprintf("[log message containing %s supporessed]", str)
		}
	}

	return msg
}

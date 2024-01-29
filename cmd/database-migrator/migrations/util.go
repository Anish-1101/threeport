package migrations

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func getGormDbFromContext(ctx context.Context) (*gorm.DB, error) {
	contextGorm := ctx.Value("gormdb")
	if contextGorm == nil {
		return nil, fmt.Errorf("could not retrieve gormdb from ctx")
	}

	var gormDb *gorm.DB
	if g, ok := contextGorm.(*gorm.DB); ok {
		gormDb = g
	} else {
		return nil, fmt.Errorf("could not type convert gormdb from ctx")
	}

	return gormDb, nil
}

package storage

import (
	"cms-project/pkg/integration/storage/s3/driver"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	driver.New,
)

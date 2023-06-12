package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

func platformID(platformID db.PlatformID) admin.PlatformID {
	switch platformID {
	case db.AndroidPlatformID:
		return admin.AndroidPlatformID
	case db.IOSPlatformID:
		return admin.IOSPlatformID
	default:
		return admin.UnknownPlatformID
	}
}

func dbPlatformID(platformID admin.PlatformID) db.PlatformID {
	switch platformID {
	case admin.AndroidPlatformID:
		return db.AndroidPlatformID
	case admin.IOSPlatformID:
		return db.IOSPlatformID
	default:
		return db.UnknownPlatformID
	}
}

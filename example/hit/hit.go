package hit

import (
	"time"

	"automatica.team/di"
	"automatica.team/di/example/db"
)

var _ = (di.D)(&Hitter{})

func New() *Hitter {
	return &Hitter{}
}

type Hitter struct {
	db *db.DB `di:"x/db"`
}

func (Hitter) Name() string {
	return "x/hit"
}

func (h Hitter) New(c di.C) (di.D, error) {
	// Automatically migrate the `hits` table.
	// `h.db` is already accessible.
	return New(), h.db.AutoMigrate(&Hit{})
}

// Hit represents a database model for server hits.
type Hit struct {
	When time.Time
	IP   string
}

// Specify a table name for GORM.
func (Hit) TableName() string {
	return "hits"
}

// Hit adds a new IP hit to the database.
func (h Hitter) Hit(ip string) error {
	return h.db.Create(&Hit{
		When: time.Now(),
		IP:   ip,
	}).Error
}

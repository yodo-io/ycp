package model

// Quota represents a quota of how many instances of a given resource a user can have
type Quota struct {
	ID      uint    `gorm:"primary_key"`
	Type    string  `                             binding:"required"`
	UserID  uint    `                             binding:"required"`
	Value   int     `                             binding:"required"`
	Catalog Catalog `gorm:"foreignkey:Type"`
}

// NewQuota creates and initialises a new quota object
func NewQuota(userID uint, tp string, val int) Quota {
	return Quota{
		UserID: userID,
		Type:   tp,
		Value:  val,
	}
}

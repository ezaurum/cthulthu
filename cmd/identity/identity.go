package identity

import (
	"github.com/ezaurum/cthulthu/database"
	"time"
	"github.com/jinzhu/gorm"
	"crypto/sha256"
	"encoding/base64"
	"github.com/labstack/gommon/random"
)

// 다른
type Identity struct {
	database.Model
	Salt string
}

func (id *Identity) GetNewToken(agentType string, agentString string, expiresAt time.Time) *IdentifyToken {
	hash := sha256.New()
	hashTarget := agentType + " " + agentString
	hashTarget += id.Salt
	hash.Write([]byte(hashTarget))
	sum := hash.Sum(nil)
	key := base64.RawURLEncoding.EncodeToString(sum)

	return &IdentifyToken{
		IdentityID:  id.ID,
		ExpiresAt:   &expiresAt,
		Token:       key,
		AgentString: agentString,
		AgentType:   agentType,
	}
}

/**
자동 로그인에 사용할 객체
 */
type IdentifyToken struct {
	database.Model
	IdentityID  int64  `gorm:"index"`
	Token       string `gorm:"index"`
	AgentType   string
	AgentString string
	ExpiresAt   *time.Time
}

func (it *IdentifyToken) CreateIfNotExist(db *gorm.DB) (interface{}, bool) {
	var tt IdentifyToken
	find := db.Where("token =?", it.Token).Find(&tt)
	if !gorm.IsRecordNotFoundError(find.Error) {
		return tt, false
	}
	r := db.Create(it)
	if nil == r.Error {
		return it, true
	}
	return nil, false
}

func CreateNewIdentity() *Identity {
	return &Identity{
		Salt:random.String(10, random.Alphabetic),
	}
}

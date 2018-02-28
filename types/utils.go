package types

import (
	"github.com/shopspring/decimal"
	ttypes "github.com/tendermint/tendermint/types"
)

type dd ttypes.EvidenceList

type ArticleContent struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt uint
	DNA       string `gorm:"size:255;unique_index"`
	Content   string `gorm:"type:longtext"`
}

type User struct {
	ID          uint            `gorm:"primary_key"`
	CreatedAt   uint
	Address     string          `gorm:"size:255;unique_index" binding:"required"`
	Name        string          `gorm:"type:text" binding:"required"`
	Extra       string          `gorm:"type:text" binding:"required"`
	Signature   string          `gorm:"type:text" binding:"required"`
	Balance     decimal.Decimal `gorm:"type:decimal(65)" binding:"-"`
	TokenBurned int             `gorm:"type:tinyint" binding:"-"`

	// Relations


	// isMember
	IsOwner     bool
}

type TokenLock struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    uint
	UserAddress  string `gorm:"size:255;index"`
	ResourceType uint `gorm:"index"`
	ResourceDNA  string `gorm:"index"`
	Amount       decimal.Decimal `gorm:"type:decimal(65)"`
	Expire       uint `gorm:"type:int unsigned;index"`
}

type PublishLogArgs struct {
	Title       []byte
	ContentHash []byte
	License     []byte
	Extras      []byte
	BlockHash   []byte
	Signature   []byte
	DNA         []byte
}

type LikeLogArgs struct {
	ArticleDNA []byte
	GroupDNA   []byte
	Signature  []byte
}

type CommentLogArgs struct {
	ArticleDNA  []byte
	GroupDNA    []byte
	ContentHash []byte
	Signature   []byte
}

type ShareLogArgs struct {
	ArticleDNA []byte
	GroupsDNA  []byte
	Signature  []byte
}
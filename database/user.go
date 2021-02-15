package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/islazy/str"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

const MinPasswordLength = 8

type User struct {
	ID           uint           `gorm:"primarykey" json:"-""`
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Email        string         `gorm:"uniqueIndex" json:"email"`
	Verification string         `gorm:"index" json:"-"`
	Verified     bool           `gorm:"index" json:"-"`
	Use2FA       bool           `gorm:"default:true" json:"use_2fa"` // i know it's not 2fa yet
	Hash         string         `json:"-"`
	Address      string         `json:"address"`
	Agents       []Agent        `json:"-"`
}

func makeRandomToken() string {
	randomShit := make([]byte, 128)
	rand.Read(randomShit)

	data := append(
		[]byte(strconv.FormatInt(time.Now().UnixNano(), 10)),
		randomShit...)

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func RegisterUser(address, email, password string) (*User, error) {
	if err := checkmail.ValidateFormat(email); err != nil {
		return nil, err
	} else if password = str.Trim(password); len(password) < MinPasswordLength {
		return nil, fmt.Errorf("minimum password length is %d", MinPasswordLength)
	}

	var found User
	if err := db.Where("email=?", email).First(&found).Error; err == nil {
		return nil, fmt.Errorf("email address already used")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("error searching email '%s': %v", email, err)
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error generating password hash: %v", err)
	}

	newUser := User{
		Email:        email,
		Verification: makeRandomToken(),
		Hash:         string(hashedPassword),
		Address:      address,
	}

	if err = db.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("error creating new user: %v", err)
	}

	return &newUser, nil
}

func VerifyUser(verification string) error {
	var found User
	if err := db.Where("verification=?", verification).First(&found).Error; err != nil {
		return err
	} else if found.Verified == true {
		return fmt.Errorf("user already verified")
	} else {
		found.Verified = true
		return db.Save(&found).Error
	}
}

func LoginUser(address, email, password string) (*User, error) {
	var found User
	if err := db.Where("email=?", email).First(&found).Error; err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else if found.Verified == false {
		return nil, fmt.Errorf("account not verified")
	} else if err = bcrypt.CompareHashAndPassword([]byte(found.Hash), []byte(password)); err != nil {
		return nil, nil
	}

	if found.Use2FA {
		found.Verification = makeRandomToken()
	}

	found.Address = address
	found.UpdatedAt = time.Now()
	if err := db.Save(&found).Error; err != nil {
		log.Error("error updating logged in user: %v", err)
	}

	return &found, nil
}

func UpdateUser(user *User, ip string, newPassword string, use2FA bool) (*User, error) {
	if newPassword = str.Trim(newPassword); newPassword != "" && len(newPassword) < MinPasswordLength {
		return nil, fmt.Errorf("minimum password length is %d", MinPasswordLength)
	}

	if newPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("error generating password hash: %v", err)
		}
		user.Hash = string(hashedPassword)
	}

	user.Use2FA = use2FA
	user.UpdatedAt = time.Now()
	user.Address = ip

	return user, db.Save(user).Error
}

func FindUserByID(id int) (*User, error) {
	var found User

	err := db.Preload("Agents", func(db *gorm.DB) *gorm.DB {
		return db.Order("agents.updated_at DESC")
	}).Find(&found, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &found, nil
}

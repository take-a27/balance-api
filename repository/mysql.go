package repository

import (
	"ARIGATOBANK/domain"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

type DB struct {
	conn *gorm.DB
}
type DBConfig struct {
	DBUser string `yaml:"DBUser"`
	DBPass string `yaml:"DBPass"`
	DBName string `yaml:"DBName"`
	DBAddr string `yaml:"DBAddr"`
	DBPort string `yaml:"DBPort"`
}

func NewMySql() *DB {
	db := new(DB)
	dbConfig := new(DBConfig)
	byteData, err := os.ReadFile("config/database.yaml")
	if err != nil {
		logrus.Fatalf("failed to read db config yaml: %s", err)
	}

	err = yaml.Unmarshal(byteData, dbConfig)
	if err != nil {
		logrus.Fatalf("failed to unmarshal db config yaml: %s", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBAddr, dbConfig.DBPort, dbConfig.DBName)
	logrus.Infof("connection string: %s", connectionString)

	dbConn, err := gorm.Open("mysql", connectionString)
	if err != nil {
		logrus.WithFields(logrus.Fields{"db_config": dbConfig}).Fatalf("failed to open db: %s", err)
	}

	db.conn = dbConn
	return db
}

func (db *DB) InsertBalanceOperation(balanceOperationTableName, userBalanceTableName string, bo *domain.BalanceOperation) error {
	tx := db.conn.Begin()
	var err error
	if err = tx.Table(balanceOperationTableName).Create(bo).Error; err != nil {
		logrus.WithFields(logrus.Fields{"balance_operation": bo}).Errorf("failed to insert balance operation: %s", err.Error())
	}

	// get current user balance to calculate amount after balance operation
	currentUb := new(domain.UserBalance)
	if err = db.conn.Table(userBalanceTableName).Where("id = ?", bo.UserId).First(currentUb).Error; err != nil {
		logrus.WithFields(logrus.Fields{"user_id": bo.UserId}).Errorf("failed to get user balance: %s", err.Error())
	}

	// insert new user balance with new amount
	updateUb := domain.NewUserBalance(bo.UserId)
	updateUb.Amount = currentUb.Amount + bo.Difference
	if err = db.conn.Table(userBalanceTableName).Model(updateUb).Where("id = ?", updateUb).Update(updateUb).Error; err != nil {
		logrus.WithFields(logrus.Fields{"user_balance": updateUb}).Errorf("failed to update user balance: %s", err.Error())
	}

	if err != nil {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
		return nil
	}
}

func (db *DB) GetBalanceOperation(balanceOperationTableName string, bo *domain.BalanceOperation) (*domain.BalanceOperation, error) {
	result := new(domain.BalanceOperation)
	if err := db.conn.Table(balanceOperationTableName).Where("user_id = ? AND timestamp = ?", bo.UserId, bo.Timestamp).First(bo).Error; err != nil {
		logrus.WithFields(logrus.Fields{"balance_operation": bo}).Errorf("failed to get balance operation: %s", err.Error())
		return nil, err
	}
	return result, nil
}

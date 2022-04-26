package repository

import (
	"ARIGATOBANK/domain"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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
		return err
	}

	// get current user balance to calculate amount after balance operation
	// if a record with this user id doesn't exist, new one will be created
	currentUb := new(domain.UserBalance)
	err = tx.Table(userBalanceTableName).Where("id = ?", bo.UserId).First(currentUb).Error
	if gorm.IsRecordNotFoundError(err) {
		currentUb = domain.NewUserBalance(bo.UserId)
		if err = tx.Table(userBalanceTableName).Create(currentUb).Error; err != nil {
			logrus.WithFields(logrus.Fields{"user_id": bo.UserId}).Errorf("failed to insert user balance: %s", err.Error())
			return err
		}
	} else if err != nil {
		logrus.WithFields(logrus.Fields{"user_id": bo.UserId}).Errorf("failed to get user balance: %s", err.Error())
		return err
	}

	// insert new user balance with new amount
	updateUb := new(domain.UserBalance)
	updateUb.Id = currentUb.Id
	updateUb.Amount = currentUb.Amount + bo.Difference

	if err = db.conn.Table(userBalanceTableName).Model(updateUb).Where("id = ?", updateUb.Id).Update(*updateUb).Error; err != nil {
		logrus.WithFields(logrus.Fields{"user_balance": updateUb}).Errorf("failed to update user balance: %s", err.Error())
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return err
}

func (db *DB) GetBalanceOperation(balanceOperationTableName string, bo *domain.BalanceOperation) (*domain.BalanceOperation, error) {
	result := new(domain.BalanceOperation)
	if err := db.conn.Table(balanceOperationTableName).Where("user_id = ? AND timestamp = ?", bo.UserId, bo.Timestamp).First(bo).Error; err != nil {
		logrus.WithFields(logrus.Fields{"balance_operation": bo}).Errorf("failed to get balance operation: %s", err.Error())
		return nil, err
	}
	return result, nil
}

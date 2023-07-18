package repository

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"errors"
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"
)

func NewAgentConnectRepository(db *gorm.DB) AgentConnectRepository {
	return &repo{db}
}

type AgentConnectRepository interface {
	Register(data model.AGCRegister) error
	Login(data model.AGCLogin) error
	ChangePassword(data model.AGCChangePassword) error
	Deposit(data model.AGCDeposit) error
	Withdraw(data model.AGCWithdraw) error
}

func (r repo) Register(data model.AGCRegister) error {

	url := fmt.Sprintf("%s/credit-auth/xregister", os.Getenv("AGENT_API"))
	result, err := helper.Post(url, data)
	if err != nil {
		return err
	}

	if result.(map[string]interface{})["Success"] == true {
		return nil
	} else {
		return errors.New(result.(map[string]interface{})["Error"].(map[string]interface{})["Message"].(string))
	}
}

func (r repo) Login(data model.AGCLogin) error {

	url := fmt.Sprintf("%s/credit-auth/login", os.Getenv("AGENT_API"))
	result, err := helper.Post(url, data)
	if err != nil {
		return err
	}

	log.Println(result)

	return nil
}

func (r repo) ChangePassword(data model.AGCChangePassword) error {

	url := fmt.Sprintf("%s/credit-auth/changepassword", os.Getenv("AGENT_API"))
	result, err := helper.Post(url, data)
	if err != nil {
		return err
	}

	log.Println(result, result.(map[string]interface{})["Success"])

	return nil
}

func (r repo) Deposit(data model.AGCDeposit) error {

	url := fmt.Sprintf("%s/credit-transfer/deposit", os.Getenv("AGENT_API"))
	result, err := helper.Post(url, data)
	if err != nil {
		return err
	}

	log.Println(result, result.(map[string]interface{})["Success"])

	return nil
}

func (r repo) Withdraw(data model.AGCWithdraw) error {

	url := fmt.Sprintf("%s/credit-transfer/withdraw", os.Getenv("AGENT_API"))
	result, err := helper.Post(url, data)
	if err != nil {
		return err
	}

	log.Println(result, result.(map[string]interface{})["Success"])

	return nil
}

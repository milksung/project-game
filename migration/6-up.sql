
CREATE Table 
    Botaccount_logs (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        external_id BIGINT NOT NULL,
        client_name VARCHAR(255) NOT NULL,
        log_type VARCHAR(255) NOT NULL,
        message TEXT NOT NULL,
        external_create_date VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Botaccount_logs`
    ADD INDEX `idx_external_id` (`external_id`);

CREATE Table 
    Botaccount_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        external_id BIGINT NOT NULL,
        bank_account_id BIGINT NOT NULL,
        bank_code VARCHAR(255) NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        date_time DATETIME NOT NULL,
        raw_date_time DATETIME NOT NULL,
        info VARCHAR(255) NOT NULL,
        channel_code VARCHAR(255) NOT NULL,
        channel_description VARCHAR(255) NOT NULL,
        txn_code VARCHAR(255) NOT NULL,
        txn_description VARCHAR(255) NOT NULL,
        checksum VARCHAR(255) NOT NULL,
        is_read BOOLEAN NOT NULL,
        external_create_date VARCHAR(255) NOT NULL,
        extermal_update_date VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Botaccount_statements`
    ADD INDEX `idx_external_id` (`external_id`),
    ADD INDEX `idx_bank_account_id` (`bank_account_id`);

ALTER TABLE `Bank_accounts`
	ADD COLUMN `external_id` BIGINT NULL AFTER `pin_code`;

ALTER TABLE `Bank_accounts`
    ADD INDEX `idx_external_id` (`external_id`);

ALTER TABLE `Bank_statements`
	ADD COLUMN `external_id` BIGINT(19) NOT NULL AFTER `account_id`,
	ADD COLUMN `from_bank_id` BIGINT(19) NULL AFTER `detail`,
	ADD COLUMN `from_account_number` VARCHAR(255) NULL AFTER `from_bank_id`;

ALTER TABLE `Bank_statements`
    ADD UNIQUE `uni_external_id` (`external_id`);

CREATE TABLE 
    `Botaccount_config` (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        config_key VARCHAR(255) NOT NULL,
        config_val VARCHAR(255) NOT NULL,
        deleted_at DATETIME NULL
    );

ALTER TABLE `Botaccount_config`
    ADD INDEX `idx_config_key` (`config_key`);


INSERT INTO `Botaccount_config` (`config_key`, `config_val`) VALUES
	('allow_create_external_account', '_all'),
	('allow_create_external_account', '_list'),
	('allow_external_account_number', 'set to list and set account number'),
    ('allow_withdraw_from_account', '_all');

ALTER TABLE `Bank_account_types`
	ADD COLUMN `allow_deposit` TINYINT NOT NULL DEFAULT 0 AFTER `limit_flag`,
	ADD COLUMN `allow_withdraw` TINYINT NOT NULL DEFAULT 0 AFTER `allow_deposit`;

UPDATE `Bank_account_types` SET `allow_deposit`=1, `allow_withdraw`=0 WHERE `limit_flag`='00001000';
UPDATE `Bank_account_types` SET `allow_deposit`=0, `allow_withdraw`=1 WHERE `limit_flag`='00000100';
UPDATE `Bank_account_types` SET `allow_deposit`=1, `allow_withdraw`=1 WHERE `limit_flag`='00001100';

ALTER TABLE `Bank_confirm_transactions`
	ADD COLUMN `credit_amount` DECIMAL(14,2) NULL DEFAULT NULL AFTER `slip_url`,
	CHANGE COLUMN `bonus_amount` `bonus_amount` DECIMAL(14,2) NOT NULL DEFAULT 0 AFTER `credit_amount`;

ALTER TABLE `Bank_confirm_transactions`
	ADD COLUMN `bank_charge_amount` DECIMAL(14,2) NOT NULL DEFAULT '0.00' AFTER `bonus_amount`;

CREATE Table 
    Bank_confirm_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        statement_id BIGINT NOT NULL,
        action_type VARCHAR(255) NOT NULL,
        user_id BIGINT NULL,
        account_id BIGINT NOT NULL,
        json_before TEXT NOT NULL,
        confirmed_at DATETIME NULL,
        confirmed_by_user_id BIGINT NULL,
        confirmed_by_username VARCHAR(255) NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_confirm_statements`
	ADD UNIQUE INDEX `uni_statement_id` (`statement_id`),
    ADD INDEX `idx_account_id` (`account_id`);

CREATE Table
    Bank_account_priorities (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        condition_type VARCHAR(255) NOT NULL DEFAULT 'or',
        min_deposit_count INT NOT NULL DEFAULT 0,
        min_deposit_total DECIMAL(14,2) NOT NULL DEFAULT 0,
        created_at DATETIME NOT NULL DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

INSERT INTO `Bank_account_priorities` (`name`, `condition_type`, `min_deposit_count`, `min_deposit_total`) VALUES
    ('ระดับ NEW ทั่วไป', 'or', 0, 0),
    ('ระดับ Gold ฝากมากกว่า 10 ครั้ง หรือ ฝากสะสมมากกว่า 10,000 บาท', 'or', 10, 10000),
    ('ระดับ Platinum ฝากมากกว่า 20 ครั้ง หรือ ฝากสะสมมากกว่า 100,000 บาท', 'or', 20, 100000),
    ('ระดับ VIP ฝากมากกว่า 30 ครั้ง หรือ ฝากสะสมมากกว่า 500,000 บาท', 'or', 30, 500000);

ALTER TABLE `Bank_accounts`
	ADD COLUMN `account_priority_id` BIGINT NULL AFTER `account_priority`,
	ADD COLUMN `auto_withdraw_credit_flag` VARCHAR(255) NOT NULL AFTER `auto_withdraw_flag`,
	ADD COLUMN `auto_withdraw_confirm_flag` VARCHAR(255) NOT NULL AFTER `auto_withdraw_credit_flag`;

ALTER TABLE `Bank_accounts`
	ADD COLUMN `is_main_withdraw` TINYINT NOT NULL DEFAULT 0 AFTER `auto_credit_flag`;

ALTER TABLE `Bank_transactions`
	CHANGE COLUMN `transfer_at` `transfer_at` DATETIME NULL AFTER `bank_charge_amount`;

ALTER TABLE `Bank_confirm_transactions`
	ADD COLUMN `action_key` VARCHAR(255) NOT NULL AFTER `id`;

UPDATE `Bank_confirm_transactions` SET `action_key`=`id` WHERE `action_key`='';

ALTER TABLE `Bank_confirm_transactions`
    ADD UNIQUE `uni_action_key` (`action_key`);

ALTER TABLE `Bank_confirm_transactions`
	DROP INDEX `uni_transaction_id`,
	ADD INDEX `uni_transaction_id` (`transaction_id`);

ALTER TABLE `Bank_accounts`
	CHANGE COLUMN `account_priority` `account_priority` VARCHAR(255) NOT NULL DEFAULT '' COLLATE 'utf8mb4_0900_ai_ci' AFTER `account_balance`;

CREATE Table 
    User_statement_types (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        code VARCHAR(255) NOT NULL,
        name VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

INSERT INTO `User_statement_types` (`code`, `name`) VALUES
    ('deposit', 'ฝากเงิน'),
    ('withdraw', 'ถอนเงิน'),
    ('bonus', 'ได้รับโบนัส'),
    ('getcreditback', 'ดึงเครดิตคืน');

CREATE Table 
    User_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        user_id BIGINT NULL,
        statement_type_id BIGINT NULL,
        Transfer_at DATETIME NULL,
        info VARCHAR(255) NOT NULL,
        before_balance DECIMAL(14,2) NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        after_balance DECIMAL(14,2) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `User_statements`
	ADD INDEX `idx_user_id` (`user_id`),
    ADD INDEX `idx_statement_type_id` (`statement_type_id`);


CREATE Table
    Banks (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        code VARCHAR(255) NOT NULL,
        icon_url VARCHAR(255) NOT NULL,
        type_flag VARCHAR(8) NOT NULL DEFAULT '00000000',
        created_at DATETIME NOT NULL DEFAULT NOW()
    );

ALTER TABLE `Banks`
	ADD UNIQUE INDEX `uni_code` (`code`);

INSERT INTO `Banks` (`name`, `code`, `icon_url`, `type_flag`) VALUES
	('กสิกรไทย', 'kbank', '', '00001111'),
	('ไทยพาณิชย์', 'scb', '', '00001111'),
	('กรุงเทพ', 'bbl', '', '00000011'),
	('กรุงศรีอยุธยา', 'bay', '', '00000011'),
	('กรุงไทย', 'ktb', '', '00000011'),
	('ทีเอ็มบีธนชาต', 'ttb', '', '00000011'),
	('ออมสิน', 'gsb', '', '00000011'),
	('ธกส', 'baac', '', '00000011'),
	('เกียรตินาคิน', 'kk', '', '00000011'),
	('อาคารสงเคราะห์', 'ghb', '', '00000011'),
	('ยูโอบี', 'uob', '', '00000011'),
	('แลนด์ แอนด์ เฮ้าส์', 'lh', '', '00000011'),
	('ซีไอเอ็มบี', 'cimb', '', '00000011'),
	('เอชเอสบีซี', 'hsbc', '', '00000011'),
	('ไอซีบีซี', 'icbc', '', '00000011');

CREATE Table
    Bank_account_types (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        limit_flag VARCHAR(8) NOT NULL DEFAULT '00000000',
        created_at DATETIME NOT NULL DEFAULT NOW()
    );

INSERT INTO `Bank_account_types` (`name`, `limit_flag`) VALUES
	('เฉพาะฝาก', '00001000'),
	('เฉพาะถอน', '00000100'),
	('ฝาก-ถอน', '00001100'),
	('พักเงิน', '00000010');

CREATE Table
    Bank_accounts (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        bank_id BIGINT NOT NULL,
        account_type_id BIGINT NOT NULL,
        account_name VARCHAR(255) NOT NULL,
        account_number VARCHAR(255) NOT NULL,
        account_balance DECIMAL(14,2) NOT NULL DEFAULT 0,
        account_priority VARCHAR(255) NOT NULL,
        account_status VARCHAR(255) NOT NULL,
        device_uid VARCHAR(255) NOT NULL,
        pin_code VARCHAR(255) NOT NULL,
        connection_status VARCHAR(255) NOT NULL,
        last_conn_update_at DATETIME NULL,
        auto_credit_flag VARCHAR(255) NOT NULL,
        auto_withdraw_flag VARCHAR(255) NOT NULL,
        auto_withdraw_max_amount VARCHAR(255) NOT NULL,
        auto_transfer_max_amount VARCHAR(255) NOT NULL,
        qr_wallet_status VARCHAR(255) NOT NULL,
        created_at DATETIME NOT NULL DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_accounts`
	ADD UNIQUE INDEX `uni_account_number` (`account_number`),
    ADD INDEX `idx_bank_id` (`bank_id`),
    ADD INDEX `idx_account_type_id` (`account_type_id`);

CREATE Table
    Bank_account_transactions (
        id INT PRIMARY KEY AUTO_INCREMENT,
        account_id BIGINT NOT NULL,
        description VARCHAR(255) NOT NULL,
        transfer_type VARCHAR(255) NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        transfer_at DATETIME NOT NULL,
        created_by_username VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_account_transactions`
    ADD INDEX `idx_account_id` (`account_id`);

CREATE Table
    Bank_account_transfers (
        id INT PRIMARY KEY AUTO_INCREMENT,
        from_account_id BIGINT NOT NULL,
        from_bank_id BIGINT NOT NULL,
        from_account_name VARCHAR(255) NOT NULL,
        from_account_number VARCHAR(255) NOT NULL,
        to_account_id BIGINT NOT NULL,
        to_bank_id BIGINT NOT NULL,
        to_account_name VARCHAR(255) NOT NULL,
        to_account_number VARCHAR(255) NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        transfer_at DATETIME NOT NULL,
        created_by_username VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
        confirmed_at DATETIME NULL,
        confirmed_by_user_id BIGINT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_account_transfers`
    ADD INDEX `idx_from_account_id` (`from_account_id`),
    ADD INDEX `idx_to_account_id` (`to_account_id`);

CREATE Table
    Webhook_logs (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        json_request TEXT NOT NULL,
        json_payload TEXT NOT NULL,
        log_type VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

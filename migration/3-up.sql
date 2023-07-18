CREATE Table
    Bank_statements (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        account_id BIGINT NOT NULL,
        amount DECIMAL(14,2) NOT NULL,
        detail VARCHAR(255) NOT NULL,
        statement_type VARCHAR(255) NOT NULL,
        transfer_at DATETIME NOT NULL,
        status VARCHAR(255) NOT NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_statements`
    ADD INDEX `idx_account_id` (`account_id`);

CREATE Table
    Bank_transactions (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        member_code VARCHAR(255) NOT NULL,
        user_id BIGINT NOT NULL,
        transfer_type VARCHAR(255) NOT NULL,
        promotion_id BIGINT NULL,
        from_account_id BIGINT NULL,
        from_bank_id BIGINT NULL,
        from_account_name VARCHAR(255) NULL,
        from_account_number VARCHAR(255) NULL,
        to_account_id BIGINT NULL,
        to_bank_id BIGINT NULL,
        to_account_name VARCHAR(255) NULL,
        to_account_number VARCHAR(255) NULL,
        credit_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        paid_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        over_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        deposit_channel VARCHAR(255) NULL,
        bonus_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        bonus_reason VARCHAR(255) NULL,
        before_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        after_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        bank_charge_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        transfer_at DATETIME NULL DEFAULT NULL,
        created_by_user_id BIGINT NOT NULL,
        created_by_username VARCHAR(255) NOT NULL,
        cancel_remark VARCHAR(255) NULL DEFAULT NULL,
        canceled_at DATETIME NULL DEFAULT NULL,
        canceled_by_user_id BIGINT NULL DEFAULT NULL,
        canceled_by_username VARCHAR(255) NULL DEFAULT NULL,
        confirmed_at DATETIME NULL DEFAULT NULL,
        confirmed_by_user_id BIGINT NULL DEFAULT NULL,
        confirmed_by_username VARCHAR(255) NULL DEFAULT NULL,
        removed_at DATETIME NULL DEFAULT NULL,
        removed_by_user_id BIGINT NULL DEFAULT NULL,
        removed_by_username VARCHAR(255) NULL DEFAULT NULL,
        status VARCHAR(255) NOT NULL,
        status_detail VARCHAR(255) NULL,
        is_auto_credit TINYINT NOT NULL DEFAULT 0,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_transactions`
    ADD INDEX `idx_user_id` (`user_id`),
    ADD INDEX `idx_from_account_id` (`from_account_id`),
    ADD INDEX `idx_to_account_id` (`to_account_id`);

CREATE Table 
    Bank_confirm_transactions (
        id BIGINT PRIMARY KEY AUTO_INCREMENT,
        transaction_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        transfer_type VARCHAR(255) NOT NULL,
        from_account_id BIGINT NULL,
        to_account_id BIGINT NULL,
        json_before TEXT NULL,
        transfer_at DATETIME NULL DEFAULT NULL,
        slip_url VARCHAR(255) NULL,
        bonus_amount DECIMAL(14,2) NOT NULL DEFAULT 0,
        confirmed_at DATETIME NULL DEFAULT NULL,
        confirmed_by_user_id BIGINT NULL,
        confirmed_by_username VARCHAR(255) NULL,
        created_at DATETIME DEFAULT NOW(),
        updated_at DATETIME NULL ON UPDATE NOW(),
        deleted_at DATETIME NULL
    );

ALTER TABLE `Bank_confirm_transactions`
	ADD UNIQUE INDEX `uni_transaction_id` (`transaction_id`),
    ADD INDEX `idx_user_id` (`user_id`);
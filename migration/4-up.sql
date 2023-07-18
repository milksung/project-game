CREATE Table
        Setting_web (
                id BIGINT PRIMARY KEY AUTO_INCREMENT,
                logo VARCHAR(255) NOT NULL,
                backgrond_color VARCHAR(8) NOT NULL,
                user_auto VARCHAR(2) NOT NULL,
                otp_register INT NOT NULL,
                tran_withdraw VARCHAR(255) NOT NULL,
                register INT NOT NULL,
                deposit_first INT NOT NULL,
                deposit_next INT NOT NULL,
                withdraw INT NOT NULL,
                line VARCHAR(255) NULL,
                url VARCHAR(255) NULL,
                opt INT NOT NULL,
                created_at DATETIME,
                updated_at DATETIME
        );
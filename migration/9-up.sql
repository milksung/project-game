ALTER TABLE `Users`
CHANGE COLUMN `member_code` `member_code` VARCHAR(255) DEFAULT NULL,
CHANGE COLUMN `username` `username` VARCHAR(255) DEFAULT NULL,
CHANGE COLUMN `password` `password` VARCHAR(255) DEFAULT NULL,
ADD COLUMN `verified_at` DATETIME DEFAULT NULL AFTER `created_at`,
ADD COLUMN `bank_id` INT NULL DEFAULT NULL AFTER `bank_account`,
ADD COLUMN `is_reset_password` TINYINT DEFAULT 0 AFTER `ip`,
ADD INDEX `idx_bank_id` (`bank_id` ASC),
ADD INDEX `idx_verified_at` (`verified_at` ASC),
ADD UNIQUE INDEX `uni_phone` (`phone` ASC);

CREATE TABLE `User_otps` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `code` varchar(10),
  `ref` varchar(10),
  `type` ENUM ('REGISTER', 'FORGET'),
  `user_id` bigint,
  `created_at` datetime DEFAULT (now()),
  `verified_at` datetime,
  `expired_at` datetime
);

CREATE INDEX `idx_code` ON `User_otps` (`code`);
CREATE INDEX `idx_ref` ON `User_otps` (`ref`);
CREATE INDEX `idx_verified_at` ON `User_otps` (`verified_at`);
CREATE INDEX `idx_created_at` ON `User_otps` (`created_at`);
CREATE INDEX `idx_expired_at` ON `User_otps` (`expired_at`);
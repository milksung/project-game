ALTER TABLE `Users` 
ADD COLUMN `partner` VARCHAR(20) DEFAULT NULL AFTER `id`,
ADD COLUMN `promotion` VARCHAR(20) DEFAULT 'ไม่รับโปรโมชั่น' AFTER `phone`,
ADD COLUMN `bankname` VARCHAR(50) DEFAULT NULL AFTER `fullname`,
ADD COLUMN `bank_account` VARCHAR(15) DEFAULT NULL AFTER `bankname`,
ADD COLUMN `channel` VARCHAR(20) DEFAULT NULL AFTER `bank_account`,
ADD COLUMN `true_wallet` VARCHAR(20) DEFAULT NULL AFTER `channel`,
ADD COLUMN `contact` VARCHAR(255) DEFAULT NULL AFTER `true_wallet`,
ADD COLUMN `note` VARCHAR(255) DEFAULT NULL AFTER `contact`,
ADD COLUMN `course` VARCHAR(50) DEFAULT NULL AFTER `note`,
ADD COLUMN `turnover_limit` INT(11) DEFAULT 0 AFTER `credit`,
ADD COLUMN `ip_registered` VARCHAR(20) DEFAULT NULL AFTER `ip`,
CHANGE COLUMN `credit` `credit` DECIMAL(14,2) NULL DEFAULT 0.00,
CHANGE COLUMN `ip` `ip` VARCHAR(20) NULL DEFAULT NULL,
CHANGE COLUMN `updated_at` `updated_at` DATETIME NULL ON UPDATE NOW();

CREATE TABLE `Scammers` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `fullname` varchar(255),
  `firstname` varchar(255) DEFAULT NULL,
  `lastname` varchar(255) DEFAULT NULL,
  `bankname` varchar(50) DEFAULT NULL,
  `bank_account` varchar(15) DEFAULT NULL,
  `phone` varchar(12) DEFAULT NULL,
  `reason` varchar(255) DEFAULT NULL,
  `created_at` datetime DEFAULT (now())
);
CREATE INDEX `idx_fullname` ON `Scammers` (`fullname`);
CREATE INDEX `idx_bankname` ON `Scammers` (`bankname`);
CREATE INDEX `idx_phone` ON `Scammers` (`phone`);
CREATE INDEX `idx_created_at` ON `Scammers` (`created_at`);

CREATE TABLE `Recommend_channels` (
  `id` bigint PRIMARY KEY AUTO_INCREMENT,
  `title` varchar(255) DEFAULT 255,
  `status` ENUM ('ACTIVE', 'DEACTIVE') DEFAULT 'ACTIVE',
  `url` varchar(255) DEFAULT 255,
  `created_at` datetime DEFAULT (now()),
  `updated_at` datetime DEFAULT NULL ON UPDATE NOW()
);
ALTER TABLE `Users`
DROP COLUMN `verified_at`,
DROP COLUMN `is_reset_password`,
DROP COLUMN `bank_id`,
DROP INDEX `idx_verified_at`,
DROP INDEX `idx_bank_id`,
DROP INDEX `uni_phone`;

DROP TABLE IF EXISTS `User_otps`;
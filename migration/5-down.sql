ALTER TABLE `Users`
DROP COLUMN `partner`,
DROP COLUMN `promotion`,
DROP COLUMN `bankname`,
DROP COLUMN `bank_account`,
DROP COLUMN `channel`,
DROP COLUMN `true_wallet`,
DROP COLUMN `contact`,
DROP COLUMN `note`,
DROP COLUMN `course`,
DROP COLUMN `turnover_limit`,
DROP COLUMN `ip_registered`;

DROP TABLE IF EXISTS `Scammers`;
DROP TABLE IF EXISTS `Recommend_channels`;
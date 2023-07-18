ALTER TABLE `Setting_web`,
DROP COLUMN `auto_withdraw`;

DROP TABLE IF EXISTS `Line_notify`;
DROP TABLE IF EXISTS `Type_notify`;
DROP TABLE IF EXISTS `Line_notifygame`;

DELETE FROM `Permissions` WHERE `permission_key` = 'deposit_list';
DELETE FROM `Permissions` WHERE `permission_key` = 'withdraw_list';

ALTER TABLE `Permissions`;
DROP INDEX `idx_position`;
DROP COLUMN `position`;

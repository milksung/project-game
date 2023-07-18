ALTER TABLE `Permissions` 
DROP INDEX `uni_permission_key`,
DROP COLUMN `main`;

ALTER TABLE `Admin_permissions` 
DROP COLUMN `is_read`,
DROP COLUMN `is_write`;

ALTER TABLE `Admin_group_permissions` 
DROP COLUMN `is_read`,
DROP COLUMN `is_write`;

ALTER TABLE `Users` 
DROP COLUMN `bank_code`;
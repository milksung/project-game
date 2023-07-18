ALTER TABLE `Setting_web`
ADD COLUMN `auto_withdraw` VARCHAR (20) DEFAULT NULL AFTER `otp_register`;

CREATE TABLE `Line_notify`
    ( `id` bigint PRIMARY KEY AUTO_INCREMENT, 
    `start_credit` DECIMAL (14, 2) NULL DEFAULT 0.00,
    `token` varchar (255) DEFAULT NULL,
    `notify_id` bigint (20), 
    status ENUM('ACTIVE', 'DEACTIVE') DEFAULT 'ACTIVE',
    `created_at` DATETIME DEFAULT current_timestamp(), 
    `updated_at` DATETIME DEFAULT NULL );

ALTER TABLE `Line_notify`
ADD INDEX `idx_notify_id` (`notify_id`);

CREATE TABLE `Type_notify`
    ( `id` bigint PRIMARY KEY AUTO_INCREMENT, 
    `name` varchar(255) DEFAULT NULL,
    `status` ENUM('ACTIVE', 'DEACTIVE') DEFAULT 'ACTIVE',
    `created_at` DATETIME DEFAULT current_timestamp(), 
    `updated_at` DATETIME DEFAULT NULL );
    
INSERT INTO `Type_notify` (`name`)
VALUES
    ( 'แจ้งเตือนเมื่อมีการสมัครสมาชิก') ,
    ('แจ้งเตือนเมื่อแอดมินล็อกอิน'),
    ('แจ้งเตือนเมื่อแอดมินล็อกเอ้าท์') , 
    ('แจ้งเตือนสรุปยอดฝากรายวัน') , 
    ('แจ้งเตือนสรุปยอดถอนรายวัน'), 
    ('แจ้งเตือนเมื่อมีการฝาก (ก่อนปรับเครดิต)'), 
    ('แจ้งเตือนเมื่อมีการฝาก (หลังปรับเครดิต)'), 
    ('แจ้งเตือนเมื่อมีการถอน (ก่อนปรับเครดิต)'),
    ('แจ้งเตือนเมื่อมีการถอน (หลังปรับเครดิต)') , 
    ('แจ้งเตือนเมื่อมีการถอน (รอโอนเงิน)');


CREATE TABLE `Line_notifygame`
    ( `id` bigint PRIMARY KEY AUTO_INCREMENT, 
    `token` varchar (255) DEFAULT NULL,
    `tag` varchar(255)  DEFAULT NULL,
    `typenotify_id` bigint DEFAULT NULL,
    `status` ENUM('ACTIVE', 'DEACTIVE') DEFAULT 'DEACTIVE',
    `created_at` DATETIME DEFAULT current_timestamp(), 
    `updated_at` DATETIME DEFAULT NULL );


UPDATE `Permissions` SET `permission_key` = 'marketing_manage' WHERE (`id` = '14');
UPDATE `Permissions` SET `permission_key` = 'marketing_manage_link' WHERE (`id` = '15');
UPDATE `Permissions` SET `permission_key` = 'marketing_manage_partner' WHERE (`id` = '16');
UPDATE `Permissions` SET `permission_key` = 'marketing_report' WHERE (`id` = '36');
UPDATE `Permissions` SET `permission_key` = 'marketing_report_return' WHERE (`id` = '37');
UPDATE `Permissions` SET `permission_key` = 'marketing_report_link' WHERE (`id` = '38');
UPDATE `Permissions` SET `permission_key` = 'marketing_report_partner' WHERE (`id` = '39');
UPDATE `Permissions` SET `permission_key` = 'activity_manage' WHERE (`id` = '17');
UPDATE `Permissions` SET `permission_key` = 'activity_manage_return' WHERE (`id` = '18');
UPDATE `Permissions` SET `permission_key` = 'activity_manage_lucky' WHERE (`id` = '19');
UPDATE `Permissions` SET `permission_key` = 'activity_manage_checkin' WHERE (`id` = '20');
UPDATE `Permissions` SET `permission_key` = 'activity_manage_coupon' WHERE (`id` = '21');
UPDATE `Permissions` SET `permission_key` = 'activity_report' WHERE (`id` = '41');
UPDATE `Permissions` SET `permission_key` = 'activity_report_lucky' WHERE (`id` = '42');
UPDATE `Permissions` SET `permission_key` = 'winlose_report' WHERE (`id` = '40');

DELETE FROM `Permissions` WHERE (`id` = '43');
DELETE FROM `Permissions` WHERE (`id` = '45');
DELETE FROM `Permissions` WHERE (`id` = '46');
DELETE FROM `Permissions` WHERE (`id` = '47');

ALTER TABLE `Permissions` 
ADD COLUMN `position` INT NULL DEFAULT NULL AFTER `name`,
ADD INDEX `idx_position` (`position` ASC);

INSERT INTO `Permissions` (`permission_key`, `main`, `name`, `position`) VALUES ('deposit_list', '1', 'รายการฝาก', 43);
INSERT INTO `Permissions` (`permission_key`, `main`, `name`, `position`) VALUES ('withdraw_list', '1', 'รายการถอน', 44);

UPDATE `Permissions` SET `position` = '0', `name` = 'คู่มือการใช้งาน' WHERE (`id` = '1');
UPDATE `Permissions` SET `position` = '1', `name` = 'ผู้ดูแลระบบ' WHERE (`id` = '2');
UPDATE `Permissions` SET `position` = '2' WHERE (`id` = '3');
UPDATE `Permissions` SET `position` = '3' WHERE (`id` = '4');
UPDATE `Permissions` SET `position` = '4' WHERE (`id` = '5');
UPDATE `Permissions` SET `position` = '5' WHERE (`id` = '6');
UPDATE `Permissions` SET `position` = '6' WHERE (`id` = '7');
UPDATE `Permissions` SET `position` = '7' WHERE (`id` = '8');
UPDATE `Permissions` SET `position` = '8' WHERE (`id` = '9');
UPDATE `Permissions` SET `position` = '9' WHERE (`id` = '10');
UPDATE `Permissions` SET `position` = '10' WHERE (`id` = '11');
UPDATE `Permissions` SET `position` = '11' WHERE (`id` = '12');
UPDATE `Permissions` SET `position` = '12' WHERE (`id` = '13');
UPDATE `Permissions` SET `position` = '13' WHERE (`id` = '14');
UPDATE `Permissions` SET `position` = '14' WHERE (`id` = '15');
UPDATE `Permissions` SET `position` = '15' WHERE (`id` = '16');
UPDATE `Permissions` SET `position` = '16' WHERE (`id` = '17');
UPDATE `Permissions` SET `position` = '17' WHERE (`id` = '18');
UPDATE `Permissions` SET `position` = '18' WHERE (`id` = '19');
UPDATE `Permissions` SET `position` = '19' WHERE (`id` = '20');
UPDATE `Permissions` SET `position` = '20' WHERE (`id` = '21');
UPDATE `Permissions` SET `position` = '21' WHERE (`id` = '22');
UPDATE `Permissions` SET `position` = '22' WHERE (`id` = '23');
UPDATE `Permissions` SET `position` = '23' WHERE (`id` = '24');
UPDATE `Permissions` SET `position` = '24' WHERE (`id` = '25');
UPDATE `Permissions` SET `position` = '25' WHERE (`id` = '26');
UPDATE `Permissions` SET `position` = '26' WHERE (`id` = '27');
UPDATE `Permissions` SET `position` = '27' WHERE (`id` = '28');
UPDATE `Permissions` SET `position` = '28' WHERE (`id` = '29');
UPDATE `Permissions` SET `position` = '29' WHERE (`id` = '30');
UPDATE `Permissions` SET `position` = '30' WHERE (`id` = '31');
UPDATE `Permissions` SET `position` = '31' WHERE (`id` = '32');
UPDATE `Permissions` SET `position` = '32' WHERE (`id` = '33');
UPDATE `Permissions` SET `position` = '33' WHERE (`id` = '34');
UPDATE `Permissions` SET `position` = '34' WHERE (`id` = '35');
UPDATE `Permissions` SET `position` = '35' WHERE (`id` = '36');
UPDATE `Permissions` SET `position` = '36' WHERE (`id` = '37');
UPDATE `Permissions` SET `position` = '37' WHERE (`id` = '38');
UPDATE `Permissions` SET `position` = '38' WHERE (`id` = '39');
UPDATE `Permissions` SET `position` = '39' WHERE (`id` = '40');
UPDATE `Permissions` SET `position` = '40' WHERE (`id` = '41');
UPDATE `Permissions` SET `position` = '41' WHERE (`id` = '42');
UPDATE `Permissions` SET `position` = '42' WHERE (`id` = '44');
UPDATE `Permissions` SET `position` = '45' WHERE (`id` = '48');
UPDATE `Permissions` SET `position` = '46' WHERE (`id` = '49');
UPDATE `Permissions` SET `position` = '47' WHERE (`id` = '50');
UPDATE `Permissions` SET `position` = '48' WHERE (`id` = '51');
UPDATE `Permissions` SET `position` = '49' WHERE (`id` = '52');
UPDATE `Permissions` SET `position` = '50' WHERE (`id` = '53');
UPDATE `Permissions` SET `position` = '51' WHERE (`id` = '54');

ALTER TABLE `Line_notifygame`
ADD COLUMN `admin_id` bigint DEFAULT 0 AFTER `token`;

ALTER TABLE `Line_notifygame`
	ADD UNIQUE INDEX `uni_admin_id` (`admin_id`);

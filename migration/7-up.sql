ALTER TABLE `Permissions` 
ADD UNIQUE INDEX `uni_permission_key` (`permission_key`),
ADD COLUMN `main` TINYINT NULL DEFAULT 0 AFTER `permission_key`;

INSERT INTO `Permissions` (`name`, `permission_key`, `main`) VALUES
('คู่มือ', 'guide', 1),
('ผู้ดูแล', 'admin', 1),
('จัดการผู้ใช้งาน', 'admin_manager', 0),
('กลุ่มผู้ใช้งาน', 'admin_group', 0),
('สรุปภาพรวม', 'summary', 1),
('จัดการเว็บเอเย่น', 'agent', 1),
('รายการเว็บเอเย่น', 'agent_list', 0),
('รายงานเพิ่ม-ลด เครดิต', 'agent_credit', 0),
('จัดการธนาคาร', 'bank', 1),
('รายการธนาคาร', 'bank_list', 0),
('รายงานธุรกรรมเงินสด', 'bank_transaction', 0),
('รายการเดินบัญชีธนาคาร', 'bank_account', 0),
('จัดการโปรโมชั่น', 'promotion', 1),
('จัดการการตลาด', 'marketing', 1),
('รายการลิ้งรับทรัพย์', 'marketing_link', 0),
('รายการพันธมิตร', 'marketing_partner', 0),
('จัดการกิจกรรม', 'activity', 1),
('คืนยอดเสีย', 'activity_return', 0),
('กงล้อนำโชค', 'activity_lucky', 0),
('เช็คอินรายวัน', 'activity_checkin', 0),
('คูปองเงินสด', 'activity_coupon', 0),
('จัดการสมาชิกเว็บ', 'member', 1),
('รายการสมาชิกเว็บ', 'member_list', 0),
('ประวัติฝาก-ถอนสมาชิก', 'member_transaction', 0),
('ตั้งค่าช่องทางที่รู้จัก', 'member_channel', 0),
('ประวัติการแก้ไขข้อมูล', 'member_history', 0),
('รายการมิจฉาชีพ', 'member_misconduct', 0),
('รายงาน', 'report', 1),
('ยอดสมาชิกผู้ใช้งาน', 'report_member', 0),
('ยอดฝาก-ถอน', 'report_deposit', 0),
('จำนวนฝาก-ถอนตามเวลา', 'report_deposittime', 0),
('รายงานการแจกโบนัส', 'report_bonus', 0),
('จำนวนสมาชิกนับเวลาบันทึก', 'report_membertime', 0),
('ยอดสมาชิกตามช่องทางที่รู้จัก', 'report_memberchannel', 0), 
('จำนวนบันทึกรายการตามผู้ใช้งาน', 'report_memberuser', 0),
('รายงานการตลาด', 'marketing1_report', 1),
('คืนยอดเสีย', 'marketing1_report_return', 0),
('ลิงค์รับทรัพย์', 'marketing1_report_link', 0),
('พันธมิตร', 'marketing1_report_partner', 0),
('รายงานข้อมูล แพ้-ชนะ', 'report1_winlose', 1),
('รายงานกิจกรรม', 'activity1_report', 1),
('กงล้อนำโชค', 'activity1_report_lucky', 0),
('รายการฝาก-ถอนเสร็จสิ้น', 'deposit_withdraw', 1),
('รายการโอนรอดำเนินการ', 'waiting_transfer', 1),
('บันทึกรายการฝาก-ถอน', 'deposit_withdraw_history', 1),
('อนุมัติฝาก(Auto)', 'auto_deposit', 1),
('อนุมัติถอน(Auto)', 'auto_withdraw', 1),
('ตั้งค่าระบบ', 'setting', 1),
('ข้อมูลเบื้องต้น', 'setting_basic', 0),
('แจ้งเตือนกลุ่ม line', 'setting_line', 0),
('PushMessage line', 'setting_sms', 0),
('แจ้งเตือน Cyber Notify', 'setting_cybernoti', 0),
('สถานะเรื่องแจ้งแก้ไข', 'status_update', 1),
('ใบแจ้งหนี้', 'invoice_notice', 1);

ALTER TABLE `Admin_permissions` 
ADD COLUMN `is_read` TINYINT NULL DEFAULT 0 AFTER `deleted_at`,
ADD COLUMN `is_write` VARCHAR(255) NULL DEFAULT 0 AFTER `is_read`;

ALTER TABLE `Admin_group_permissions` 
ADD COLUMN `is_read` TINYINT NULL DEFAULT 0 AFTER `deleted_at`,
ADD COLUMN `is_write` TINYINT NULL DEFAULT 0 AFTER `is_read`;

ALTER TABLE `Users` 
ADD COLUMN `bank_code` VARCHAR(10) NULL DEFAULT NULL AFTER `bankname`;
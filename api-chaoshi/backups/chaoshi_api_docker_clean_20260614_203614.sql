DROP TABLE IF EXISTS `activities`;
CREATE TABLE `activities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `image` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `link_type` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `link_value` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sort` int unsigned NOT NULL DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_activities_type` (`type`),
  KEY `idx_activities_status` (`status`),
  KEY `idx_activities_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `admin_profiles`;
CREATE TABLE `admin_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `contact_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_admin_profiles_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `admin_profiles` VALUES (1,'Head Office','Admin','13800138000',1,'2026-06-14 18:46:55','2026-06-14 18:46:55');
DROP TABLE IF EXISTS `admin_users`;
CREATE TABLE `admin_users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `admin_profile_id` bigint unsigned NOT NULL,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `role` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'operator',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `last_login_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_admin_users_username` (`username`),
  KEY `idx_admin_users_profile_id` (`admin_profile_id`),
  CONSTRAINT `fk_admin_users_profile` FOREIGN KEY (`admin_profile_id`) REFERENCES `admin_profiles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `admin_users` VALUES (1,1,'admin','$2a$10$FD5KRv2ER/jSpSCkEU/PMeMA4kRjqdv0tN6RlmRCY2ReaNFibmPre','Admin','13800138000','admin',1,'2026-06-14 18:50:17','2026-06-14 18:46:55','2026-06-14 18:50:17');
DROP TABLE IF EXISTS `announcements`;
CREATE TABLE `announcements` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_announcements_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sort` int unsigned NOT NULL DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_categories_merchant_id` (`merchant_id`),
  KEY `idx_categories_sort` (`merchant_id`,`sort`),
  CONSTRAINT `fk_categories_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `categories` VALUES (1,3,'AA',0,1,'2026-06-14 19:06:41','2026-06-14 19:06:41'),(2,3,'BB',0,1,'2026-06-14 19:06:44','2026-06-14 19:06:44');
DROP TABLE IF EXISTS `cloud_printers`;
CREATE TABLE `cloud_printers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `brand` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `device_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `api_key` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `api_url` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `feie_user` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `feie_ukey` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `feie_sn` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `print_types` json DEFAULT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `auto_print` tinyint(1) NOT NULL DEFAULT '0',
  `is_default` tinyint(1) NOT NULL DEFAULT '0',
  `print_count` int NOT NULL DEFAULT '0',
  `last_print_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_cloud_printers_merchant_id` (`merchant_id`),
  CONSTRAINT `fk_cloud_printers_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `coupon_records`;
CREATE TABLE `coupon_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `coupon_id` bigint unsigned NOT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '0',
  `used_at` datetime DEFAULT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_coupon_records_user_id` (`user_id`),
  KEY `idx_coupon_records_coupon_id` (`coupon_id`),
  KEY `idx_coupon_records_order_id` (`order_id`),
  CONSTRAINT `fk_coupon_records_coupon` FOREIGN KEY (`coupon_id`) REFERENCES `coupons` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_coupon_records_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `coupons`;
CREATE TABLE `coupons` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `discount_amount` decimal(10,2) DEFAULT NULL,
  `min_order_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `total_count` int NOT NULL,
  `remaining_count` int NOT NULL,
  `per_user_limit` int NOT NULL DEFAULT '1',
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_coupons_merchant_id` (`merchant_id`),
  KEY `idx_coupons_status` (`status`),
  CONSTRAINT `fk_coupons_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `merchant_delivery_settings`;
CREATE TABLE `merchant_delivery_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL COMMENT '鍟嗗ID(涓€瀵逛竴)',
  `enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '鏄惁鍚敤閰嶉€佽垂瑙勫垯: true=鍚敤 false=鍋滅敤(浠呮帶鍒惰垂鐢ㄨ鍒?閰嶉€佸紑鍏冲彇merchants.takeout_enabled)',
  `base_fee` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '鍩虹閰嶉€佽垂(鍏?',
  `free_delivery_amount` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '婊″噺鍏嶉厤閫佽垂閲戦(0=涓嶅惎鐢?',
  `max_distance` bigint unsigned NOT NULL DEFAULT '10' COMMENT '鏈€澶ч厤閫佽窛绂?km)',
  `distance_rules` json DEFAULT NULL COMMENT '璺濈闃舵璁¤垂瑙勫垯JSON [{min_distance,max_distance,fee}]',
  `created_at` datetime(3) DEFAULT NULL COMMENT '鍒涘缓鏃堕棿',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '鏇存柊鏃堕棿',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchant_delivery_settings_merchant_id` (`merchant_id`),
  UNIQUE KEY `merchant_id` (`merchant_id`),
  UNIQUE KEY `idx_merchant_delivery_settings_merchant_id` (`merchant_id`),
  CONSTRAINT `fk_merchant_delivery_settings_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `merchant_delivery_settings` VALUES (1,1,1,5.00,50.00,10,'[{\"fee\": 0, \"max_distance\": 2, \"min_distance\": 0}, {\"fee\": 3, \"max_distance\": 5, \"min_distance\": 2}, {\"fee\": 6, \"max_distance\": 10, \"min_distance\": 5}]','2026-06-14 18:46:55.000','2026-06-14 18:46:55.000');
DROP TABLE IF EXISTS `merchant_full_reduction_rules`;
CREATE TABLE `merchant_full_reduction_rules` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL COMMENT '鎵€灞炲晢瀹禝D',
  `threshold_amount` decimal(10,2) NOT NULL COMMENT '婊″噺闂ㄦ閲戦(鍏?',
  `discount_amount` decimal(10,2) NOT NULL COMMENT '鍑忓厤閲戦(鍏?',
  `sort` bigint NOT NULL DEFAULT '0' COMMENT '鎺掑簭鍊?,
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '鐘舵€? 1=鍚敤 0=鍋滅敤',
  `created_at` datetime(3) DEFAULT NULL COMMENT '鍒涘缓鏃堕棿',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '鏇存柊鏃堕棿',
  PRIMARY KEY (`id`),
  KEY `idx_full_reduction_merchant_id` (`merchant_id`),
  KEY `idx_merchant_full_reduction_rules_merchant_id` (`merchant_id`),
  CONSTRAINT `fk_full_reduction_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `merchant_pickup_points`;
CREATE TABLE `merchant_pickup_points` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `address` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `lat` decimal(10,6) NOT NULL,
  `lng` decimal(10,6) NOT NULL,
  `is_default` tinyint(1) NOT NULL DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `sort` int unsigned NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_merchant_pickup_points_merchant_id` (`merchant_id`),
  KEY `idx_merchant_pickup_points_default` (`merchant_id`,`is_default`),
  CONSTRAINT `fk_merchant_pickup_points_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `merchant_pickup_points` VALUES (1,1,'盲赂禄氓潞鈥斆モ€奥嵜ヂ徛懊ㄢ€÷β徛惷р€毬?,'氓艗鈥斆ぢ郝ヂ糕€毭ε撀澝┧溌趁ヅ捖好ヂ宦好モ€郝矫仿?8氓聫路盲赂鈧ヂ扁€毭ε撀嵜ヅ犅∶ヂ徛?,39.908823,116.397470,1,1,1,'2026-06-14 18:59:47','2026-06-14 18:59:47'),(2,1,'氓聬沤茅鈥斅绰ε韭睹ㄢ€÷β徛惷р€毬?,'氓艗鈥斆ぢ郝ヂ糕€毭ε撀澝┧溌趁ヅ捖好ヂ宦好モ€郝矫仿?8氓聫路氓聬沤茅鈥斅ヂ柯ヂ忊€撁ヅ捖?,39.909120,116.398020,0,1,2,'2026-06-14 18:59:47','2026-06-14 18:59:47'),(3,2,'盲赂艙茅鈥斅ㄢ€÷β徛惷р€毬?,'茅鈥斅ヂ衡€? 盲赂艙茅鈥斅モ€βッヂ徛Ｃヂ徛趁ぢ韭ヂ柯ヂ忊€撁ε概?,31.230610,121.473700,1,1,1,'2026-06-14 18:59:47','2026-06-14 18:59:47'),(4,2,'氓聛艙猫陆娄氓艙潞猫鈥÷β徛惷р€毬?,'茅鈥斅ヂ衡€? 氓艙掳盲赂鈥姑ヂ伵撁铰γヅ撀?B1 氓驴芦忙聧路氓聫鈥撁绰ヅ捖?,31.229980,121.474250,0,1,2,'2026-06-14 18:59:47','2026-06-14 18:59:47'),(5,3,'AAAA','AAAA',23.129240,113.264360,1,1,1,'2026-06-14 18:59:47','2026-06-14 19:00:14'),(6,3,'BBBB','BBBB',23.128760,113.265040,0,1,2,'2026-06-14 18:59:47','2026-06-14 19:00:19');
DROP TABLE IF EXISTS `merchant_staffs`;
CREATE TABLE `merchant_staffs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL,
  `username` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `push_openid` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `role` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'staff',
  `notify_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `browse_notify_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `last_login_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_merchant_staffs_merchant_username` (`merchant_id`,`username`),
  KEY `idx_merchant_staffs_push_openid` (`push_openid`),
  KEY `idx_merchant_staffs_phone` (`phone`),
  CONSTRAINT `fk_merchant_staffs_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `merchant_staffs` VALUES (1,1,'merchant','$2a$10$mP89UzDWaHy0LVxdDqWhheUJ/UN4tVkArcEhTqW7kqScW7lk.558W','Merchant Admin','13900139000',NULL,'owner',1,1,1,NULL,'2026-06-14 18:46:55','2026-06-14 18:46:55'),(2,2,'13000000001','$2a$10$djcuR/VgdAq9fz2TkFYFSOfGURPftEuRxFJz6CNejQmsDOSsBqtMS','閲戦櫟閰掑簵','13000000001','','owner',1,1,1,NULL,'2026-06-14 18:50:51','2026-06-14 18:50:51'),(3,3,'13000000002','$2a$10$.wkjic7yjV9eZ8J193icWOvcRLhWE67xJjHNaDDAaMCNRXUunT0wW','灏忕澘瓒呭競','13000000002','','owner',1,1,1,NULL,'2026-06-14 18:51:20','2026-06-14 18:51:20');
DROP TABLE IF EXISTS `merchants`;
CREATE TABLE `merchants` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `logo` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cover_image` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_email` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `lat` decimal(10,6) DEFAULT NULL,
  `lng` decimal(10,6) DEFAULT NULL,
  `business_category` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `business_hours` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `announcement` text COLLATE utf8mb4_unicode_ci,
  `min_order_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `takeout_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `dine_in_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `pickup_enabled` tinyint(1) NOT NULL DEFAULT '1',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `rating` decimal(2,1) NOT NULL DEFAULT '5.0',
  `sales_count` int unsigned NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_merchants_status` (`status`),
  KEY `idx_merchants_location` (`lat`,`lng`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `merchants` VALUES (1,'Demo Store','https://example.com/images/merchant_logo_1.jpg',NULL,'Li Si','13900139000','lisi@example.com','88 Jianguo Road, Chaoyang District, Beijing',39.908823,116.407470,'Retail','08:00-22:00','Welcome to our store',20.00,1,1,1,1,5.0,0,'2026-06-14 18:46:55','2026-06-14 18:46:55'),(2,'閲戦櫟閰掑簵','','','閲戦櫟閰掑簵','13000000001','','',0.000000,0.000000,'','','',0.00,1,1,1,1,5.0,0,'2026-06-14 18:50:51','2026-06-14 18:50:51'),(3,'灏忕澘瓒呭競','https://xiaorui.huaiangaoxin.top/uploads/common/2026/06/14/1781434284885163312_f4ae2567-35b0-4428-ba4a-41df544de072_-_鍓湰.jpg','','灏忕澘瓒呭競','13000000002','','',0.000000,0.000000,'','','',0.00,1,1,1,1,5.0,0,'2026-06-14 18:51:20','2026-06-14 18:51:25');
DROP TABLE IF EXISTS `order_items`;
CREATE TABLE `order_items` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint unsigned NOT NULL,
  `merchant_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL,
  `product_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `image` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `price` decimal(10,2) NOT NULL,
  `quantity` int unsigned NOT NULL DEFAULT '1',
  `spec_info` json DEFAULT NULL,
  `subtotal` decimal(10,2) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_items_order_id` (`order_id`),
  KEY `idx_order_items_merchant_id` (`merchant_id`),
  KEY `idx_order_items_product_id` (`product_id`),
  CONSTRAINT `fk_order_items_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_order_items_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_order_items_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `order_items` VALUES (1,1,3,1,'AA','https://xiaorui.huaiangaoxin.top/uploads/common/2026/06/14/1781435220467657977_f4ae2567-35b0-4428-ba4a-41df544de072.jpg',112.00,1,'1',112.00,'2026-06-14 19:07:30');
DROP TABLE IF EXISTS `orders`;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `merchant_id` bigint unsigned NOT NULL,
  `total_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `delivery_fee` decimal(10,2) NOT NULL DEFAULT '0.00',
  `discount_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `pay_amount` decimal(10,2) NOT NULL DEFAULT '0.00',
  `delivery_type` tinyint unsigned NOT NULL DEFAULT '1',
  `delivery_distance` decimal(5,2) DEFAULT NULL,
  `delivery_address` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `contact_phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pickup_point_id` bigint unsigned DEFAULT NULL,
  `pickup_point_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pickup_point_address` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pickup_point_lat` decimal(10,6) DEFAULT NULL,
  `pickup_point_lng` decimal(10,6) DEFAULT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `remark` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `verify_code` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `transaction_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `paid_at` datetime DEFAULT NULL,
  `pay_notify_payload` json DEFAULT NULL,
  `completed_at` datetime DEFAULT NULL,
  `completed_by_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cancelled_at` datetime DEFAULT NULL,
  `refunded_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orders_order_no` (`order_no`),
  KEY `idx_orders_user_id` (`user_id`),
  KEY `idx_orders_merchant_id` (`merchant_id`),
  KEY `idx_orders_status` (`status`),
  KEY `idx_orders_created_at` (`created_at`),
  KEY `idx_orders_paid_at` (`paid_at`),
  KEY `idx_orders_verify_code` (`verify_code`),
  KEY `idx_orders_pickup_point_id` (`pickup_point_id`),
  CONSTRAINT `fk_orders_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE RESTRICT,
  CONSTRAINT `fk_orders_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `orders` VALUES (1,'202606141907293132232',1,3,112.00,0.00,0.00,112.00,3,0.00,'','','',6,'BBBB','BBBB',23.128760,113.265040,4,'','381525','',NULL,NULL,NULL,'','2026-06-14 19:08:04',NULL,'2026-06-14 19:07:30','2026-06-14 19:08:04');
DROP TABLE IF EXISTS `print_logs`;
CREATE TABLE `print_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL,
  `printer_id` bigint unsigned NOT NULL,
  `order_id` bigint unsigned DEFAULT NULL,
  `type` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '0',
  `error_message` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_print_logs_merchant_id` (`merchant_id`),
  KEY `idx_print_logs_printer_id` (`printer_id`),
  KEY `idx_print_logs_order_id` (`order_id`),
  CONSTRAINT `fk_print_logs_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_print_logs_printer` FOREIGN KEY (`printer_id`) REFERENCES `cloud_printers` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `product_specs`;
CREATE TABLE `product_specs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint unsigned NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `options` json DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_product_specs_product_id` (`product_id`),
  CONSTRAINT `fk_product_specs_product` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `product_specs` VALUES (1,1,'1','[{\"name\": \"1\", \"price\": 1}, {\"name\": \"2\", \"price\": 2}]','2026-06-14 19:07:11','2026-06-14 19:07:11');
DROP TABLE IF EXISTS `products`;
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL,
  `category_id` bigint unsigned DEFAULT NULL,
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text COLLATE utf8mb4_unicode_ci,
  `images` json DEFAULT NULL,
  `price` decimal(10,2) NOT NULL,
  `original_price` decimal(10,2) DEFAULT NULL,
  `stock` int unsigned NOT NULL DEFAULT '0',
  `unit` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'item',
  `sales` int unsigned NOT NULL DEFAULT '0',
  `sort` int unsigned NOT NULL DEFAULT '0',
  `status` tinyint unsigned NOT NULL DEFAULT '1',
  `deleted_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_products_merchant_id` (`merchant_id`),
  KEY `idx_products_category_id` (`category_id`),
  KEY `idx_products_status` (`status`),
  KEY `idx_products_sales` (`merchant_id`,`sales`),
  KEY `idx_products_sort` (`merchant_id`,`sort`),
  CONSTRAINT `fk_products_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_products_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `products` VALUES (1,3,1,'AA','AA','[\"https://xiaorui.huaiangaoxin.top/uploads/common/2026/06/14/1781435220467657977_f4ae2567-35b0-4428-ba4a-41df544de072.jpg\"]',111.00,1111.00,110,'浠?,1,0,1,NULL,'2026-06-14 19:07:11','2026-06-14 19:07:30');
DROP TABLE IF EXISTS `refunds`;
CREATE TABLE `refunds` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `order_id` bigint unsigned NOT NULL,
  `refund_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `refund_amount` decimal(10,2) NOT NULL,
  `refund_reason` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` tinyint unsigned NOT NULL DEFAULT '0',
  `refund_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `refunded_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refunds_refund_no` (`refund_no`),
  KEY `idx_refunds_order_id` (`order_id`),
  KEY `idx_refunds_status` (`status`),
  CONSTRAINT `fk_refunds_order` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `user_addresses`;
CREATE TABLE `user_addresses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `province` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `city` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `district` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `lat` decimal(10,6) DEFAULT NULL,
  `lng` decimal(10,6) DEFAULT NULL,
  `is_default` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_addresses_user_id` (`user_id`),
  CONSTRAINT `fk_user_addresses_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

DROP TABLE IF EXISTS `user_behavior_events`;
CREATE TABLE `user_behavior_events` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint unsigned NOT NULL COMMENT '鍟嗗ID',
  `user_id` bigint unsigned NOT NULL COMMENT '鐢ㄦ埛ID',
  `openid` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '寰俊OpenID',
  `event_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '浜嬩欢绫诲瀷: page_view=椤甸潰娴忚 product_view=鍟嗗搧鏌ョ湅 submit_order=鎻愪氦璁㈠崟 pay_success=鏀粯鎴愬姛',
  `page` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '椤甸潰鏍囪瘑(濡俿tore_home/store_product)',
  `product_id` bigint unsigned DEFAULT NULL COMMENT '鍏宠仈鍟嗗搧ID(鍟嗗搧鏌ョ湅浜嬩欢)',
  `order_id` bigint unsigned DEFAULT NULL COMMENT '鍏宠仈璁㈠崟ID(涓嬪崟/鏀粯浜嬩欢)',
  `source` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '浜嬩欢鏉ユ簮: scan=鎵爜 direct=鐩存帴杩涘叆',
  `payload` json DEFAULT NULL COMMENT '浜嬩欢闄勫姞鏁版嵁JSON',
  `created_at` datetime(3) DEFAULT NULL COMMENT '鍒涘缓鏃堕棿',
  PRIMARY KEY (`id`),
  KEY `idx_user_behavior_events_merchant_id` (`merchant_id`),
  KEY `idx_user_behavior_events_user_id` (`user_id`),
  KEY `idx_user_behavior_events_openid` (`openid`),
  KEY `idx_user_behavior_events_event_type` (`event_type`),
  KEY `idx_user_behavior_events_product_id` (`product_id`),
  KEY `idx_user_behavior_events_order_id` (`order_id`),
  KEY `idx_user_behavior_events_created_at` (`created_at`),
  KEY `idx_user_behavior_events_open_id` (`openid`)
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `user_behavior_events` VALUES (1,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 18:49:24.547'),(2,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 18:49:24.558'),(3,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:49:52.472'),(4,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:49:52.488'),(5,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:49:53.613'),(6,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:49:53.625'),(7,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:50:03.225'),(8,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:50:03.238'),(9,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:52:44.607'),(10,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:52:44.622'),(11,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:52:52.690'),(12,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:52:52.707'),(13,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:53:14.527'),(14,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:53:14.540'),(15,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:53:47.080'),(16,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:53:47.095'),(17,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:53:49.041'),(18,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:53:49.055'),(19,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:53:49.475'),(20,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:53:49.500'),(21,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:53:59.571'),(22,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:53:59.595'),(23,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:54:02.637'),(24,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:54:02.650'),(25,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:54:09.590'),(26,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:54:09.603'),(27,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:54:21.236'),(28,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:54:21.246'),(29,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:54:24.912'),(30,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:54:24.924'),(31,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:54:32.908'),(32,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:54:32.920'),(33,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'dd','{\"page\": \"store_home\"}','2026-06-14 18:55:05.211'),(34,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'dd','{\"source\": \"dd\"}','2026-06-14 18:55:05.223'),(35,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'dd','{\"page\": \"store_home\"}','2026-06-14 18:55:14.364'),(36,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'dd','{\"source\": \"dd\"}','2026-06-14 18:55:14.378'),(37,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 18:55:24.215'),(38,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 18:55:24.226'),(39,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 18:55:51.562'),(40,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 18:55:51.576'),(41,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:55:56.770'),(42,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:55:56.782'),(43,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 18:55:57.883'),(44,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 18:55:57.892'),(45,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 19:07:16.889'),(46,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 19:07:16.900'),(47,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_confirm',NULL,NULL,'scan','{\"page\": \"store_confirm\"}','2026-06-14 19:07:23.469'),(48,3,1,'','submit_order','store_confirm',NULL,1,'store','{\"pay_amount\": 112, \"delivery_type\": 3, \"delivery_distance\": 0}','2026-06-14 19:07:29.866'),(49,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','pay_success','store_payment_result',NULL,1,'store','{\"amount\": 112, \"order_id\": 1}','2026-06-14 19:07:29.886'),(50,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 19:07:43.241'),(51,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 19:07:43.253'),(52,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 19:07:45.178'),(53,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 19:07:45.199'),(54,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 19:07:45.984'),(55,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 19:07:46.008'),(56,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 19:07:46.713'),(57,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 19:07:46.724'),(58,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 19:07:55.506'),(59,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 19:07:55.532'),(60,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 19:08:00.756'),(61,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 19:08:00.768'),(62,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 19:08:06.714'),(63,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 19:08:06.724'),(64,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 19:15:17.253'),(65,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 19:15:17.262'),(66,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 19:48:45.413'),(67,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 19:48:45.434'),(68,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'xcx_shell','{\"page\": \"store_home\"}','2026-06-14 19:57:25.889'),(69,3,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'xcx_shell','{\"source\": \"xcx_shell\"}','2026-06-14 19:57:25.905'),(70,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','page_view','store_home',NULL,NULL,'scan','{\"page\": \"store_home\"}','2026-06-14 19:57:26.206'),(71,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','store_visit','store_home',NULL,NULL,'scan','{\"source\": \"scan\"}','2026-06-14 19:57:26.221');
DROP TABLE IF EXISTS `user_visits`;
CREATE TABLE `user_visits` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL COMMENT '鐢ㄦ埛ID',
  `merchant_id` bigint unsigned NOT NULL COMMENT '鍟嗗ID',
  `openid` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '寰俊OpenID',
  `visit_time` datetime(3) DEFAULT NULL COMMENT '璁块棶鏃堕棿',
  `source` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '璁块棶鏉ユ簮: scan=鎵爜 direct=鐩存帴杩涘叆',
  PRIMARY KEY (`id`),
  KEY `idx_user_visits_user_id` (`user_id`),
  KEY `idx_user_visits_merchant_id` (`merchant_id`),
  KEY `idx_user_visits_openid` (`openid`),
  KEY `idx_user_visits_open_id` (`openid`),
  CONSTRAINT `fk_user_visits_merchant` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_visits_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `user_visits` VALUES (1,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:49:24.546','xcx_shell'),(2,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:49:52.472','scan'),(3,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:49:53.612','scan'),(4,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:50:03.224','scan'),(5,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:52:44.606','scan'),(6,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:52:52.689','scan'),(7,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:53:14.526','scan'),(8,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:53:47.079','scan'),(9,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:53:49.040','scan'),(10,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:53:49.475','scan'),(11,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:53:59.570','scan'),(12,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:54:02.636','scan'),(13,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:54:09.590','scan'),(14,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:54:21.235','scan'),(15,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:54:24.911','scan'),(16,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:54:32.907','scan'),(17,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:55:05.210','dd'),(18,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:55:14.364','dd'),(19,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:55:24.214','xcx_shell'),(20,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:55:51.562','xcx_shell'),(21,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:55:56.769','scan'),(22,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 18:55:57.882','scan'),(23,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:07:16.889','scan'),(24,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:07:43.240','scan'),(25,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:07:45.177','scan'),(26,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:07:45.983','scan'),(27,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:07:46.712','xcx_shell'),(28,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:07:55.505','scan'),(29,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:08:00.755','xcx_shell'),(30,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:08:06.713','xcx_shell'),(31,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:15:17.251','xcx_shell'),(32,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:48:45.411','xcx_shell'),(33,1,3,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:57:25.888','xcx_shell'),(34,1,1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','2026-06-14 19:57:26.205','scan');
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `openid` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '寰俊OpenID(鐢ㄦ埛鍞竴鏍囪瘑)',
  `union_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '寰俊UnionID(璺ㄥ皬绋嬪簭鍞竴)',
  `nickname` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '鐢ㄦ埛鏄电О(榛樿寰俊鐢ㄦ埛)',
  `avatar` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '鐢ㄦ埛澶村儚URL',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '鐢ㄦ埛鎵嬫満鍙?,
  `status` tinyint unsigned NOT NULL DEFAULT '1' COMMENT '鐘舵€? 1=姝ｅ父 0=绂佺敤',
  `first_visit_at` datetime(3) DEFAULT NULL COMMENT '棣栨璁块棶鏃堕棿',
  `last_visit_at` datetime(3) DEFAULT NULL COMMENT '鏈€鍚庤闂椂闂?,
  `visit_count` bigint unsigned NOT NULL DEFAULT '1' COMMENT '绱璁块棶娆℃暟',
  `has_ordered` tinyint(1) NOT NULL DEFAULT '0' COMMENT '鏄惁涓嬭繃鍗? true=鏄?false=鍚?,
  `total_orders` bigint unsigned NOT NULL DEFAULT '0' COMMENT '绱璁㈠崟鏁?,
  `total_spent` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '绱娑堣垂閲戦(鍏?',
  `has_paid` tinyint(1) NOT NULL DEFAULT '0' COMMENT '鏄惁瀹屾垚杩囨敮浠? true=鏄?false=鍚?,
  `first_paid_at` datetime(3) DEFAULT NULL COMMENT '棣栨鏀粯鏃堕棿',
  `created_at` datetime(3) DEFAULT NULL COMMENT '鍒涘缓鏃堕棿',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '鏇存柊鏃堕棿',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_users_openid` (`openid`),
  UNIQUE KEY `openid` (`openid`),
  UNIQUE KEY `idx_users_open_id` (`openid`),
  KEY `idx_users_union_id` (`union_id`),
  KEY `idx_users_phone` (`phone`),
  KEY `idx_users_first_visit_at` (`first_visit_at`),
  KEY `idx_users_has_paid` (`has_paid`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `users` VALUES (1,'o3GoB7ii-JMuOprBdYg54Wg5dT6E','','寰俊鐢ㄦ埛','','',1,'2026-06-14 18:49:21.808','2026-06-14 19:57:26.205',39,1,1,112.00,1,'2026-06-14 19:07:29.870','2026-06-14 18:49:21.803','2026-06-14 19:57:26.206');



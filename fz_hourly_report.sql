/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 90300
 Source Host           : localhost:3306
 Source Schema         : release_atd

 Target Server Type    : MySQL
 Target Server Version : 90300
 File Encoding         : 65001

 Date: 12/02/2026 16:41:20
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for fz_hourly_report
-- ----------------------------
DROP TABLE IF EXISTS `fz_hourly_report`;
CREATE TABLE `fz_hourly_report` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `media` varchar(50) NOT NULL COMMENT '媒体简称',
  `media_adv_id` varchar(100) NOT NULL COMMENT '媒体账户ID',
  `media_adv_name` varchar(200) NOT NULL COMMENT '媒体账户名称',
  `report_date` int NOT NULL COMMENT '报表日期，格式：YYYYMMDD',
  `cost` decimal(12,2) DEFAULT '0.00' COMMENT '消耗',
  `convert_dp` bigint DEFAULT '0' COMMENT '拉活数',
  `dp_app_order_nums` bigint DEFAULT '0' COMMENT '订单数',
  `click` bigint DEFAULT '0' COMMENT '点击数',
  `expose` bigint DEFAULT '0' COMMENT '曝光数',
  `convert_dp_price` decimal(12,2) DEFAULT '0.00' COMMENT '拉活成本',
  `dp_app_order_price` decimal(12,2) DEFAULT '0.00' COMMENT '订单成本',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_date_media_adv` (`report_date`,`media`,`media_adv_id`) COMMENT '唯一约束：防止重复数据',
  KEY `idx_report_date` (`report_date`) COMMENT '按日期查询索引',
  KEY `idx_media_adv` (`media`,`media_adv_id`) COMMENT '按媒体账户查询索引',
  KEY `idx_create_time` (`create_time`) COMMENT '创建时间索引'
) ENGINE=InnoDB AUTO_INCREMENT=125 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='飞猪时报数据表';

-- ----------------------------
-- Records of fz_hourly_report
-- ----------------------------
BEGIN;
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (1, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:51', '2026-02-11 15:37:51');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (2, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', 20260211, 0.00, 0, 0, 0, 12, 0.00, 0.00, '2026-02-11 15:37:51', '2026-02-11 15:37:51');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (3, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', 20260211, 0.00, 0, 0, 0, 3, 0.00, 0.00, '2026-02-11 15:37:51', '2026-02-11 15:37:51');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (4, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:51', '2026-02-11 15:37:51');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (5, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', 20260211, 0.00, 0, 0, 0, 1, 0.00, 0.00, '2026-02-11 15:37:51', '2026-02-11 15:37:51');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (6, 'oppo', '1000521019', '飞猪酒店', 20260211, 300000.00, 5014, 57, 12111, 384945, 59.83, 5263.16, '2026-02-11 15:37:52', '2026-02-11 15:47:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (7, 'oppo', '1000521020', '飞猪-酒店', 20260211, 0.00, 21, 0, 0, 2126, 0.00, 0.00, '2026-02-11 15:37:52', '2026-02-11 15:47:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (8, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', 20260211, 0.00, 2, 0, 0, 230, 0.00, 0.00, '2026-02-11 15:37:52', '2026-02-11 15:47:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (9, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', 20260211, 20300.00, 419, 5, 962, 57948, 48.45, 4060.00, '2026-02-11 15:37:52', '2026-02-11 15:47:19');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (10, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', 20260211, 0.00, 0, 0, 0, 28, 0.00, 0.00, '2026-02-11 15:37:52', '2026-02-11 15:37:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (11, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', 20260211, 0.00, 0, 0, 0, 41, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:47:19');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (12, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', 20260211, 0.00, 4, 0, 0, 500, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:40:23');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (13, 'oppo', '1001197524', '美数-飞猪-拉活-10', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:37:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (14, 'oppo', '1001197530', '美数-飞猪-拉活-11', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:37:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (15, 'oppo', '1001197536', '美数-飞猪-拉活-12', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:37:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (16, 'oppo', '1001197541', '美数-飞猪-拉活-13', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:37:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (17, 'oppo', '1001197547', '美数-飞猪-拉活-14', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:37:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (18, 'oppo', '1001197552', '美数-飞猪-拉活-15', 20260211, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-11 15:37:53', '2026-02-11 15:37:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (55, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:22', '2026-02-12 11:46:22');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (56, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', 20260212, 0.00, 0, 0, 0, 22, 0.00, 0.00, '2026-02-12 11:46:22', '2026-02-12 16:30:17');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (57, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', 20260212, 0.00, 0, 0, 0, 1, 0.00, 0.00, '2026-02-12 11:46:22', '2026-02-12 16:30:17');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (58, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:22', '2026-02-12 11:46:22');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (59, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:22', '2026-02-12 11:46:22');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (60, 'oppo', '1000521019', '飞猪酒店', 20260212, 350000.00, 6259, 66, 14095, 483301, 55.92, 5303.03, '2026-02-12 11:46:23', '2026-02-12 16:30:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (61, 'oppo', '1000521020', '飞猪-酒店', 20260212, 0.00, 7, 0, 0, 1262, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 16:30:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (62, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', 20260212, 0.00, 0, 0, 0, 116, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 16:30:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (63, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', 20260212, 0.00, 7, 0, 0, 630, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 16:30:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (64, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', 20260212, 0.00, 0, 0, 0, 29, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 16:30:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (65, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', 20260212, 0.00, 0, 0, 0, 24, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 16:30:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (66, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', 20260212, 0.00, 3, 0, 0, 297, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 16:30:18');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (67, 'oppo', '1001197524', '美数-飞猪-拉活-10', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 11:46:23');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (68, 'oppo', '1001197530', '美数-飞猪-拉活-11', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 11:46:23');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (69, 'oppo', '1001197536', '美数-飞猪-拉活-12', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 11:46:23');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (70, 'oppo', '1001197541', '美数-飞猪-拉活-13', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 11:46:23');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (71, 'oppo', '1001197547', '美数-飞猪-拉活-14', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 11:46:23');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (72, 'oppo', '1001197552', '美数-飞猪-拉活-15', 20260212, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-02-12 11:46:23', '2026-02-12 11:46:23');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (73, 'xiaomi', '1261482', '杭州淘美航空服务有限公司（飞猪买量非商店104）', 20260212, 139002.00, 1407, 14, 2579, 90697, 99.00, 9929.00, '2026-02-12 11:46:25', '2026-02-12 16:30:20');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (74, 'xiaomi', '1261492', '杭州淘美航空服务有限公司（飞猪买量非商店105）', 20260212, 140.00, 1, 0, 3, 264, 140.00, 0.00, '2026-02-12 11:46:25', '2026-02-12 16:30:20');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (75, 'xiaomi', '1261502', '杭州淘美航空服务有限公司（飞猪买量非商店106）', 20260212, 621.00, 9, 0, 24, 143, 69.00, 0.00, '2026-02-12 11:46:25', '2026-02-12 16:30:21');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (98, 'adn', '1073770254_348117', '飞猪-拉活', 20260212, 48035.00, 1253, 1, 6288, 26220, 38.00, 48035.00, '2026-02-12 15:49:59', '2026-02-12 15:58:59');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (99, 'adn', '1073770254', '飞猪-拉活', 20260212, 8317.00, 1188, 0, 2762, 4655, 7.00, 0.00, '2026-02-12 15:49:59', '2026-02-12 15:58:59');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

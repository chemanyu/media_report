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

 Date: 11/05/2026 14:30:43
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for cainiao_advertiser
-- ----------------------------
DROP TABLE IF EXISTS `cainiao_advertiser`;
CREATE TABLE `cainiao_advertiser` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `media_adv_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '媒体账户ID（如：巨量引擎平台账户ID）',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_unique_record` (`media_adv_id`) COMMENT '唯一约束：防止重复数据'
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜鸟账户表';

-- ----------------------------
-- Records of cainiao_advertiser
-- ----------------------------
BEGIN;
INSERT INTO `cainiao_advertiser` (`id`, `media_adv_id`, `create_time`, `update_time`) VALUES (1, '95790767', '2026-02-06 16:21:08', '2026-02-06 16:21:08');
INSERT INTO `cainiao_advertiser` (`id`, `media_adv_id`, `create_time`, `update_time`) VALUES (2, '95790768', '2026-02-06 16:22:52', '2026-02-06 16:22:52');
INSERT INTO `cainiao_advertiser` (`id`, `media_adv_id`, `create_time`, `update_time`) VALUES (3, '95790769', '2026-02-06 16:22:57', '2026-02-06 16:22:57');
INSERT INTO `cainiao_advertiser` (`id`, `media_adv_id`, `create_time`, `update_time`) VALUES (4, '95790770', '2026-02-06 16:23:01', '2026-02-06 16:23:01');
INSERT INTO `cainiao_advertiser` (`id`, `media_adv_id`, `create_time`, `update_time`) VALUES (5, '95790765', '2026-02-06 16:23:05', '2026-02-06 16:23:05');
INSERT INTO `cainiao_advertiser` (`id`, `media_adv_id`, `create_time`, `update_time`) VALUES (6, '95790766', '2026-02-06 16:23:09', '2026-02-06 16:23:09');
INSERT INTO `cainiao_advertiser` (`id`, `media_adv_id`, `create_time`, `update_time`) VALUES (7, '95790764', '2026-02-06 16:23:13', '2026-02-06 16:23:13');
COMMIT;

-- ----------------------------
-- Table structure for cainiao_cardinality
-- ----------------------------
DROP TABLE IF EXISTS `cainiao_cardinality`;
CREATE TABLE `cainiao_cardinality` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `cardinality` decimal(3,1) NOT NULL COMMENT '菜鸟基数：如1.4、4.1、1.0等，小数点后一位',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_cardinality` (`cardinality`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='菜鸟基数配置表';

-- ----------------------------
-- Records of cainiao_cardinality
-- ----------------------------
BEGIN;
INSERT INTO `cainiao_cardinality` (`id`, `cardinality`, `update_time`, `create_time`) VALUES (1, 1.0, '2026-02-09 15:28:50', '2026-02-09 15:17:24');
COMMIT;

-- ----------------------------
-- Table structure for elm_hc_media_report
-- ----------------------------
DROP TABLE IF EXISTS `elm_hc_media_report`;
CREATE TABLE `elm_hc_media_report` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `performance_id` int NOT NULL COMMENT '客户表id',
  `media_adv_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '媒体账户ID（如：巨量引擎平台账户ID）',
  `media_adv_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '媒体账户名称（如：巨量引擎平台账户名称）',
  `huichuan_adv_id` bigint NOT NULL COMMENT '汇川账户ID',
  `redirect_num` int NOT NULL DEFAULT '0' COMMENT '调起数',
  `pay_num` int NOT NULL DEFAULT '0' COMMENT '付费数',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_unique_record` (`media_adv_id`,`huichuan_adv_id`) COMMENT '唯一约束：防止重复数据'
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='汇川饿了么账户报表';

-- ----------------------------
-- Records of elm_hc_media_report
-- ----------------------------
BEGIN;
INSERT INTO `elm_hc_media_report` (`id`, `performance_id`, `media_adv_id`, `media_adv_name`, `huichuan_adv_id`, `redirect_num`, `pay_num`, `create_time`, `update_time`) VALUES (1, 1, '1826264272898059', 'RD-饿了么-19', 211571796, 0, 0, '2026-02-04 22:43:06', '2026-02-04 22:43:06');
INSERT INTO `elm_hc_media_report` (`id`, `performance_id`, `media_adv_id`, `media_adv_name`, `huichuan_adv_id`, `redirect_num`, `pay_num`, `create_time`, `update_time`) VALUES (2, 1, '1826264271670538', 'RD-饿了么-17', 211571795, 0, 0, '2026-02-04 22:46:34', '2026-02-04 22:46:34');
INSERT INTO `elm_hc_media_report` (`id`, `performance_id`, `media_adv_id`, `media_adv_name`, `huichuan_adv_id`, `redirect_num`, `pay_num`, `create_time`, `update_time`) VALUES (3, 1, '1852910715285515', 'RD-饿了么-41', 211571785, 0, 0, '2026-02-04 22:46:49', '2026-02-04 22:46:49');
INSERT INTO `elm_hc_media_report` (`id`, `performance_id`, `media_adv_id`, `media_adv_name`, `huichuan_adv_id`, `redirect_num`, `pay_num`, `create_time`, `update_time`) VALUES (4, 1, '1826263935546762', 'RD-饿了么-11', 211554240, 1098, 6, '2026-02-04 22:47:05', '2026-02-05 11:23:47');
INSERT INTO `elm_hc_media_report` (`id`, `performance_id`, `media_adv_id`, `media_adv_name`, `huichuan_adv_id`, `redirect_num`, `pay_num`, `create_time`, `update_time`) VALUES (5, 1, '1826264272231820', 'RD-饿了么-18', 211557314, 1274, 0, '2026-02-04 22:47:17', '2026-02-05 11:23:56');
COMMIT;

-- ----------------------------
-- Table structure for elm_hc_performance_report
-- ----------------------------
DROP TABLE IF EXISTS `elm_hc_performance_report`;
CREATE TABLE `elm_hc_performance_report` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `customer_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户名称（如：拉扎斯网络科技（上海）有限公司）',
  `customer_short` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户简称（如：淘宝闪购）',
  `agent_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代理名称（如：北京美数信息科技有限公司）',
  `agent_short` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代理简称（如：美数）',
  `media_platform_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '媒体平台名称（如：巨量引擎）',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_unique_record` (`media_platform_name`) COMMENT '唯一约束：防止重复数据'
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='汇川饿了么数据报表';

-- ----------------------------
-- Records of elm_hc_performance_report
-- ----------------------------
BEGIN;
INSERT INTO `elm_hc_performance_report` (`id`, `customer_name`, `customer_short`, `agent_name`, `agent_short`, `media_platform_name`, `create_time`, `update_time`) VALUES (1, '拉扎斯网络科技（上海）有限公司', '淘宝闪购', '北京美数信息科技有限公司', '美数', '巨量引擎', '2026-02-04 22:34:11', '2026-02-05 10:53:47');
COMMIT;

-- ----------------------------
-- Table structure for error_statistics
-- ----------------------------
DROP TABLE IF EXISTS `error_statistics`;
CREATE TABLE `error_statistics` (
  `id` int NOT NULL AUTO_INCREMENT,
  `query_log_id` int DEFAULT NULL,
  `error_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `error_count` int DEFAULT '0',
  `account_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `query_log_id` (`query_log_id`),
  CONSTRAINT `error_statistics_ibfk_1` FOREIGN KEY (`query_log_id`) REFERENCES `query_logs` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of error_statistics
-- ----------------------------
BEGIN;
COMMIT;

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
) ENGINE=InnoDB AUTO_INCREMENT=343 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='飞猪时报数据表';

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
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (125, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:00', '2026-03-02 14:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (126, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', 20260302, 0.00, 0, 0, 0, 9, 0.00, 0.00, '2026-03-02 14:01:00', '2026-03-02 14:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (127, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', 20260302, 0.00, 0, 0, 0, 1, 0.00, 0.00, '2026-03-02 14:01:00', '2026-03-02 14:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (128, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (129, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (130, 'oppo', '1000521019', '飞猪酒店', 20260302, 148411.00, 3135, 6, 7564, 177560, 47.34, 24735.17, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (131, 'oppo', '1000521020', '飞猪-酒店', 20260302, 0.00, 0, 0, 0, 43, 0.00, 0.00, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (132, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', 20260302, 0.00, 0, 0, 0, 5, 0.00, 0.00, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (133, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (134, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', 20260302, 0.00, 0, 0, 0, 1, 0.00, 0.00, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (135, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:01', '2026-03-02 14:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (136, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', 20260302, 205222.00, 4062, 38, 9227, 491046, 50.52, 5400.58, '2026-03-02 14:01:02', '2026-03-02 14:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (137, 'oppo', '1001197524', '美数-飞猪-拉活-10', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:02', '2026-03-02 14:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (138, 'oppo', '1001197530', '美数-飞猪-拉活-11', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:02', '2026-03-02 14:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (139, 'oppo', '1001197536', '美数-飞猪-拉活-12', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:02', '2026-03-02 14:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (140, 'oppo', '1001197541', '美数-飞猪-拉活-13', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:02', '2026-03-02 14:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (141, 'oppo', '1001197547', '美数-飞猪-拉活-14', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:02', '2026-03-02 14:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (142, 'oppo', '1001197552', '美数-飞猪-拉活-15', 20260302, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-02 14:01:02', '2026-03-02 14:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (143, 'xiaomi', '1261482', '杭州淘美航空服务有限公司（飞猪买量非商店104）', 20260302, 13426.00, 150, 0, 300, 5243, 90.00, 0.00, '2026-03-02 14:01:03', '2026-03-02 14:01:03');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (144, 'xiaomi', '1261492', '杭州淘美航空服务有限公司（飞猪买量非商店105）', 20260302, 18528.00, 247, 5, 575, 37161, 75.00, 3706.00, '2026-03-02 14:01:05', '2026-03-02 14:01:05');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (145, 'xiaomi', '1261502', '杭州淘美航空服务有限公司（飞猪买量非商店106）', 20260302, 0.00, 1, 0, 3, 918, 0.00, 0.00, '2026-03-02 14:01:06', '2026-03-02 14:01:06');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (146, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:00', '2026-03-20 18:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (147, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', 20260320, 0.00, 0, 0, 0, 14, 0.00, 0.00, '2026-03-20 18:01:00', '2026-03-20 18:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (148, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', 20260320, 0.00, 0, 0, 0, 1, 0.00, 0.00, '2026-03-20 18:01:00', '2026-03-20 18:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (149, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:00', '2026-03-20 18:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (150, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', 20260320, 0.00, 0, 0, 0, 1, 0.00, 0.00, '2026-03-20 18:01:00', '2026-03-20 18:01:00');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (151, 'oppo', '1000521019', '飞猪酒店', 20260320, 0.00, 11, 0, 0, 1727, 0.00, 0.00, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (152, 'oppo', '1000521020', '飞猪-酒店', 20260320, 0.00, 0, 0, 0, 40, 0.00, 0.00, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (153, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', 20260320, 0.00, 0, 0, 0, 1, 0.00, 0.00, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (154, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (155, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (156, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (157, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', 20260320, 50000.00, 1204, 7, 3469, 91623, 41.53, 7142.86, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (158, 'oppo', '1001197524', '美数-飞猪-拉活-10', 20260320, 202.00, 19, 0, 202, 6396, 10.63, 0.00, '2026-03-20 18:01:01', '2026-03-20 18:01:01');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (159, 'oppo', '1001197530', '美数-飞猪-拉活-11', 20260320, 2674.00, 20, 1, 58, 6544, 133.70, 2674.00, '2026-03-20 18:01:02', '2026-03-20 18:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (160, 'oppo', '1001197536', '美数-飞猪-拉活-12', 20260320, 323510.00, 754, 76, 2659, 906692, 429.06, 4256.71, '2026-03-20 18:01:02', '2026-03-20 18:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (161, 'oppo', '1001197541', '美数-飞猪-拉活-13', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:02', '2026-03-20 18:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (162, 'oppo', '1001197547', '美数-飞猪-拉活-14', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:02', '2026-03-20 18:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (163, 'oppo', '1001197552', '美数-飞猪-拉活-15', 20260320, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-03-20 18:01:02', '2026-03-20 18:01:02');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (164, 'xiaomi', '1261482', '杭州淘美航空服务有限公司（飞猪买量非商店104）', 20260320, 37552.00, 350, 6, 761, 51673, 107.00, 6259.00, '2026-03-20 18:01:03', '2026-03-20 18:01:03');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (165, 'xiaomi', '1261492', '杭州淘美航空服务有限公司（飞猪买量非商店105）', 20260320, 24303.00, 193, 4, 397, 6274, 126.00, 6076.00, '2026-03-20 18:01:03', '2026-03-20 18:01:03');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (166, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:34');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (167, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', 20260423, 0.00, 0, 0, 0, 6, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:34');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (168, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:34');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (169, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:34');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (170, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:34');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (171, 'oppo', '1000521019', '飞猪酒店', 20260423, 0.00, 0, 0, 0, 252, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (172, 'oppo', '1000521020', '飞猪-酒店', 20260423, 0.00, 0, 0, 0, 3, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (173, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (174, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (175, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (176, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (177, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', 20260423, 50000.00, 1205, 3, 3715, 60201, 41.49, 16666.67, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (178, 'oppo', '1001197524', '美数-飞猪-拉活-10', 20260423, 64965.00, 85, 11, 257, 73093, 764.29, 5905.91, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (179, 'oppo', '1001197530', '美数-飞猪-拉活-11', 20260423, 320063.00, 517, 38, 1831, 683519, 619.08, 8422.71, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (180, 'oppo', '1001197536', '美数-飞猪-拉活-12', 20260423, 50867.00, 93, 10, 365, 169811, 546.96, 5086.70, '2026-04-23 14:20:28', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (181, 'oppo', '1001197541', '美数-飞猪-拉活-13', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:29', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (182, 'oppo', '1001197547', '美数-飞猪-拉活-14', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:29', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (183, 'oppo', '1001197552', '美数-飞猪-拉活-15', 20260423, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-23 14:20:29', '2026-04-23 14:25:35');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (202, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (203, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', 20260424, 0.00, 0, 0, 0, 6, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (204, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (205, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (206, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (207, 'oppo', '1000521019', '飞猪酒店', 20260424, 0.00, 0, 0, 0, 156, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (208, 'oppo', '1000521020', '飞猪-酒店', 20260424, 0.00, 0, 0, 0, 4, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (209, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:52', '2026-04-24 10:31:52');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (210, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (211, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (212, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (213, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', 20260424, 49872.00, 1199, 2, 4065, 67012, 41.59, 24936.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (214, 'oppo', '1001197524', '美数-飞猪-拉活-10', 20260424, 7831.00, 31, 1, 108, 22377, 252.61, 7831.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (215, 'oppo', '1001197530', '美数-飞猪-拉活-11', 20260424, 15653.00, 62, 5, 419, 56769, 252.47, 3130.60, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (216, 'oppo', '1001197536', '美数-飞猪-拉活-12', 20260424, 68653.00, 201, 21, 604, 193613, 341.56, 3269.19, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (217, 'oppo', '1001197541', '美数-飞猪-拉活-13', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (218, 'oppo', '1001197547', '美数-飞猪-拉活-14', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (219, 'oppo', '1001197552', '美数-飞猪-拉活-15', 20260424, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-24 10:31:53', '2026-04-24 10:31:53');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (220, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:36');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (221, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', 20260429, 0.00, 0, 0, 0, 9, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:36');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (222, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:36');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (223, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:36');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (224, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:36');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (225, 'oppo', '1000521019', '飞猪酒店', 20260429, 0.00, 1, 0, 0, 303, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (226, 'oppo', '1000521020', '飞猪-酒店', 20260429, 0.00, 0, 0, 0, 4, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (227, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (228, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (229, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (230, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:52:59', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (231, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', 20260429, 50000.00, 1591, 4, 5072, 80089, 31.43, 12500.00, '2026-04-29 13:52:59', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (232, 'oppo', '1001197524', '美数-飞猪-拉活-10', 20260429, 1444.00, 12, 0, 47, 6323, 120.33, 0.00, '2026-04-29 13:53:00', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (233, 'oppo', '1001197530', '美数-飞猪-拉活-11', 20260429, 65074.00, 227, 18, 960, 251457, 286.67, 3615.22, '2026-04-29 13:53:00', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (234, 'oppo', '1001197536', '美数-飞猪-拉活-12', 20260429, 217990.00, 565, 59, 1745, 558429, 385.82, 3694.75, '2026-04-29 13:53:00', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (235, 'oppo', '1001197541', '美数-飞猪-拉活-13', 20260429, 0.00, 0, 0, 0, 2, 0.00, 0.00, '2026-04-29 13:53:00', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (236, 'oppo', '1001197547', '美数-飞猪-拉活-14', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:53:00', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (237, 'oppo', '1001197552', '美数-飞猪-拉活-15', 20260429, 0.00, 0, 0, 0, 0, 0.00, 0.00, '2026-04-29 13:53:00', '2026-04-29 16:42:37');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (292, 'honor', '424696', '荣耀账户424696', 20260429, 32498.08, 326, 6, 925, 45668, 0.00, 0.00, '2026-04-29 16:24:16', '2026-04-29 16:42:38');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (293, 'honor', '424699', '荣耀账户424699', 20260429, 76258.80, 438, 18, 968, 69097, 0.00, 0.00, '2026-04-29 16:24:16', '2026-04-29 16:42:38');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (294, 'honor', '424773', '荣耀账户424773', 20260429, 487087.29, 2780, 109, 6201, 512634, 0.00, 0.00, '2026-04-29 16:24:16', '2026-04-29 16:42:38');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (295, 'honor', '424776', '荣耀账户424776', 20260429, 100000.00, 629, 4, 843, 63196, 0.00, 0.00, '2026-04-29 16:24:16', '2026-04-29 16:42:38');
INSERT INTO `fz_hourly_report` (`id`, `media`, `media_adv_id`, `media_adv_name`, `report_date`, `cost`, `convert_dp`, `dp_app_order_nums`, `click`, `expose`, `convert_dp_price`, `dp_app_order_price`, `create_time`, `update_time`) VALUES (296, 'honor', '424766', '荣耀账户424766', 20260429, 12014.40, 180, 0, 749, 18962, 0.00, 0.00, '2026-04-29 16:24:16', '2026-04-29 16:42:38');
COMMIT;

-- ----------------------------
-- Table structure for fz_media_advertiser
-- ----------------------------
DROP TABLE IF EXISTS `fz_media_advertiser`;
CREATE TABLE `fz_media_advertiser` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `media` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '媒体简称，oppo, xiaomi, adn',
  `media_adv_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '媒体账户ID',
  `media_adv_name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '媒体账户名称',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `client_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'OAuth2 Client ID（honor专用）',
  `client_secret` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'OAuth2 Client Secret（honor专用）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_unique_record` (`media_adv_id`) COMMENT '唯一约束：防止重复数据'
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='飞猪账户报表';

-- ----------------------------
-- Records of fz_media_advertiser
-- ----------------------------
BEGIN;
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (1, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', '2026-02-11 15:16:06', '2026-02-11 15:16:06', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (2, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', '2026-02-11 15:16:21', '2026-02-11 15:16:21', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (3, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', '2026-02-11 15:16:29', '2026-02-11 15:16:29', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (4, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', '2026-02-11 15:16:36', '2026-02-11 15:16:36', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (5, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', '2026-02-11 15:16:43', '2026-02-11 15:16:43', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (6, 'oppo', '1000521019', '飞猪酒店', '2026-02-11 15:16:51', '2026-02-11 15:16:51', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (7, 'oppo', '1000521020', '飞猪-酒店', '2026-02-11 15:16:58', '2026-02-11 15:16:58', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (8, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', '2026-02-11 15:17:05', '2026-02-11 15:17:05', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (9, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', '2026-02-11 15:17:14', '2026-02-11 15:17:14', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (10, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', '2026-02-11 15:17:22', '2026-02-11 15:17:22', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (11, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', '2026-02-11 15:17:28', '2026-02-11 15:17:28', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (12, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', '2026-02-11 15:17:35', '2026-02-11 15:17:35', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (13, 'oppo', '1001197524', '美数-飞猪-拉活-10', '2026-02-11 15:17:41', '2026-02-11 15:17:41', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (14, 'oppo', '1001197530', '美数-飞猪-拉活-11', '2026-02-11 15:17:48', '2026-02-11 15:17:48', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (15, 'oppo', '1001197536', '美数-飞猪-拉活-12', '2026-02-11 15:17:53', '2026-02-11 15:17:53', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (16, 'oppo', '1001197541', '美数-飞猪-拉活-13', '2026-02-11 15:17:59', '2026-02-11 15:17:59', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (17, 'oppo', '1001197547', '美数-飞猪-拉活-14', '2026-02-11 15:18:05', '2026-02-11 15:18:05', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (18, 'oppo', '1001197552', '美数-飞猪-拉活-15', '2026-02-11 15:18:11', '2026-02-11 15:18:11', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (19, 'xiaomi', '1261482', '杭州淘美航空服务有限公司（飞猪买量非商店104）', '2026-02-12 11:45:49', '2026-02-12 11:45:49', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (20, 'xiaomi', '1261492', '杭州淘美航空服务有限公司（飞猪买量非商店105）', '2026-02-12 11:46:00', '2026-02-12 11:46:00', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (21, 'xiaomi', '1261502', '杭州淘美航空服务有限公司（飞猪买量非商店106）', '2026-02-12 11:46:11', '2026-02-12 11:46:11', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (22, 'adn', '1073770254', '飞猪-拉活', '2026-02-12 16:01:00', '2026-02-12 16:01:00', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (23, 'adn', '1073772364', '飞猪', '2026-04-23 13:36:02', '2026-04-23 13:36:02', '', '');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (24, 'honor', '424696', '荣耀账户424696', '2026-04-29 13:51:08', '2026-04-29 13:51:08', '2049052778130702336', 'kj6HTJKiFl2BDDwVGS9d1fZ9C9GND3Ew');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (25, 'honor', '424699', '荣耀账户424699', '2026-04-29 13:51:08', '2026-04-29 13:51:08', '2049321968009740288', 'FP88pt7KyqPMHgHrrcMF3HQQZ48LWsiF');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (26, 'honor', '424773', '荣耀账户424773', '2026-04-29 13:51:08', '2026-04-29 13:51:08', '2049322050954199040', 'ma3yj0TNPRo07yevhBgNLFc8IX1k6jvj');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (27, 'honor', '424776', '荣耀账户424776', '2026-04-29 13:51:08', '2026-04-29 13:51:08', '2049322125440843776', 'nSBnkJRwRF0xthBBsaeMWqlevhAhvFmu');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`, `client_id`, `client_secret`) VALUES (28, 'honor', '424766', '荣耀账户424766', '2026-04-29 13:51:08', '2026-04-29 13:51:08', '2049322205832937472', '6TvGdmZolSZXby7MXcyR3PtVOhSuumNY');
COMMIT;

-- ----------------------------
-- Table structure for media_token
-- ----------------------------
DROP TABLE IF EXISTS `media_token`;
CREATE TABLE `media_token` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `media` varchar(64) DEFAULT '' COMMENT '媒体名称',
  `token` text COMMENT '媒体token',
  `refresh_token` varchar(62) DEFAULT '' COMMENT '媒体刷新token',
  `agent_id` varchar(64) DEFAULT '' COMMENT '代理商id',
  `advertiser_id` varchar(64) DEFAULT '' COMMENT '账户ID',
  `del_flag` tinyint(1) DEFAULT '0' COMMENT '删除(0:正常;1:删除)',
  `create_time` int DEFAULT '0' COMMENT '创建时间',
  `update_time` int DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni` (`media`,`del_flag`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb3;

-- ----------------------------
-- Records of media_token
-- ----------------------------
BEGIN;
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (1, 'kuaishou', '39f82389a557774b86a40bd42f42836f', 'de08d1d9909783e164fc0c1ff9344be5', '62633', '95790770', 0, 1769149660, 1769161632);
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (2, 'juliang_pachong', 'unionuuid=V2_blRFCEsCRUUlABYDfh4MATcCFA9KAEJCdFpOUHwdWFUIABNeRlZAFnIIT1R6G1lqZAASQkFWRwp1DENLexhY; cud=e11b32607ffdc4e6273ff1d1006c31b5; __jdu=1740564634752297692458; shshshfpa=1f665855-d40d-4573-911e-1817e2ec065a-1740564636; shshshfpx=1f665855-d40d-4573-911e-1817e2ec065a-1740564636; jcap_dvzw_fp=KeaWLFuo8pcS4hUTlKEg5XpWkm238VoTRLkArQdfBtIo7yDw2_OTR5m71rAh1t7xyWt-FgBS4B2tqOm1VbBEuA==; xapieid=jdd03BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRIIAAAAMW7MTH2RAAAAAADF5UMG6JW5HCCYX; b_webp=1; b_avif=1; guid=506306d32ecf8b1aeb57f2f4eb93e8c05249956125234fb26f74336d6487b8d6; mba_muid=1740564634752297692458; x-rp-evtoken=mGW9U4qbzsaBdCMe70m9pJbXbFXe6C8Jsw3mKiteWCHGLPn9LYcODlL17XjUR_UaWKgtKm3eF7qnYvc5KxUGpw%3D%3D; webp=1; visitkey=4831893135550762302; jdd69fo72b8lfeoe=P7NFMPEEEQ67KDC3QBNPJ6H7WGBY5MCBVTVXAM3OIM2E2NEXEANYMYQWZ5SLM7FUGYNCNOYOEL4KZMS4PEIFMY74SE; focus-login-switch=saas; DeviceSeq=51d062e2a052495380a8ffada20493f9; login=true; cn=1; ceshi3.com=000; MNoticeId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80---%E8%AF%9D%E8%B4%B9%E5%88%B8=657; app_id=jdsaas; MNoticeIdjd_778c9c130e0cf=657; focus-token-type=3; MMsgIdjd_778c9c130e0cf=20704488; idtUserSession=69fdacfd-d1c7-4fbf-9bb2-b21c66052078; o2State=; MNoticeId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80-%E5%A4%8D%E6%8A%95=663; appCode=ms0ca95114; __wga=1774441698473.1774441698473.1758537421714.1758537421714.1.2; CSID=EGo%2fSyRVCFtcQVEMWUNbT1A0dnF3L19YSQoHUUBbDQFqeHZ5cXl1dxhUVy1TWlRTWndlYA5TRBduZBh4XFleOkVaQFxOX1lde3dw; qid_uid=cb17f32b-3136-4cdf-bd6b-993afda7b4f9; qid_fs=1774441722233; sidebarStatus=1; qid_ls=1774485116932; qid_ts=1774488195631; qid_vis=6; qid_evord=16; me_fp=ab439f323d83d7aba0cc1db3fabf1c4d; 3AB9D23F7A4B3C9B=BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRII; MNoticeId%E4%BC%99%E4%BC%B4%E8%AE%A1%E5%88%92--%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80=664; MNoticeId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80---JDZT=664; __jdv=209449046|direct|-|none|-|1776234206565; MMsgId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80---JDZT=20805994; me_js_token=jdd03BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRIIAAAAM5TISIDPIAAAAAD4OZWYKCHTVAPQX; focus-team-id=ve_9jazjdntr; focus-client=WEB; me_saas_userInfo=U2FsdGVkX198Fq2VcGntmJ3FZ5/+BEDcJGfWa2QOEg1AbOGQuw/8YP9cEN6LFiayfC8liA3HgN9RHijTBT+/oA==; b_dw=1052; b_dh=1084; b_dpr=1; __jd_ref_cls=Babel_H5FirstClick; MMsgId%E4%BC%99%E4%BC%B4%E8%AE%A1%E5%88%92--%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80=20807936; cvt=168; _t=pFcEIYnNjw/bt+k27JjqcrC0MH9MLNxTgWR5jVwmhSM=; 3AB9D23F7A4B3CSS=jdd03BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRIIAAAAM5V7DW7ZIAAAAACP2SDAF7W7ER3MX; JSESSIONID=B3660D0EDADABD9E8EA3C2E62EE462D2.s1; mp=%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80--JDZT%E5%B9%BF%E7%82%B9%E9%80%9A; TrackID=1KqdrHDX5kCUWzhqImt_JPk1DervEMYeScncwvQqiK0r-0L_SuFjvsTXzDgKmHQMuP98Btg7MNIEiZtq7QfB_x3qix996C3wuTRxWAeM2iyjWZZ6WkOe3SmKJJIgKIdJG9qyWI7RlSPeRZ5mY2kQh4A; light_key=AASBKE7rOxgWQziEhC_QY6yaupcAcImsUZIkA0Ztre_Sj_XnEicRA9C4K9KOmP_K2TxBh1Xrr5CzLMpmOVF1gh-GYiEQtg; pinId=7ivgiUrWDZoa8zJwJl0Nexk38Oip3RighCT6qKiQXUg; pin=%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80--JDZT%E5%B9%BF%E7%82%B9%E9%80%9A; unick=niu4z076oja7vm; _tp=rebo07rdgrnFn%2FXZRur6MN%2By2pABxq3dOhgEpgXHtdQ6uHMOuMFDZcMJGmv2vR9wIlH9ZX%2FsruUzTsD0RNeZbRqDDqBD3qZNdNcsX%2F8xCEc%3D; _pst=%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80--JDZT%E5%B9%BF%E7%82%B9%E9%80%9A; __jda=209449046.1740564634752297692458.1740564634.1776752808.1776770553.166; __jdc=209449046; shshshfpb=BApXWrgjPrPhAkPZFJkqrxyPXcWBZKreABgIXQRlX9xJ1d5l-KsP0t1-zqW7fZcMqJvB65uSihYcUcK086rh7484jQk3qryqOOZDrE35qEF3D7-AVcI3WGqywvzQgsH3OZMv8m-k8zYx6K5PhUGNdhapjHUhgqz7F7acKPBU-8_w; MNoticeId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80--JDZT%E5%B9%BF%E7%82%B9%E9%80%9A=665; MMsgId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80--JDZT%E5%B9%BF%E7%82%B9%E9%80%9A=20813510; thor=4C81D599AF6638A10932C8F36003092E36E6DE411AFF049061A43311450758EC642BB1D49D8D7121CAE7FA0276F534F58FC1283769F84DB02BA072F42FD032063B2C8126E7D657E90F929A6271B3D6E328118A6E9B69529C364B101C7CAF2F3E7676E1A540B4B412D2D5ED753F1890B731743C75E007A9EF8F4D7D9318A1362CFFAA00C47321E691F79F1313F6341E72; flash=3_kRoG1WXTYCdWjEM83dxj2nRwp_hmeZNCa0TERtdQSD9W7TTwjgkw__yZYe-tqRjHUnIbNsKN2opbPyKX8eKuDZWKWTJyrHP9yUGYoZbIPzb62GIxXb0ujXSgPDi0KqKd_7oM2XwmjATxFxCrNuyyzmfcvFUF6GLPxFiiag7RepDAI5x26z2UiNAEXQEUjHrpjrdp; sdtoken=AAbEsBpEIOVjqTAKCQtvQu17BLM6THU_a7r2Cst09SLOQHPzN_RxuEpZVbUGkhEdsjBmQyUORP-RqWFbPWHakWJYjnbNADeFUpURqjaiuDVtbdPtNwGUvSMeU8iK2j9Xl4mg0gWibZevyXHS-RbBpA0Z3V-C5u9c7e6O0NND7vI', '', '', '', 0, 0, 1777452880);
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (3, 'jingcheng_pachong', 'jd_eid=jdd03BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRIIAAAAMWINJI4OAAAAAACWZUXJ5GJBDHXEX; unionuuid=V2_blRFCEsCRUUlABYDfh4MATcCFA9KAEJCdFpOUHwdWFUIABNeRlZAFnIIT1R6G1lqZAASQkFWRwp1DENLexhY; cud=e11b32607ffdc4e6273ff1d1006c31b5; __jdu=1740564634752297692458; o2State={%22webp%22:true%2C%22avif%22:true}; shshshfpa=1f665855-d40d-4573-911e-1817e2ec065a-1740564636; shshshfpx=1f665855-d40d-4573-911e-1817e2ec065a-1740564636; jcap_dvzw_fp=KeaWLFuo8pcS4hUTlKEg5XpWkm238VoTRLkArQdfBtIo7yDw2_OTR5m71rAh1t7xyWt-FgBS4B2tqOm1VbBEuA==; guid=77170d2da975835a8b2e9fa0c2f7202eb3d4e41907061534005232f5f2d2bc32; qrsc=3; xapieid=jdd03BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRIIAAAAMW7MTH2RAAAAAADF5UMG6JW5HCCYX; b_dw=1502; b_dh=728; b_dpr=2; b_webp=1; b_avif=1; guid=506306d32ecf8b1aeb57f2f4eb93e8c05249956125234fb26f74336d6487b8d6; __jda=238391251.1740564634752297692458.1740564634.1758248626.1758279059.22; mba_muid=1740564634752297692458; x-rp-evtoken=mGW9U4qbzsaBdCMe70m9pJbXbFXe6C8Jsw3mKiteWCHGLPn9LYcODlL17XjUR_UaWKgtKm3eF7qnYvc5KxUGpw%3D%3D; webp=1; visitkey=4831893135550762302; __wga=1758537421714.1758537421714.1758537421714.1758537421714.1.1; jdd69fo72b8lfeoe=P7NFMPEEEQ67KDC3QBNPJ6H7WGBY5MCBVTVXAM3OIM2E2NEXEANYMYQWZ5SLM7FUGYNCNOYOEL4KZMS4PEIFMY74SE; focus-login-switch=saas; DeviceSeq=51d062e2a052495380a8ffada20493f9; login=true; cart-main=xx; cn=1; user-key=eea59422-e9ef-4976-9a35-194857cc0799; idtUserSession=bf040a08-7179-4b84-81ba-98a8514a9856; ceshi3.com=000; MNoticeId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80---%E8%AF%9D%E8%B4%B9%E5%88%B8=657; app_id=jdsaas; sidebarStatus=1; MNoticeIdjd_778c9c130e0cf=657; me_fp=fd24709377ac93f9bbbf30ad594895f7; 3AB9D23F7A4B3C9B=BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRII; __jdv=209449046|direct|-|none|-|1768790198414; MNoticeId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80---JDZT=659; MMsgId%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80---JDZT=20705801; me_js_token=jdd03BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRIIAAAAM4A6U2Q3IAAAAACPX2U2U77IPEK4X; me_saas_userInfo=U2FsdGVkX1/Grztr/QtZzOfIoUxcjY8YTZQ7piipEHlJ3qCVTrOqGU6qAoRspMg29qzEOd0WRUYtRQTARugIM5I2tT5duiYK96EoP3HYogU=; focus-token-type=3; focus-team-id=$z9tKQlZAxTS0xOn1WWPz3; js_user_track=; MMsgIdjd_778c9c130e0cf=20704488; _t=HoIjMb/LzW3rpxDOuiDMqr0l0pTF7x3z+HHPB2yMK7Q=; JSESSIONID=115B6E0E313142767D97FE7E19FF99A9.s1; mp=%E4%BC%99%E4%BC%B4%E8%AE%A1%E5%88%92--%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80; TrackID=1PNlUsnH4YqoknM3IUfnlDeCJfmPJgeaipvixrDCOy8mwGH5rwCaxkxagZs_iuRQqzmVmmc7USbBj2vMpFLz7nzXfEqEiaZibzL10fT-xkN_ibMTN624jzZVq_LBZkSbJhMzbV34HnqlsxU5EXyYm3Q; light_key=AASBKE7rOxgWQziEhC_QY6yawCa-wXDgLajck1ZHNsGbo6nBEu46npFrzqIya8YWmCgq_r8H4TWoWQiNHo-WLnH4clEg1w; pinId=C4k73xwU8lX3I4S8SxuoI4N_gKxl9BLfnRytSPX6E9U; pin=%E4%BC%99%E4%BC%B4%E8%AE%A1%E5%88%92--%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80; unick=14pvh8vw2t7f9i; _tp=Zy0qZGk%2F4LH%2B9ChrciUT9XK0E%2B8dfHfw1Sa7EXGdnhfoQ4EfVxjVOq9LyhsnctfEnfX%2FpjfOsyO4wSl1GGI1Oc6dDxrwtTNjzuCmxS990X0%3D; _pst=%E4%BC%99%E4%BC%B4%E8%AE%A1%E5%88%92--%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80; MMsgId%E4%BC%99%E4%BC%B4%E8%AE%A1%E5%88%92--%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80=20710033; thor=B34E7578B0967B8420965A11126652FF608A80FA2406D997FE6015FC3945F1B74DCD10FB3633145AB924322428A2157C578B1CBF6FD6BCC00FC2B8D997EDBAFFA526AB3B0B4E76DE124F33F043E8168EB96C90C0DE4D7725F268A0B1D352D755F73F0B455997ACEE204E2DF04255D7D5CA51F332C2C565165DAF2A247FA12843CB11527B3F94FEFFB0EE71DEE093A1B2; __jda=209449046.1740564634752297692458.1740564634.1769751643.1770024792.112; __jdc=209449046; cvt=118; 3AB9D23F7A4B3CSS=jdd03BZGES66BKGLDCNSLHTOXJZ7MJQDWBY6YJS3ECPTGI4IMWM6ICIG4W7GVCSKVRUQQ5IUIXAMLW3EECRVYLFESMNVRIIAAAAM4DWZPBJYAAAAACRRU4WJTSHGKKMX; MNoticeId%E4%BC%99%E4%BC%B4%E8%AE%A1%E5%88%92--%E7%BE%8E%E6%95%B0%E7%A7%91%E6%8A%80=659; shshshfpb=BApXW-Z-7HvlAkPZFJkqrxyPXcWBZKreABgIXQRlC9xJ1eJl-KcP0sV-y1W7fZ8NRUPB64uSiiocUAq1E6rh74M4jO03rrSqOVpGGYhlqAEDE84BhE5GmCMu4qFRVrBm8eNf1jokzqSjTsMzRz4XTgJsCIkFjrmnAvPpabJke31NaaxrXKQ; flash=3_p97Zbnxsel8V9uz8Dj6u0WFz6IZ1AzBH2laT1cW2RXaVIMmv-a8GyeuB8NnfcMFC68Gnq1L5nQK-1B4X_U7LOlrNpYJlM648VyBTvYbE-a4Q4diIJWukDAJLSoEt8WwY6txsnAUcsfzqY0as_kXdWLZinbiGTPu4iNDrH8tIaBsAAUUTePBgfCnkJlS9CwTNq7P*; sdtoken=AAbEsBpEIOVjqTAKCQtvQu173o9sKo5k2gQTI0tVI8OGsWkCaaYnzL5ZMpj4uf0uSUtefKVWdmJexPiWz5Jooz37fdnuABv-izO8x4kCv78W6hNBg9KtCUrOQ9C6b9SgnKEiDcrYJ1iFXQrcepSuh-GItIuExcY8TlRMBtReIxIGys4QSWw', '', '', '', 0, 0, 1770029703);
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (4, 'juliang_dls', '103bd43e2da9cbb866c9e36e3842db7c1f0f8e12', 'c505f639cd35a84edfa261a7284d522cc2689021', '', '', 0, 0, 1770261849);
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (5, 'juliang_kh', '99ba1b22e27530bf54da2902c4f8f36126d8f3b2', '8d01ecab2dbca125ad9f6de3a6be9be4c8b614aa', '', '', 0, 0, 1770261849);
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (6, 'zfb_pachong', 'JSESSIONID=1AA78E15D72120E047E1158E6FE0D4A6; riskMobileBankSendTime=-1; riskMobileAccoutSendTime=-1; riskMobileCreditSendTime=-1; riskCredibleMobileSendTime=-1; riskOriginalAccountMobileSendTime=-1; cna=0qo5IIV7enACAd3Nmml5AUIB; umdata_=G80628DD3C029F4C421059938224C807CFB739C; apay_aid=1762410937217208; session.cookieNameId=ALIPAYJSESSIONID; userId=2088931782944304; spanner=8x0W1onuQ4WblTWWlGIyR3RbZTW+hylfXt2T4qEYgj0=; entportal_access_channel=standardportal; LoginForm=alipay_login_auth; CLUB_ALIPAY_COM=2088512069383737; ali_apache_tracktmp=\"uid=2088512069383737\"; ALI_PAMIR_SID=\"U73x8EJVfXwoUBnqrrdZF3GtTcz#+xxQWKsoQ1WUYc9n/ridmzcz\"; __TRACERT_COOKIE_bucUserId=2088931782944304; mustAddPartitionedTag=noNeedToAdd; bs_n_lang=zh_CN; iw.userid=\"K1iSL1z1pqWRGn21dzIomg==\"; __compass_session_id=3f056da4-a65b-4265-b8fa-29b93bbb4bb9; _uab_collina=177061996534093291024565; mobileSendTime=-1; credibleMobileSendTime=-1; ctuMobileSendTime=-1; _umdata=G80628DD3C029F4C421059938224C807CFB739C; apay_id=712785577.1739766994986652.1770615931845472.1770619965715022.61; jsh_t_c_e=jsh_t_0.48457809763965654; zone=RZ55B; auth_jwt=e30.eyJleHAiOjE3NzA2MjA5MjE4NjAsInJsIjoiNSwwLDI3LDE5LDI4LDMwLDEzLDEwIiwic2N0IjoiaVdQQ25ZeG9xK3A4a3ZNZytwbU5ybzNGTG9LYWtIZWplMDE5MmVWIiwidWlkIjoiMjA4ODUxMjA2OTM4MzczNyJ9.tUo6D8aDsQTToLm7kCoA4PuIO2CmienwAieRWNCvsjk; apay_sid=357553119.1770619965715022.1770620322.283460.9.177127.9; ALIPAYJSESSIONID=RZ42Eks6tnMCzxtlxwGTrr5xRK1Vy2authRZ55GZ00; rtk=Ezd/ECNpHJiZZUHJGzuRDVVJMncENhWfjwMYQBoIXL2TNhJiXuK; ALIPAYBUMNGJSESSIONID=GZ00XoGGewyD4p1zK5DnQjuaC20eacantbuserviceGZ00; ctoken=jJhIUeVPd81u7otw', '', '', '', 0, 0, 1770636685);
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (7, 'tanx_pachong', 'arms_uid=dfe9b128-f159-42aa-843c-3f849a72a3a6; __itrace_wid=95ae536e-98c9-4b4d-80b3-e28c2a186e6a; cna=VWfNIPe5zgoCASvziD7zfIGi; t_alimama=ee2cffdec21995b8105c3f08a73ab6c9; lego2_cna=4501C0MUK255Y4YCE08YK85T; wk_cookie2=167f0c4d672ac968e0a93b99b0d0321b; v_alimama=0; lgc=; login=true; cancelledSubSites=empty; sn=; t=aad3d47a3693faf2be2983e66f9366c7; cookie2_alimama=15bd10f3b6205075f98cb01e4f8bce85; XSRF-TOKEN=2091b594-181c-4ce1-952d-ff2abfa64050; cookie2=1db157d62a47d1600098fae31701049f; _tb_token_=eaf4a76971871; __wpkreporterwid_=cfd94cff-5125-4946-ab3d-deccfc9949a5; _l_g_=Ug%3D%3D; env_bak=FM%2BgywHD7Unoz8SbIc%2FnUBfFQSvsP66Kkma08fliLj6Y; xlly_s=1; dnk=%5Cu4E0A%5Cu6D77%5Cu7F8E%5Cu6570; uc1=cookie21=V32FPkk%2FhSt4&cookie16=WqG3DMC9UpAPBHGz5QBErFxlCA%3D%3D&cookie14=UoYZbjdFkkD63A%3D%3D&cookie15=VFC%2FuZ9ayeYq2g%3D%3D&pas=0&existShop=false; tracknick=%5Cu4E0A%5Cu6D77%5Cu7F8E%5Cu6570; lid=%E4%B8%8A%E6%B5%B7%E7%BE%8E%E6%95%B0; havana_lgc_exp=1805361636966; unb=2207447452491; cookie1=BvXjtrOZ5cHkMZPXVjPH0yv0jKITKPPKC7y20aZJsfA%3D; cookie17=UUphzWModBWn%2FYeb%2BQ%3D%3D; _nk_=%5Cu4E0A%5Cu6D77%5Cu7F8E%5Cu6570; sgcookie=E100ztcOnntJeIBp7fOcZyCMsQ6TjV7Q6QQpEinnwagXVSGuHkqcppIs59w59tXA5RHEjFw4Rydskql7VJHjXdPZ%2FNw0xELytE7dzdV5RrEKouo%3D; sg=%E6%95%B018; csg=8df32ad0; wk_unb=UUphzWModBWn%2FYeb%2BQ%3D%3D; isg=BEdHq9M34h3jb2cUcDO42Iq31v0RTBsuKuZqZhk0Y1b9iGdKIRyrfoUbLEjWYPOm; rurl=aHR0cHM6Ly9wdWIuYWxpbWFtYS5jb20vP2ZvcndhcmQ9aHR0cCUzQSUyRiUyRnB1Yi5hbGltYW1hLmNvbSUyRnBvcnRhbCUyRnYyJTJGcGFnZXMlMkZhY3Rpdml0eSUyRm9mZmljaWFsJTJGaW5kZXguaHRtJTNGc3BtJTNEYTIxOXQuMTE4MTY5OTUuY2Y3Njg3NzU0LjEuMWM1NzZhMTVFWDJ6UkU%3D; JSESSIONID=0FB4322E4D97339A6EB46B08CD01B17E; tfstk=gMtSaVTzJBKqxDH-97k4hHCSRCsQAxoaAJ6pIpEzpgIJpkdCgYjezLSBAB524_JP28YD_dEy44jUAgjhvflZbcukEMjLIV5Kh_vvn9HVeCppNSjhvfl4ukQofMAWjr21JKHfL9zR9BQdDKCcdzCp9_QYMO6Gv6dp9ieAL9F8JMU8HtCcpMCp9MHXHsXCv6ddvx9x6cul3J1DFvFR41Ym7_JRGkEpDDb51LUUvkKfFatpeswYHn65P1QK2nlWfdpJj9tnBy1BIUdPnBisGiKWN3QfADZc4dTvNZTSOo_DkdxdlEc_SdvJNFQ92mMcl3vN2stnKkfJzpKFR3hTL_Kkan_Mx7GVTFJ22ZOtgoOhRULO2HGsDgy0b1szgywfSk6f_xMb-yxgFuFE5h6hhaBc3JHjhPrlytXf_xMb-ybRnt1Zhxa3-', '', '', '', 0, 0, 1774507422);
INSERT INTO `media_token` (`id`, `media`, `token`, `refresh_token`, `agent_id`, `advertiser_id`, `del_flag`, `create_time`, `update_time`) VALUES (8, 'dhh_pachong', 't=aad3d47a3693faf2be2983e66f9366c7; wk_cookie2=167f0c4d672ac968e0a93b99b0d0321b; cookie2=1db157d62a47d1600098fae31701049f; tkSid=1760177613326_556240329_0.0; _samesite_flag_=true; _tb_token_=eaf4a76971871; thw=cn; sdkSilent=1773670236484; havana_sdkSilent=1773670236484; mt=ci=0_0; tkmb=e=okcO90WRuM9wv7nGTZ2_-rX3DTSVrV-TdPW7X9DGPergabDxBccRKcvcnO8iiwnXiiFPB_IQFOVWLajJKI63f9Z9rtbmhnxUteAIBJBhSwPYLuVcajQge_SdNqaVWClOIyCYESgjl5dSVi5Q0yOZYXsbgVor69P7hI1aZ-3ZD_ykjrecxAMSRzo84ttZ8tvjQ9C_kswJLB1Wmw4VzoKVGYqmY1JsNQIob6uvwItC6OcwEmTz_MnbQUcI4VbMweAAe4wMclroAkFVbHpdHLbehIuSgxhL99vOtDuIKwcMY0VKEw17FpewGr-9tCpAtRQVRFcP7nDXPACqoPhMO46pXbiq_Nw6iOfVmQEnT2zBn7-ybhShgiPNswthg69zVkf5DcbShJUCG-rx0nrjo33587ZzzLDMP5b_KdSAomdgdXdgYh-m7jOlAx_awAQ-_g97lGpcJ2moA8OYVhifpb4Tzl3erkpkUHy64oAlkR1vQ3A0pJThyex4uFbZd5s_ymM0NgWDEEGi-3w&iv=0&et=1774000790&tk_cps_param=874030133; XSRF-TOKEN=1b912238-b177-4f14-a0e0-b8abf9aefc02; wk_unb=UUphzWModBWn%2FYeb%2BQ%3D%3D; cna=1/pwIGSZkE0CASRwzCE/RY3I; lgc=%5Cu4E0A%5Cu6D77%5Cu7F8E%5Cu6570; cancelledSubSites=empty; dnk=%5Cu4E0A%5Cu6D77%5Cu7F8E%5Cu6570; tracknick=%5Cu4E0A%5Cu6D77%5Cu7F8E%5Cu6570; sn=; _hvn_lgc_=0; havana_lgc2_0=eyJoaWQiOjIyMDc0NDc0NTI0OTEsInNnIjoiMjUzYjY1MTFhMWRmY2IyMGZhMGU0NGYxYmFmNDhkZDYiLCJzaXRlIjowLCJ0b2tlbiI6IjF4d2R4bDI4NHB2V0R4M0hOOGw4SEtRIn0; 3PcFlag=1774579361139; xlly_s=1; fastSlient=1774853771569; unb=2207447452491; uc1=existShop=false&cookie15=VFC%2FuZ9ayeYq2g%3D%3D&cookie21=W5iHLLyFfoWm&cookie14=UoYZbj3Y8wXIog%3D%3D&cookie16=VFC%2FuZ9az08KUQ56dCrZDlbNdA%3D%3D&pas=0; uc3=nk2=qiAr7FspCes%3D&lg2=UIHiLt3xD8xYTw%3D%3D&id2=UUphzWModBWn%2FYeb%2BQ%3D%3D&vt3=F8dD29QGc6mJBx6pv3g%3D; csg=aaff2a41; ultraCookieBase=1k6S45BQHTO3QxAuph9cvIwZMzArTGwFkYCfMRSZJhGKIRWk2KHNCEaoEITzCMasvhYNJp2CFK0RbzA2HhHgF%2FWT9zn68GGKIjEBiRXXsjL1CsLlS%2F2kwMfobBxzk6HbHH0ZvjkFYE0%2BTsyheh1umXbk9G5T%2FDPMwiYTTmTLgIvCOe%2BZmk1mSh%2BLSaqpgdMyaglz%2FKm4GfGXfIOvs0dX0bpCRlbGZBVqLQNPF6VpAFDSh%2Bzh9eRIYUdX7dLM01pbx%2FMKDGP%2F7pA%2BdaX3%2BAz7gzTzHzA%3D%3D; cookie17=UUphzWModBWn%2FYeb%2BQ%3D%3D; skt=7e82f5187faf4fab; existShop=MTc3NDg1Mzc3Mw%3D%3D; uc4=nk4=0%40qBqkhc2koMYJ6LPmYM%2Bji%2Fl63Q%3D%3D&id4=0%40U2grFnyVW3yldGfZWWJnA8tPxZuODu%2B4; _cc_=URm48syIZQ%3D%3D; _l_g_=Ug%3D%3D; sg=%E6%95%B018; _nk_=%5Cu4E0A%5Cu6D77%5Cu7F8E%5Cu6570; cookie1=BvXjtrOZ5cHkMZPXVjPH0yv0jKITKPPKC7y20aZJsfA%3D; sgcookie=E100s%2BUq44BApdrlCKV5C1d%2BWLjtAwq5uVWN4gG3eqV9c5cbEouuDGRze2t0Q0A6niqiXDkqtnVNNbunUfphHQ1Wp4dgb2jD%2BxGnzq7nT5WbscM%3D; havana_lgc_exp=1805957773381; isg=BDQ0Y06H4a16intBclkx3LgJBfSmDVj3xfcZ_86V179COdSD9hlahjh6uXHhwZBP; tfstk=gLKm7fqgz_NI9sPrx5SjyTbW2YuJhis1eCEO6GCZz_5WXCK9QOSMaCHfMdlXshfP1s-AWxsGQIppksJxlOrMIIDfXiubbOAvL1stl1UMsdAul1dODGXGc3LYXshfjdRdjXhKvDpXhijZ9XQUidqAbOoO_4QVGWXAYXhKv0pXhGsZ9KLCbXvPC_WVQsrZE7WdaoW2_s7za9XCb1Rw_gWPK96a_ZSa4855QG5w_GkkU_6Gb1RNbYvr6QzV9h8e4X5UiqN2OwKlnZfeE2ZajEEdo6R5ZlleqtfNTK5ubl5nyoIyQClgXM8BZhbe1cq5DpYP4tYoQz5V-ERAFCm30s-k3n5v0bEAgn9yJhdrQl5Mze7dgUGSfZ8BMe7k2jq5TUJkBaYSCrC9ROR1RhhQcsJwdHTdxc2GEevPqg7Uzyl7ll6rB3z_5ZW5E6FOTkQI6oGAnYDuJd_VF91pXY4k1ZW5FDMoEyFFuT6vi', '1b912238-b177-4f14-a0e0-b8abf9aefc02', '', '', 0, 0, 1774858411);
COMMIT;

-- ----------------------------
-- Table structure for qczj_config
-- ----------------------------
DROP TABLE IF EXISTS `qczj_config`;
CREATE TABLE `qczj_config` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `total_uv` bigint NOT NULL DEFAULT '110000',
  `ratio` decimal(5,4) NOT NULL DEFAULT '0.4000',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of qczj_config
-- ----------------------------
BEGIN;
INSERT INTO `qczj_config` (`id`, `total_uv`, `ratio`, `update_time`, `create_time`) VALUES (1, 110000, 0.4000, '2026-03-06 17:12:59', '2026-03-06 17:12:59');
COMMIT;

-- ----------------------------
-- Table structure for qczj_report_data
-- ----------------------------
DROP TABLE IF EXISTS `qczj_report_data`;
CREATE TABLE `qczj_report_data` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `report_date` int NOT NULL,
  `view` bigint NOT NULL DEFAULT '0',
  `click` bigint NOT NULL DEFAULT '0',
  `action` bigint NOT NULL DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_report_date` (`report_date`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of qczj_report_data
-- ----------------------------
BEGIN;
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (6, 2026030600, 88386, 7830, 4645, '2026-03-06 18:43:45', '2026-03-06 18:43:45');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (7, 2026030601, 15664, 4033, 2480, '2026-03-06 18:43:45', '2026-03-06 18:43:45');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (8, 2026030602, 13117, 2786, 1727, '2026-03-06 18:43:45', '2026-03-06 18:43:45');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (9, 2026030603, 10243, 2515, 1558, '2026-03-06 18:43:45', '2026-03-06 18:43:45');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (10, 2026030604, 12576, 3414, 2082, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (11, 2026030605, 20511, 5571, 3443, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (12, 2026030606, 43024, 13433, 8563, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (13, 2026030607, 71758, 23162, 14838, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (14, 2026030608, 66103, 21611, 13871, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (15, 2026030609, 59129, 19610, 12352, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (16, 2026030610, 58306, 18987, 11986, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (17, 2026030611, 62723, 20302, 13015, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (18, 2026030612, 57419, 18047, 11614, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (19, 2026030613, 47657, 14616, 9348, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (20, 2026030614, 44722, 13600, 8741, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (21, 2026030615, 56276, 14541, 9129, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (22, 2026030616, 68076, 15329, 9739, '2026-03-06 18:43:46', '2026-03-06 18:43:46');
INSERT INTO `qczj_report_data` (`id`, `report_date`, `view`, `click`, `action`, `create_time`, `update_time`) VALUES (24, 2026030617, 72152, 16911, 10705, '2026-03-06 18:44:38', '2026-03-06 18:44:38');
COMMIT;

-- ----------------------------
-- Table structure for query_logs
-- ----------------------------
DROP TABLE IF EXISTS `query_logs`;
CREATE TABLE `query_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `query_date` varchar(8) COLLATE utf8mb4_unicode_ci NOT NULL,
  `total_request_count` int DEFAULT '0',
  `total_error_count` int DEFAULT '0',
  `query_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `query_logs_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of query_logs
-- ----------------------------
BEGIN;
INSERT INTO `query_logs` (`id`, `user_id`, `query_date`, `total_request_count`, `total_error_count`, `query_time`) VALUES (1, 1, '20250819', 636639, 467486, '2025-08-20 10:48:05');
INSERT INTO `query_logs` (`id`, `user_id`, `query_date`, `total_request_count`, `total_error_count`, `query_time`) VALUES (2, 1, '20250819', 636639, 467486, '2025-08-20 11:05:45');
INSERT INTO `query_logs` (`id`, `user_id`, `query_date`, `total_request_count`, `total_error_count`, `query_time`) VALUES (3, 1, '20250819', 636639, 467486, '2025-08-20 11:20:38');
INSERT INTO `query_logs` (`id`, `user_id`, `query_date`, `total_request_count`, `total_error_count`, `query_time`) VALUES (4, 1, '20250819', 636639, 467486, '2025-08-20 11:29:24');
INSERT INTO `query_logs` (`id`, `user_id`, `query_date`, `total_request_count`, `total_error_count`, `query_time`) VALUES (5, 1, '20250819', 636639, 467486, '2025-08-20 11:29:33');
INSERT INTO `query_logs` (`id`, `user_id`, `query_date`, `total_request_count`, `total_error_count`, `query_time`) VALUES (6, 1, '20250819', 636639, 467486, '2025-08-20 11:30:19');
INSERT INTO `query_logs` (`id`, `user_id`, `query_date`, `total_request_count`, `total_error_count`, `query_time`) VALUES (7, 1, '20250819', 636639, 467486, '2025-08-20 11:39:48');
COMMIT;

-- ----------------------------
-- Table structure for rebate
-- ----------------------------
DROP TABLE IF EXISTS `rebate`;
CREATE TABLE `rebate` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `subject` varchar(50) NOT NULL COMMENT '主体：新杰、晴川、美数、魔米、归流',
  `port` varchar(50) NOT NULL COMMENT '端口：优居、至也、各界',
  `rebate_rate` decimal(5,3) NOT NULL COMMENT '返点率：如0.025表示2.5%',
  `subject_type` tinyint NOT NULL DEFAULT '1' COMMENT '主体类型：1-京东主体, 2-三方主体（用于区分规则）',
  `remark` varchar(255) DEFAULT '' COMMENT '备注',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_subject_port` (`subject`,`port`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='返点配置表';

-- ----------------------------
-- Records of rebate
-- ----------------------------
BEGIN;
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (2, '新杰', '优居', 0.000, 1, '', '2026-01-28 16:40:44', '2026-01-27 15:46:37');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (4, '晴川', '优居', 0.000, 1, '', '2026-01-27 16:42:42', '2026-01-27 16:42:42');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (5, '美数', '优居', 0.025, 2, '', '2026-01-28 16:41:47', '2026-01-28 16:41:47');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (6, '震来', '优居', 0.025, 2, '', '2026-01-28 16:42:11', '2026-01-28 16:42:11');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (7, '归流', '优居', 0.025, 2, '', '2026-01-28 16:42:44', '2026-01-28 16:42:44');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (8, '新杰', '至也', 0.000, 1, '', '2026-01-28 16:43:13', '2026-01-28 16:43:13');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (9, '晴川', '至也', 0.000, 1, '', '2026-01-28 16:43:29', '2026-01-28 16:43:29');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (10, '美数', '至也', 0.040, 2, '', '2026-01-28 16:43:47', '2026-01-28 16:43:47');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (11, '震来', '至也', 0.040, 2, '', '2026-01-28 16:44:01', '2026-01-28 16:44:01');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (12, '归流', '至也', 0.040, 2, '', '2026-01-28 16:44:12', '2026-01-28 16:44:12');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (13, '新杰', '谷界', 0.000, 1, '', '2026-01-28 16:44:37', '2026-01-28 16:44:37');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (14, '晴川', '谷界', 0.000, 1, '', '2026-01-28 16:44:50', '2026-01-28 16:44:50');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (15, '美数', '谷界', 0.035, 2, '', '2026-01-28 16:45:10', '2026-01-28 16:45:10');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (16, '震来', '谷界', 0.035, 2, '', '2026-01-28 16:45:26', '2026-01-28 16:45:26');
INSERT INTO `rebate` (`id`, `subject`, `port`, `rebate_rate`, `subject_type`, `remark`, `update_time`, `create_time`) VALUES (17, '归流', '谷界', 0.035, 2, '', '2026-01-28 16:45:44', '2026-01-28 16:45:44');
COMMIT;

-- ----------------------------
-- Table structure for service_fee
-- ----------------------------
DROP TABLE IF EXISTS `service_fee`;
CREATE TABLE `service_fee` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `service_provider` varchar(50) NOT NULL COMMENT '服务商名称：通途、蚁行、创效、凯旋、星河、云谷、美数',
  `fee_rate` decimal(5,3) NOT NULL COMMENT '服务费率：如0.04表示4%',
  `remark` varchar(255) DEFAULT '' COMMENT '备注',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_provider` (`service_provider`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='服务费配置表';

-- ----------------------------
-- Records of service_fee
-- ----------------------------
BEGIN;
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (1, '通途', 0.040, '', '2026-01-27 16:42:57', '2026-01-27 16:42:57');
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (2, '蚁行', 0.040, '', '2026-01-28 16:46:23', '2026-01-28 16:46:23');
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (3, '创效', 0.040, '', '2026-01-28 16:46:41', '2026-01-28 16:46:41');
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (4, '凯旋', 0.030, '', '2026-01-28 16:46:53', '2026-01-28 16:46:53');
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (5, '星河', 0.030, '', '2026-01-28 16:47:03', '2026-01-28 16:47:03');
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (6, '云谷', 0.030, '', '2026-01-28 16:47:14', '2026-01-28 16:47:14');
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (7, '美数', 0.000, '', '2026-01-28 16:47:25', '2026-01-28 16:47:25');
INSERT INTO `service_fee` (`id`, `service_provider`, `fee_rate`, `remark`, `update_time`, `create_time`) VALUES (8, '鑫聚量', 0.040, '', '2026-01-29 18:23:29', '2026-01-29 18:23:29');
COMMIT;

-- ----------------------------
-- Table structure for tanx_monitor
-- ----------------------------
DROP TABLE IF EXISTS `tanx_monitor`;
CREATE TABLE `tanx_monitor` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ds` varchar(32) NOT NULL COMMENT '日期',
  `pid` varchar(64) NOT NULL COMMENT '广告位',
  `adzone_name` varchar(255) DEFAULT NULL COMMENT '广告名称',
  `qingqiupv` varchar(32) DEFAULT '' COMMENT 'tanx有效请求',
  `active_ratio_df` varchar(32) DEFAULT '' COMMENT '东风手淘换端率-同步点击',
  `tanx_effect_pv` varchar(32) DEFAULT '' COMMENT '有效展现pv',
  `tanx_clk` varchar(32) DEFAULT '' COMMENT '有效点击',
  `dongfeng_ef` varchar(32) DEFAULT '' COMMENT '有效媒体消耗结算',
  `create_time` int NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_ds_pid` (`ds`,`pid`)
) ENGINE=InnoDB AUTO_INCREMENT=937 DEFAULT CHARSET=utf8mb3 COMMENT='阿里妈妈展点唤报表';

-- ----------------------------
-- Records of tanx_monitor
-- ----------------------------
BEGIN;
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (904, '2026-02-08', 'mm_1902210064_2348000105_111662900182', '佳投联盟-信息流-ios', '973273025', '93.86%', '860413', '138433', '21654.99', 1770635350);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (905, '2026-02-08', 'mm_1902210064_2348000105_111663100185', '佳投联盟-信息流-安卓', '3453206492', '70.19%', '1902781', '118338', '26442.23', 1770635350);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (906, '2026-02-08', 'mm_1902210064_2348000105_113763000348', '佳投_激励视频_ios', '157479746', '76.97%', '138874', '100236', '5156.12', 1770635351);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (907, '2026-02-08', 'mm_1902210064_2348000105_113765750160', '佳投_激励视频_Android', '706141326', '96.48%', '210765', '121405', '7349.56', 1770635351);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (908, '2026-02-08', 'mm_3447365382_2749500311_114553400117', '有境_信息流_ios', '456355322', '96.53%', '334876', '53494', '7195.64', 1770635352);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (909, '2026-02-08', 'mm_3447365382_2749500311_114552150188', '有境_信息流_Android', '1440960277', '97.64%', '563550', '122786', '14978.68', 1770635352);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (910, '2026-02-08', 'mm_1873810155_2320450209_111384500336', '新数DSP_信息流_Android', '1408619861', '79.48%', '1308598', '360743', '33084.76', 1770635352);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (911, '2026-02-08', 'mm_1873810155_2320450209_111953150429', '新数DSP-信息流独占流量-Android', '1146823908', '0.00%', '15192786', '1380014', '11409.63', 1770635353);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (912, '2026-02-08', 'mm_1873810155_2320450209_111953700467', '新数DSP-信息流独占流量-ios', '129173939', '96.24%', '240438', '23840', '2056.48', 1770635353);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (913, '2026-02-08', 'mm_1873810155_2320450209_115792250490', '新数联盟-视频信息流-IOS', '213472724', '92.31%', '862993', '128145', '12539.58', 1770635354);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (914, '2026-02-08', 'mm_1873810155_2320450209_115890350442', '新数DSP_信息流/推荐流/图集_Android', '1403888277', '78.09%', '3579603', '269991', '14720.15', 1770635354);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (915, '2026-02-08', 'mm_1873810155_2320450209_115894200243', '新数DSP_信息流/推荐流/图集_IOS', '612419426', '90.42%', '2608906', '96265', '13504.68', 1770635355);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (916, '2026-02-08', 'mm_1873810155_2320450209_115797000138', '新数联盟-视频信息流-安卓', '70682997', '92.48%', '26535', '7771', '683.03', 1770635355);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (917, '2026-02-08', 'mm_1873810155_2320450209_111388400121', '新数DSP_信息流_IOS', '144654093', '0.00%', '3772266', '613265', '2479.38', 1770635355);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (918, '2026-02-08', 'mm_1861850082_2364600045_115796950045', '快友联盟-信息流-IOS(新）', '734891633', '98.38%', '365585', '42935', '5341.34', 1770635356);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (919, '2026-02-08', 'mm_1861850082_2364600045_111952850176', '快友_信息流独占流量_Android', '139020507', '61.42%', '72631', '5467', '634.24', 1770635356);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (920, '2026-02-08', 'mm_1861850082_2364600045_111952350174', '快友_信息流独占流量_ios', '197084742', '90.95%', '177866', '20657', '3087.36', 1770635357);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (921, '2026-02-08', 'mm_1861850082_2364600045_111484350168', '快友_联盟信息流_Android', '697303003', '63.52%', '192333', '11708', '2007.35', 1770635357);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (922, '2026-02-08', 'mm_1562690005_2184850029_111650700372', '多盟-信息流-安卓', '3456633150', '84.03%', '3544897', '369149', '26853.7', 1770635358);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (923, '2026-02-08', 'mm_1562690005_2184850029_111651000378', '多盟-信息流-ios', '717435603', '84.85%', '2406701', '99243', '26222.24', 1770635358);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (924, '2026-02-08', 'mm_1562690005_2184850029_112164850011', '多盟-视频资源-ios', '215532144', '100.44%', '396868', '274887', '13740.9', 1770635358);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (925, '2026-02-08', 'mm_1562690005_2184850029_112159050469', '多盟-视频资源-Android', '952159092', '94.53%', '1163753', '576593', '26448.95', 1770635359);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (926, '2026-02-08', 'mm_6899417631_3131950139_115762350460', '浩睿科技-信息流-IOS', '238408446', '83.56%', '496345', '26114', '5818.62', 1770635359);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (927, '2026-02-08', 'mm_6899417631_3131950139_115765600056', '浩睿科技-信息流-安卓', '1007990211', '109.38%', '415670', '70944', '6452.5', 1770635360);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (928, '2026-02-08', 'mm_6899417631_3131950139_115768000140', '浩睿科技-开屏-IOS', '250237998', '91.79%', '375524', '108013', '7813.15', 1770635360);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (929, '2026-02-08', 'mm_6899417631_3131950139_115840550019', '浩睿联盟-信息流2-IOS', '31916585', '93.83%', '28122', '1703', '210.41', 1770635361);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (930, '2026-02-08', 'mm_6899417631_3131950139_115834100467', '浩睿联盟-信息流2-安卓', '536916591', '87.90%', '147132', '8364', '1417.04', 1770635361);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (931, '2026-02-08', 'mm_2279650033_2540650473_111835400327', '美数_信息流_Android', '3685629213', '93.16%', '3179341', '351011', '50498.13', 1770635363);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (932, '2026-02-08', 'mm_2279650033_2540650473_111835250296', '美数_信息流_ios', '428712268', '89.89%', '4468698', '170055', '59079.73', 1770635363);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (933, '2026-02-08', 'mm_2279650033_2540650473_111835500278', '美数_开屏_Android', '48736186', '0.00%', '10141', '1713', '311.42', 1770635363);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (934, '2026-02-08', 'mm_2279650033_2540650473_111835750274', '美数_开屏_ios', '308342482', '83.77%', '1794642', '366812', '38421.64', 1770635364);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (935, '2026-02-08', 'mm_2279650033_2540650473_114534000022', '美数_激励视频_Android', '42755666', '5.05%', '53806', '31764', '1666.86', 1770635364);
INSERT INTO `tanx_monitor` (`id`, `ds`, `pid`, `adzone_name`, `qingqiupv`, `active_ratio_df`, `tanx_effect_pv`, `tanx_clk`, `dongfeng_ef`, `create_time`) VALUES (936, '2026-02-08', 'mm_2279650033_2540650473_114532100161', '美数_激励视频_ios', '94535729', '89.90%', '658721', '92886', '10881.65', 1770635365);
COMMIT;

-- ----------------------------
-- Table structure for task_type
-- ----------------------------
DROP TABLE IF EXISTS `task_type`;
CREATE TABLE `task_type` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '任务名称，如：app、首购',
  `code` varchar(20) NOT NULL COMMENT '任务编码',
  `settlement_price` decimal(10,2) NOT NULL COMMENT '结算单价',
  `media` varchar(50) NOT NULL COMMENT '媒体平台，如：巨量',
  `status` tinyint NOT NULL DEFAULT '1' COMMENT '状态：1-启用, 0-停用',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_media_code` (`media`,`code`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='任务类型表';

-- ----------------------------
-- Records of task_type
-- ----------------------------
BEGIN;
INSERT INTO `task_type` (`id`, `name`, `code`, `settlement_price`, `media`, `status`, `create_time`, `update_time`) VALUES (1, 'APP', 'app', 42.00, '巨量', 1, '2026-01-27 16:43:22', '2026-01-27 16:43:22');
INSERT INTO `task_type` (`id`, `name`, `code`, `settlement_price`, `media`, `status`, `create_time`, `update_time`) VALUES (2, '首购', 'shougou', 58.00, '巨量', 1, '2026-01-29 12:06:58', '2026-01-29 12:06:58');
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` (`id`, `username`, `password_hash`, `email`, `created_at`, `is_active`) VALUES (1, 'admin', '0192023a7bbd73250516f069df18b500', 'admin@example.com', '2025-08-20 10:41:32', 1);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

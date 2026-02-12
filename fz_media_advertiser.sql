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

 Date: 12/02/2026 16:41:29
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

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
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_unique_record` (`media_adv_id`) COMMENT '唯一约束：防止重复数据'
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='飞猪账户报表';

-- ----------------------------
-- Records of fz_media_advertiser
-- ----------------------------
BEGIN;
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (1, 'oppo', '1000336285', 'SMBBJ-飞猪酒店-XXL', '2026-02-11 15:16:06', '2026-02-11 15:16:06');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (2, 'oppo', '1000391889', 'SMBBJ-飞猪-XXL', '2026-02-11 15:16:21', '2026-02-11 15:16:21');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (3, 'oppo', '1000391890', 'SMBBJ-飞猪酒店-XXL', '2026-02-11 15:16:29', '2026-02-11 15:16:29');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (4, 'oppo', '1000397947', 'SMBBJ-飞猪push-XXL', '2026-02-11 15:16:36', '2026-02-11 15:16:36');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (5, 'oppo', '1000485531', 'SMBBJ-飞猪酒店-XXL', '2026-02-11 15:16:43', '2026-02-11 15:16:43');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (6, 'oppo', '1000521019', '飞猪酒店', '2026-02-11 15:16:51', '2026-02-11 15:16:51');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (7, 'oppo', '1000521020', '飞猪-酒店', '2026-02-11 15:16:58', '2026-02-11 15:16:58');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (8, 'oppo', '1000765351', '美数-飞猪-拉活-春运交通会场', '2026-02-11 15:17:05', '2026-02-11 15:17:05');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (9, 'oppo', '1001156858', '美数-飞猪-拉活-2026度假春节会场', '2026-02-11 15:17:14', '2026-02-11 15:17:14');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (10, 'oppo', '1001156859', '美数-飞猪-拉活-2026酒店春节会场', '2026-02-11 15:17:22', '2026-02-11 15:17:22');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (11, 'oppo', '1001156860', '美数-飞猪-拉活-春节爆款酒店', '2026-02-11 15:17:28', '2026-02-11 15:17:28');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (12, 'oppo', '1001156862', '美数-飞猪-拉活-出境机票春节提前购', '2026-02-11 15:17:35', '2026-02-11 15:17:35');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (13, 'oppo', '1001197524', '美数-飞猪-拉活-10', '2026-02-11 15:17:41', '2026-02-11 15:17:41');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (14, 'oppo', '1001197530', '美数-飞猪-拉活-11', '2026-02-11 15:17:48', '2026-02-11 15:17:48');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (15, 'oppo', '1001197536', '美数-飞猪-拉活-12', '2026-02-11 15:17:53', '2026-02-11 15:17:53');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (16, 'oppo', '1001197541', '美数-飞猪-拉活-13', '2026-02-11 15:17:59', '2026-02-11 15:17:59');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (17, 'oppo', '1001197547', '美数-飞猪-拉活-14', '2026-02-11 15:18:05', '2026-02-11 15:18:05');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (18, 'oppo', '1001197552', '美数-飞猪-拉活-15', '2026-02-11 15:18:11', '2026-02-11 15:18:11');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (19, 'xiaomi', '1261482', '杭州淘美航空服务有限公司（飞猪买量非商店104）', '2026-02-12 11:45:49', '2026-02-12 11:45:49');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (20, 'xiaomi', '1261492', '杭州淘美航空服务有限公司（飞猪买量非商店105）', '2026-02-12 11:46:00', '2026-02-12 11:46:00');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (21, 'xiaomi', '1261502', '杭州淘美航空服务有限公司（飞猪买量非商店106）', '2026-02-12 11:46:11', '2026-02-12 11:46:11');
INSERT INTO `fz_media_advertiser` (`id`, `media`, `media_adv_id`, `media_adv_name`, `create_time`, `update_time`) VALUES (22, 'adn', '1073770254', '飞猪-拉活', '2026-02-12 16:01:00', '2026-02-12 16:01:00');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;

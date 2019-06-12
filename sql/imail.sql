/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 50562
 Source Host           : localhost:3306
 Source Schema         : imail

 Target Server Type    : MySQL
 Target Server Version : 50562
 File Encoding         : 65001

 Date: 12/06/2019 17:19:58
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for im_mail
-- ----------------------------
DROP TABLE IF EXISTS `im_mail`;
CREATE TABLE `im_mail` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `mail_from` varchar(255) NOT NULL DEFAULT '' COMMENT '邮件来源',
  `mail_to` varchar(255) NOT NULL DEFAULT '' COMMENT '发送给谁',
  `content` text NOT NULL COMMENT '邮件内容',
  `size` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '内容大小',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
  `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for im_user
-- ----------------------------
DROP TABLE IF EXISTS `im_user`;
CREATE TABLE `im_user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(32) NOT NULL COMMENT '密码',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
  `update_time` bigint(20) unsigned NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for im_user_box
-- ----------------------------
DROP TABLE IF EXISTS `im_user_box`;
CREATE TABLE `im_user_box` (
  `id` bigint(255) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `mid` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '邮件ID',
  `size` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '邮件字节大小',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '类型， 0:接收邮件;1:发送邮件',
  `create_time` bigint(20) unsigned NOT NULL COMMENT '创建时间',
  `update_time` bigint(20) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;

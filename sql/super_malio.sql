drop table if exists `server_info`;
create table `server_info`(
	`key` varchar(128) not null default '' comment '',
	`value` varchar(128) not null default '' comment '',
	key(`key`)
) engine=innodb default character set=utf8 collate=utf8_general_ci;

drop table if exists `account`;
create table `account`(
	`name` varchar(128) not null comment '账号',
	`salt` varbinary(20) not null comment '',
	`vkey` varbinary(256) not null comment '',
    `guid` int unsigned not null comment '角色ID',
	key(`name`)
) engine=innodb default character set=utf8 collate=utf8_general_ci;

drop table if exists `role_simple`;
create table `role_simple`(
    `guid` int unsigned not null comment '角色ID',
    `occupation` tinyint unsigned not null default 0 comment '职业',
    `gender` tinyint unsigned not null default 0 comment '性别',
    `hp` int unsigned not null default 0 comment '血量',
    `mp` int unsigned not null default 0 comment '蓝量',
    `xp` int unsigned not null default 0 comment '经验',
    `mapid` smallint unsigned not null default 0 comment '当前地图',
    `mapx` tinyint unsigned not null default 0 comment '当前地图坐标x',
    `mapy` tinyint unsigned not null default 0 comment '当前地图坐标y',
    `level` smallint unsigned not null default 0 comment '等级',
    `name` char(32) not null default '' comment '角色名字',
    primary key(`guid`)
) engine=innodb default character set=utf8 collate=utf8_general_ci;


drop table if exists `item`;
create table `item`(
    `role_id` int unsigned not null comment '角色ID',
    `item_id` int unsigned not null comment '物品ID',
    `level` int unsigned not null default 0 comment '等级',
    `data` varchar(128) not null default '' comment '额外数据',
    key(`role_id`)
) engine=innodb default character set=utf8 collate=utf8_general_ci;

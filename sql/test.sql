delete from `server_info`;
delete from `account`;
delete from `role_simple`;
delete from `item`;

insert into `server_info` values('platform', 'test');
insert into role_simple(guid, `level`, name) values(1, 1, "agan"); 
insert into item(role_id, item_id, level, data) values(1, 1, 1, "wapen");

mysql
create database mydb default character set utf8 collate utf8_bin;
grant all on mydb.* to 'cloud'@'%' identified by 'passwd' with grant option;
flush privileges;
use mydb
create table t_employee 
(id int primary key auto_increment, 
name varchar(50) not null, 
address varchar(100) not null,
age int not null);
insert into t_employee values(NULL, 'zhangsan', 'jilin', 23);
insert into t_employee values(NULL, 'lisi', 'beijing', 25);
insert into t_employee values(NULL, 'wangwu', 'shanghai', 24);


go get -u github.com/jteeuwen/go-bindata/...
go-bindata.exe -o asset/asset.go -pkg=asset static/... view/...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o rest restdemo

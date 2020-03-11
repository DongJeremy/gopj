use mydb
create table t_user (
    id int primary key auto_increment, 
    name varchar(50) not null, 
    address varchar(100) not null,
    age int not null);
insert into t_user values(NULL, 'zhangsan', 'jilin', 23);
insert into t_user values(NULL, 'lisi', 'beijing', 25);
insert into t_user values(NULL, 'wangwu', 'shanghai', 24);

conn oracle/ora123@orcl
create table t_user (  
    id         number(10)   primary key,  
    name       varchar2(50),  
    address    varchar2(100), 
    age        number(10)
);
create sequence t_user1_id_seq increment by 1 start with 1 maxvalue 9999999999 cycle;
create or replace trigger t_user_trigger1 
before insert on t_user 
for each row when(new.id is null)
begin
select t_user1_id_seq.nextval into:NEW.ID from dual;
end;
/

insert into t_user values(null, 'Peter', 'jilin', 20);
insert into t_user values(null, 'Wang', 'shanghai', 30);
select * from t_user;
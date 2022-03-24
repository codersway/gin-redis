-- auto-generated definition
create table test
(
    id int auto_increment
        primary key,
    v  int null,
    constraint test_id_uindex
        unique (id)
);

INSERT INTO mtest.test (id, v) VALUES (100, 16);
INSERT INTO mtest.test (id, v) VALUES (101, 1111);

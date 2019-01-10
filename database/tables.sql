drop table antichistes;

create table antichistes (
    id			serial	primary key,
    first_part	text	not null,
    second_part text	not null,
    votes		integer not null default 0,
    public      boolean not null default false
);

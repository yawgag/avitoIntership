CREATE EXTENSION IF NOT EXISTS pgcrypto;

create table role(
	id serial primary key,
	name text not null);

create table users(
	id serial primary key,
	email text not null,
	password text not null,
	roleId int not null references role(id));
	
create table sessions(
        sessionId text primary key,
        userId int not null,
        userRole text not null,
        expireAt timeStamp not null);

create table cities (
    id   serial primary key,
    name text not null unique);

create table reception_statuses (
    id   serial primary key,
    name text not null unique);

create table product_types (
    id   serial primary key,
    name text not null unique);

create table pvzs (
    id       UUID primary key default gen_random_uuid(),
    reg_date date not null default CURRENT_DATE,
    city_id  int not null references cities(id) ON DELETE RESTRICT);

create table receptions (
	id UUID primary key default gen_random_uuid(),
	reception_start_datetime TIMESTAMPTZ not null default now(),
	pvz_id UUID not null references pvzs(id)ON DELETE RESTRICT,
	status_id int not null default 1 references reception_statuses(id) ON DELETE RESTRICT);

create table products (
	id UUID primary key default gen_random_uuid(),
	added_at TIMESTAMPTZ not null default now(),
	type_id int  not null references product_types(id) ON DELETE RESTRICT);

create table reception_products (
    reception_id UUID not null references receptions(id) ON DELETE CASCADE,
    product_id   UUID not null references products(id) ON DELETE CASCADE,
    primary key (reception_id, product_id));



insert into cities(name)
values ('Москва'),
	('Санкт-Петербург'),
	('Казань');

insert into reception_statuses(name)
values ('in_progress'),
	('close');

insert into product_types(name)
values ('обувь'),
	('одежда'),
	('электроника');


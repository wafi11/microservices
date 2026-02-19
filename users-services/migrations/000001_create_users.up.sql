create table users (
    id serial primary key,
    full_name varchar(100),
    username varchar(50),
    email varchar(200) unique,
    password text,
    phone_number varchar(20),
    is_active boolean,
    is_deleted boolean default false,
    deleted_at timestamp,
    created_at timestamp default current_timestamp,
    updated_at timestamp default  current_timestamp
);

create index idx_users_username on users(username);
create unique index idx_users_email on users(email) where is_deleted = false;
create unique index idx_users_phone_number on users(phone_number);
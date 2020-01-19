create table users (
    created_at timestamp default now(),
    updated_at timestamp default now(),
    id varchar(64) primary key,
    username varchar(255) not null unique,
    password_hash varchar(255) not null,
    admin bool not null default false
);

create trigger trg_set_update_ts_users before update on users for each row execute procedure set_update_ts();
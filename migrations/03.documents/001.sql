create table documents (
    created_at timestamp default now(),
    updated_at timestamp default now(),
    id serial primary key,
    path varchar not null unique,
    project_id int references projects(id) on delete cascade,
    title varchar(256),
    body text
);

create trigger trg_set_update_ts_documents before update on documents for each row execute procedure set_update_ts();
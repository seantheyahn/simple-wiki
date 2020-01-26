create table projects (
    created_at timestamp default now(),
    updated_at timestamp default now(),
    id serial primary key,
    title varchar(256),
    description text
);

create trigger trg_set_update_ts_projects before update on projects for each row execute procedure set_update_ts();

--projects and users junction table
create table projects__users (
    project_id int references projects(id) on delete cascade, 
    user_id varchar(64) references users(id) on delete cascade,
    can_write bool not null default true,
    constraint pk_projects__users primary key (project_id, user_id)
);
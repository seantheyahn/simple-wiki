alter table documents
    add column sort_order int not null default 0,
    drop column path;
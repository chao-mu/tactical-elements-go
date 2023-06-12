create table player (
    id serial primary key,
    gold int not null default 0,
    experience int not null default 0,
    health int not null default 3,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP
);

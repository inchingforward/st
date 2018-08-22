create table story (
    id bigserial primary key,
    title text not null,
    uuid text not null,
    authors text not null,
    private boolean default false not null,
    started_at timestamp with time zone default now(),
    published boolean default false not null,
    published_at timestamp with time zone
);

create table story_part (
    id bigserial primary key,
    story_id bigint references story(id),
    part_num int not null,
    part_text text not null,
    written_by text not null,
    written_at timestamp with time zone not null default now()
);

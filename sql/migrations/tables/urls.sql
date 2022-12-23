CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists urls (
    "key" varchar(4) not null primary key
    , "url" varchar(1024)
    , "visits" int8 not null default 0
    , "issuer" varchar(255)
);

create index if not exists urls_key_idx on "urls" ("key");
create index if not exists urls_issuer_idx on "urls" ("issuer");
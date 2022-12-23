alter table urls add column if not exists sub uuid;
create index if not exists urls_sub_idx on "urls" ("sub");
alter table tweets add column parent_id UUID references tweets(id) on delete cascade;
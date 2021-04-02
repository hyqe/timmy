

create table congregations (
	_id uuid primary key,
	_created timestamptz default now(),
    _modified timestamptz default now(),
	name text not null,
	circuit text not null default '',
	street text not null default '',
	city text not null default '',
	state text not null default '',
	url text not null default '',
	phone text not null default ''
);

create table Publishers (
    _id uuid primary key,
	_created timestamptz default now(),
);
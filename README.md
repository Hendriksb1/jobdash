go run insert.go -firm="" -type_job="" -result="pending" -url=""

litecli findajob.db
.tables


SCHEMA:
CREATE TABLE job_types (
    id integer primary key,
    type_nameTEXT Not null)                                                                                                                                                                                   |
CREATE TABLE "openings" (
    id integer primary key,
    firm TEXT Not null, type_job integer,
    result integer,
    application_date date default (date('now')), 
    url TEXT,

    foreign key (type_job) references job_types(id),
    foreign key (result) references results(id)
) 

CREATE TABLE results (
    id integer primary key,
    result_name TEXT Not null
    )
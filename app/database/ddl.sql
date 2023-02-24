create table if not exists sources
(
    id          INTEGER primary key autoincrement,
    label       varchar(256) default ''                not null,
    title       varchar(256) default ''                not null,
    author      varchar(256) default ''                not null,
    age         varchar(64)  default ''                not null,
    contents    text         default ''                not null,
    sort        int          default 0                 not null,
    create_time timestamp    default CURRENT_TIMESTAMP not null
);

create table if not exists words
(
    id                    INTEGER primary key autoincrement,
    source_id             int         default 0                 not null,
    word                  varchar(32) default ''                not null,
    word_length           tinyint     default 0                 not null,
    pinyin                varchar(32) default ''                not null,
    pinyin_tone           varchar(32) default ''                not null,
    last_pinyin           varchar(16) default ''                not null,
    last_pinyin_tone      varchar(16) default ''                not null,
    last_pinyin_tone_name tinyint     default 0                 not null,
    create_time           timestamp   default CURRENT_TIMESTAMP not null
);

create index if not exists source_id_index on words (source_id);
create index if not exists pinyin_index on words (word_length, last_pinyin, last_pinyin_tone_name);
create index if not exists pinyin_tone_index on words (word_length, last_pinyin_tone, last_pinyin_tone_name);

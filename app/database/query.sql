-- name: CreateWord :exec
insert into words (source_id, word, pinyin, pinyin_tone, last_pinyin, last_pinyin_tone, last_pinyin_tone_name,
                   word_length)
values (?, ?, ?, ?, ?, ?, ?, ?);

-- name: SearchRhymeWords :many
select w.*,
       sources.label,
       sources.contents
from (select words.word,
             words.source_id,
             words.pinyin,
             words.pinyin_tone,
             count(words.word) as use_count
      from words
      where word_length = ?
        and last_pinyin = ?
      group by word
      having use_count > 10
      limit ? offset ?) w
         join sources on sources.id = w.source_id
order by use_count desc;

-- name: GetSearchRhymeWordsCount :one
select count(1)
from (select count(words.word) as use_count
      from words
      where word_length = ?
        and last_pinyin = ?
      group by word
      having use_count > 10) w;

-- name: SearchRhymeToneWords :many
select w.*,
       sources.label,
       sources.contents
from (select words.word,
             words.source_id,
             words.pinyin,
             words.pinyin_tone,
             count(words.word) as use_count
      from words
      where word_length = ?
        and last_pinyin_tone = ?
      group by word
      having use_count > 10
      limit ? offset ?) w
         join sources on sources.id = w.source_id
order by use_count desc;

-- name: GetSearchRhymeToneWordsCount :one
select count(1)
from (select count(words.word) as use_count
      from words
      where word_length = ?
        and last_pinyin_tone = ?
      group by word
      having use_count > 10) w;

-- name: GetRandomWords :many
select w.*,
       sources.label,
       sources.contents
from (select words.word,
             words.source_id,
             words.pinyin,
             words.pinyin_tone,
             count(words.word) as use_count
      from words
      where word_length = ?
      group by word
      limit ? offset ?) w
         join sources on sources.id = w.source_id;

-- name: GetWordsCount :one
select count(1)
from (select count(1)
      from words
      where word_length = ?
      group by word) w;

-- name: CreateSource :one
insert into sources
    (label, age, title, author, contents, sort)
values (?, ?, ?, ?, ?, ?)
RETURNING id;

-- name: GetSources :many
select *
from sources
order by id desc;

-- name: GetSourcesByWord :many
select sources.*
from (select id,
             source_id
      from words
      where word = ?) w
         inner join sources on w.source_id = sources.id
limit ? offset ?;

-- name: GetSourcesByWordCount :one
select count(1)
from (select id, source_id
      from words
      where word = ?) w
         inner join sources on w.source_id = sources.id;

-- name: FindSource :one
select *
from sources
where id = ?;

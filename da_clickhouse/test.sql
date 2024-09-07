CREATE TABLE IF NOT EXISTS users (
    id UInt32,
    nickname String NOT NULL,
    email String NOT NULL,
    hash_password String NOT NULL,
    registered_at DateTime64 NOT NULL
) ENGINE = ReplacingMergeTree(id)
ORDER BY (registered_at);



INSERT INTO users 
SELECT * 
FROM postgresql('forum-db-1:5430', 'forumdb', 'users', 'postgres', 'some_password') as postgr
WHERE
(SELECT count(*) FROM users) = 0
    OR 
postgr.registered_at > (SELECT max(registered_at) FROM users)
AND postgr.id != (SELECT id FROM users WHERE registered_at = (SELECT max(registered_at) FROM users));




SELECT email from users

OPTIMIZE TABLE users 


SELECT max(registered_at) FROM users


drop table users 
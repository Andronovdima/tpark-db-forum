
SET SYNCHRONOUS_COMMIT = 'off';


DROP TRIGGER IF EXISTS on_post_insert ON posts;
DROP TRIGGER IF EXISTS on_thread_insert ON threads;

DROP FUNCTION IF EXISTS fn_update_thread_votes_ins();


CREATE TABLE IF NOT EXISTS users (
      nickname varchar primary key,
      about varchar,
      fullname varchar,
      email varchar not null unique
    );


CREATE TABLE IF NOT EXISTS forums (
     slug varchar not null primary key,
     posts bigint not null,
     threads bigint not null,
    title varchar not null,
    author varchar not null references users
    );


CREATE TABLE IF NOT EXISTS threads (
     id bigserial not null primary key,
    author varchar not null references users,
    created timestamptz DEFAULT now(),
    forum varchar not null references forums,
    message varchar not null,
    slug varchar,
    title varchar not null,
    votes int
);

CREATE TABLE IF NOT EXISTS posts (
    id bigserial not null primary key,
    author varchar not null references users,
    created timestamptz DEFAULT now(),
    forum varchar not null references forums,
    isEdited boolean not null DEFAULT false,
    message varchar not null,
    parent bigint not null DEFAULT NULL,
    thread bigint references threads,
    path bigint[]  NOT NULL DEFAULT '{0}'
);


CREATE TABLE IF NOT EXISTS votes (
    id bigserial not null primary key,
    nickname varchar references users,
    thread bigint references threads,
    voice int
);

CREATE TABLE IF NOT EXISTS forum_users
(
    forum    varchar not null references forums,
    nickname varchar not null references users
);

CREATE OR REPLACE FUNCTION forum_users_update()
    RETURNS TRIGGER AS '
    BEGIN
        INSERT INTO forum_users (forum, nickname) VALUES ((SELECT slug FROM forums WHERE LOWER(NEW.forum) = LOWER(slug)),
                                                          (SELECT users.nickname FROM users WHERE LOWER(NEW.author) = LOWER(nickname)));
        RETURN NULL;
    END;
' LANGUAGE plpgsql;

CREATE TRIGGER on_post_insert
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE forum_users_update();

CREATE TRIGGER on_thread_insert
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE forum_users_update();

SET SYNCHRONOUS_COMMIT = 'off';



DROP INDEX IF EXISTS idx_users_email_index;
DROP INDEX IF EXISTS idx_users_nickname_index;
DROP INDEX IF EXISTS idx_forums_slug_index;
DROP INDEX IF EXISTS idx_threads_slug;
DROP INDEX IF EXISTS idx_threads_forum;
DROP INDEX IF EXISTS idx_posts_forum;
DROP INDEX IF EXISTS idx_posts_parent;
DROP INDEX IF EXISTS idx_posts_path;
DROP INDEX IF EXISTS idx_posts_thread;
DROP INDEX IF EXISTS idx_posts_thread_id;

DROP TRIGGER IF EXISTS on_vote_insert ON votes;
DROP TRIGGER IF EXISTS on_vote_update ON votes;

DROP FUNCTION IF EXISTS fn_update_thread_votes_ins();

DROP TABLE IF EXISTS forum_users CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE ;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;

CREATE TABLE IF NOT EXISTS users (
    nickname varchar  COLLATE "POSIX" primary key ,
    about varchar,
    fullname varchar,
    email varchar not null unique
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_index
    ON users (LOWER(email));
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_nickname_index
    ON users (LOWER(nickname));

CREATE INDEX IF NOT EXISTS idx_users_pok
    ON users (nickname, email, fullname, about, LOWER(email), LOWER(nickname));

CREATE TABLE IF NOT EXISTS forums (
    slug varchar not null primary key,
    posts bigint not null,
    threads bigint not null,
    title varchar not null,
    author varchar not null references users
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_slug_index
    ON forums (LOWER(slug));
CREATE UNIQUE INDEX IF NOT EXISTS idx_forums_userNick_unique
    ON forums (LOWER(author));



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

CREATE INDEX IF NOT EXISTS idx_threads_slug
    ON threads (LOWER(slug));

CREATE INDEX IF NOT EXISTS idx_threads_forum
    ON threads (LOWER(forum));

CREATE INDEX IF NOT EXISTS idx_threads_pok
    ON threads (id, forum, author, slug, created, title, message, votes);

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

CREATE INDEX IF NOT EXISTS idx_posts_path ON posts USING GIN (path);
CREATE INDEX IF NOT EXISTS idx_posts_thread ON posts (thread);
CREATE INDEX IF NOT EXISTS idx_posts_forum ON posts (forum);
CREATE INDEX IF NOT EXISTS idx_posts_parent ON posts (parent);
CREATE INDEX IF NOT EXISTS idx_posts_thread_id ON posts (thread, id);
CREATE INDEX IF NOT EXISTS idx_posts_pok
    ON posts (id, parent, thread, forum, author, created, message, isedited, path);


CREATE TABLE IF NOT EXISTS votes (
    id bigserial not null primary key,
    nickname varchar references users,
    thread bigint references threads,
    voice smallint
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_votes_nickname_thread_unique
    ON votes (LOWER(nickname), thread);


CREATE TABLE IF NOT EXISTS forum_users
(
    forum    varchar not null references forums,
    nickname varchar not null references users
);

CREATE INDEX idx_forum_users_user_id
    ON forum_users(forum);

CREATE INDEX idx_forum_users_forum_id
    ON forum_users(nickname);

CREATE INDEX idx_forum_users_user_id_forum_id
    ON forum_users (forum, nickname);


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
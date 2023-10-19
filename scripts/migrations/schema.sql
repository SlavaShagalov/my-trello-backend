CREATE TABLE IF NOT EXISTS users
(
    id              serial    NOT NULL PRIMARY KEY,
    username        text      NOT NULL UNIQUE,
    hashed_password bytea     NOT NULL,
    email           varchar   NOT NULL,
    name            varchar   NOT NULL,
    avatar          varchar   NULL,
    created_at      timestamp NOT NULL DEFAULT now(),
    updated_at      timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS workspaces
(
    id          serial    NOT NULL PRIMARY KEY,
    user_id     int       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title       varchar   NOT NULL,
    description varchar   NULL,
    created_at  timestamp NOT NULL DEFAULT now(),
    updated_at  timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS boards
(
    id           serial    NOT NULL PRIMARY KEY,
    workspace_id int       NOT NULL REFERENCES workspaces (id) ON DELETE CASCADE,
    title        varchar   NOT NULL DEFAULT '',
    description  varchar   NULL,
    background   varchar   NULL,
    created_at   timestamp NOT NULL DEFAULT now(),
    updated_at   timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS lists
(
    id         serial    NOT NULL PRIMARY KEY,
    board_id   int       NOT NULL REFERENCES boards (id) ON DELETE CASCADE,
    title      varchar   NOT NULL DEFAULT '',
    position   int       NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS cards
(
    id         serial    NOT NULL PRIMARY KEY,
    list_id    int       NOT NULL REFERENCES lists (id) ON DELETE CASCADE,
    title      varchar   NOT NULL DEFAULT '',
    content    varchar   NULL,
    position   int       NOT NULL,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

-- Update positions after list was deleted
CREATE OR REPLACE FUNCTION on_list_delete() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE lists
    SET position = position - 1
    WHERE board_id = old.board_id
      AND position > old.position;

    RETURN NULL;
END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER list_delete
    AFTER DELETE
    ON lists
    FOR EACH ROW
EXECUTE PROCEDURE on_list_delete();

CREATE OR REPLACE PROCEDURE update_list_positions(new_position int, list_id int) AS
$$
DECLARE
    old_position int;
    boardID      int;
BEGIN
    SELECT l.position, l.board_id
    INTO old_position, boardID
    FROM lists l
    WHERE id = list_id;

    IF new_position > old_position THEN
        UPDATE lists l
        SET position = position - 1
        WHERE l.board_id = boardID
          AND position > old_position
          AND position <= new_position;
        --         RAISE NOTICE 'old_position: %', old_position;
--         RAISE NOTICE 'new_position: %', new_position;
    ELSIF new_position < old_position THEN
        UPDATE lists l
        SET position = position + 1
        WHERE l.board_id = boardID
          AND position >= new_position
          AND position < old_position;
    END IF;
END
$$ LANGUAGE plpgsql;

-- Update positions after card was deleted
CREATE OR REPLACE FUNCTION on_card_delete() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE cards
    SET position = position - 1
    WHERE list_id = old.list_id
      AND position > old.position;

    RETURN NULL;
END
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER card_delete
    AFTER DELETE
    ON cards
    FOR EACH ROW
EXECUTE PROCEDURE on_card_delete();

-- Update positions after card position was updated
CREATE OR REPLACE PROCEDURE update_cards_positions(new_position int, card_id int) AS
$$
DECLARE
    old_position int;
    listID       int;
BEGIN
    SELECT c.position, c.list_id
    INTO old_position, listID
    FROM cards c
    WHERE id = card_id;

    IF new_position > old_position THEN
        UPDATE cards c
        SET position = position - 1
        WHERE c.list_id = listID
          AND position > old_position
          AND position <= new_position;
    ELSIF new_position < old_position THEN
        UPDATE cards c
        SET position = position + 1
        WHERE c.list_id = listID
          AND position >= new_position
          AND position < old_position;
    END IF;
END
$$ LANGUAGE plpgsql;

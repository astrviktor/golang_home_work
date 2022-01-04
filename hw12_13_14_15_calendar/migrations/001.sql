CREATE DATABASE calendar;
CREATE USER "user" WITH ENCRYPTED PASSWORD 'password123';
GRANT ALL PRIVILEGES ON DATABASE calendar TO "user";

CREATE SCHEMA calendar;

CREATE TABLE calendar.event (
    id uuid,
    title text NOT NULL,
    date_start timestamp with time zone NOT NULL,
    date_end timestamp with time zone NOT NULL,
    description text NOT NULL,
    user_id integer,
    time_to_notification integer
);

COMMENT ON TABLE calendar.event IS 'События календаря';

COMMENT ON COLUMN calendar.event.id IS 'Уникальный идентификатор события';
COMMENT ON COLUMN calendar.event.title IS 'Заголовок - короткий';
COMMENT ON COLUMN calendar.event.date_start IS 'Дата и время старта события';
COMMENT ON COLUMN calendar.event.date_end IS 'Дата и время завершения события';
COMMENT ON COLUMN calendar.event.description IS 'Описание события - длинный текст';
COMMENT ON COLUMN calendar.event.user_id IS 'ID пользователя, владельца события';
COMMENT ON COLUMN calendar.event.time_to_notification IS 'За сколько минут высылать уведомление';

CREATE INDEX event_id_index ON calendar.event (id);
CREATE INDEX event_title_index ON calendar.event (title);
CREATE INDEX event_date_start_index ON calendar.event (date_start);
CREATE INDEX event_date_end_index ON calendar.event (date_end);
CREATE INDEX event_description_index ON calendar.event (description);
CREATE INDEX event_user_id_index ON calendar.event (user_id);
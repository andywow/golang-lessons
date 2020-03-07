create table "event" (
  uuid serial primary key,
  start_time timestamp not null,
  duration integer not null,
  header text not null,
  description text not null,
  username text not null,
  notification_period integer not null
);

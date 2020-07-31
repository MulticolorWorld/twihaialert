create table if not exists user
(
    id                int auto_increment primary key,
    created_at        datetime not null,
    last_login        datetime not null,
    last_notify       datetime null,
    last_notify_count int      not null,
    dm_notification   int      not null
);
create table if not exists twitter_account
(
    id                  int auto_increment primary key,
    user_id             int         not null,
    twitter_id          bigint      not null unique,
    screen_name         varchar(50) not null unique,
    access_token        varchar(50) not null,
    access_token_secret varchar(50) not null,
    created_at          datetime    not null,
    updated_at          datetime    not null,
    last_tweet_id       bigint      not null,
    daily_count         int         not null,
    daily_count_rt      int         not null,
    count_update        datetime    null,
    foreign key fk_user_id (user_id) references user (id)
)
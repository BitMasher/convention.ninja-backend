create table if not exists asset_categories
(
    id   bigint auto_increment
    primary key,
    name varchar(256) not null,
    constraint asset_categories_id_uindex
    unique (id)
    );

create table if not exists asset_manufacturers
(
    id   bigint auto_increment
    primary key,
    name varchar(256) not null,
    constraint asset_manufacturers_id_uindex
    unique (id)
    );

create table if not exists asset_models
(
    id              bigint auto_increment
    primary key,
    name            varchar(256) not null,
    manufacturer_id bigint       null,
    category_id     bigint       not null,
    constraint asset_models_id_uindex
    unique (id),
    constraint asset_models_asset_categories_id_fk
    foreign key (category_id) references asset_categories (id),
    constraint asset_models_asset_manufacturers_id_fk
    foreign key (manufacturer_id) references asset_manufacturers (id)
    );

create table if not exists asset_tags
(
    id       bigint auto_increment
    primary key,
    asset_id bigint       not null,
    tag_code varchar(512) null,
    constraint asset_tags_id_uindex
    unique (id)
    );

create table if not exists users
(
    id           bigint       not null
    primary key,
    name         varchar(256) null,
    display_name varchar(32)  null,
    constraint users_id_uindex
    unique (id)
    );

create table if not exists organizations
(
    id       bigint       not null
    primary key,
    name     varchar(256) not null,
    owner_id bigint       not null,
    constraint organizations_id_uindex
    unique (id),
    constraint organizations_users_id_fk
    foreign key (owner_id) references users (id)
    );

create table if not exists assets
(
    id              bigint       not null
    primary key,
    model_id        bigint       not null,
    serial_number   varchar(512) null,
    organization_id bigint       not null,
    constraint assets_id_uindex
    unique (id),
    constraint assets_asset_models_id_fk
    foreign key (model_id) references asset_models (id),
    constraint assets_organizations_id_fk
    foreign key (organization_id) references organizations (id)
    );


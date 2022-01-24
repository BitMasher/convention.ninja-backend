create table ods.users
(
    id             bigint primary key,
    created_at     datetime(3) null,
    updated_at     datetime(3) null,
    deleted_at     datetime(3) null,
    name           longtext null,
    display_name   longtext null,
    email          longtext null,
    email_verified tinyint(1) null,
    firebase_id    longtext null
);

create table ods.organizations
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    name            longtext null,
    normalized_name longtext null,
    owner_id        bigint null,
    constraint fk_organizations_owner
        foreign key (owner_id) references ods.users (id)
);

create table ods.categories
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    name            longtext null,
    organization_id bigint null,
    constraint fk_categories_organization
        foreign key (organization_id) references ods.organizations (id)
);

create index idx_categories_deleted_at
    on ods.categories (deleted_at);

create table ods.manufacturers
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    name            longtext null,
    organization_id bigint null,
    constraint fk_manufacturers_organization
        foreign key (organization_id) references ods.organizations (id)
);

create index idx_manufacturers_deleted_at
    on ods.manufacturers (deleted_at);

create table ods.models
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    name            longtext null,
    manufacturer_id bigint null,
    category_id     bigint null,
    organization_id bigint null,
    constraint fk_models_category
        foreign key (category_id) references ods.categories (id),
    constraint fk_models_manufacturer
        foreign key (manufacturer_id) references ods.manufacturers (id),
    constraint fk_models_organization
        foreign key (organization_id) references ods.organizations (id)
);

create table ods.assets
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    model_id        bigint null,
    serial_number   longtext null,
    organization_id bigint null,
    constraint fk_assets_model
        foreign key (model_id) references ods.models (id),
    constraint fk_assets_organization
        foreign key (organization_id) references ods.organizations (id)
);

create table ods.asset_tags
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    tag_id          longtext null,
    asset_id        bigint null,
    organization_id bigint null,
    constraint fk_assets_asset_tags
        foreign key (asset_id) references ods.assets (id)
);

create index idx_asset_tags_deleted_at
    on ods.asset_tags (deleted_at);

create index idx_assets_deleted_at
    on ods.assets (deleted_at);

create index idx_models_deleted_at
    on ods.models (deleted_at);

create index idx_organizations_deleted_at
    on ods.organizations (deleted_at);

create index idx_users_deleted_at
    on ods.users (deleted_at);


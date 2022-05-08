create table ods.venues
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    name            longtext    null,
    address         longtext    null,
    organization_id bigint,
    constraint fk_venues_organization
        foreign key (organization_id) references ods.organizations (id)
);

create table ods.venue_rooms
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    name            longtext    null,
    venue_id        bigint,
    organization_id bigint,
    constraint fk_venue_rooms_organization
        foreign key (organization_id) references ods.organizations (id),
    constraint fk_venue_rooms_venue
        foreign key (venue_id) references ods.venues (id)
);

create table ods.events
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    name            longtext    null,
    organization_id bigint,
    constraint fk_events_organization
        foreign key (organization_id) references ods.organizations (id)
);

create table ods.event_schedules
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    start_date      datetime(3) null,
    end_date        datetime(3) null,
    event_id        bigint,
    organization_id bigint,
    constraint fk_event_schedules_organization
        foreign key (organization_id) references ods.organizations (id),
    constraint fk_event_schedules_event
        foreign key (event_id) references ods.events (id)
);

create table ods.event_schedule_entries
(
    id                bigint primary key,
    created_at        datetime(3) null,
    updated_at        datetime(3) null,
    deleted_at        datetime(3) null,
    name              longtext    null,
    start_date        datetime(3) null,
    end_date          datetime(3) null,
    venue_room_id     bigint,
    event_schedule_id bigint,
    organization_id   bigint,
    constraint fk_event_schedule_entries_organization
        foreign key (organization_id) references ods.organizations (id),
    constraint fk_event_schedule_entries_event_schedule
        foreign key (event_schedule_id) references ods.event_schedules (id),
    constraint fk_event_schedule_entries_venue_room
        foreign key (venue_room_id) references ods.venue_rooms (id)
);

create table ods.manifests
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    ship_date       datetime(3) null,
    creator_id      bigint,
    organization_id bigint,
    constraint fk_manifests_organization
        foreign key (organization_id) references ods.organizations (id),
    constraint fk_manifests_user
        foreign key (creator_id) references ods.users (id)
);

create table ods.manifest_entries
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    manifest_id     bigint,
    asset_id        bigint,
    organization_id bigint,
    constraint fk_manifest_entries_manifest
        foreign key (manifest_id) references ods.manifests (id),
    constraint fk_manifest_entries_asset
        foreign key (asset_id) references ods.assets (id),
    constraint fk_manifest_entries_organization
        foreign key (organization_id) references ods.organizations (id)
);

create index idx_venues_deleted_at
    on ods.venues (deleted_at);

create index idx_venue_rooms_deleted_at
    on ods.venue_rooms (deleted_at);

create index idx_events_deleted_at
    on ods.events (deleted_at);

create index idx_event_schedules_deleted_at
    on ods.event_schedules (deleted_at);

create index idx_event_schedule_entries_deleted_at
    on ods.event_schedule_entries (deleted_at);

create index idx_manifests_deleted_at
    on ods.manifests (deleted_at);

create index idx_manifest_entries_deleted_at
    on ods.manifest_entries (deleted_at);
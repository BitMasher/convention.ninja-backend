alter table ods.manifests
    add column room_id bigint null after deleted_at;

alter table ods.manifests
    add constraint fk_manifests_venue_room
        foreign key (room_id) references ods.venue_rooms (id);

alter table ods.manifests
    add column responsible_user_id bigint null after room_id;

alter table ods.manifests
    add constraint fk_manifests_responsible_user
        foreign key (responsible_user_id) references ods.users (id);

alter table ods.manifests
    add column responsible_external_party json null after responsible_user_id;

create table ods.external_party_identifiers
(
    id              bigint primary key,
    created_at      datetime(3) null,
    updated_at      datetime(3) null,
    deleted_at      datetime(3) null,
    fields          longtext,
    organization_id bigint,
    constraint fk_external_party_identifiers_organization
        foreign key (organization_id) references ods.organizations (id)
)
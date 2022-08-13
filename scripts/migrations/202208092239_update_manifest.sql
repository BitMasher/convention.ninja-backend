alter table ods.manifests
    drop foreign key fk_manifests_venue_room;

alter table ods.manifests
    modify room_id varchar(255) null;
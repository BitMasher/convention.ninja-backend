alter table ods.assets
drop foreign key assets_venue_rooms_id_fk;

drop index assets_venue_rooms_id_fk on ods.assets;

alter table ods.assets
    modify room_id varchar(255) null;


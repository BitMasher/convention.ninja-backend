alter table ods.assets
    add room_id bigint null after serial_number;

alter table ods.assets
    add constraint assets_venue_rooms_id_fk
        foreign key (room_id) references ods.venue_rooms (id);


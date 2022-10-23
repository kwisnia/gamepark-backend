insert into age_rating_organizations
    (name, id)
values ('ESRB', 1),
       ('PEGI', 2),
       ('CERO', 3),
       ('USK', 4),
       ('GRAC', 5),
       ('CLASS_IND', 6),
       ('ACB', 7);

insert into age_ratings
    (name, id)
values ('Three', 1),
       ('Seven', 2),
       ('Twelve', 3),
       ('Sixteen', 4),
       ('Eighteen', 5),
       ('RP', 6),
       ('EC', 7),
       ('E', 8),
       ('E10', 9),
       ('T', 10),
       ('M', 11),
       ('AO', 12),
       ('CERO_A', 13),
       ('CERO_B', 14),
       ('CERO_C', 15),
       ('CERO_D', 16),
       ('CERO_Z', 17),
       ('USK_0', 18),
       ('USK_6', 19),
       ('USK_12', 20),
       ('USK_16', 21),
       ('USK_18', 22),
       ('GRAC_ALL', 23),
       ('GRAC_Twelve', 24),
       ('GRAC_Fifteen', 25),
       ('GRAC_Eighteen', 26),
       ('GRAC_TESTING', 27),
       ('CLASS_IND_L', 28),
       ('CLASS_IND_Ten', 29),
       ('CLASS_IND_Twelve', 30),
       ('CLASS_IND_Fourteen', 31),
       ('CLASS_IND_Sixteen', 32),
       ('CLASS_IND_Eighteen', 33),
       ('ACB_G', 34),
       ('ACB_PG', 35),
       ('ACB_M', 36),
       ('ACB_MA15', 37),
       ('ACB_R18', 38),
       ('ACB_RC', 39);

insert into game_categories
    (name, id)
values ('main_game', 0),
       ('dlc_addon', 1),
       ('expansion', 2),
       ('bundle', 3),
       ('standalone_expansion', 4),
       ('mod', 5),
       ('episode', 6),
       ('season', 7),
       ('remake', 8),
       ('remaster', 9),
       ('expanded_game', 10),
       ('port', 11),
       ('fork', 12);

insert into external_categories
    (name, id)
values ('steam', 1),
       ('gog', 5),
       ('youtube', 10),
       ('microsoft', 11),
       ('apple', 13),
       ('twitch', 14),
       ('android', 15),
       ('amazon_asin', 20),
       ('amazon_luna', 22),
       ('amazon_adg', 23),
       ('epic_game_store', 26),
       ('oculus', 28),
       ('utomik', 29),
       ('itch_io', 30),
       ('xbox_marketplace', 31),
       ('kartridge', 32),
       ('playstation_store_us', 36),
       ('focus_entertainment', 37);

insert into release_date_categories
    (name, id)
values ('YYYYMMMMDD', 0),
       ('YYYYMMMM', 1),
       ('YYYY', 2),
       ('YYYYQ1', 3),
       ('YYYYQ2', 4),
       ('YYYYQ3', 5),
       ('YYYYQ4', 6),
       ('TBD', 7);

insert into release_regions
    (name, id)
values ('europe', 1),
       ('north_america', 2),
       ('australia', 3),
       ('new_zealand', 4),
       ('japan', 5),
       ('china', 6),
       ('asia', 7),
       ('worldwide', 8),
       ('korea', 9),
       ('brazil', 10);

insert into game_completions
    (name, id)
values ('main_story', 0),
       ('main_and_extras', 1),
       ('completionist', 2),
       ('dropped', 3);
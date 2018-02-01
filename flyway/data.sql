-- Adminer 4.5.0 PostgreSQL dump

INSERT INTO "authority" ("id", "version", "name") VALUES
(1,	'0',	'ADMINISTRATOR');

INSERT INTO "authority_event" ("authority_id", "event_id") VALUES
(1,	3),
(1,	5);

INSERT INTO "event" ("id", "created_by", "created_at", "type", "description") VALUES
(1,	1,	'2018-02-01 10:11:26.902248+00',	'CREATE_USER',	'User "erik@eies.land" created.'),
(2,	1,	'2018-02-01 10:24:56.896291+00',	'CREATE_ROLE',	'Role "ADMINISTRATOR" created.'),
(3,	1,	'2018-02-01 10:25:40.598196+00',	'CREATE_AUTHORITY',	'Authority "ADMINISTRATOR" created.'),
(4,	1,	'2018-02-01 10:30:37.329862+00',	'USER_ROLE_ASSIGNED',	'User "erik@eies.land" assigned role "ADMINISTRATOR".'),
(5,	1,	'2018-02-01 10:34:08.762837+00',	'ROLE_AUTHORITY_ASSIGNED',	'Role "ADMINISTRATOR" was assigned authority "ADMINISTRATOR".');

INSERT INTO "event_type" ("name", "description") VALUES
('ROLE_AUTHORITY_ASSIGNED',	'A role was assigned an authority.'),
('USER_ROLE_ASSIGNED',	'A user was assigned a role.'),
('CREATE_USER',	'A new user was created.'),
('CREATE_AUTHORITY',	'A new authority was created.'),
('CREATE_ROLE',	'A new role was created.');

INSERT INTO "role" ("id", "version", "name") VALUES
(1,	'0',	'ADMINISTRATOR');

INSERT INTO "role_authority" ("role_id", "authority_id") VALUES
(1,	1);

INSERT INTO "role_event" ("role_id", "event_id") VALUES
(1,	2),
(1,	4),
(1,	5);

INSERT INTO "user" ("id", "version", "first_name", "last_name", "email", "password", "enabled", "deleted") VALUES
(1,	'0',	'Erik',	'Eiesland',	'erik@eies.land',	'$2y$10$rW.o8LqVIyXjnDtEQ.eLvurFs7BPfdqIZexZPaFgTMn.Qms6HMvyi',	'1',	'');

INSERT INTO "user_event" ("user_id", "event_id") VALUES
(1,	1),
(1,	4);

INSERT INTO "user_role" ("user_id", "role_id") VALUES
(1,	1);

-- 2018-02-01 10:39:47.950554+00
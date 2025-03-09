INSERT INTO `roles` (created_at, name)
VALUES (CURRENT_TIMESTAMP, 'superadmin');
INSERT INTO `roles` (created_at, name)
VALUES (CURRENT_TIMESTAMP, 'admin');

INSERT INTO `domain_roles` (created_at, name)
VALUES (CURRENT_TIMESTAMP, 'admin');
INSERT INTO `domain_roles` (created_at, name)
VALUES (CURRENT_TIMESTAMP, 'record_admin');
INSERT INTO `domain_roles` (created_at, name)
VALUES (CURRENT_TIMESTAMP, 'readonly');

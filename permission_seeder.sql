USE erajaya_be_tech_test;

-- User / Admin Permissions
DELETE FROM
	user_permissions
WHERE
	user_id = "3f7b39b5-b647-4c3d-8a2e-aa4e6ec31801";

INSERT INTO
	user_permissions (id, user_id, permission_id)
SELECT
	UUID(),
	"3f7b39b5-b647-4c3d-8a2e-aa4e6ec31801",
	id
FROM
	permissions
WHERE permissions.package = "WebsiteAdmin";

CREATE ROLE reader;

ALTER ROLE reader PASSWORD '1111';
ALTER ROLE reader LOGIN;

GRANT SELECT ON users TO reader;
GRANT SELECT ON workspaces TO reader;
GRANT SELECT ON boards TO reader;
GRANT SELECT ON lists TO reader;
GRANT SELECT ON cards TO reader;

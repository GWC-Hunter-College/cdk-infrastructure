/* put the query of:

GET /me/clubs/eboard

in this folder */ 

-- the ? denotes the current user's id

SELECT ?, club_id FROM club_members cm
INNER JOIN clubs ON ? = clubs.id
WHERE cm.is_eboard = TRUE OR cm.is_owner = TRUE


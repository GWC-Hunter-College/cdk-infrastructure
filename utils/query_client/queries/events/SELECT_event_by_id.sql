SELECT 
	ec.fk_club_id AS club_id, 
	ec.club_is_event_owner AS is_owner,
    i.object_key AS object_key,
	e.*
FROM events e
INNER JOIN events_to_clubs ec ON e.id = ec.fk_event_id
INNER JOIN clubs c ON ec.fk_club_id = c.id
LEFT JOIN images i ON c.fk_logo_id = i.id
WHERE status LIKE ?
    AND e.id = ?;
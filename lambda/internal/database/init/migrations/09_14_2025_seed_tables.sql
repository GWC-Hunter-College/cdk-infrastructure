-- seed students
INSERT INTO students (id, email) VALUES
  ('11111111-1111-1111-1111-111111111111', 'kelly@example.com'),
  ('22222222-2222-2222-2222-222222222222', 'kyle@example.com'),
  ('33333333-3333-3333-3333-333333333333', 'shohruz@example.com'),
  ('44444444-4444-4444-4444-444444444444', 'anthony@example.com'),
  ('55555555-5555-5555-5555-555555555555', 'maria@example.com'),
  ('66666666-6666-6666-6666-666666666666', 'james@example.com'),
  ('77777777-7777-7777-7777-777777777777', 'sophia@example.com'),
  ('88888888-8888-8888-8888-888888888888', 'li@example.com'),
  ('99999999-9999-9999-9999-999999999999', 'david@example.com'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'fatima@example.com');

-- seed student_info (FK to students.id)
INSERT INTO student_info (fk_student_id, username, first_name, last_name) VALUES
  ('11111111-1111-1111-1111-111111111111', 'kelly', 'Kelly', 'Smith'),
  ('22222222-2222-2222-2222-222222222222', 'kyle', 'Kyle', 'Jones'),
  ('33333333-3333-3333-3333-333333333333', 'shohruz', 'Shohruz', 'Aliyev'),
  ('44444444-4444-4444-4444-444444444444', 'anth', 'Anthony', 'Brown'),
  ('55555555-5555-5555-5555-555555555555', 'maria', 'Maria', 'Lopez'),
  ('66666666-6666-6666-6666-666666666666', 'james', 'James', 'Smith'),
  ('77777777-7777-7777-7777-777777777777', 'soph', 'Sophia', 'Wang'),
  ('88888888-8888-8888-8888-888888888888', 'li', 'Li', 'Zhang'),
  ('99999999-9999-9999-9999-999999999999', 'david', 'David', 'Kim'),
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'fat', 'Fatima', 'Hussein');

-- create cooking club (no logo yet)
INSERT INTO clubs (name, fk_logo_id) VALUES ('Cooking Club', NULL);
INSERT INTO club_info (fk_club_id, website_url)
SELECT id, 'https://www.example.com/cooking' FROM clubs WHERE name='Cooking Club';

-- add maria as owner of cooking club
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '55555555-5555-5555-5555-555555555555', id, TRUE, TRUE FROM clubs WHERE name='Cooking Club';

-- add anthony, james, kelly as members
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '44444444-4444-4444-4444-444444444444', id, FALSE, FALSE FROM clubs WHERE name='Cooking Club';
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '66666666-6666-6666-6666-666666666666', id, TRUE, FALSE FROM clubs WHERE name='Cooking Club';
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '11111111-1111-1111-1111-111111111111', id, FALSE, FALSE FROM clubs WHERE name='Cooking Club';

-- make GWC club
INSERT INTO clubs (name, fk_logo_id) VALUES ('Girls Who Code @ Hunter', NULL);
INSERT INTO club_info (fk_club_id, website_url)
SELECT id, 'https://www.example.com/gwc' FROM clubs WHERE name='Girls Who Code @ Hunter';

-- assign kelly as owner
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '11111111-1111-1111-1111-111111111111', id, TRUE, TRUE FROM clubs WHERE name='Girls Who Code @ Hunter';

-- add sophia and anthony, anthony becomes eboard
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '77777777-7777-7777-7777-777777777777', id, FALSE, FALSE FROM clubs WHERE name='Girls Who Code @ Hunter';
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '44444444-4444-4444-4444-444444444444', id, TRUE, FALSE FROM clubs WHERE name='Girls Who Code @ Hunter';

-- add kyle and li
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '22222222-2222-2222-2222-222222222222', id, TRUE, FALSE FROM clubs WHERE name='Girls Who Code @ Hunter';
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '88888888-8888-8888-8888-888888888888', id, TRUE, FALSE FROM clubs WHERE name='Girls Who Code @ Hunter';

-- create hunter cs club
INSERT INTO clubs (name, fk_logo_id) VALUES ('Hunter CS Club', NULL);
INSERT INTO club_info (fk_club_id, website_url)
SELECT id, 'https://www.example.com/huntercs' FROM clubs WHERE name='Hunter CS Club';

-- shohruz is owner
INSERT INTO club_members (fk_student_id, fk_club_id, member_is_eboard, member_is_owner)
SELECT '33333333-3333-3333-3333-333333333333', id, TRUE, TRUE FROM clubs WHERE name='Hunter CS Club';


INSERT INTO events (
  fk_author_id, fk_thumbnail_id, title, location, rsvp_link, status,
  start_date, end_date, timezone, created_at, updated_at
)
SELECT '11111111-1111-1111-1111-111111111111',
       NULL,
       'First Kickoff Meeting','Room 101','some link','posted',
       '2025-09-15 17:00:00','2025-09-15 19:00:00','America/New_York',NOW(),NOW();

SET @event_id = LAST_INSERT_ID();

INSERT INTO event_descriptions (fk_event_id, description)
VALUES(@event_id,'meow');

INSERT INTO events_to_clubs (fk_event_id,fk_club_id,club_is_event_owner)
SELECT @event_id,id,TRUE FROM clubs WHERE name='Girls Who Code @ Hunter';

-- verify clubs
INSERT INTO verified_clubs (fk_club_id) SELECT id FROM clubs WHERE name='Hunter CS Club';
INSERT INTO verified_clubs (fk_club_id) SELECT id FROM clubs WHERE name='Girls Who Code @ Hunter';

-- admins
INSERT INTO admins (fk_student_id) VALUES
  ('22222222-2222-2222-2222-222222222222'), -- kyle
  ('44444444-4444-4444-4444-444444444444'); -- anthony
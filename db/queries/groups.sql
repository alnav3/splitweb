-- name: GroupsCountByUserId :one
SELECT count(id)
FROM group_members
WHERE user_id = $1;

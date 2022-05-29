package repo

import (
	"gorm.io/gorm"
	"presentio-server-posts/src/v0/models"
)

type TagsRepo struct {
	db *gorm.DB
}

func CreateTagsRepo(db *gorm.DB) TagsRepo {
	return TagsRepo{
		db,
	}
}

func (r *TagsRepo) BulkInsert(tags []models.Tag) error {
	sql := "INSERT INTO tags (name) VALUES "
	first := true
	erasedArgs := make([]interface{}, 0, len(tags))

	for i := 0; i < len(tags); i++ {
		if !first {
			sql += ", "
		}

		sql += "(?)"
		first = false
		erasedArgs = append(erasedArgs, tags[i].Name)
	}

	sql += " ON CONFLICT DO NOTHING"

	return r.db.Exec(sql, erasedArgs...).Error
}

func (r *TagsRepo) BulkInsertRelation(tags []models.Tag, postId int64) error {
	sql := "INSERT INTO post_tags (post_id, tag_id) VALUES "
	erasedArgs := make([]interface{}, 0, len(tags)*2)
	first := true

	for i := 0; i < len(tags); i++ {
		if !first {
			sql += ","
		}

		first = false
		sql += "(?, (SELECT id FROM tags WHERE name = ?))"
		erasedArgs = append(erasedArgs, postId, tags[i].Name)
	}

	return r.db.Exec(sql, erasedArgs...).Error
}

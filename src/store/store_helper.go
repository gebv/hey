package store

import (
	"github.com/golang/glog"
	"models"
	"fmt"
	"github.com/jackc/pgx"
)

func FindModel(model models.ModelAbstractInterface, db *pgx.Conn, tx *pgx.Tx, fields ...string) error {
	fieldNames, fieldValues := model.Fields(fields...)

	query := SqlSelect(model.TableName(), fieldNames)
	

	where := fmt.Sprintf(" WHERE %s = ?", model.PrimaryName())

	query = FormateToPQuery(query+where)
	args := []interface{}{model.PrimaryValue()}

	var err error

	if tx != nil {
		err = tx.QueryRow(query, args...).Scan(fieldValues...)
	} else {
		err = db.QueryRow(query, args...).Scan(fieldValues...)
	}

	if err == pgx.ErrNoRows {

		return models.ErrNotFound
	}

	if err != nil {

		return err
	}

	return nil
}


func CreateModel(model models.ModelAbstractInterface, db *pgx.Conn, tx *pgx.Tx, fields ...string) (err error) {
	if len(model.PrimaryValue()) == 0 {
		glog.Errorf("Create '%T'. Primary key = nil", model)
		return models.ErrNotValid
	}

	_fields := models.StringArray{}
	_fields.AddAsArray(fields)
	_fields.Del("removed_at")
	_fields.Add("created_at")
	_fields.Add("updated_at")
	_fields.Add(model.PrimaryName())

	return updateOrCreateModel(model, db, tx, true, false, _fields...)
}

func UpdateModel(model models.ModelAbstractInterface, db *pgx.Conn, tx *pgx.Tx, fields ...string) (err error) {
	if len(model.PrimaryValue()) == 0 {
		glog.Errorf("Update '%T'. Primary key = nil", model)
		return models.ErrNotValid
	}

	_fields := models.StringArray{}
	_fields.AddAsArray(fields)
	_fields.Del(model.PrimaryName())
	_fields.Del("removed_at")
	_fields.Del("created_at")
	_fields.Add("updated_at")

	return updateOrCreateModel(model, db, tx, false, false, _fields...)
}

func DeleteModel(model models.ModelAbstractInterface, db *pgx.Conn, tx *pgx.Tx) (err error) {
	if len(model.PrimaryValue()) == 0 {
		glog.Errorf("Delete '%T'. Primary key = nil", model)
		return models.ErrNotValid
	}

	return updateOrCreateModel(model, db, tx, false, true, "updated_at", "removed_at", "is_removed")
}

func updateOrCreateModel(model models.ModelAbstractInterface, db *pgx.Conn, tx *pgx.Tx, isNew, isRemove bool, fields ...string) error {
	if isNew {
		model.BeforeCreate()
	}

	if isRemove {
		model.BeforeDelete()
	}

	model.BeforeSave()

	fieldNames, fieldValues := model.Fields(fields...)

	var query string = SqlUpdate(model.TableName(), fieldNames)

	if isNew {
		query = SqlInsert(model.TableName(), fieldNames)
	}

	where := fmt.Sprintf(" WHERE %[1]s = ? RETURNING %[1]s", model.PrimaryName())
	if isNew {
		where = fmt.Sprintf(" RETURNING %s", model.PrimaryName())
	} else {
		fieldValues = append(fieldValues, model.PrimaryValue())
	}

	query = FormateToPQuery(query+where)

	var err error

	if tx != nil {
		err = tx.QueryRow(query, fieldValues...).Scan(model.Maps()[model.PrimaryName()])

	} else {
		err = db.QueryRow(query, fieldValues...).Scan(model.Maps()[model.PrimaryName()])
	}

	if err != nil {

		return err
	}

	return nil
}
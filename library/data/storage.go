package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/appcontext"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ErrNotEnough declare specific error for Not Enough
// ErrExisted declare specific error for data already exist
var (
	ErrNotFound     = fmt.Errorf("data is not found")
	ErrAlreadyExist = fmt.Errorf("data already exists")
)

// GenericStorage represents the generic Storage
// for the domain models that matches with its database models
type GenericStorage interface {
	Single(ctx *gin.Context, elem interface{}, where string, arg map[string]interface{}) error
	Where(ctx *gin.Context, elems interface{}, where string, arg map[string]interface{}) error
	SinglePOSTEMP(ctx *gin.Context, elem interface{}, where string, arg map[string]interface{}) error
	WherePOSTEMP(ctx *gin.Context, elems interface{}, where string, arg map[string]interface{}) error
	SelectWithQuery(ctx *gin.Context, elem interface{}, query string, args map[string]interface{}) error
	FindByID(ctx *gin.Context, elem interface{}, id interface{}) error
	FindAll(ctx *gin.Context, elems interface{}, page int, limit int, isAsc bool) error
	Insert(ctx *gin.Context, elem interface{}) (*sql.Result, error)
	InsertNoTrail(ctx *gin.Context, elem interface{}) (*sql.Result, error)
	InsertMany(ctx *gin.Context, elem interface{}) error
	InsertManyWithTime(ctx *gin.Context, elem interface{}, created_at time.Time) error
	Update(ctx *gin.Context, elem interface{}) error
	UpdateNoTrail(ctx *gin.Context, elem interface{}) error
	UpdateMany(ctx *gin.Context, elems interface{}) error
	Delete(ctx *gin.Context, id interface{}) error
	DeleteMany(ctx *gin.Context, ids interface{}) error
	CountAll(ctx *gin.Context, count interface{}) error
	HardDelete(ctx *gin.Context, id interface{}) error
	ExecQuery(ctx *gin.Context, query string, args map[string]interface{}) error
	SelectFirstWithQuery(ctx *gin.Context, elem interface{}, query string, args map[string]interface{}) error
	InsertTrail(ctx *gin.Context, id string) (*sql.Result, error)
	UpdateTrail(ctx *gin.Context, existingElem interface{}, elem interface{}, id interface{}) (*sql.Result, error)
	UpdateStatus(ctx *gin.Context, id string, status_code string) error
}

// ImmutableGenericStorage represents the immutable generic Storage
// for the domain models that matches with its database models.
// The immutable generic Storage provides only the find & insert methods.
type ImmutableGenericStorage interface {
	Single(ctx *gin.Context, elem interface{}, where string, arg map[string]interface{}) error
	Where(ctx *gin.Context, elems interface{}, where string, arg map[string]interface{}) error
	FindByID(ctx *gin.Context, elem interface{}, id interface{}) error
	FindAll(ctx *gin.Context, elems interface{}, page int, limit int, isAsc bool) error
	Insert(ctx *gin.Context, elem interface{}) error
	DeleteMany(ctx *gin.Context, ids interface{}) error
}

// MySQLStorage is the postgres implementation of generic Storage
type MySQLStorage struct {
	db                  Queryer
	tableName           string
	elemType            reflect.Type
	isImmutable         bool
	selectFields        string
	insertFields        string
	insertParams        string
	updateSetFields     string
	updateManySetFields string
	logStorage          LogStorage
}

// LogStorage storage for logs
type LogStorage struct {
	db           Queryer
	logName      string
	elemType     reflect.Type
	insertFields string
	insertParams string
}

// MysqlConfig represents the configuration for the postgres Storage.
type MysqlConfig struct {
	IsImmutable bool
}

// Single queries an element according to the query & argument provided
func (r *MySQLStorage) Single(ctx *gin.Context, elem interface{}, where string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	// if !r.isImmutable {
	// 	where = fmt.Sprintf(`"deletedAt" IS NULL AND %s`, where)
	// }

	statement, err := db.PrepareNamed(fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", r.selectFields, r.tableName, where))
	if err != nil {
		return err
	}
	defer statement.Close()

	err = statement.Get(elem, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// SinglePOSTEMP queries an element according to the query & argument provided
func (r *MySQLStorage) SinglePOSTEMP(ctx *gin.Context, elem interface{}, where string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}
	// if !r.isImmutable {
	// 	where = fmt.Sprintf(`"deletedAt" IS NULL AND %s`, where)
	// }

	statement, err := db.PrepareNamed(fmt.Sprintf("SELECT %s FROM `%s` WHERE %s",
		r.selectFields, r.tableName, where))
	if err != nil {
		return err
	}
	defer statement.Close()

	err = statement.Get(elem, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// Where queries the elements according to the query & argument provided
func (r *MySQLStorage) Where(ctx *gin.Context, elems interface{}, where string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	// if !r.isImmutable {
	// 	where = fmt.Sprintf(`"deletedAt" IS NULL AND %s`, where)
	// }

	query := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", r.selectFields, r.tableName, where)

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)

	err = db.Select(elems, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// WherePOSTEMP queries the elements according to the query & argument provided
func (r *MySQLStorage) WherePOSTEMP(ctx *gin.Context, elems interface{}, where string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	// if !r.isImmutable {
	// 	where = fmt.Sprintf(`"deletedAt" IS NULL AND %s`, where)
	// }

	query := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", r.selectFields, r.tableName, where)
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)

	err = db.Select(elems, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// SelectWithQuery Customizable Query for Select
func (r *MySQLStorage) SelectWithQuery(ctx *gin.Context, elems interface{}, query string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)

	err = db.Select(elems, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// FindByID finds an element by its id
// it's defined in this project context that
// the element id column in the db should be "id"
func (r *MySQLStorage) FindByID(ctx *gin.Context, elem interface{}, id interface{}) error {
	where := `id = :id`

	err := r.Single(ctx, elem, where, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return err
	}

	return nil
}

// FindAll finds all elements from the database.
func (r *MySQLStorage) FindAll(ctx *gin.Context, elems interface{}, page int, limit int, isAsc bool) error {
	where := `TRUE`
	where = fmt.Sprintf(`%s ORDER BY id`, where)

	if !isAsc {
		where = fmt.Sprintf(`%s %s`, where, "DESC")
	}

	where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)

	err := r.Where(ctx, elems, where, map[string]interface{}{
		"limit":  limit,
		"offset": (page - 1) * limit,
	})

	if err != nil {
		return err
	}

	return nil
}

func interfaceConversion(i interface{}) (map[string]interface{}, error) {
	resJSON, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	err = json.Unmarshal(resJSON, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Insert inserts a new element into the database.
// It assumes the primary key of the table is "id" with serial type.
// It will set the "owner" field of the element with the current account in the context if exists.
// It will set the "created_at" and "updated_at" fields with current time.
// If immutable set true, it won't insert the updated_at
func (r *MySQLStorage) Insert(ctx *gin.Context, elem interface{}) (*sql.Result, error) {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    INSERT INTO %s(%s)
    VALUES (%s)`, r.tableName, r.insertFields, r.insertParams))
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	dbArgs := r.insertArgs(*currentUserID, elem, 0)
	result, err := statement.Exec(dbArgs)
	if err != nil {
		return nil, err
	}

	// lastID, _ := (result).LastInsertId()

	// _, err = r.InsertTrail(ctx, fmt.Sprintf("%d", lastID))
	// if err != nil {
	// 	return nil, err
	// }

	// Assuming id is pre-generated before this
	lastID := r.findID(elem)

	_, err = r.InsertTrail(ctx, lastID.(string))
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *MySQLStorage) InsertTrail(ctx *gin.Context, id string) (*sql.Result, error) {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	created_at := library.UTCPlus7().Format("2006-01-02 15:04:05")
	query := fmt.Sprintf(`
  INSERT INTO user_actions(id, user_id, table_name, action, created_at, ref_id)
  VALUES (UUID(), '%s', '%s', 'Create', :created_at, '%s')`, *currentUserID, r.tableName, id)
	statement, err := db.PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	//	dbArgs := r.insertArgs(*currentUserID, nil, 0)
	//	fmt.Println(dbArgs)
	dbArgs := make(map[string]interface{})
	dbArgs["created_at"] = created_at
	result, err := statement.Exec(dbArgs)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *MySQLStorage) insertArgs(currentUserID string, elem interface{}, index int) map[string]interface{} {
	res := map[string]interface{}{
		"created_at": library.UTCPlus7(),
		"created_by": currentUserID,
	}

	if !r.isImmutable {
		res["updated_at"] = library.UTCPlus7()
		res["updated_by"] = currentUserID
	}

	var v reflect.Value
	if reflect.TypeOf(elem) == reflect.TypeOf(reflect.Value{}) {
		data := elem.(reflect.Value)
		v = reflect.Indirect(data)
	} else {
		v = reflect.ValueOf(elem).Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			var typeMapString map[string]interface{}
			var val interface{}
			if v.Field(i).Type() == reflect.TypeOf(typeMapString) {
				metadataBytes, err := json.Marshal(v.Field(i).Interface())
				if err != nil {
					val = "{}"
				} else {
					val = string(metadataBytes)
				}
			} else {
				val = v.Field(i).Interface()
			}
			res[dbTag] = val
		}
	}

	if index != 0 {
		s := strconv.Itoa(index)
		res = renamingKey(res, s)
	}

	return res
}

func (r *LogStorage) insertArgs(currentUserID int, elem interface{}, index int) map[string]interface{} {
	res := map[string]interface{}{
		"created_at": library.UTCPlus7(),
		"created_by": currentUserID,
	}

	var v reflect.Value

	if reflect.TypeOf(elem) == reflect.TypeOf(reflect.Value{}) {
		data := elem.(reflect.Value)
		v = reflect.Indirect(data)
	} else {
		v = reflect.ValueOf(elem).Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			var typeMapString map[string]interface{}
			var val interface{}
			if v.Field(i).Type() == reflect.TypeOf(typeMapString) {
				metadataBytes, err := json.Marshal(v.Field(i).Interface())
				if err != nil {
					val = "{}"
				} else {
					val = string(metadataBytes)
				}
			} else {
				val = v.Field(i).Interface()
			}
			res[dbTag] = val
		}
	}

	if index != 0 {
		s := strconv.Itoa(index)
		res = renamingKey(res, s)
	}

	return res
}

// InsertMany is function for creating many datas into specific table in database.
func (r *MySQLStorage) InsertMany(ctx *gin.Context, elem interface{}) error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	sqlStr := fmt.Sprintf(`
  INSERT INTO %s(%s)
  VALUES `, r.tableName, r.insertFields)

	var dbArgs map[string]interface{}

	datas := reflect.ValueOf(elem)

	insertFields := strings.Split(r.insertFields, ",")
	limit := 60000 / len(insertFields)
	indexData := 0
	if datas.Kind() == reflect.Slice {
		for i := 0; i < datas.Len(); i++ {
			sqlStr += fmt.Sprintf("(%s),", insertParams(r.elemType, r.isImmutable, i+1))
			arg := r.insertArgs(*currentUserID, datas.Index(i), i+1)
			if indexData == 0 {
				dbArgs = arg
			} else {
				for k, v := range arg {
					dbArgs[k] = v
				}
			}
			indexData++
			if indexData == limit {
				err := r.insertData(ctx, sqlStr, dbArgs)
				if err != nil {
					return err
				}

				indexData = 0
				sqlStr = fmt.Sprintf(`
        INSERT INTO "%s"(%s)
        VALUES `, r.tableName, r.insertFields)
				dbArgs = map[string]interface{}{}
			}
		}
	}

	if datas.Kind() == reflect.Map {
		for key, element := range datas.MapKeys() {
			sqlStr += fmt.Sprintf("(%s),", insertParams(r.elemType, r.isImmutable, key+1))
			arg := r.insertArgs(*currentUserID, datas.MapIndex(element), key+1)
			if indexData == 0 {
				dbArgs = arg
			} else {
				for k, v := range arg {
					dbArgs[k] = v
				}
			}
			indexData++
			if indexData == limit {
				err := r.insertData(ctx, sqlStr, dbArgs)
				if err != nil {
					return err
				}

				indexData = 0
				sqlStr = fmt.Sprintf(`
        INSERT INTO "%s"(%s)
        VALUES `, r.tableName, r.insertFields)
				dbArgs = map[string]interface{}{}
			}
		}
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")

	statement, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(dbArgs)
	if err != nil {
		return err
	}

	return nil
}

// InsertManyWithTime is function for creating many datas into specific table in database with specific created_at.
func (r *MySQLStorage) InsertManyWithTime(ctx *gin.Context, elem interface{}, created_at time.Time) error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	sqlStr := fmt.Sprintf(`
  INSERT INTO "%s"(%s)
  VALUES `, r.tableName, r.insertFields)

	var dbArgs map[string]interface{}

	datas := reflect.ValueOf(elem)

	a := strings.Split(r.insertFields, ",")
	limit := 60000 / len(a)
	indexData := 0
	if datas.Kind() == reflect.Slice {
		for i := 0; i < datas.Len(); i++ {
			sqlStr += fmt.Sprintf("(%s),", insertParams(r.elemType, r.isImmutable, i+1))

			arg := r.insertArgs(*currentUserID, datas.Index(i), i+1)
			arg[fmt.Sprintf("created_at%d", i+1)] = created_at
			if indexData == 0 {
				dbArgs = arg
			} else {
				for k, v := range arg {
					dbArgs[k] = v
				}
			}
			indexData++
			if indexData == limit {
				err := r.insertData(ctx, sqlStr, dbArgs)
				if err != nil {
					return err
				}

				indexData = 0
				sqlStr = fmt.Sprintf(`
        INSERT INTO "%s"(%s)
        VALUES `, r.tableName, r.insertFields)
				dbArgs = map[string]interface{}{}
			}
		}
	}

	if datas.Kind() == reflect.Map {
		for key, element := range datas.MapKeys() {
			sqlStr += fmt.Sprintf("(%s),", insertParams(r.elemType, r.isImmutable, key+1))
			arg := r.insertArgs(*currentUserID, datas.MapIndex(element), key+1)
			arg[fmt.Sprintf("created_at%d", key+1)] = created_at
			if indexData == 0 {
				dbArgs = arg
			} else {
				for k, v := range arg {
					dbArgs[k] = v
				}
			}
			indexData++
			if indexData == limit {
				err := r.insertData(ctx, sqlStr, dbArgs)
				if err != nil {
					return err
				}

				indexData = 0
				sqlStr = fmt.Sprintf(`
        INSERT INTO "%s"(%s)
        VALUES `, r.tableName, r.insertFields)
				dbArgs = map[string]interface{}{}
			}
		}
	}

	err := r.insertData(ctx, sqlStr, dbArgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLStorage) insertData(ctx *gin.Context, sqlStr string, dbArgs map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")

	statement, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(dbArgs)
	if err != nil {
		return err
	}

	return nil
}

// RenamingKey is function for renaming key for map
func renamingKey(m map[string]interface{}, add string) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range m {
		newKey := fmt.Sprint(k, add)
		newMap[newKey] = v
	}
	return newMap
}

func (r *MySQLStorage) findChanges(existingElem interface{}, elem interface{}) map[string]interface{} {
	diff := map[string]interface{}{}
	ev := reflect.ValueOf(existingElem).Elem()
	v := reflect.ValueOf(elem).Elem()
	for i := 0; i < ev.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			val1 := ev.Field(i).Interface()
			val2 := v.Field(i).Interface()
			if !reflect.DeepEqual(val1, val2) {

				singleDiff := make([]interface{}, 2)
				singleDiff[0] = val1
				singleDiff[1] = val2
				diff[dbTag] = singleDiff
			}
		}
	}
	return diff
}

// Update updates the element in the database.
// It will update the "updated_at" field.
func (r *MySQLStorage) Update(ctx *gin.Context, elem interface{}) error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	id := r.findID(elem)
	existingElem := reflect.New(r.elemType).Interface()
	err := r.FindByID(ctx, existingElem, id)
	if err != nil {
		return err
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    UPDATE %s SET %s WHERE id = :id`,
		r.tableName,
		r.updateSetFields))
	if err != nil {
		return err
	}
	defer statement.Close()

	updateArgs := r.updateArgs(*currentUserID, existingElem, elem)
	updateArgs["id"] = id

	_, err = statement.Exec(updateArgs)
	if err != nil {
		return err
	}

	_, err = r.UpdateTrail(ctx, existingElem, elem, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *MySQLStorage) UpdateTrail(ctx *gin.Context, existingElem interface{}, elem interface{}, id interface{}) (*sql.Result, error) {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}
	created_at := library.UTCPlus7().Format("2006-01-02 15:04:05")
	statement, err := db.PrepareNamed(fmt.Sprintf(`
    INSERT INTO user_actions(id, user_id, table_name,action,created_at, ref_id)
    VALUES (UUID(), '%s', '%s', '%s', :created_at, :id)`, *currentUserID, r.tableName, "Update"))
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	dbArgs := make(map[string]interface{})
	dbArgs["created_at"] = created_at
	dbArgs["id"] = id
	result, err := statement.Exec(dbArgs)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *MySQLStorage) UpdateStatus(ctx *gin.Context, id string, status_code string) error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	updated_at := library.UTCPlus7().Format("2006-01-02 15:04:05")

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    UPDATE %s set status_id = '%s', updated_at = :updated_at, updated_by = '%s' where id = '%s'`, r.tableName, status_code, *currentUserID, id))
	if err != nil {
		return err
	}
	defer statement.Close()

	dbArgs := make(map[string]interface{})
	dbArgs["updated_at"] = updated_at
	_, err = statement.Exec(dbArgs)
	if err != nil {
		return err
	}

	statementTrail, err := db.PrepareNamed(fmt.Sprintf(`
    INSERT INTO user_actions(id, user_id, table_name, action, action_value, created_at, ref_id)
    VALUES (UUID(), '%s', '%s', '%s', '%s', :updated_at, '%s')`, *currentUserID, r.tableName, "Update Status", status_code, id))
	if err != nil {
		return err
	}
	defer statementTrail.Close()

	dbArgs = make(map[string]interface{})
	dbArgs["updated_at"] = updated_at
	_, err = statementTrail.Exec(dbArgs)
	if err != nil {
		return err
	}

	return nil
}

// UpdateMany updates the element in the database.
// It will update the "updated_at" field.
func (r *MySQLStorage) UpdateMany(ctx *gin.Context, elems interface{}) error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	dbArgs := map[string]interface{}{}

	sqlStr := fmt.Sprintf(`
  UPDATE "%s" as "currentTable"
  SET
    %s
  FROM (VALUES
  `, r.tableName, r.updateManySetFields)

	datas := reflect.ValueOf(elems)

	limit := 2000
	indexData := 0
	if datas.Kind() == reflect.Slice {
		for i := 0; i < datas.Len(); i++ {
			sqlStrIndex, arg := r.updateManyParams(*currentUserID, datas.Index(i), i+1)
			sqlStr += sqlStrIndex
			if indexData == 0 {
				dbArgs = arg
			} else {
				for k, v := range arg {
					dbArgs[k] = v
				}
			}
			indexData++
			if indexData == limit {
				err := r.updated_ata(ctx, sqlStr, dbArgs)
				if err != nil {
					return err
				}

				indexData = 0
				sqlStr = fmt.Sprintf(`
        UPDATE "%s" as "currentTable"
        SET
          %s
        FROM (VALUES
        `, r.tableName, r.updateManySetFields)
				dbArgs = map[string]interface{}{}
			}
		}
	}

	if datas.Kind() == reflect.Map {
		for key, element := range datas.MapKeys() {
			sqlStrIndex, arg := r.updateManyParams(*currentUserID, datas.MapIndex(element), key+1)
			sqlStr += sqlStrIndex
			if indexData == 0 {
				dbArgs = arg
			} else {
				for k, v := range arg {
					dbArgs[k] = v
				}
			}
			indexData++
			if indexData == limit {
				err := r.updated_ata(ctx, sqlStr, dbArgs)
				if err != nil {
					return err
				}

				indexData = 0
				sqlStr = fmt.Sprintf(`
        UPDATE "%s" as "currentTable"
        SET
          %s
        FROM (VALUES
        `, r.tableName, r.updateManySetFields)
				dbArgs = map[string]interface{}{}
			}
		}
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr = fmt.Sprintf(`%s
  ) as "updatedTable"("updated_at", %s)
  where cast("currentTable".id as int) = cast("updatedTable".id as int)
  `, sqlStr, r.selectFields)

	statement, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(dbArgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLStorage) updated_ata(ctx *gin.Context, sqlStr string, dbArgs map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr = fmt.Sprintf(`%s
  ) as "updatedTable"("updated_at",  %s)
  where cast("currentTable".id as int) = cast("updatedTable".id as int)
  `, sqlStr, r.selectFields)

	statement, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(dbArgs)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLStorage) updateManyParams(currentUserID string, elem interface{}, index int) (string, map[string]interface{}) {
	sqlStr := fmt.Sprintf(`(cast(:updated_at%d as timestamp),%d,`, index, currentUserID)

	var v reflect.Value

	res := map[string]interface{}{
		"updated_at": library.UTCPlus7().Format(time.RFC3339),
		"updated_by": currentUserID,
	}

	if reflect.TypeOf(elem) == reflect.TypeOf(reflect.Value{}) {
		data := elem.(reflect.Value)
		v = reflect.Indirect(data)
	} else {
		v = reflect.ValueOf(elem).Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if !emptyTag(dbTag) {
			var typeMapString types.Metadata
			var val interface{}
			var typeTime time.Time
			var field reflect.Value
			if v.Field(i).Kind() == reflect.Ptr {
				field = v.Field(i).Elem()
			} else {
				field = v.Field(i)
			}
			if field.Type() == reflect.TypeOf(typeMapString) {
				metadataBytes, err := json.Marshal(field.Interface())
				if err != nil {
					val = "{}"
				} else {
					val = string(metadataBytes)
				}
			} else {
				val = field.Interface()
			}

			if dbTag == "created_at" || dbTag == "updated_at" || field.Type() == reflect.TypeOf(typeTime) {
				valTime := val.(time.Time)
				val = valTime.Format(time.RFC3339)
				sqlStr += fmt.Sprintf(`cast(:%s%d as timestamp),`, dbTag, index)
				res[dbTag] = val
			} else if field.Type() == reflect.TypeOf(typeMapString) {
				sqlStr += fmt.Sprintf(`cast(:%s%d as jsonb),`, dbTag, index)
				res[dbTag] = val
			} else {
				switch field.Kind() {
				case reflect.String:
					if r.elemType.Field(i).Tag.Get("cast") != "" {
						sqlStr += fmt.Sprintf(`cast('%s' as %s),`, val, r.elemType.Field(i).Tag.Get("cast"))
					} else {
						sqlStr += fmt.Sprintf(`'%s',`, val)
					}
					break
				default:
					sqlStr += fmt.Sprint(val, ",")
					break
				}
			}
		}
	}

	sqlStr = strings.TrimSuffix(sqlStr, ",")

	sqlStr += "),"

	if index != 0 {
		s := strconv.Itoa(index)
		res = renamingKey(res, s)
	}

	return sqlStr, res
}

// it assumes the id column named "id"
func (r *MySQLStorage) findID(elem interface{}) interface{} {
	v := reflect.ValueOf(elem).Elem()
	for i := 0; i < v.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if idTag(dbTag) {
			return v.Field(i).Interface()
		}
	}
	return nil
}

func (r *MySQLStorage) updateArgs(currentUserID string, existingElem interface{}, elem interface{}) map[string]interface{} {
	res := map[string]interface{}{
		"updated_at": library.UTCPlus7(),
		"updated_by": currentUserID,
	}

	v := reflect.ValueOf(elem).Elem()
	ev := reflect.ValueOf(existingElem).Elem()
	for i := 0; i < ev.NumField(); i++ {
		dbTag := r.elemType.Field(i).Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			var typeMapString map[string]interface{}
			var val interface{}

			if v.Field(i).Type() == reflect.TypeOf(typeMapString) {
				metadataBytes, err := json.Marshal(v.Field(i).Interface())
				if err != nil {
					val = "{}"
				} else {
					val = string(metadataBytes)
				}
			} else {
				val = v.Field(i).Interface()
			}
			res[dbTag] = val
		}
	}
	return res
}

// Delete deletes the elem from database.
// Delete not really deletes the elem from the db, but it will set the
// "deletedAt" column to current time.
func (r *MySQLStorage) Delete(ctx *gin.Context, id interface{}) error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    UPDATE %s SET deletedAt = :deletedAt, deletedBy = :deletedBy WHERE id = :id`, r.tableName))
	if err != nil {
		return err
	}
	defer statement.Close()

	deleteArgs := map[string]interface{}{
		"id":        id,
		"deletedAt": library.UTCPlus7(),
		"deletedBy": currentUserID,
	}
	_, err = statement.Exec(deleteArgs)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMany delete elems from database.
// DeleteMany not really delete elems from the db, but it will set the
// "deletedAt" column to current time.
func (r *MySQLStorage) DeleteMany(ctx *gin.Context, ids interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	// Check if interface is type of slices
	datas := reflect.ValueOf(ids)
	if datas.Kind() != reflect.Slice {
		return fmt.Errorf("ids data should be slices")
	}

	if r.isImmutable {
		query := fmt.Sprintf(`DELETE FROM %s WHERE id IN (:ids)`, r.tableName)
		query, args, err := sqlx.Named(query, map[string]interface{}{
			"ids": ids,
		})
		if err != nil {
			return err
		}

		query, args, err = sqlx.In(query, args...)
		if err != nil {
			return err
		}

		query = db.Rebind(query)
		db.MustExec(query, args...)
		return nil
	}

	var queryParam string
	var payloads = map[string]interface{}{
		"deletedAt": library.UTCPlus7(),
	}

	for i := 0; i < datas.Len(); i++ {
		queryParam += fmt.Sprintf(":%s%d,", "id", i+1)
		payloads[fmt.Sprintf("%s%d", "id", i+1)] = datas.Index(i).Interface()
	}

	queryParam = strings.TrimSuffix(queryParam, ",")

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    UPDATE "%s" SET "deletedAt" = :deletedAt WHERE "id" in (%s) RETURNING %s
  `, r.tableName, queryParam, r.selectFields))
	if err != nil {
		return err
	}

	_, err = statement.Exec(payloads)
	if err != nil {
		return err
	}
	return nil
}

// CountAll is function to count all row datas in specific table in database
func (r *MySQLStorage) CountAll(ctx *gin.Context, count interface{}) error {
	var where string
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	if !r.isImmutable {
		where = fmt.Sprintf(`"deletedAt" IS NULL`)
	}

	q := fmt.Sprintf(`SELECT COUNT(*) FROM "%s" WHERE %s`, r.tableName, where)

	err := db.Get(count, q)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// HardDelete is function to hard deleting data into specific table in database
func (r *MySQLStorage) HardDelete(ctx *gin.Context, id interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    DELETE FROM %s WHERE id = :id
  `, r.tableName))
	if err != nil {
		return err
	}
	defer statement.Close()

	deleteArgs := map[string]interface{}{
		"id": id,
	}
	_, err = statement.Exec(deleteArgs)
	if err != nil {
		return err
	}

	return nil
}

// ExecQuery is function to only execute raw query into database
func (r *MySQLStorage) ExecQuery(ctx *gin.Context, query string, args map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	statement, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(args)
	if err != nil {
		return err
	}

	return nil
}

// SelectFirstWithQuery Customizable Query for Select only take the first row
func (r *MySQLStorage) SelectFirstWithQuery(ctx *gin.Context, elems interface{}, query string, arg map[string]interface{}) error {
	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}

	query = db.Rebind(query)

	err = db.Get(elems, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}

	return nil
}

// ActivityLog log for transactions (insert, update, delete)
type ActivityLog struct {
	ID              int                    `db:"id"`
	UserID          int                    `db:"userId"`
	TableName       string                 `db:"tableName"`
	ReferenceID     int                    `db:"referenceId"`
	Metadata        map[string]interface{} `db:"metadata"`
	ValueBefore     map[string]interface{} `db:"valueBefore"`
	ValueAfter      map[string]interface{} `db:"valueAfter"`
	TransactionTime *time.Time             `db:"transactionTime"`
	TransactionType string                 `db:"transactionType"`
	created_at      *time.Time             `db:"created_at"`
}

func getContextVariables(ctx *gin.Context) *string {
	return appcontext.UserID(ctx)
}

func determineUser(ctx *gin.Context) string {
	userID := getContextVariables(ctx)
	var resUserID string
	if userID != nil {
		resUserID = *userID
	}

	return resUserID
}

// NewLogStorage creates a logStorage
func NewLogStorage(db *sqlx.DB, logName string) *LogStorage {
	logType := reflect.TypeOf(ActivityLog{})
	return &LogStorage{
		db:           db,
		logName:      logName,
		elemType:     logType,
		insertFields: insertFields(logType, true),
		insertParams: insertParams(logType, true, 0),
	}
}

// NewMySQLStorage creates a new generic postgres Storage
func NewMySQLStorage(db *sqlx.DB, tableName string, elem interface{}, cfg MysqlConfig) *MySQLStorage {
	elemType := reflect.TypeOf(elem)
	return &MySQLStorage{
		db:                  db,
		tableName:           tableName,
		elemType:            elemType,
		isImmutable:         cfg.IsImmutable,
		selectFields:        selectFields(elemType),
		insertFields:        insertFields(elemType, cfg.IsImmutable),
		insertParams:        insertParams(elemType, cfg.IsImmutable, 0),
		updateSetFields:     updateSetFields(elemType),
		updateManySetFields: updateManySetFields(elemType),
	}
}

func selectFields(elemType reflect.Type) string {
	dbFields := []string{}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag != "" && dbTag != "-" {
			dbFields = append(dbFields, fmt.Sprintf("`%s`", dbTag))
		}
	}
	return strings.Join(dbFields, ",")
}

func insertFields(elemType reflect.Type, isImmutable bool) string {
	dbFields := []string{"`created_by`", "`created_at`"}
	if !isImmutable {
		dbFields = append(dbFields, "`updated_by`", "`updated_at`")
	}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			dbFields = append(dbFields, fmt.Sprintf("`%s`", dbTag))
		}
	}
	return strings.Join(dbFields, ",")
}

func insertParams(elemType reflect.Type, isImmutable bool, index int) string {
	dbParams := []string{":created_by", ":created_at"}
	if !isImmutable {
		dbParams = append(dbParams, ":updated_by", ":updated_at")
	}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			dbParams = append(dbParams, fmt.Sprintf(":%s", dbTag))
		}
	}

	if index != 0 {
		s := strconv.Itoa(index)
		for i, v := range dbParams {
			dbParams[i] = fmt.Sprint(v, s)
		}
	}

	return strings.Join(dbParams, ",")
}

func updateSetFields(elemType reflect.Type) string {
	setFields := []string{"`updated_at` = :updated_at", "`updated_by` = :updated_by"}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			setFields = append(setFields, fmt.Sprintf("`%s` = :%s", dbTag, dbTag))
		}
	}
	return strings.Join(setFields, ",")
}

func updateManySetFields(elemType reflect.Type) string {
	setManyFields := []string{"`updated_at` = `updatedTable`.`updated_at`"}
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		if !readOnlyTag(dbTag) && !emptyTag(dbTag) {
			setManyFields = append(setManyFields, fmt.Sprintf("`%s` = `updatedTable`.`%s`", dbTag, dbTag))
		}
	}

	return strings.Join(setManyFields, ",")
}

func idTag(dbTag string) bool {
	return dbTag == "id"
}

func emptyTag(dbTag string) bool {
	emptyTags := []string{"", "-"}
	for _, t := range emptyTags {
		if dbTag == t {
			return true
		}
	}
	return false
}

func readOnlyTag(dbTag string) bool {
	readOnlyTags := []string{"created_at", "updated_at", "deletedAt"}
	for _, t := range readOnlyTags {
		if dbTag == t {
			return true
		}
	}
	return false
}

// Insert inserts a new element into the database.
// It assumes the primary key of the table is "id" with serial type.
// It will set the "owner" field of the element with the current account in the context if exists.
// It will set the "created_at" and "updated_at" fields with current time.
// If immutable set true, it won't insert the updated_at
func (r *MySQLStorage) InsertNoTrail(ctx *gin.Context, elem interface{}) (*sql.Result, error) {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    INSERT INTO %s(%s)
    VALUES (%s)`, r.tableName, r.insertFields, r.insertParams))
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	dbArgs := r.insertArgs(*currentUserID, elem, 0)
	result, err := statement.Exec(dbArgs)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Update updates the element in the database.
// It will update the "updated_at" field.
func (r *MySQLStorage) UpdateNoTrail(ctx *gin.Context, elem interface{}) error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	db := r.db
	tx, ok := TxFromContext(ctx)
	if ok {
		db = tx
	}

	id := r.findID(elem)
	existingElem := reflect.New(r.elemType).Interface()
	err := r.FindByID(ctx, existingElem, id)
	if err != nil {
		return err
	}

	statement, err := db.PrepareNamed(fmt.Sprintf(`
    UPDATE %s SET %s WHERE id = :id`,
		r.tableName,
		r.updateSetFields))
	if err != nil {
		return err
	}
	defer statement.Close()

	updateArgs := r.updateArgs(*currentUserID, existingElem, elem)
	updateArgs["id"] = id

	_, err = statement.Exec(updateArgs)
	if err != nil {
		return err
	}

	return nil
}

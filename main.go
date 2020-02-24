package main

import (
	"errors"
	"reflect"
	"strings"
)

const (
	sep   = "."
	dbTag = "db"
)

// Foo is foo
type Foo struct {
	ID    uint  `db:"id"`
	Title uint  `db:"title"`
	Boos  []Boo `db:"boos"`
}

// Boo is boo
type Boo struct {
	ID   uint `db:"id"`
	Name uint `db:"name"`
}

func main() {
	cols := []string{"id", "title", "boos.id", "boos.name"}
	rows := []map[string]interface{}{
		map[string]interface{}{"id": 1, "title": "title1", "boos.id": 101, "boos.name": "name1"},
		map[string]interface{}{"id": 2, "title": "title2", "boos.id": 102, "boos.name": "name2"},
	}
	foo := []Foo{}

	Map(&foo, rows, cols)
}

// Map is map
func Map(dest interface{}, rows []map[string]interface{}, cols []string) error {
	// instead of to record data to dest directlry,
	// how about making an intermediate data structure & transforming it to dest at last?

	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer")
	}

	// assume that dest is a pointer of slice
	// dereferences reflect value
	drv := rv.Elem()

	// element reflect type
	ert := drv.Type().Elem()

	// for convenience
	// this is just a copy of the dest's element on which I currently work
	// ex) { "id": 1, "name": "name1", "boos": reflect.Value(slice.Of(Boo)), "boos.doos": reflect.Value(slice.Of(Coo)), "coos": reflect.Value(slice.Of(Coo))}
	dic := map[string]reflect.Value{}

	for _, row := range rows {
		// First, create a new one from (row, cols)
		created := createNewValue(ert, row, cols)

		// find matched element from dest to created
		// drv.Len()
	}
	return nil
}

func getHierarchies(cols []string) [][]string {
	// I think I should use a map as a store hierarchy
	// Not able to handle the case with 2-dimensional string array like - ["id", "boos.id", "coos.id", "boos.doos.id", "coos.eoos.id"]
	result := [][]string{}
	deepest := cols[0]
	maxSepCnt := 1
	for _, col := range cols {
		parts := strings.Split(col, sep)
		if len(parts) < 1 {
			continue
		}
		for i, part := range parts {
			if i > 0 {

			}
		}
	}
	return result
}

func createNewValue(t reflect.Type, row map[string]interface{}, cols []string) (reflect.Value, map[string]reflect.Value) {
	// record map representaion at the same time
	dic := map[string]reflect.Value{}

	createdPtr := reflect.New(t)
	created := createdPtr.Elem()
	for _, col := range cols {
		v := row[col]
		rv := reflect.ValueOf(v)
		parts := strings.Split(col, sep)
		// dereference to manipulate T, not Ptr(T)
		cursor := created
		for _, part := range parts {
			found, ok := getFieldByTagValue(cursor, part)
			if !ok {
				continue
			}
			cursor = found
		}
		if cursor.CanSet() {
			// cursor's kind must be the one of (Slice, basic type(int, bool...))
			if cursor.Kind() == reflect.Slice {
				lastPart := parts[len(parts)-1]
				// assume that ts is slice, not slice ptr
				createdEmbeded := reflect.New(cursor.Elem().Type()).Elem()
				matchedField, _ := getFieldByTagValue(createdEmbeded, lastPart)
				matchedField.Set(rv)
				cursor.Set(reflect.Append(cursor, createdEmbeded))
				// is it possible the below?
				// reflect.Append(cursor, createdEmbeded)
			} else {
				cursor.Set(rv)
			}
		}
	}
	return created
}

func findMatchedFromDestToCreated(destSlice reflect.Value, valueToFind reflect.Value, cols []string) (reflect.Value, bool) {
	l := destSlice.Len()
	if l < 1 {
		return reflect.ValueOf((*int)(nil)), false
	}

	for i := 0; i < l; i++ {
		el := destSlice.Index(i)

	}
}

func getFieldByTagValue(v reflect.Value, tagValueToFind string) (reflect.Value, bool) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return reflect.ValueOf(0), false
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tagValue := t.Field(i).Tag.Get(dbTag)
		if tagValue == tagValueToFind {
			return v.Field(i), true
		}
	}
	return reflect.ValueOf(0), false
}

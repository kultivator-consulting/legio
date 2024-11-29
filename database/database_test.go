package database

import (
	"context"
	"cortex_api/database/db_gen"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"testing"
)

func TestApiDatabase(t *testing.T) {
	var tests []struct {
		name string
		want *Model
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ApiDatabase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApiDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModel_Close(t *testing.T) {
	type fields struct {
		DbPool *pgxpool.Pool
	}
	var tests []struct {
		name    string
		fields  fields
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Model{
				DbPool: tt.fields.DbPool,
			}
			if err := db.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModel_Open(t *testing.T) {
	type fields struct {
		DbPool *pgxpool.Pool
	}
	var tests []struct {
		name    string
		fields  fields
		want    *pgxpool.Pool
		want1   *db_gen.Queries
		want2   context.Context
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Model{
				DbPool: tt.fields.DbPool,
			}
			got, got1, got2, err := db.Open()
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Open() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Open() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

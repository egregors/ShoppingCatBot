package main

import (
	"reflect"
	"sync"
	"testing"
)

func TestInmem_GetAll(t *testing.T) {
	type fields struct {
		items map[int64][]string
		mu    *sync.Mutex
	}
	type args struct {
		chatID int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   [][]string
	}{
		{
			"Missing key",
			fields{
				items: map[int64][]string{0: {"foo", "bar"}},
				mu:    &sync.Mutex{},
			},
			args{chatID: 1},
			[][]string(nil),
		},
		{
			"2 items",
			fields{
				items: map[int64][]string{0: {"foo", "bar"}},
				mu:    &sync.Mutex{},
			},
			args{chatID: 0},
			[][]string{{"foo", "bar"}},
		},
		{
			"11 items",
			fields{
				items: map[int64][]string{0: {"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}},
				mu:    &sync.Mutex{},
			},
			args{chatID: 0},
			[][]string{{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, {"11", "10"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := &Inmem{
				items: tt.fields.items,
				mu:    tt.fields.mu,
			}
			if got := db.GetAll(tt.args.chatID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

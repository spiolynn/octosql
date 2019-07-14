package execution

import (
	"testing"
	"time"

	"github.com/cube2222/octosql"
)

func TestOrderBy_Get(t *testing.T) {
	now := time.Now()

	type args struct {
		stream RecordStream
		fields []OrderField
	}
	tests := []struct {
		name    string
		args    args
		want    RecordStream
		wantErr bool
	}{
		{
			name: "simple order - one column int ascending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "age"},
						[]interface{}{1, 7}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "age"},
						[]interface{}{2, 10}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"id", "age"},
						[]interface{}{3, 2}),
				}),
				fields: []OrderField{
					{
						ColumnName: "age",
						Direction:  Ascending,
					},
				},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "age"},
					[]interface{}{3, 2}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "age"},
					[]interface{}{1, 7}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"id", "age"},
					[]interface{}{2, 10}),
			}),
			wantErr: false,
		},
		{
			name: "simple order - one column string descending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"b", 7}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"c", 10}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"a", 2}),
				}),
				fields: []OrderField{
					{
						ColumnName: "name",
						Direction:  Descending,
					},
				},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"c", 10}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"b", 7}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"a", 2}),
			}),
			wantErr: false,
		},
		{
			name: "simple order - one column time descending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "birth"},
						[]interface{}{"b", now}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "birth"},
						[]interface{}{"c", now.Add(time.Hour)}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "birth"},
						[]interface{}{"a", now.Add(-1 * time.Hour)}),
				}),
				fields: []OrderField{
					{
						ColumnName: "birth",
						Direction:  Descending,
					},
				},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "birth"},
					[]interface{}{"c", now.Add(time.Hour)}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "birth"},
					[]interface{}{"b", now}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "birth"},
					[]interface{}{"a", now.Add(-1 * time.Hour)}),
			}),
			wantErr: false,
		},
		{
			name: "complex order - string ascending then int descending",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"a", 7}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"d", 19}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"a", -2}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"c", 1}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"d", 17}),
				}),
				fields: []OrderField{
					{
						ColumnName: "name",
						Direction:  Ascending,
					},
					{
						ColumnName: "age",
						Direction:  Descending,
					},
				},
			},
			want: NewInMemoryStream([]*Record{
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"a", 7}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"a", -2}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"c", 1}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"d", 19}),
				NewRecordFromSliceWithNormalize(
					[]octosql.VariableName{"name", "age"},
					[]interface{}{"d", 17}),
			}),
			wantErr: false,
		},

		{
			name: "failed - missing field",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"a", 7}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age?"},
						[]interface{}{"d", 19}),
				}),
				fields: []OrderField{
					{
						ColumnName: "age",
						Direction:  Descending,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},

		{
			name: "failed - type mismatch",
			args: args{
				stream: NewInMemoryStream([]*Record{
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"a", 7}),
					NewRecordFromSliceWithNormalize(
						[]octosql.VariableName{"name", "age"},
						[]interface{}{"d", 19.5}),
				}),
				fields: []OrderField{
					{
						ColumnName: "age",
						Direction:  Descending,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ordered, err := createOrderedStream(tt.args.fields, tt.args.stream)
			if err != nil && !tt.wantErr {
				t.Errorf("Error in create stream: %v", err)
				return
			} else if err != nil {
				return
			}

			equal, err := AreStreamsEqual(tt.want, ordered)
			if err != nil {
				t.Errorf("Error in AreStreamsEqual(): %v", err)
				return
			}

			if !equal {
				t.Errorf("Streams don't match")
				return
			}
		})
	}
}

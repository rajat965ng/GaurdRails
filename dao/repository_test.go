package dao

import (
	"findings/model"
	"testing"

	"gorm.io/gorm"
)

func Test_repository_FindRepositoryByStatus(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		status model.StatusType
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "find repos by status Queued",
			args: args{
				model.QUEUED,
			},
			fields: fields{
				NewDatabaseInstance().GetConnection(),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dao := &repository{
				db: tt.fields.db,
			}
			got, err := dao.FindRepositoryByStatus(tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.FindRepositoryByStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == tt.want {
				t.Logf("repository.FindRepositoryByStatus() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

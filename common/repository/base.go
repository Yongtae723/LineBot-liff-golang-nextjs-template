package repository

import (
	"fmt"
	"sync"

	"github.com/supabase-community/supabase-go"
)

// BaseRepo is the base struct for Supabase repositories
type BaseRepo struct {
	Client          *supabase.Client
	supabaseRoleKey string
}

var (
	once     sync.Once
	baseRepo *BaseRepo
)

func InitSupabase(supabaseURL, supabaseKey string) error {
	var err error
	once.Do(func() {
		client, err := supabase.NewClient(supabaseURL, supabaseKey, &supabase.ClientOptions{})
		if err != nil {
			fmt.Println("cannot initalize client", err)
		}
		baseRepo = &BaseRepo{Client: client, supabaseRoleKey: supabaseKey}
	})
	return err
}

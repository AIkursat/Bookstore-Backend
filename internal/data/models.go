package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3 // If db access takes longer than 3 seconds, cancel it

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User:  User{},
		Token: Token{},
	}
}

type Models struct {
	User  User
	Token Token
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname, omitempty"`
	LastName  string    `json:"lastname, omitempty"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     Token     `json:"token"`
}

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}
	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil{
		return nil, err
	}
	return &user, nil
}

func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, created_at, updated_at from users where email = $1`

	var user User
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil{
		return nil, err
	}
	return &user, nil
}

func (u *User) Update() error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// statement, because it's not query
	stmt := `update users set 
	email = $1,
	first_name = $2,
    last_name = $3,
	updated_at = $4,
	where id = $5 

	`

	_, err := db.ExecContext(ctx, stmt,
	u.Email,
	u.FirstName,
    u.LastName,
	time.Now(),
	u.ID,
	)

	if err != nil{
		return err
	}

	return nil 
}

func (u *User) Delete() error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    stmt := `delete from users where id = $1`

	_, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil{
		return err
	}

	return nil 

}


func(u *User) Insert(user User) (int, error){ // because we return a id
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12) // default is 10, but used 12 for hash.
    if err != nil{
		return 0, err
	}

	// If that pass that.

	var newID int

	stmt := `insert into users(email, first_name, last_name, password, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id		
		`
        
		// we are using the all values for replacement
		err = db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
        hashedPassword,
		time.Now(),
		time.Now(),
		).Scan(&newID)

		if err != nil{
			return 0, err
		}
		return newID, nil

}

func (u *User) ResetPassword(password string) error{
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
    
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12) // default is 10, but used 12 for hash.
    if err != nil{
		return err
	}

	stmt := `update users set password = $1 where id =$2`
    
	_, err = db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err != nil{
		return  err
	}
     return nil
}

func(u *User) PasswordMatches(plainText string) (bool, error) {
	// Here we compare the password from the db and written password
   err := bcrypt.CompareHashAndPassword([]byte(u.Password), [](plainText))

   if err != nil{
	switch{
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword): // invalid password
       return false, nil
	default:
		return false, err
	}
   }
    return true, nil
}



type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash []byte    `json:"-"` // because not getting send that
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Expiry    time.Time `json:"expiry"`
}

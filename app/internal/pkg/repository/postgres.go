package repository

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net"
	"pehks1980/be2_hw81/internal/app/endpoint"
	"pehks1980/be2_hw81/internal/pkg/model"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

const ( DDL = `
		DROP TABLE IF EXISTS catalog;
        DROP TABLE IF EXISTS environments;
        DROP TABLE IF EXISTS users;

        CREATE TABLE IF NOT EXISTS users
        (
			id uuid NOT NULL CONSTRAINT users_pk PRIMARY KEY,
			name varchar(150) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS environments
		(
			id uuid NOT NULL CONSTRAINT environments_pk PRIMARY KEY,
			title varchar(150) NOT NULL,
			text text NOT NULL
		);

		CREATE TABLE IF NOT EXISTS catalog
		(
			id uuid NOT NULL PRIMARY KEY,
			title varchar(150) NOT NULL,
			user_id uuid NOT NULL CONSTRAINT users_id_fk REFERENCES users ON DELETE CASCADE,
			environment_id uuid NOT NULL CONSTRAINT environments_id_fk REFERENCES environments ON DELETE CASCADE
		);
				

        INSERT INTO public.users (id, name) VALUES
			('b29f95a2-499a-4079-97f5-ff55c3854fcb', 'usr1');
        INSERT INTO public.users (id, name) VALUES
			('b6dede74-ad09-4bb7-a036-997ab3ab3130', 'usr2');

        INSERT INTO public.environments (id, title, text) VALUES
			('e4e12c87-88d8-413c-8ab6-57bfa4e953a8', 'env11', 'RUSSIA');
		INSERT INTO public.environments (id, title, text) VALUES
			('68792339-715c-4823-a4d5-a85cefec8d36', 'env12', 'USA');
        INSERT INTO public.environments (id, title, text) VALUES
			('e095e3a2-5b8e-4bc8-b793-bc3606c4fdd5', 'env13', 'CHINA');

		INSERT INTO public.catalog (id, title, user_id, environment_id ) VALUES
			('02486fe5-787d-45c8-b89e-34ff888f5ea8', 'ENVIRONMENTS USERS CATALOG', 
				'b29f95a2-499a-4079-97f5-ff55c3854fcb', 'e4e12c87-88d8-413c-8ab6-57bfa4e953a8');
		INSERT INTO public.catalog (id, title, user_id, environment_id ) VALUES
			('bc4d018c-c74a-4169-86d9-45f2f2a17e55', 'ENVIRONMENTS USERS CATALOG', 
				'b6dede74-ad09-4bb7-a036-997ab3ab3130', 'e4e12c87-88d8-413c-8ab6-57bfa4e953a8');

		INSERT INTO public.catalog (id, title, user_id, environment_id ) VALUES
			('11666af9-5b8c-4fe0-b6be-f2a3ce1b010e', 'ENVIRONMENTS USERS CATALOG', 
				'b29f95a2-499a-4079-97f5-ff55c3854fcb', '68792339-715c-4823-a4d5-a85cefec8d36');
	`
	FINDUSER = `select id from users where name = $1;`
	FINDENVIRONMENT = `select id, text from environments where title = $1;`
	USERMEMBEROF = `select title from environments where id IN
						(select environment_id from catalog where user_id 
						= (select id from users where name = $1));
	`
	ENVUSERSLIST = `select name from users where id IN
						(select user_id from catalog where environment_id 
						= (select id from environments where title = $1));
	`
	// OTHER QUERYES SO ON.
)

// PgRepo - init pg go struct holds connex to db
type PgRepo struct {
	URL    string
	DBPool *pgxpool.Pool
}

func (pgr *PgRepo) FindGood(ctx context.Context, key string) ([]model.Good, error) {
	panic("implement me")
}

// New Init of pg driver
func (pgr *PgRepo) New(ctx context.Context, filename, filename1 string) endpoint.RepoIf {
	// Строка для подключения к базе данных
	url := filename //"postgres://postuser:postpassword@192.168.1.204:5432/a4"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}
	// Pool соединений обязательно ограничивать сверху
	cfg.MaxConns = 8
	cfg.MinConns = 4
	// HealthCheckPeriod - частота пингования соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute
	// MaxConnLifetime - сколько времени будет жить соединение.
	//можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour
	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	cfg.MaxConnIdleTime = 30 * time.Minute
	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	//init DDL
	_, err = dbpool.Exec(ctx, DDL)
	if err != nil {
		log.Fatal(err)
	}

	return &PgRepo{
		URL:    url,
		DBPool: dbpool,
	}
}

func (pgr *PgRepo) AddUpdGood(ctx context.Context, good model.Good) (string, error) {
	panic("implement me")
}

func (pgr *PgRepo) GetGood(ctx context.Context, title string) (model.Good, error) {
	panic("implement me")
}

func (pgr *PgRepo) DelGood(ctx context.Context, id uuid.UUID) error {
	panic("implement me")
}

// CloseConn - close db connex when server quit
func (pgr *PgRepo) CloseConn() {
	pgr.DBPool.Close()
}



func (pgr *PgRepo) AuthUser(ctx context.Context, user model.User) (string, error) {
	panic("implement me")
}

func (pgr *PgRepo) GetUserEnvs(ctx context.Context, name string) (model.Envs, error) {
	//todo implement sql USERMEMBEROF
	envs := model.Envs{}

	return envs, nil
}

func (pgr *PgRepo) GetEnvUsers(ctx context.Context, title string) (model.Users, error) {
	//todo implement sql ENVUSERSLIST
	users := model.Users{}
	return users, nil
}

func (pgr *PgRepo) GetUser(ctx context.Context, name string) (model.User, error) {
	//todo implement sql FINDUSER
	user := model.User{}
	return user, nil
}

func (pgr *PgRepo) GetEnv(ctx context.Context, title string) (model.Environment, error) {
	//todo implement sql FINDENVIRONMENT
	env := model.Environment{}
	return env, nil
}

func (pgr *PgRepo) AddUpdEnv(ctx context.Context, env model.Environment) (string, error) {
	//todo implement sql AddUpdENVIRONMENT
	env_id := ""
	return env_id, nil
}

func (pgr *PgRepo) AddUpdUser(ctx context.Context, user model.User) (string, error) {
	//todo implement sql AddUpdUSER
	user_id := ""
	return user_id, nil
}

func (pgr *PgRepo) DelUser(ctx context.Context, id uuid.UUID) (error) {
	//todo implement sql delUSER
	return nil
}

func (pgr *PgRepo) DelEnv(ctx context.Context, id uuid.UUID) (error) {
	//todo implement sql delENV
	return nil
}
package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Project-TETI/AgriSense/backend/internal/domain"
	"github.com/Project-TETI/AgriSense/backend/internal/repository"
	"github.com/Project-TETI/AgriSense/backend/internal/service"
	"github.com/Project-TETI/AgriSense/backend/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type authTestDB struct {
	db     *gorm.DB
	admin  *gorm.DB
	schema string
}

func newAuthTestDB(t *testing.T) authTestDB {
	t.Helper()
	dsn := os.Getenv("AUTH_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("set AUTH_TEST_DATABASE_DSN to an isolated PostgreSQL test database")
	}
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	adminSQL := stdlib.OpenDB(*cfg)
	admin, err := gorm.Open(postgres.New(postgres.Config{Conn: adminSQL}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		_ = adminSQL.Close()
		t.Fatal(err)
	}
	schema := "auth_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := admin.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		_ = adminSQL.Close()
		t.Fatal(err)
	}
	testCfg := cfg.Copy()
	testCfg.RuntimeParams["search_path"] = schema
	testCfg.RuntimeParams["application_name"] = schema
	testSQL := stdlib.OpenDB(*testCfg)
	testSQL.SetMaxOpenConns(8)
	t.Cleanup(func() {
		_ = testSQL.Close()
		if err := admin.Exec("DROP SCHEMA " + schema + " CASCADE").Error; err != nil {
			t.Errorf("clean test schema: %v", err)
		}
		_ = adminSQL.Close()
	})
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: testSQL}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&domain.User{}, &domain.RefreshToken{}, &domain.Plant{}, &domain.Disease{}, &domain.Inspection{}); err != nil {
		t.Fatal(err)
	}
	return authTestDB{db: db, admin: admin, schema: schema}
}

func authServiceFor(db *gorm.DB) service.AuthService {
	return service.NewAuthService(repository.NewUserRepository(db), repository.NewRefreshTokenRepository(db), repository.NewAuthTransaction(db), strings.Repeat("x", 32), 15*time.Minute, 7*24*time.Hour)
}

func seedRefresh(t *testing.T, db *gorm.DB, expiry time.Time) (domain.User, domain.RefreshToken, string) {
	t.Helper()
	hash, err := utils.HashPassword("test-password")
	if err != nil {
		t.Fatal(err)
	}
	user := domain.User{Name: "Test User", Email: uuid.NewString() + "@example.com", PasswordHash: hash}
	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}
	raw, tokenHash, err := utils.GenerateRefreshToken()
	if err != nil {
		t.Fatal(err)
	}
	token := domain.RefreshToken{UserID: user.ID, TokenHash: tokenHash, ExpiresAt: expiry}
	if err := db.Create(&token).Error; err != nil {
		t.Fatal(err)
	}
	return user, token, raw
}

func rejectRefreshInserts(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec(`CREATE FUNCTION reject_refresh_insert() RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'injected refresh insert failure'; END $$;
CREATE TRIGGER fail_refresh_insert BEFORE INSERT ON refresh_tokens FOR EACH ROW EXECUTE FUNCTION reject_refresh_insert();`).Error; err != nil {
		t.Fatal(err)
	}
}

func TestAuthRegisterRollsBackWhenSessionInsertFails(t *testing.T) {
	testDB := newAuthTestDB(t)
	rejectRefreshInserts(t, testDB.db)
	session, err := authServiceFor(testDB.db).Register(context.Background(), domain.RegisterRequest{Name: "Test User", Email: "register@example.com", Password: "test-password"})
	if err == nil || session != nil {
		t.Fatal("expected failed registration without returned session")
	}
	var users, tokens int64
	if err := testDB.db.Model(&domain.User{}).Count(&users).Error; err != nil {
		t.Fatal(err)
	}
	if err := testDB.db.Model(&domain.RefreshToken{}).Count(&tokens).Error; err != nil {
		t.Fatal(err)
	}
	if users != 0 || tokens != 0 {
		t.Fatalf("registration left partial data: users=%d tokens=%d", users, tokens)
	}
}

func TestAuthRegisterConcurrentEmailConflict(t *testing.T) {
	testDB := newAuthTestDB(t)
	svc := authServiceFor(testDB.db)
	start := make(chan struct{})
	results := make(chan error, 2)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for i := 0; i < 2; i++ {
		go func() {
			<-start
			_, err := svc.Register(ctx, domain.RegisterRequest{Name: "Test User", Email: "same@example.com", Password: "test-password"})
			results <- err
		}()
	}
	close(start)
	var successes, conflicts int
	for i := 0; i < 2; i++ {
		err := <-results
		if err == nil {
			successes++
		} else if errors.Is(err, domain.ErrEmailAlreadyRegistered) {
			conflicts++
		} else {
			t.Fatalf("unexpected registration error: %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("successes=%d conflicts=%d", successes, conflicts)
	}
	var users, tokens int64
	testDB.db.Model(&domain.User{}).Count(&users)
	testDB.db.Model(&domain.RefreshToken{}).Count(&tokens)
	if users != 1 || tokens != 1 {
		t.Fatalf("unexpected registration data: users=%d tokens=%d", users, tokens)
	}
	duplicate := domain.User{Email: "same@example.com", Name: "Duplicate", PasswordHash: "unused"}
	if err := repository.NewUserRepository(testDB.db).Create(ctx, &duplicate); !errors.Is(err, domain.ErrEmailAlreadyRegistered) {
		t.Fatalf("database conflict not classified: %v", err)
	}
}

func waitForRefreshLocks(t *testing.T, db authTestDB, want int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		var count int64
		err := db.admin.Raw(`SELECT count(*) FROM pg_stat_activity WHERE application_name = ? AND wait_event_type = 'Lock' AND query LIKE '%FOR UPDATE%'`, db.schema).Scan(&count).Error
		if err != nil {
			t.Fatal(err)
		}
		if count >= int64(want) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("refresh requests did not reach the held PostgreSQL row lock")
}

func TestAuthRefreshConcurrentSingleUse(t *testing.T) {
	testDB := newAuthTestDB(t)
	expiry := time.Now().UTC().Add(time.Hour).Truncate(time.Microsecond)
	_, token, raw := seedRefresh(t, testDB.db, expiry)
	holder := testDB.db.Begin()
	if holder.Error != nil {
		t.Fatal(holder.Error)
	}
	defer holder.Rollback()
	if err := holder.Exec("SELECT id FROM refresh_tokens WHERE id = ? FOR UPDATE", token.ID).Error; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	results := make(chan error, 2)
	svc := authServiceFor(testDB.db)
	for i := 0; i < 2; i++ {
		go func() { _, err := svc.RefreshToken(ctx, raw); results <- err }()
	}
	waitForRefreshLocks(t, testDB, 2)
	if err := holder.Commit().Error; err != nil {
		t.Fatal(err)
	}
	var successes, rejections int
	for i := 0; i < 2; i++ {
		err := <-results
		if err == nil {
			successes++
		} else if errors.Is(err, domain.ErrInvalidRefreshToken) {
			rejections++
		} else {
			t.Fatalf("unexpected refresh error: %v", err)
		}
	}
	if successes != 1 || rejections != 1 {
		t.Fatalf("successes=%d rejections=%d", successes, rejections)
	}
	var rows []domain.RefreshToken
	if err := testDB.db.Find(&rows).Error; err != nil {
		t.Fatal(err)
	}
	var active int
	for _, row := range rows {
		if row.RevokedAt == nil {
			active++
			if !row.ExpiresAt.Equal(expiry) {
				t.Fatal("rotation extended refresh expiry")
			}
		}
	}
	if len(rows) != 2 || active != 1 {
		t.Fatalf("unexpected tokens: total=%d active=%d", len(rows), active)
	}
}

func TestAuthRefreshRollsBackRevocationOnInsertFailure(t *testing.T) {
	testDB := newAuthTestDB(t)
	_, token, raw := seedRefresh(t, testDB.db, time.Now().UTC().Add(time.Hour))
	rejectRefreshInserts(t, testDB.db)
	session, err := authServiceFor(testDB.db).RefreshToken(context.Background(), raw)
	if err == nil || session != nil {
		t.Fatal("expected refresh insert failure without session")
	}
	var stored domain.RefreshToken
	if err := testDB.db.First(&stored, "id = ?", token.ID).Error; err != nil {
		t.Fatal(err)
	}
	if stored.RevokedAt != nil {
		t.Fatal("failed rotation revoked the original token")
	}
	var count int64
	if err := testDB.db.Model(&domain.RefreshToken{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("unexpected tokens after rollback: %d", count)
	}
}

func TestAuthRefreshChecksExpiryAfterWaitingForLock(t *testing.T) {
	testDB := newAuthTestDB(t)
	_, token, raw := seedRefresh(t, testDB.db, time.Now().UTC().Add(time.Hour))
	holder := testDB.db.Begin()
	if holder.Error != nil {
		t.Fatal(holder.Error)
	}
	defer holder.Rollback()
	if err := holder.Exec("SELECT id FROM refresh_tokens WHERE id = ? FOR UPDATE", token.ID).Error; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() { _, err := authServiceFor(testDB.db).RefreshToken(ctx, raw); result <- err }()
	waitForRefreshLocks(t, testDB, 1)
	if err := holder.Model(&domain.RefreshToken{}).Where("id = ?", token.ID).Update("expires_at", time.Now().UTC().Add(-time.Second)).Error; err != nil {
		t.Fatal(err)
	}
	if err := holder.Commit().Error; err != nil {
		t.Fatal(err)
	}
	if err := <-result; !errors.Is(err, domain.ErrInvalidRefreshToken) {
		t.Fatalf("expired token accepted after waiting: %v", err)
	}
}

func TestAuthSessionFlowAndAbsoluteExpiry(t *testing.T) {
	testDB := newAuthTestDB(t)
	svc := authServiceFor(testDB.db)
	req := domain.RegisterRequest{Name: "Test User", Email: "flow@example.com", Password: "test-password"}
	registered, err := svc.Register(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Login(context.Background(), domain.LoginRequest{Email: req.Email, Password: "incorrect"}); !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatal("wrong password accepted")
	}
	loggedIn, err := svc.Login(context.Background(), domain.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		t.Fatal(err)
	}
	if loggedIn.Response.User.ID != registered.Response.User.ID {
		t.Fatal("login changed user")
	}
	if _, err := svc.GetProfile(context.Background(), loggedIn.Response.User.ID); err != nil {
		t.Fatal(err)
	}
	if err := svc.Logout(context.Background(), loggedIn.RefreshToken); err != nil {
		t.Fatal(err)
	}
	if err := svc.Logout(context.Background(), loggedIn.RefreshToken); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RefreshToken(context.Background(), loggedIn.RefreshToken); !errors.Is(err, domain.ErrInvalidRefreshToken) {
		t.Fatal("revoked token refreshed")
	}
	_, _, raw := seedRefresh(t, testDB.db, time.Now().UTC().Add(2*time.Minute).Truncate(time.Microsecond))
	rotated, err := svc.RefreshToken(context.Background(), raw)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := utils.ValidateAccessToken(rotated.Response.AccessToken, strings.Repeat("x", 32))
	if err != nil {
		t.Fatal(err)
	}
	if claims.ExpiresAt.Unix() != rotated.RefreshExpiresAt.Unix() || rotated.Response.ExpiresIn > 120 {
		t.Fatal("access token outlives absolute session expiry")
	}
	encoded, err := json.Marshal(rotated)
	if err != nil || strings.Contains(string(encoded), rotated.RefreshToken) {
		t.Fatal("refresh secret leaked through JSON")
	}
	duplicate := domain.RefreshToken{UserID: registered.Response.User.ID, TokenHash: utils.HashToken(registered.RefreshToken), ExpiresAt: time.Now().Add(time.Hour)}
	if err := testDB.db.Create(&duplicate).Error; err == nil {
		t.Fatal("duplicate refresh hash was accepted")
	}
}

func TestAuthMalformedRefreshDoesNotUseDatabase(t *testing.T) {
	svc := service.NewAuthService(nil, nil, nil, "", time.Minute, time.Hour)
	for _, raw := range []string{"", "short", strings.Repeat("z", 64)} {
		if session, err := svc.RefreshToken(context.Background(), raw); !errors.Is(err, domain.ErrInvalidRefreshToken) || session != nil {
			t.Fatal("malformed refresh accepted")
		}
	}
}

func TestAuthMigrationRejectsDuplicateHashesWithoutDeletingData(t *testing.T) {
	testDB := newAuthTestDB(t)
	user, token, _ := seedRefresh(t, testDB.db, time.Now().UTC().Add(time.Hour))
	if err := testDB.db.Exec("DROP INDEX idx_refresh_tokens_token_hash").Error; err != nil {
		t.Fatal(err)
	}
	duplicate := domain.RefreshToken{UserID: user.ID, TokenHash: token.TokenHash, ExpiresAt: token.ExpiresAt}
	if err := testDB.db.Create(&duplicate).Error; err != nil {
		t.Fatal(err)
	}
	script, err := os.ReadFile("../../migrations/001_refresh_token_hash_unique.sql")
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := testDB.db.DB()
	if err != nil {
		t.Fatal(err)
	}
	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if _, err := conn.ExecContext(context.Background(), string(script)); err == nil {
		t.Fatal("migration accepted duplicate hashes")
	}
	if _, err := conn.ExecContext(context.Background(), "ROLLBACK"); err != nil {
		t.Fatal(err)
	}
	var count int64
	if err := testDB.db.Model(&domain.RefreshToken{}).Count(&count).Error; err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatal("migration deleted conflicting records")
	}
}

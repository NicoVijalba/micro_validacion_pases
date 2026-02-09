package integration

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/validacion-pases/internal/domain"
	"github.com/example/validacion-pases/internal/repository/mysql"
	_ "github.com/go-sql-driver/mysql"
	mysqltc "github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestInsertRecordIntegration(t *testing.T) {
	ctx := context.Background()

	container, err := mysqltc.Run(ctx,
		"mysql:8.4",
		mysqltc.WithDatabase("validacion_pases"),
		mysqltc.WithUsername("app"),
		mysqltc.WithPassword("app"),
		mysqltc.WithScripts(filepath.Join("..", "..", "migrations", "000001_init.up.sql")),
		mysqltc.WithWaitStrategy(wait.ForLog("port: 3306  MySQL Community Server").WithStartupTimeout(2*time.Minute)),
	)
	if err != nil {
		t.Fatalf("mysql container failed: %v", err)
	}
	defer func() { _ = container.Terminate(ctx) }()

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "3306")
	dsn := fmt.Sprintf("app:app@tcp(%s:%s)/validacion_pases?parseTime=true", host, port.Port())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := mysql.NewRecordRepository(db)
	now := time.Now().UTC()
	id, err := repo.Insert(ctx, domain.Record{
		Emision:             now,
		Nave:                "NAVE 1",
		Viaje:               "VJ1",
		Cliente:             "CLIENTE 1",
		Booking:             "BK1",
		Rama:                "internacional",
		Contenedor:          "ABCU1234567",
		PuertoDescargue:     "Balboa",
		LibreRetencionHasta: now.AddDate(0, 0, 2),
		DiasLibre:           2,
		Transportista:       "",
		TituloTerminal:      "TERMINAL PACIFICO - BALBOA",
		UsuarioFirma:        "tester",
		CreatedAt:           now,
	})
	if err != nil {
		t.Fatalf("insert failed: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected id > 0")
	}
}

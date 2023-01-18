package store

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kokdot/go-musthave-devops/internal/metricsserver"
	"github.com/kokdot/go-musthave-devops/internal/repo"
	// "github.com/kokdot/go-musthave-devops/internal/repo"
)

var zeroG Gauge = 0
var zeroC Counter = 0
// type Metrics = repo.Metrics
type DbStorage struct {
	StoreMap      *StoreMap
	restore       bool
    storeFile   string
	storeInterval time.Duration
	key           string
	url           string
	dataBaseDSN   string
	dbconn        *sql.DB
}
func NewDbStorage(storeInterval time.Duration, storeFile string, restore bool, url string, key string, dataBaseDSN string) (*DbStorage, error){
    dbconn, err := sql.Open("pgx", dataBaseDSN)
	if err != nil {
		return nil, err
	}
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = dbconn.PingContext(ctx); err != nil {
		return nil, err
	}
    var sm = make(StoreMap, 0)
    var dbStorage =   DbStorage{
        StoreMap: &sm,
        storeInterval: storeInterval, 
		restore: restore,
        storeFile: storeFile,
		url: url,
		key: key,
		dataBaseDSN: dataBaseDSN,
        dbconn: dbconn,
    }
    if err = dbStorage.createStorage(); err != nil {
        return nil, err
    }

    return &dbStorage , nil
}

func (d DbStorage) SaveByBatch(sm []repo.Metrics) (*repo.StoreMap, error) {
// func (d DbStorage) SaveByBatch(sm *repo.StoreMap) (*repo.StoreMap, error) {
    smtx := make(repo.StoreMap)
    for _, val := range sm {
        mtx, err := d.Save(&val)
        if err != nil {
            return nil, err
        }
        smtx[val.ID] = *mtx
    }
    return &smtx, nil
}

func (d DbStorage) Save(mtxNew *Metrics) (*Metrics, error) {
    // var mtxOld *Metrics
    mtxOld, err := d.Get(mtxNew.ID)
    if err == nil && mtxNew.MType == "counter" {
            *mtxNew.Delta += *mtxOld.Delta
	}
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `INSERT INTO Metrics
    (
        ID, 
        MType, 
        Delta, 
        Value, 
        Hash
    ) values($1, $2, $3, $4, $5) ON CONFLICT (ID) DO UPDATE SET 
    ID = Metrics.ID,
    MType = Metrics.MType,
    Delta = EXCLUDED.Delta + Metrics.Delta, 
    Value = EXCLUDED.Value,
    Hash = EXCLUDED.Hash;
    `
    // ON CONFLICT (ID) DO UPDATE SET ID = Metrics.ID;

    // ON CONFLICT (a) DO UPDATE SET c = tablename.c + 1;
    // INSERT INTO tablename (a, b, c) values (1, 2, 10)
    _, err = d.dbconn.ExecContext(ctx, query, mtxNew.ID, mtxNew.MType, mtxNew.Delta, mtxNew.Value, mtxNew.Hash)
    if err != nil {
		return mtxNew, fmt.Errorf("не удалось выполнить запрос создания записи в таблице Metrics: %v", err)
	}
    // var mtxOld *Metrics
    mtxOld, err = d.Get(mtxNew.ID)
    if err != nil {
		return mtxNew, fmt.Errorf("не удалось выполнить запрос получения записи в таблице Metrics: %v", err)
	}
    return mtxOld, nil
}

func (d DbStorage) createStorage() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if !d.restore {
        query := `DROP TABLE IF EXISTS Metrics;`
        _, err := d.dbconn.ExecContext(ctx, query)
        if err != nil {
            return fmt.Errorf("не удалось выполнить запрос удаления таблицы Metrics: %v", err)
        }
    }
    query := `
		CREATE TABLE  IF NOT EXISTS  Metrics
        (
            ID VARCHAR(255) NOT NULL UNIQUE,
            MType VARCHAR(10) NOT NULL,
            Delta BIGINT,
            Value double precision,
            Hash VARCHAR(255)
        );
	`
    _, err := d.dbconn.ExecContext(ctx, query)
    if err != nil {
		return fmt.Errorf("не удалось выполнить запрос создания таблицы Metrics: %v", err)
	}
    return nil
}
func (d DbStorage) ReadStorage() error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `
        SELECT ID, MType, Delta, Value, Hash FROM Metrics;
    `
	rows, err := d.dbconn.QueryContext(ctx, query)
     if err != nil {
		return fmt.Errorf("не удалось выполнить запрос на полуыенин таблицы Metrics: %v", err)
	}
    defer rows.Close()
    var mtx Metrics
    var sm = make(StoreMap, 0)
    var delta sql.NullInt64
    var hash sql.NullString
    var value sql.NullFloat64
    for rows.Next() {
        err = rows.Scan(&mtx.ID, &mtx.MType, &delta, &value, &hash)
        if err != nil {
		    return fmt.Errorf("не удалось отсканировать строку запроса GetTable: %v", err)
	    }
        if value.Valid {
            value1 := Gauge(value.Float64)
            mtx.Value = &value1 
        } else {
            mtx.Value = &zeroG
        }
        if delta.Valid {
            delta1 := Counter(delta.Int64)  
            mtx.Delta = &delta1
        } else {
            mtx.Delta = &zeroC
        }
        if hash.Valid {
            hash1 := hash.String
            mtx.Hash = hash1  
        } else {
            mtx.Hash = ""
        }
        sm[mtx.ID] = mtx
    }
    *(d.StoreMap) = sm
    return nil   
}
func (d DbStorage) Get(id string) (*Metrics, error) {
     ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `
        SELECT ID, MType, Delta, Value FROM Metrics
        WHERE ID=$1
       `
    row := d.dbconn.QueryRowContext(ctx, query, id)

    var mtx Metrics
    var delta sql.NullInt64
    // var hash sql.NullString
    var value sql.NullFloat64
    err := row.Scan(&mtx.ID, &mtx.MType, &delta, &value)
    if err != nil {
        return nil, fmt.Errorf("не удалось отсканировать строку запроса GetMtx: %v", err)
    }
    if value.Valid {
        value1 := Gauge(value.Float64)
        mtx.Value = &value1 
    } else {
        mtx.Value = &zeroG
    }
    if delta.Valid {
        delta1 := Counter(delta.Int64)  
        mtx.Delta = &delta1
    } else {
        mtx.Delta = &zeroC
    }
    if d.key != "" {
        mtx.Hash = metricsserver.Hash(&mtx, d.key)
    } else {
        mtx.Hash = ""
    }
    err = row.Err()
    if err != nil {
        return nil, err
    }
    return &mtx, nil
}

func (d DbStorage) GetGaugeValue(id string) (Gauge, error) {
    mtxNew, err := d.Get(id)
    if err != nil {
        return zeroG, err
    }
    return *mtxNew.Value, nil
}
func (d DbStorage) GetCounterValue(id string) (Counter, error) {
    mtxNew, err := d.Get(id)
    if err != nil {
        return zeroC, err
    }
    return *mtxNew.Delta, nil
}
func (d DbStorage) GetDataBaseDSN() string {
	return d.dataBaseDSN
}
func (d DbStorage) GetStoreFile() string {
    return d.storeFile
}
func (d DbStorage) GetURL() string {
	return d.url
}
func (d DbStorage) GetRestore() bool {
	return d.restore
}

func (d DbStorage) GetKey() string {
	return d.key
}
func (d DbStorage) GetStoreInterval() time.Duration {
	return d.storeInterval
}

func (d DbStorage) GetPing() (bool, error) {
	// urlExample := "postgres://postgres:postgres@localhost:5432/postgres"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := d.dbconn.PingContext(ctx); err != nil {
		return false, err
	}
	fmt.Println("Ping Ok")
	return true, nil
}
func (d DbStorage) SaveCounterValue(name string, counter Counter) (Counter, error) {
    mtx := metricsserver.NewMetrics(name, "counter")
    mtxNew, err :=(d.Save(&mtx)) //Save(mtx)
    if err != nil {
        return counter, fmt.Errorf("%s", err)
    }
    return *mtxNew.Delta, nil
}
func (d DbStorage) SaveGaugeValue(name string, gauge Gauge) error {
    mtx := metricsserver.NewMetrics(name, "gauge")
    _, err :=(d.Save(&mtx)) 
    if err != nil {
        return fmt.Errorf("%s", err)
    }
    return nil
}
func (d DbStorage) GetAllValues() string {
    _, _ = d.GetAll()
    return repo.StoreMapToString(d.StoreMap)
}
func (d DbStorage) GetAll() (StoreMap, error) {
    err := d.ReadStorage()
    if err != nil {
        return nil, fmt.Errorf("%s", err)
    }
    return *d.StoreMap, nil
}
func (d DbStorage) WriteStorage() error {
    for _, val := range *d.StoreMap  {
        _, err := d.Save(&val)
        if err != nil {
            return fmt.Errorf("%s", err)
        }
    }
    return nil
}
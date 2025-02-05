package store

import (
	"context"
	"fmt"
	"time"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kokdot/go-musthave-devops/internal/metricsserver"
	"github.com/kokdot/go-musthave-devops/internal/repo"
)

var zeroG Gauge = 0
var zeroC Counter = 0
type DBStorage struct {
	StoreMap      *StoreMap
	restore       bool
    storeFile   string
	storeInterval time.Duration
	key           string
	url           string
	dataBaseDSN   string
	dbconn        *sql.DB
}
func NewDBStorage(storeInterval time.Duration, storeFile string, restore bool, url string, key string, dataBaseDSN string) (*DBStorage, error){
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
    var dbStorage =   DBStorage{
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
func (d DBStorage) SaveByBatch1(sm *repo.StoreMap) (*repo.StoreMap, error) {
    logg.Print("--------------------------------------------SaveByBatch----------------------------start-----------------------------------")
        // шаг 1 — объявляем транзакцию
    tx, err := d.dbconn.Begin()
    if err != nil {
        logg.Print("--------------------------------------------SaveByBatch----------------------------1-----------------------------------")
        return nil, err
    }
    // шаг 1.1 — если возникает ошибка, откатываем изменения
    defer tx.Rollback()

    // шаг 2 — готовим инструкцию
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
     query := `INSERT INTO Metrics
    (
        ID, 
        MType, 
        Delta, 
        Value
    ) values($1, $2, $3, $4) ON CONFLICT (ID) DO UPDATE SET 
    ID = Metrics.ID,
    MType = Metrics.MType,
    Delta = EXCLUDED.Delta + Metrics.Delta, 
    Value = EXCLUDED.Value
    `
    stmt, err := tx.PrepareContext(ctx, query)
    if err != nil {
        logg.Print("--------------------------------------------SaveByBatch----------------------------2-----------------------------------")
        return nil, err
    }
    // шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
    defer stmt.Close()

    for _, v := range *sm {
        // шаг 3 — указываем, что каждое видео будет добавлено в транзакцию
        if _, err = stmt.ExecContext(ctx, v.ID, v.MType, v.Delta, v.Value); err != nil {
            logg.Print("--------------------------------------------SaveByBatch----------------------------3-----------------------------------")
            return nil, err
        }
    }
    // шаг 4 — сохраняем изменения
    err = tx.Commit()
    if err != nil {
        logg.Print("--------------------------------------------SaveByBatch----------------------------4-----------------------------------")
        return nil, err
    }
    smtx := make(repo.StoreMap)
    for _, val := range *sm {
        mtx, err := d.Get(val.ID)
        if err != nil {
            logg.Print("--------------------------------------------SaveByBatch----------------------------5-----------------------------------")
           return nil, err
        }
        smtx[val.ID] = *mtx
    }
    logg.Print("--------------------------------------------SaveByBatch----------------------------finish-----------------------------------")
    return &smtx, nil
}

func (d DBStorage) SaveByBatch(sm []repo.Metrics) (*[]repo.Metrics, error) {
        // шаг 1 — объявляем транзакцию
    tx, err := d.dbconn.Begin()
    if err != nil {
        return nil, err
    }
    // шаг 1.1 — если возникает ошибка, откатываем изменения
    defer tx.Rollback()

    // шаг 2 — готовим инструкцию
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
     query := `INSERT INTO Metrics
    (
        ID, 
        MType, 
        Delta, 
        Value
    ) values($1, $2, $3, $4) ON CONFLICT (ID) DO UPDATE SET 
    ID = Metrics.ID,
    MType = Metrics.MType,
    Delta = EXCLUDED.Delta + Metrics.Delta, 
    Value = EXCLUDED.Value
    `
    stmt, err := tx.PrepareContext(ctx, query)
    if err != nil {
        return nil, err
    }
    // шаг 2.1 — не забываем закрыть инструкцию, когда она больше не нужна
    defer stmt.Close()

    for _, v := range sm {
        // шаг 3 — указываем, что каждое видео будет добавлено в транзакцию
        if _, err = stmt.ExecContext(ctx, v.ID, v.MType, v.Delta, v.Value); err != nil {
            return nil, err
        }
    }
    // шаг 4 — сохраняем изменения
    err = tx.Commit()
    if err != nil {
        return nil, err
    }
    smNew := make([]repo.Metrics, 0)
    for _, val := range sm {
        mtx, err := d.Get(val.ID)
        if err != nil {
           return nil, err
        }
        smNew = append(smNew, *mtx)
    }
    return &smNew, nil
}
func (d DBStorage) SaveByBatchOld(sm []repo.Metrics) (*repo.StoreMap, error) {
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

func (d DBStorage) Save(mtxNew *Metrics) (*Metrics, error) {
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
    _, err := d.dbconn.ExecContext(ctx, query, mtxNew.ID, mtxNew.MType, mtxNew.Delta, mtxNew.Value, mtxNew.Hash)
    if err != nil {
		return mtxNew, fmt.Errorf("не удалось выполнить запрос создания записи в таблице Metrics: %v", err)
	}
    var mtxOld *Metrics
    mtxOld, err = d.Get(mtxNew.ID)
    if err != nil {
		return mtxNew, fmt.Errorf("не удалось выполнить запрос получения записи в таблице Metrics: %v", err)
	}
    return mtxOld, nil
}

func (d DBStorage) createStorage() error {
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
func (d DBStorage) ReadStorage() error {
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
    err = rows.Err()
        if err != nil {
            return err
    }
    *(d.StoreMap) = sm
    return nil   
}
func (d DBStorage) Get(id string) (*Metrics, error) {
     ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
    query := `
        SELECT ID, MType, Delta, Value FROM Metrics
        WHERE ID=$1
       `
    row := d.dbconn.QueryRowContext(ctx, query, id)

    var mtx Metrics
    var delta sql.NullInt64
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
    if mtx.MType == "counter" {
        mtx.Value = nil
    } else {
        mtx.Delta = nil
    }
    err = row.Err()
    if err != nil {
        return nil, err
    }
    return &mtx, nil
}

func (d DBStorage) GetGaugeValue(id string) (Gauge, error) {
    mtxNew, err := d.Get(id)
    if err != nil {
        return zeroG, err
    }
    return *mtxNew.Value, nil
}
func (d DBStorage) GetCounterValue(id string) (Counter, error) {
    mtxNew, err := d.Get(id)
    if err != nil {
        return zeroC, err
    }
    return *mtxNew.Delta, nil
}
func (d DBStorage) GetDataBaseDSN() string {
	return d.dataBaseDSN
}
func (d DBStorage) GetStoreFile() string {
    return d.storeFile
}
func (d DBStorage) GetURL() string {
	return d.url
}
func (d DBStorage) GetRestore() bool {
	return d.restore
}

func (d DBStorage) GetKey() string {
	return d.key
}
func (d DBStorage) GetStoreInterval() time.Duration {
	return d.storeInterval
}

func (d DBStorage) GetPing() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := d.dbconn.PingContext(ctx); err != nil {
		return false, err
	}
	logg.Print("Ping Ok")
	return true, nil
}
func (d DBStorage) SaveCounterValue(name string, counter Counter) (Counter, error) {
    logg.Printf("Couunter: %v", counter)
    mtx := metricsserver.NewMetrics(name, "counter")
    mtx.Delta = &counter
    logg.Printf("mtx: %#v, ; Delta: %d", mtx, *mtx.Delta)
    mtxNew, err :=(d.Save(&mtx)) //Save(mtx)
    if err != nil {
        return counter, fmt.Errorf("%s", err)
    }
    logg.Printf("mtxNew: %#v, ; Delta: %d", mtxNew, *mtxNew.Delta)
    return *mtxNew.Delta, nil
}
func (d DBStorage) SaveGaugeValue(name string, gauge Gauge) error {
    mtx := metricsserver.NewMetrics(name, "gauge")
    mtx.Value = &gauge
    _, err :=(d.Save(&mtx)) 
    if err != nil {
        return fmt.Errorf("%s", err)
    }
    return nil
}
func (d DBStorage) GetAllValues() string {
    _, _ = d.GetAll()
    return repo.StoreMapToString(d.StoreMap)
}
func (d DBStorage) GetAll() (StoreMap, error) {
    err := d.ReadStorage()
    if err != nil {
        return nil, fmt.Errorf("%s", err)
    }
    return *d.StoreMap, nil
}
func (d DBStorage) WriteStorage() error {
    for _, val := range *d.StoreMap  {
        _, err := d.Save(&val)
        if err != nil {
            return fmt.Errorf("%s", err)
        }
    }
    return nil
}
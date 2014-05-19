package db

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/cache"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"logger"
	"os"
	"sync"
)

const (
	dbVersion     int = 2
	dbMaxTransCnt     = 20000
	dbMaxTransMem     = 64 * 1024 * 1024 // 64 MB
)

var (
	TxShaMissing   = errors.New("Requested transaction does not exist")
	DuplicateSha   = errors.New("Duplicate insert attempted")
	DbDoesNotExist = errors.New("Non-existent database")
	DbUnknownType  = errors.New("Non-existent database type")
	DbUnOpen       = errors.New("Non-open database")
)

var gpdb *LevelDb

type LevelDb struct {
	// lock preventing multiple entry
	dbLock sync.Mutex

	// leveldb pieces
	lDb *leveldb.DB
	ro  *opt.ReadOptions
	wo  *opt.WriteOptions

	lbatch *leveldb.Batch
}

func init() {
	var err error
	gpdb, err = CreateDB("tx.db")
	if err != nil {
		fmt.Println(err)
		logger.Infoln("init db failed!")
	}
}

// parseArgs parses the arguments from the btcdb Open/Create methods.
func parseArgs(funcName string, args ...interface{}) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Invalid arguments to ldb.%s -- "+
			"expected database path string", funcName)
	}
	dbPath, ok := args[0].(string)
	if !ok {
		return "", fmt.Errorf("First argument to ldb.%s is invalid -- "+
			"expected database path string", funcName)
	}
	return dbPath, nil
}

// OpenDB opens an existing database for use.
func OpenDB(args ...interface{}) (*LevelDb, error) {
	dbpath, err := parseArgs("OpenDB", args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(dbpath)
	db, err := openDB(dbpath, false)
	if err != nil {
		return nil, err
	}

	return db, nil
}

var CurrentDBVersion int32 = 1

func openDB(dbpath string, create bool) (pbdb *LevelDb, err error) {
	var db LevelDb
	var tlDb *leveldb.DB
	var dbversion int32

	defer func() {
		if err == nil {
			db.lDb = tlDb

			pbdb = &db
		}
	}()

	if create == true {
		err = os.Mkdir(dbpath, 0750)
		if err != nil {
			logger.Errorf("mkdir failed %v %v", dbpath, err)
			//return
		}
	} else {
		_, err = os.Stat(dbpath)
		if err != nil {
			err = DbDoesNotExist
			return
		}
	}

	needVersionFile := false
	verfile := dbpath + ".ver"
	fi, ferr := os.Open(verfile)
	if ferr == nil {
		defer fi.Close()

		ferr = binary.Read(fi, binary.LittleEndian, &dbversion)
		if ferr != nil {
			dbversion = ^0
		}
	} else {
		if create == true {
			needVersionFile = true
			dbversion = CurrentDBVersion
		}
	}

	myCache := cache.NewEmptyCache()
	opts := &opt.Options{
		BlockCache:   myCache,
		MaxOpenFiles: 256,
		Compression:  opt.NoCompression,
	}

	switch dbversion {
	case 0:
		opts = &opt.Options{}
	case 1:
		// uses defaults from above
	default:
		err = fmt.Errorf("unsupported db version %v", dbversion)
		return
	}

	tlDb, err = leveldb.OpenFile(dbpath, opts)
	if err != nil {
		return
	}

	// If we opened the database successfully on 'create'
	// update the
	if needVersionFile {
		fo, ferr := os.Create(verfile)
		if ferr != nil {
			// TODO(design) close and delete database?
			//err = ferr
			//return
		}
		defer fo.Close()
		err = binary.Write(fo, binary.LittleEndian, dbversion)
		if err != nil {
			return
		}
	}

	return
}

// CreateDB creates, initializes and opens a database for use.
func CreateDB(args ...interface{}) (*LevelDb, error) {
	dbpath, err := parseArgs("Create", args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(dbpath)
	// No special setup needed, just OpenBB
	db, err := openDB(dbpath, true)
	return db, err
}

func (db *LevelDb) close() {
	db.lDb.Close()
}

// Sync verifies that the database is coherent on disk,
// and no outstanding transactions are in flight.
func (db *LevelDb) Sync() {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	// while specified by the API, does nothing
	// however does grab lock to verify it does not return until other operations are complete.
}

// Close cleanly shuts down database, syncing all data.
func (db *LevelDb) Close() {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	db.close()
}

func int64ToKey(keyint int64) []byte {
	key := fmt.Sprintf("%d", keyint)
	return []byte(key)
}

func (db *LevelDb) lBatch() *leveldb.Batch {
	if db.lbatch == nil {
		db.lbatch = new(leveldb.Batch)
	}
	return db.lbatch
}

func (db *LevelDb) RollbackClose() {
	db.dbLock.Lock()
	defer db.dbLock.Unlock()

	db.close()
}

func GetTx(refid string) (cmd, id string, timestamp int64, amount, price string, err error) {
	if gpdb != nil {
		return gpdb.getTx(refid)
	} else {
		err = DbUnOpen
		return
	}
}

func SetTx(cmd, id string, timestamp int64, amount, price string) (err error) {
	if gpdb != nil {
		return gpdb.setTx(cmd, id, timestamp, amount, price)
	} else {
		err = DbUnOpen
		return
	}
}

func (db *LevelDb) getTx(refid string) (cmd, id string, timestamp int64,
	amount, price string, err error) {
	var buf []byte

	key := []byte(refid)
	buf, err = db.lDb.Get(key, db.ro)
	if err != nil {
		fmt.Println(err)
		fmt.Println("get id failed", key)
		return
	}

	fmt.Println(buf)

	dr := bytes.NewBuffer(buf)

	err = binary.Read(dr, binary.LittleEndian, &cmd)
	if err != nil {
		err = fmt.Errorf("Db Corrupt 1")
		return
	}
	err = binary.Read(dr, binary.LittleEndian, &id)
	if err != nil {
		err = fmt.Errorf("Db Corrupt 2")
		return
	}
	err = binary.Read(dr, binary.LittleEndian, &timestamp)
	if err != nil {
		err = fmt.Errorf("Db Corrupt 3")
		return
	}

	err = binary.Read(dr, binary.LittleEndian, &amount)
	if err != nil {
		err = fmt.Errorf("Db Corrupt 4")
		return
	}

	err = binary.Read(dr, binary.LittleEndian, &price)
	if err != nil {
		err = fmt.Errorf("Db Corrupt 5")
		return
	}

	return cmd, id, timestamp, amount, price, nil
}

func (db *LevelDb) setTx(cmd, id string, timestamp int64,
	amount, price string) error {

	var txW bytes.Buffer

	err := binary.Write(&txW, binary.LittleEndian, []byte(cmd))
	if err != nil {
		fmt.Println(err)
		err = fmt.Errorf("Write fail 1 ")
		return err
	}

	err = binary.Write(&txW, binary.LittleEndian, []byte(id))
	if err != nil {
		fmt.Println(err)
		err = fmt.Errorf("Write fail 2 ")
		return err
	}

	err = binary.Write(&txW, binary.LittleEndian, timestamp)
	if err != nil {
		err = fmt.Errorf("Write fail 3")
		return err
	}

	err = binary.Write(&txW, binary.LittleEndian, []byte(amount))
	if err != nil {
		err = fmt.Errorf("Write fail 4")
		return err
	}

	err = binary.Write(&txW, binary.LittleEndian, []byte(price))
	if err != nil {
		err = fmt.Errorf("Write fail 5")
		return err
	}

	key := []byte(id)
	db.lDb.Put(key, txW.Bytes(), db.wo)

	return nil
}

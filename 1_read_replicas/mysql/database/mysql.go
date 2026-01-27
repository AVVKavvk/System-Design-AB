package database

import (
	"hash/crc32"
	"log"
	"sync"

	"github.com/AVVKavvk/mysql-replicas/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	PrimaryDB *gorm.DB
	Replicas  []*gorm.DB // Slice to hold R1, R2, R3, ...
	once      sync.Once
)

func InitDB() {

	once.Do(func() {
		var err error
		// 1. Connection Strings
		dsnPrimary := "root:primary_password@tcp(localhost:3306)/replica_test?charset=utf8mb4&parseTime=True&loc=Local"
		// Ensure these are in a consistent order!
		replicaDSNs := []string{
			"root:replica_password@tcp(localhost:3307)/replica_test?charset=utf8mb4&parseTime=True&loc=Local", // Index 0 (R1)
			"root:replica_password@tcp(localhost:3308)/replica_test?charset=utf8mb4&parseTime=True&loc=Local", // Index 1 (R2)
			"root:replica_password@tcp(localhost:3309)/replica_test?charset=utf8mb4&parseTime=True&loc=Local", // Index 2 (R3)
		}

		// 2. Connect to Primary
		PrimaryDB, err = gorm.Open(mysql.Open(dsnPrimary), &gorm.Config{})
		if err != nil {
			log.Fatal("Failed to connect to primary:", err)
		}

		for _, dsn := range replicaDSNs {
			conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatalf("Failed to connect to replica %s: %v", dsn, err)
			}
			Replicas = append(Replicas, conn)
		}
		log.Println("Database connection established. Manual Sharding enabled.")

		// Migration on Primary
		if err := PrimaryDB.AutoMigrate(&models.User{}); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("Database connection established with Read/Write splitting enabled.")
	})
}

// func InitDB() {
// 	once.Do(func() {
// 		var err error

// 		// 1. Define Connection Strings
// 		// Primary (Write) - Uses Root (or a user with Write privs)
// 		dsnPrimary := "root:primary_password@tcp(localhost:3306)/replica_test?charset=utf8mb4&parseTime=True&loc=Local"

// 		// Replicas (Read)
// 		dsnReplica1 := "root:replica_password@tcp(localhost:3307)/replica_test?charset=utf8mb4&parseTime=True&loc=Local"
// 		dsnReplica2 := "root:replica_password@tcp(localhost:3308)/replica_test?charset=utf8mb4&parseTime=True&loc=Local"
// 		dsnReplica3 := "root:replica_password@tcp(localhost:3309)/replica_test?charset=utf8mb4&parseTime=True&loc=Local"

// 		// 2. Open Initial Connection to Primary
// 		MysqlDB, err = gorm.Open(mysql.Open(dsnPrimary), &gorm.Config{})
// 		if err != nil {
// 			log.Fatal("Failed to connect to primary database:", err)
// 		}

// 		// 3. Register DBResolver (The Magic Part)
// 		err = MysqlDB.Use(dbresolver.Register(dbresolver.Config{
// 			// Sources = Primary (Writes)
// 			Sources: []gorm.Dialector{mysql.Open(dsnPrimary)},

// 			// Replicas = Read Replicas (Reads)
// 			Replicas: []gorm.Dialector{
// 				mysql.Open(dsnReplica1),
// 				mysql.Open(dsnReplica2),
// 				mysql.Open(dsnReplica3),
// 			},

// 			// Policy: Randomly load balance between replicas
// 			Policy: dbresolver.RandomPolicy{},
// 		}))

// 		if err != nil {
// 			log.Fatal("Failed to register dbresolver:", err)
// 		}
// 		// 4. Auto Migrate (Only runs on Primary automatically)
// 		err = MysqlDB.AutoMigrate(&models.User{})
// 		if err != nil {
// 			log.Fatal("Migration failed:", err)
// 		}

// 		log.Println("Database connection established with Read/Write splitting enabled.")
// 	})
// }

func CheckDatabaseConnection(db *gorm.DB) {
	var serverID string

	// 1. Check READ Connection (Should be Replica ID: 2, 3, or 4)
	db.Raw("SELECT @@server_id").Scan(&serverID)
	log.Printf("READ Query handled by Server ID: %s", serverID)
}

// GetReplicaByKey selects a replica based on a hash of the input string
func GetReplicaByKey(key string) *gorm.DB {
	if len(Replicas) == 0 {
		return PrimaryDB // Fallback if no replicas
	}

	// 1. Hash the key (CRC32 is fast and good enough for distribution)
	hash := crc32.ChecksumIEEE([]byte(key))

	// 2. Modulo to find the index (0, 1, or 2)
	index := hash % uint32(len(Replicas))

	log.Println("Replica Index:", index)
	// 3. Return the specific replica
	return Replicas[index]
}

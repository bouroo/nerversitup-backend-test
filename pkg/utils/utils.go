package utils

import (
	"os"
	"strconv"

	"github.com/influxdata/influxdb/v2/pkg/snowflake"
	"golang.org/x/crypto/bcrypt"
)

// sortable unique IDs generator
var SnowflakeGen *snowflake.Generator

func init() {
	if SnowflakeGen == nil {
		// mock machine ID with process ID
		SnowflakeGen = snowflake.New(os.Getpid() % 1023)
	}
}

func SnowflakeId() uint64 {
	return SnowflakeGen.Next()
}

func SnowflakeIdStr() string {
	return SnowflakeGen.NextString()
}

func HashPassword(passwd string) (hashed []byte, err error) {
	hashed, err = bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	return
}

func CheckPassword(passwd, hashed string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(passwd))
	return
}

func LevelFromLocals(localsLevel interface{}) (level int) {
	strLevel, ok := localsLevel.(string)
	if !ok {
		return 0
	}
	level, _ = strconv.Atoi(strLevel)
	return
}

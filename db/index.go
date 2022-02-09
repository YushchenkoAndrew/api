package db

import "api/interfaces"

func Init(tables []interfaces.Table) {
	ConnectToDB(tables)
	ConnectToRedis(tables)
}

package main

import jobtest "gin-redis/extend-datatype/lock/jobs/src"

func main() {
	jobtest.Run()
	select {}
}

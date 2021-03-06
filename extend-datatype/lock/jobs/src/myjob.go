package jobtest

import (
	"fmt"
	"gin-redis/extend-datatype/lock"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func Run() {
	job("job1")
	job("job2")
}

// 一个经典场景：模拟分布式任务单点执行
func job(jobname string) {
	c := cron.New(cron.WithSeconds())
	id, err := c.AddFunc("0/5 * * * * *", func() {
		defer func() {
			if e := recover(); e != nil {
				log.Println(jobname, "执行失败:", e)
			}
		}()
		locker := lock.NewLockerWithTTL("job", time.Second*5).Lock()
		defer locker.Unlock()
		time.Sleep(time.Second * 2)
		db := GormDB.Exec("update test set v=v+1 where id=100")
		if db.Error != nil {
			log.Println(db.Error)
		} else {
			log.Println(jobname, "任务执行完毕")
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s任务ID是:%d 启动\n", jobname, id)
	c.Start()
}

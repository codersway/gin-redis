package chapter5

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	Conn *redis.Client
}

var (
	ctx = context.Background()
)

func NewClient(conn *redis.Client) *Client {
	return &Client{Conn: conn}
}

func (this Client) ImportIpsToRedis(filename string) {
	res := csvReader(filename)
	pipe := this.Conn.Pipeline()
	for count, row := range res {
		var (
			startIp string
			resIP   int64
		)
		if len(row) == 0 {
			startIp = ""
		} else {
			startIp = row[0]
		}
		if strings.Contains(strings.ToLower(startIp), "i") {
			continue
		}
		if strings.Contains(startIp, ".") {
			resIP = this.IpToScore(startIp)
		} else {
			var err error
			resIP, err = strconv.ParseInt(startIp, 10, 64)
			if err != nil {
				continue
			}
		}
		cityID := row[2] + "_" + strconv.Itoa(count)
		pipe.ZAdd(ctx, "ip2cityid:", &redis.Z{Member: cityID, Score: float64(resIP)})
		if (count+1)%1000 == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				log.Println("pipeline err in ImportIpsToRedis: ", err)
				return
			}
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("pipeline err in ImportIpsToRedis: ", err)
		return
	}
}

type CityInfo struct {
	CityID  string
	Country string
	Region  string
	City    string
}

func (this Client) ImportCityToRedis(filename string) {
	res := csvReader(filename)
	pipe := this.Conn.Pipeline()
	for count, row := range res {
		if len(row) < 4 || !IsDigit(row[0]) {
			continue
		}

		city := CityInfo{
			CityID:  row[0],
			Country: row[1],
			Region:  row[2],
			City:    row[3],
		}

		value, err := json.Marshal(city)
		if err != nil {
			log.Println("marshal json failed, err: ", err)
		}

		pipe.HSet(ctx, "cityid2city:", city.CityID, value)
		if (count+1)%1000 == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				log.Println("pipeline err in ImportCityToRedis: ", err)
				return
			}
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("pipeline err in ImportCityToRedis: ", err)
		return
	}
}

func (this Client) IpToScore(ip string) int64 {
	var score int64 = 0
	for _, v := range strings.Split(ip, ".") {
		n, _ := strconv.ParseInt(v, 10, 0)
		score = score*256 + n
	}
	return score
}

func (this Client) FindCityByIp(ip string) string {
	ipAddr := strconv.Itoa(int(this.IpToScore(ip)))
	res := this.Conn.ZRangeByScore(ctx, "ip2cityid:", &redis.ZRangeBy{
		Max:    ipAddr,
		Min:    "0",
		Offset: 0,
		Count:  2,
	}).Val()
	if len(res) == 0 {
		return ""
	}
	cityId := strings.Split(res[0], "_")[0]
	var result CityInfo
	if err := json.Unmarshal([]byte(this.Conn.HGet(ctx, "cityid2city:", cityId).Val()), &result); err != nil {
		log.Fatalln("unmarshal err: ", err)
	}
	return strings.Join([]string{result.CityID, result.City, result.Country, result.Region}, "")
}

func csvReader(filename string) [][]string {
	csvFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("open file fault, filename: %s, err: %v", filename, err)
	}
	file := csv.NewReader(csvFile)
	var res [][]string
	for {
		record, err := file.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("read csv fault, err: ", err)
		}
		res = append(res, record)
	}
	return res
}

func IsDigit(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return false
	}
	return true
}

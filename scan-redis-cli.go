package main

import (
	"flag"
	"fmt"
	"time"

	"strconv"

	"os"

	"github.com/garyburd/redigo/redis"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func main() {
	start := makeTimestamp()
	v := scanKeys()
	end := makeTimestamp()
	fmt.Printf("cost %v ms, total scan key = %v\n", (end - start), len(v))
}

var (
	server    string
	passwd    string
	keyLength int
	keyCount  int
	pattern   string
	showHelp  bool
	redisPool *redis.Pool
)

func scanKeys() []string {
	con := redisPool.Get()
	if con.Err() != nil {
		panic(con.Err())
	}
	defer con.Close()
	//SCAN cursor [MATCH pattern] [COUNT count]
	keys := make([]string, 0, 5000000)
	cursor := 0
	for {
		v, e := con.Do("scan", strconv.FormatInt(int64(cursor), 10), "match", pattern, "count", "10000")
		if e != nil {
			panic(e)
		}
		arr, e := redis.MultiBulk(v, e)
		if e != nil {
			panic(e)
		}
		cursor, _ = redis.Int(arr[0], nil)
		members, _ := redis.Strings(arr[1], nil)
		if len(members) > 0 {
			counter := 0
			for _, m := range members {
				isLengthOK := keyLength <= 0 || keyLength == len(m)
				isCountOK := keyCount <= 0 || counter <= keyCount
				if isLengthOK && isCountOK {
					fmt.Printf("key=>%v\n", m)
					counter++
					keys = append(keys, m)
				}
			}
		}
		if cursor == 0 {
			break
		}
	}
	return keys
}

func init() {
	flag.StringVar(&server, "h", "127.0.0.1:6379", "-h=IP地址:端口，默认为 127.0.0.1:6379")
	flag.StringVar(&passwd, "a", "", "-a=密码，默认为空")
	flag.StringVar(&pattern, "p", "*", "-p=匹配符，默认为*")
	flag.IntVar(&keyLength, "l", 0, "-l=固定key的长度，默认值为0，即不限")
	flag.IntVar(&keyCount, "c", 0, "-c=最大key个数，默认为0，即不限")

	flag.BoolVar(&showHelp, "help", false, "-help 显示该帮助")
	flag.Parse()

	fmt.Printf("连接的主机为 %v, 匹配模式为 %v, 固定的 key 长度为 %v, key结果个数限制为 %v\n", server, pattern, keyLength, keyCount)

	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	redisPool = newPool(server, passwd)

}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     500,
		MaxActive:   5000,
		IdleTimeout: 5 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}

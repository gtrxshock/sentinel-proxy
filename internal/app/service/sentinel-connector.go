package service

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/gtrxshock/sentinel-proxy/internal/app/core"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type SentinelConnector struct {
	sentinelClientList []*redis.SentinelClient
	requestTimeout     time.Duration
}

var AllSentinelsBrokenErr = errors.New("all sentinels unavailable")

func NewSentinelConnector(sentinelAddrList []string, requestTimeout int) *SentinelConnector {
	sentinelConnector := &SentinelConnector{requestTimeout: time.Duration(requestTimeout)}

	for _, sentinelAddr := range sentinelAddrList {
		sentinelConnector.AddSentinelClient(sentinelAddr)
	}

	return sentinelConnector
}

func (sc *SentinelConnector) AddSentinelClient(sentinelAddr string) {
	sentinelClient := redis.NewSentinelClient(&redis.Options{
		Addr: sentinelAddr,
	})

	sc.sentinelClientList = append(sc.sentinelClientList, sentinelClient)
}

func (sc *SentinelConnector) GetActualRedisAddr(dbName string) (string, error) {
	actualRedisMasterList := make(map[string]int)
	for _, sentinelClient := range sc.sentinelClientList {
		ctx, _ := context.WithTimeout(context.Background(), sc.requestTimeout*time.Second)
		res, err := sentinelClient.GetMasterAddrByName(ctx, dbName).Result()
		if err != nil {
			core.GetLogger().Errorf("get master redis addr from sentinel, db: %s, error: %s", dbName, err)
			// mark sentinel as broken
			continue
		}

		addr := strings.Join(res, ":")

		if _, ok := actualRedisMasterList[addr]; !ok {
			actualRedisMasterList[addr] = 1
		} else {
			actualRedisMasterList[addr]++
		}
	}

	if len(actualRedisMasterList) == 0 {
		return "", AllSentinelsBrokenErr
	}

	if len(actualRedisMasterList) > 1 {
		core.GetLogger().Errorf("cluster synchronization broken, quorum result: %v", actualRedisMasterList)

		var preferredMasterAddr string
		preferredMasterCnt := 0
		for addr, cnt := range actualRedisMasterList {
			if cnt > preferredMasterCnt {
				preferredMasterAddr = addr
			}
		}

		return preferredMasterAddr, nil
	}

	masterAddr := ""
	for addr := range actualRedisMasterList {
		masterAddr = addr
	}

	return masterAddr, nil
}

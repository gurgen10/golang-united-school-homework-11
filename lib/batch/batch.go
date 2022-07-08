package batch

import (
	"time"
	"sync"
	"fmt"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	var wg sync.WaitGroup
	var mx sync.Mutex
	sem := make(chan struct{}, pool)
	var i int64 = 0

	for ; i < n/pool; i++ {
		wg.Add(1)
		sem <- struct{}{}
		
			go func(id int64) {
				mx.Lock()
				user := getOne(id)
				res = append(res, user)
				mx.Unlock()
				<-sem
				wg.Done()
				}(int64(i))	
			
		}
		for i = n/pool; i < n; i++ {
			wg.Add(1)
			sem <- struct{}{}
			
				go func(id int64) {
					mx.Lock()
					user := getOne(id)
					res = append(res, user)
					mx.Unlock()
					<-sem
					wg.Done()
					}(int64(i))	
				
			}
			
		wg.Wait()
	return res
}

func GetUsers(n int64, pool int64) {
	u := getBatch(n, pool)

	fmt.Println(u)
}

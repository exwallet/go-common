/*
 * @Author: kidd
 * @Date: 1/14/19 5:52 PM
 */

package taskManager

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/exwallet/go-common/gologger"
	"os"
	"sync"
	"time"
)

/*
任务管理器
*/
var (
	signalInterrupt = make(chan os.Signal, 1)
	waitSegment     = time.Millisecond * 500
	log             = gologger.GetLogger()
)

type TaskManager struct {
	isRunning   bool                             // 运行状态
	tasksMap    map[*TaskInterface]chan struct{} // 用于通知任务退出信号
	stopWG      *sync.WaitGroup                  // 停止等待通知
	initWaitSec time.Duration                    // 运行前等待时间,默认为0
	tasksArr    *arraylist.List
}

func New(waitSec ...time.Duration) *TaskManager {
	tm := &TaskManager{
		isRunning:   false,
		tasksMap:    make(map[*TaskInterface]chan struct{}, 0),
		stopWG:      new(sync.WaitGroup),
		initWaitSec: 0,
		tasksArr:    arraylist.New(),
	}
	if waitSec != nil && len(waitSec) != 0 {
		tm.initWaitSec = time.Second * waitSec[0]
	}
	return tm
}

func (tm *TaskManager) IsRunning() bool {
	return tm.isRunning
}

func (tm *TaskManager) ListTasks() []string {
	var s []string
	for k, _ := range tm.tasksMap {
		s = append(s, (*k).String())
	}
	return s
}

func (th *TaskManager) HasTask(t string) bool {
	return th.tasksArr.Contains(t)
}

// 增加任务
func (tm *TaskManager) Add(task interface{}) bool {
	t, ok := task.(TaskInterface)
	if !ok {
		panic("添加任务失败: ")
	}
	if tm.tasksArr.Contains(t.String()) {
		log.Info("已存在重复任务:%s", t.String())
		return false
	}
	tm.tasksArr.Add(t.String())
	tm.tasksMap[&t] = make(chan struct{})

	return true
}

// 停止任务集
func (tm *TaskManager) Stop() {
	if tm.isRunning != true {
		return
	}
	for _, c := range tm.tasksMap {
		c <- struct{}{}
	}
	tm.isRunning = false
	// 等待所有task destroy完成
	tm.stopWG.Wait()
}

// 开始前重设变量
func (tm *TaskManager) resetBeforeStart() {
	//tm.stopChans = make([]chan int, len(tm.tasks))
}

// 开始任务集
func (tm *TaskManager) Start() {
	if tm.isRunning {
		return
	}
	tm.isRunning = true
	tm.resetBeforeStart()
	go func() {
		time.Sleep(tm.initWaitSec)
		wg := new(sync.WaitGroup)
		defer func() {
			log.Info("--> taskManager协程结束, 等待所有子任务结束...\n")
			wg.Wait()
		}()
		wg.Add(len(tm.tasksMap))        // Start goru等待信号
		tm.stopWG.Add(len(tm.tasksMap)) // Stop  goru等待信号
		for task, stopC := range tm.tasksMap {
			//tf := reflect.Indirect(reflect.ValueOf(task)).Type()
			// 运行任务
			go func(wg *sync.WaitGroup, task TaskInterface, stopC chan struct{}) {
				// 如果发生错误, 强制中断任务.
				task.Init()
				waitCount := task.DelayTime().Nanoseconds() / waitSegment.Nanoseconds()
				//log.Info("waitCount:%d", waitCount)
				for {
					select {
					case <-stopC:
						task.Destroy()
						wg.Done()
						tm.stopWG.Done()
						return
					default:
						e := task.Run()
						if e != nil { // 发生error强制退出
							task.Destroy()
							wg.Done()
							tm.stopWG.Done()
							log.Error("任务%s 发生错误导致退出", task)
							return
						}
						// 等待下一轮RUN时也要看有没有退出信号
						if waitCount == 0 {
							select {
							case <-stopC:
								task.Destroy()
								wg.Done()
								tm.stopWG.Done()
								return
							default:
								<-time.After(task.DelayTime())
							}
						} else {
							for i := int64(1); i < waitCount; i++ {
								select {
								case <-stopC:
									task.Destroy()
									wg.Done()
									tm.stopWG.Done()
									return
								default:
									<-time.After(waitSegment)
								}
							}
						}
					}
				}
			}(wg, *task, stopC)
		}
	}()

}

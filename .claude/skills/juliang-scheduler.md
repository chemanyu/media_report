---
name: juliang-scheduler
description: "定时发送巨量日报数据"
triggers:
  - "巨量日报"
  - "juliang"
  - "定时任务"
---

# 现阶段代码逻辑
1. ./service/api/etc/media-api.yaml 是配置文件，  
    JuliangDayReportCron: "10 12 * * *"  这个是巨量日报表配置
2. ./service/api/internal/script/report_scheduler.go 文件是我启动定时任务的文件
    config.Schedule.JuliangReportCron 是时报任务
    config.Schedule.JuliangDayReportCron 是日报任务
3. ./service/api/internal/script/juliang_scheduler.go 是巨量的定时获取任务文件
    executeJuliangReportJob 方法是主要逻辑

# 需要修改的方案
1.  config.Schedule.JuliangReportCron ,config.Schedule.JuliangDayReportCron 。这两个任务都需要调  用 executeJuliangReportJob 方法
2. 但是时报任务，executeJuliangReportJob 中 startTime 和 endTime 得是今天的时间
3. 日报任务，startTime 和 endTime 得是昨天的时间
4. 在 juliang_scheduler.go 的 第 340，341 行
    cost := parseNumber(account.StatCost)         // 消耗
	cashCost := parseNumber(account.StatCashCost) // 现金消耗
    这两个分别是 消耗和现金消耗
    在之后的 计算返点消耗 和 计算服务费 中，时报得用 cost，日报得用 cashCost
5. generateAndUploadExcelReport 生成文件的逻辑，看看分时报日报有没有影响
6. sendJuliangDingTalkNotification 发送钉钉消息的这个接口，如果是日报
    "#### 巨量时报  \n---\n" 这里得写 巨量日报

# 帮我改写代码

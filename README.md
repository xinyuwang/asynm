# asynm
Async Mission based on redis

### Redis Key Mission

#### mis_{id} 存储任务结构映射（hash）

submis_{no}, JSON<start, end, data, error>
count_all 总任务（int）
count_cur 已完成（int） （INCR key）
create_time 创建时间
finish_time 完成时间
expire_time 任务超时时间
mission_state 任务状态
mission_data 任务内容，JSON<自定义>

#### mis_{id}_record_log 任务记录（json)

key, JSON<start, end, data, num（拆分任务数）, expire（任务超时时间）, state（状态，进行中，已完成，已关闭，不存在）>

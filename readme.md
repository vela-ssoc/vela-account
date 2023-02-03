# 用户信息
> 获取系统账户信息 

## 账户信息
> vela.account <br />
> 支持:windows 和 linux

#### 账户字段
- name &emsp;&emsp;用户名
- home &emsp;&emsp;用户目录
- uid  &emsp;&emsp;用户uid
- gid  &emsp;&emsp;用户gid
- status &emsp;账户状态 启用或者禁用
- [ps(cnd)](获取该账进程) 获取账户对应进程
```lua
    local v = vela.account.all()
    for i=1,v.size do
        v[i].ps("exe re *bash*")
    end 
```

#### 扩展方法
- [vela.account.all(cnd)](#####全部账号) &emsp;&emsp;&emsp;获取全部账户
- [vela.account.current()](#####当前账号) &emsp;&emsp;&emsp;获取当前用户
- [vela.account.snapshot()](#####监控快照) &emsp;&emsp;&emsp;获取账户快照

##### 全部账号
> []v = vela.account.all(cnd) 返回Slice
```lua
    local v = vela.account.all()
    local v = vela.account.all("name = root")

    for i=1,v.size do
       print(v[i]) 
    end 
```

##### 当前账号
> v = vela.account.current() 返回当前用户account
```lua
    local v = vela.account.current()
    print(v)
```

##### 账户快照
> v = vela.account.snapshot() 生成快照 <br />
> 用来监控账户的新增、修改、删除动作 <br />
> 参数: account
- v.sync() &emsp;同步数据到管理中心
- v.on_delete(pipe) &emsp;删除账号
- v.on_create(pipe) &emsp;新建账户
- v.on_update(pipe) &emsp;账户更新
- v.poll(time)   &emsp;监控周期 单位:秒
```lua
    local s = vela.account.snapshot()
    s.sync()
    v.on_delete(function(account) end)
    v.on_create(function(account) end)
    v.on_update(function(account) end)
    v.poll(4)
```

## 组信息
> vela.group <br />
> 跟上面的接口使用类似 支持:windows 和 linux

#### 组字段
- name &emsp;&emsp;名称
- gid  &emsp;&emsp;id

#### 扩展方法
- [vela.group.all(cnd)](#####全部组) &emsp;&emsp;&emsp;获取组
- [vela.group.snapshot()](#####组快照) &emsp;&emsp;&emsp;获取组快照

##### 全部组
> []v = vela.group.all(cnd) 返回Slice
```lua
    local v = vela.group.all()
    local v = vela.group.all("name = root")

    for i=1,v.size do
       print(v[i]) 
    end 
```

#### 组快照
> v = vela.group.snapshot()
```lua
    local s = vela.group.snapshot()
    s.sync()
    v.on_delete(function(group) end)
    v.on_create(function(group) end)
    v.on_update(function(group) end)
    v.poll(4)
```

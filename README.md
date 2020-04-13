### GameTools
游戏中常用的一些工具
- excel配置表导出为纯文本格式的配置表，以及对应的.d.ts文件
- 合并纯文本格式的配置表
- 配置表自动生成.d.ts文件
- 将文件的文件名转成全小写
- 可视化操作

## 使用
##### 1.配置表导出
1.将assets文件夹下的app.exe放到任意位置，打开app.exe,新建一个项目，选择excel配置表所在目录，选择数据导出目录，选择定义文件导出目录。

2.点击导出，导出成功后，在指定的目录可以找到生成的文件

3.将导出的文件放到h5游戏项目中,将assets文件夹下的configs.ts放到h5游戏项目源码中。

4.示例代码片段如下：
```javascript 1.8
await configs.load("resource/out/config.txt");
console.log("activity#1.name", ConfigData.activity[1].name);
console.log("playerInit.hp", ConfigData.playerInit.hp);
```


#### 项目环境配置
配置环境变量GOPATH\
值配置当前项目的路径, 当前项目就会处于活跃状态，导入的依赖包才会下载到当前项目中

#### 项目目录

|--GameTools
- |--assets 测试目录，源文件和测试文件都放在这个目录下
- |--src golang源码文件夹
- |--|--app 可视化工具源码
  


---
注意点
- 如果使用命令行运行某个golang文件，则运行时的目录是在该文件的同一级
- 如果使用idea运行某个golang文件，则运行时的目录是idea所配置的目录下
- 拥有main方法的文件只能在main包下运行
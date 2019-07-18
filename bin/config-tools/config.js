/**
 * create by wangcheng on 2019/4/15 20:51
 * 这里文件的解析方法不依赖任何引擎
 */
var zero;
(function (zero) {
    /**
     * 解析配置表文件中的数据
     * 合并的配置表是将每个单独的配置表文件内容合并，并通过2个换行符分隔构成的
     * 单独的配置表格式：
     * tableName
     * nameTypeDef
     * dataRow
     *
     * @param data 配置表文件中的原始数据
     * @param collector 解析后的数据的集合，所有配置表结构体将放到这个对象中
     */
    function parseConfigData(data, collector) {
        var configs = data.split(/\n\n/);
        if (configs.length <= 0) {
            console.error("加载配置表失败,行数不正确");
        }
        for (var _i = 0, configs_1 = configs; _i < configs_1.length; _i++) {
            var c = configs_1[_i];
            var ret = parseConfig(c);
            if (ret) {
                collector[ret.filename] = ret.config;
            }
            else {
                console.error("某张配置表无法正常加载,已经忽略");
            }
        }
    }
    zero.parseConfigData = parseConfigData;
    /**
     * 解析一份配置表数据
     * 每份配置表文件中，第一行是文件名,第二行是配置表对象的属性名和类型,接下来的行是数据
     * @param data 配置表内容
     */
    function parseConfig(data) {
        var lines = data.split("\n");
        var numRow = lines.length;
        if (numRow < 3) {
            return;
        }
        //文件名
        var filename = lines[0];
        //属性|类型 定义行
        var lineNameType = lines[1];
        //数据起始行
        var start = 2;
        //解析 属性名和类型
        var types = [];
        var names = [];
        var nameTypes = lineNameType.split(" ");
        for (var _i = 0, nameTypes_1 = nameTypes; _i < nameTypes_1.length; _i++) {
            var strNT = nameTypes_1[_i];
            var nt = strNT.split("|");
            names.push(nt[0]);
            types.push(nt[1]);
        }
        //将配置表的每一行保存到对象中，通过第一列的值映射每一行
        var config = {};
        var numColumn = names.length;
        for (var i = start; i < numRow; i++) {
            var data_1 = lines[i].split(" ");
            var obj = {};
            for (var j = 0; j < numColumn; j++) {
                obj[names[j]] = convert(data_1[j], types[j]);
            }
            var key = obj[names[0]];
            config[key] = obj;
        }
        return {
            config: config,
            filename: filename
        };
    }
    /**
     * 数据类型转换
     * @param value 字符串类型的数据
     * @param type 类型
     */
    function convert(value, type) {
        if (type == 'number')
            return +value;
        if (type == 'boolean')
            return !!+value;
        value = decodeURIComponent(value);
        if (type == "string")
            return value;
        if (type == "number[]")
            return value.split(",").map(function (v) { return +v; });
        if (type == "boolean[]")
            return value.split(",").map(function (v) { return v == 'true'; });
        if (type == "string[]")
            return value.split(",");
    }
    // function getFileNameFromPath(path: string): string {
    //     return path.match(/\/(\w+)\.\w+/)[1];
    // }
})(zero || (zero = {}));

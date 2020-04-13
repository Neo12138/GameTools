namespace configs {
    export async function load(url: string, namespace: string = "ConfigData") {
        let res = await fetch(url);
        let data = await res.text();

        window[namespace] = {};
        parseConfigData(decode(data), window[namespace]);
        console.log("data: ", window[namespace]);
    }

    function decode(data: string): string {
        return data.split("").map(v => {
            let code = v.charCodeAt(0);
            if (code <= 128) {
                code--;
            }
            return String.fromCharCode(code)
        }).join("");
    }


    /**
     * 解析配置表文件中的数据
     * 合并的配置表是将每个单独的配置表文件内容合并，并通过2个换行符分隔构成的
     * 单独的配置表格式：
     * tableName
     * nameDef
     * typeDef
     * dataRow
     *
     * @param data 配置表文件中的原始数据
     * @param collector 解析后的数据的集合，所有配置表结构体将放到这个对象中
     */
    function parseConfigData(data: string, collector: object): void {
        let configs = data.split(/\n\n/);
        configs = configs.filter(v => !!v);
        if (configs.length <= 0) {
            console.error("加载配置表失败,行数不正确")
        }
        for (let c of configs) {
            let ret = parseConfig(c);
            if (ret) {
                collector[ret.filename] = ret.config;
            } else {
                console.error("某张配置表无法正常加载,已经忽略");
            }
        }
    }

    /**
     * 解析一份配置表数据
     * 每份配置表文件中，第一行是文件名,第二行是配置表对象的属性名和类型,接下来的行是数据
     * @param data 配置表内容
     */
    function parseConfig(data: string): { config: object, filename: string } {
        let lines = data.split("\n");
        let numRow = lines.length;
        if (numRow < 1) {
            console.error("数据格式不正确 filename:" + lines[0]);
            return;
        }
        //文件名 类型
        let [filename, type] = lines[0].split(" ");
        let config: object;
        if (type == "object") {
            config = parseObjectConfig(lines);
        } else if (type == "object_array") {
            config = parseObjectArrayConfig(lines)
        }
        if (!config) return;

        return {
            config: config,
            filename: filename
        };
    }

    function parseObjectConfig(lines: string[]): object {
        let numRow = lines.length;
        //数据起始行
        let start: number = 1;

        //将配置表的每一行保存到对象中，通过第一列的值映射每一行
        let config = {};
        for (let i = start; i < numRow; i++) {
            let data = lines[i].split(" ");
            let type = data[0];
            let name = data[1];
            let value = data[2];
            config[name] = convert(value, type);
        }
        return config;
    }

    function parseObjectArrayConfig(lines: string[]): object {
        let numRow = lines.length;
        if (numRow < 4) {
            console.error("数据格式不正确 filename:" + lines[0]);
            return;
        }
        //属性类型
        let types: string[] = lines[1].split(" ");
        //属性名
        let names: string[] = lines[2].split(" ");
        //数据起始行
        let start: number = 3;

        //将配置表的每一行保存到对象中，通过第一列的值映射每一行
        let config = {};
        let numColumn = names.length;
        for (let i = start; i < numRow; i++) {
            let data = lines[i].split(" ");
            let obj = {};
            for (let j = 0; j < numColumn; j++) {
                obj[names[j]] = convert(data[j], types[j]);
            }
            let key = obj[names[0]];
            config[key] = obj;
        }
        return config
    }

    /**
     * 数据类型转换
     * @param value 字符串类型的数据
     * @param type 类型
     */
    function convert(value: string, type: string): any {
        if (type == 'number') return +value;
        if (type == 'boolean') return !!+value;
        value = decodeURIComponent(value);
        if (type == "string") return value;
        if (type == "number[]") return value.split("|").map(v => +v);
        if (type == "boolean[]") return value.split("|").map(v => v == 'true');
        if (type == "string[]") return value.split("|");
    }
}

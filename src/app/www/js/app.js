function checkPath(path) {
    return /^[A-z]:[^:*?"\|<>]+$/.test(path)
}

function generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
        let r = Math.random() * 16 | 0, v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

function getTimestamp() {
    let d = new Date();
    return `${paddingLeft(d.getHours())}:${paddingLeft(d.getMinutes())}:${paddingLeft(d.getSeconds())}.${d.getMilliseconds()}`;
}

function paddingLeft(n) {
    return n < 10 ? "0" + n : "" + n;
}

async function doGoMethod(method, ...args) {
    if (window[method]) return await window[method](...args)
}

let v = new Vue({
    el: "#app",
    data: {
        popup: false,
        projects: [],

        tempProject: {
            type: 0,
            id: "",
            name: "新建项目",
            source: "",
            dataOutDir: "",
            dataFileName: "config",
            defineOutDir: "",
            defineFileName: "config",
            defineNamespace: "ConfigData",
            createTime: 0,
            isExporting: false
        },
        hasNullValue: false,
        toast: false,
        isToastOut: false,
        toastMsg: "",
        log: false,
        logs: [],
        operateLog: ""
    },
    async mounted() {
        let data = await doGoMethod("loadData");
        if (data && data.length) {
            data.forEach(v => this.projects.push(v))
        }
        this.projects.push({
            type: 0,
            id: "",
            name: "新建项目",
            source: "",
            dataOutDir: "",
            dataFileName: "config",
            defineOutDir: "",
            defineFileName: "config",
            defineNamespace: "ConfigData",
            createTime: 0,
            isExporting: false
        });

    },
    watch: {
        tempProject: function () {
            let hasNullValue = false;
            for (let k in this.tempProject) {
                if (this.tempProject[k] === "") {
                    hasNullValue = true;
                    break;
                }
            }
            this.hasNullValue = hasNullValue;
        }
    },
    methods: {
        async exportExcel(index) {
            this.projects[index].isExporting = true;
            this.addLog(index, "点击导出");
            let ret = await doGoMethod("exportExcel", this.projects[index]);

            if (ret && ret.code) {
                this.addLog(index, `${ret.msg}`);
                this.showToast(`${this.projects[index].name} 导出失败`);
            } else {
                this.addLog(index, `导出成功`);
                this.showToast(`${this.projects[index].name} 导出成功`);
            }
            this.projects[index].isExporting = false;
            this.$set(this.projects, "" + index, this.projects[index]);
        },
        setting(index) {
            this.addLog(index, `点击设置`);
            this.tempProject = {};
            let proj = this.projects[index];
            for (let k in proj) {
                this.tempProject[k] = proj[k];
            }
            this.popup = true;
        },
        explore(index) {
            console.log(`${this.projects[index].name} 浏览文件夹`);
            this.log = true;
            this.operateLog = this.logs[index];
        },
        async deleteProject(index) {
            this.addLog(index, `点击删除`);
            let ret = await doGoMethod("deleteProject", this.tempProject.id);

            if (ret && ret.code) {
                this.addLog(index, `${ret.msg}`);
                this.showToast(`${this.projects[index].name} 删除失败`);
            } else {
                this.showToast(`${this.projects[index].name} 删除成功`);
                this.projects.splice(index, 1);
                this.logs.splice(index, 1)
            }
        },
        closePopup() {
            this.popup = false;
        },
        async applyModify() {
            if (this.tempProject.id === "") {
                this.tempProject.type = 1;
                this.tempProject.id = generateUUID();
                this.tempProject.createTime = Date.now();
                let ret = await doGoMethod("addProject", this.tempProject);
                if (ret && ret.code) {
                    this.showToast(`${this.tempProject.name} 创建失败`);
                } else {
                    this.showToast(`${this.tempProject.name} 创建成功`);
                    this.projects.push(this.tempProject);
                }
            } else {
                let i = this.projects.findIndex(v => v.id === this.tempProject.id);
                let ret = await doGoMethod("modifyProject", this.tempProject);
                if (ret && ret.code) {
                    this.addLog(i, `${ret.msg}`);
                    this.showToast(`${this.tempProject.name} 修改失败`);
                } else {
                    this.addLog(i, `项目修改成功`);
                    this.showToast(`${this.tempProject.name} 修改成功`);
                    this.projects[i] = this.tempProject;
                }
            }
            this.projects.sort((a, b) => {
                if (a.type === b.type) return a.createTime - b.createTime;
                return b.type - a.type;
            });
            this.closePopup();
        },
        checkInputPath(str) {
            if (str !== "" && !checkPath(str)) {
                this.showToast("路径不正确")
            }
        },
        showToast(msg, delay) {
            if (delay === void 0) delay = 2000;
            this.isToastOut = false;
            this.toastMsg = msg;
            this.toast = true;
            setTimeout(() => {
                this.isToastOut = true
            }, delay)
        },
        closeLogPanel() {
            this.log = false;
        },
        addLog(index, msg) {
            let newLog = `${getTimestamp()} ${msg}\n`;
            if (this.logs[index]) {
                this.logs[index] += newLog;
            } else {
                this.logs[index] = newLog;
            }
        }
    }
})
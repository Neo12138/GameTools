//由工具自动生成，请勿手动修改
declare namespace config {
	interface Iactivity {
		/** ID */
		readonly id: number;
		/** 活动名称 */
		readonly name: string;
		/** 是否生效 */
		readonly take_effect: number;
		/** 平台生效 */
		readonly platform: number;
		/** 显示位置 */
		readonly position: number;
		/** 活动图标/banner */
		readonly icon: string[];
		/** 新玩家是否可见 */
		readonly delay: number;
		/** 活动类型 */
		readonly type: number;
		/** 开启时间 */
		readonly start: number;
		/** 结束时间 */
		readonly end: number;
		/** 显示顺序 */
		readonly order: number;
		/** HOT标识 */
		readonly hot: number;
		/** 功能参数 */
		readonly parameter: string;
		/** 功能参数说明 */
		readonly des: string;
	}

	interface Ilanding {
		/** ID */
		readonly id: number;
		/** 天数 */
		readonly day: number;
		/** 奖励 */
		readonly gift: number;
	}

	interface Iseven_days {
		/** ID */
		readonly id: number;
		/** 类型 */
		readonly type: number;
		/** 天数 */
		readonly day: number;
		/** 奖励 */
		readonly gift: number;
	}

	interface Iai_battle {
		/** ID */
		readonly id: number;
		/** 未命中 */
		readonly miss: number;
		/** 人-头 */
		readonly head: number;
		/** 人-身体 */
		readonly body: number;
		/** 马-头 */
		readonly horse_head: number;
		/** 马-身体 */
		readonly horse_body: number;
		/** 弱点 */
		readonly weak: number;
		/** 拉满弓概率 */
		readonly full: number;
		/** 额外瞄准时间最小 */
		readonly time_min: number;
		/** 额外瞄准时间最大 */
		readonly time_max: number;
	}

	interface Ibattle {
		/** ID */
		readonly id: number;
		/** 部位 */
		readonly type: number;
		/** 上限 */
		readonly max: number;
		/** 下限 */
		readonly min: number;
	}

	interface Ibox {
		/** ID */
		readonly id: number;
		/** 图标 */
		readonly icon: string;
		/** 开启时间 */
		readonly time: number;
		/** 掉落ID */
		readonly drop_id: number;
		/** 奖励数量 */
		readonly number: number;
		/** 是否重复掉落 */
		readonly repeat: number;
	}

	interface IchangeServer {
		/** ID */
		readonly id: number;
		/** 名字 */
		readonly name: string;
		/** 地址 */
		readonly ip: string;
		/** 端口 */
		readonly port: string;
	}

	let activity: { [key: number]: Iactivity };
	let landing: { [key: number]: Ilanding };
	let seven_days: { [key: number]: Iseven_days };
	let ai_battle: { [key: number]: Iai_battle };
	let battle: { [key: number]: Ibattle };
	let box: { [key: number]: Ibox };
	let changeServer: { [key: number]: IchangeServer };
}
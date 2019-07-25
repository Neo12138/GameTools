//由工具自动生成，请勿手动修改
declare namespace ConfigData {
	interface IActivity {
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
		readonly icon: string;
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

	interface IAIBattle {
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

	interface IBattle {
		/** ID */
		readonly id: number;
		/** 部位 */
		readonly type: number;
		/** 上限 */
		readonly max: number;
		/** 下限 */
		readonly min: number;
	}

	interface IBox {
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

	interface IChangeServer {
		/** ID */
		readonly id: number;
		/** 名字 */
		readonly name: string;
		/** 地址 */
		readonly ip: string;
		/** 端口 */
		readonly port: string;
	}

	interface ILanding {
		/** ID */
		readonly id: number;
		/** 天数 */
		readonly day: number;
		/** 奖励 */
		readonly gift: number;
	}

	interface ISevenDays {
		/** ID */
		readonly id: number;
		/** 类型 */
		readonly type: number;
		/** 天数 */
		readonly day: number;
		/** 奖励 */
		readonly gift: number;
	}

	let Activity: { [key: number]: IActivity };
	let AIBattle: { [key: number]: IAIBattle };
	let Battle: { [key: number]: IBattle };
	let Box: { [key: number]: IBox };
	let ChangeServer: { [key: number]: IChangeServer };
	let Landing: { [key: number]: ILanding };
	let SevenDays: { [key: number]: ISevenDays };
}
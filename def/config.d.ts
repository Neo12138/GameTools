//由工具自动生成，请勿手动修改
declare namespace ConfigData {
	interface IActivity {
		/** ID */
		readonly id: number;
		/** 活动名称 */
		readonly name: string;
		/** 是否生效 */
		readonly takeEffect: number;
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

	let activity: { [key: number]: IActivity};

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

	let battle: { [key: number]: IBattle};

	interface IPlayerInit {
		/** 初始生命值 */
		readonly hp: number;
		/** 初始攻击力 */
		readonly attack: number;
		/** 初始法术强度 */
		readonly spellPower: number;
		/** 初始魔法值 */
		readonly mp: number;
		/** 初始移速 */
		readonly speed: number;
	}

	let playerInit: IPlayerInit;

}
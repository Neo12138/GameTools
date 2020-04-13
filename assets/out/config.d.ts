//由工具自动生成，请勿手动修改
declare namespace ConfigData1 {
	interface IActivity {
		/** ID */
		readonly id: number;
		readonly name: string;
		readonly takeEffect: number;
		readonly platform: number;
		readonly position: number;
		readonly icon: string;
		readonly delay: number;
		readonly type: number;
		readonly start: number;
		readonly end: number;
		readonly order: number;
		readonly hot: number;
		readonly parameter: string;
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
/*
 * @Author: cloudyi.li
 * @Date: 2023-06-07 11:57:56
 * @LastEditTime: 2023-06-27 15:08:49
 * @LastEditors: cloudyi.li
 * @FilePath: /chatserver-api/pkg/tiktoken/tiktoken_test.go
 */
package tiktoken

import (
	"fmt"
	"testing"
)

func TestNumTokensSingleString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		// wantNum_tokens int
	}{
		// TODO: Add test cases.
		{
			name: "tset1",
			args: args{
				str: `李强是一名中国政治人物，现任中共中央政治局常委、国务院总理、党组书记。他曾担任浙江省人民政府省长，中共江苏省委书记、省人大常委会主任，中共上海市委书记等职务。李强在地方工作期间，被认为是浙江新一轮政府改革的主要设计者和推动者，其主导创立的浙江政务服务网成为了全国“互联网+政务服务”样本。在国家层面上，他也负责推动中国大陆走出疫情封控困难的问题。同时，他与现任中共中央总书记习近平长期密切共事，也因此被外界视为习近平“之江新军”的主要人物。
				`,
			},
			// wantNum_tokens: 707,
		},
		{
			name: "tset3",
			args: args{
				str: `北京高考696分以上考生104人，前20名成绩暂不公布
				6月25日上午，北京开放2023年高考成绩查询入口，并公布高招录取分数线及高考考生分数分布表。澎湃新闻（www.thepaper.cn）注意到，今年北京普通本科录取控制分数线为448分，较去年（425分）上涨了23分；特殊类型招生控制分数线为527分，较去年（518分）上涨了9分；普通专科录取控制分数线与去年相同，为120分（语数外三科总分）。同时，今年北京继续采取暂不公布排名前20名考生成绩的做法，不影响其志愿填报和录取。相关考生进行成绩查询时，将提示“祝贺你高考取得全市前20名的优异成绩！”。2023年北京市高考考生分数分布表显示，今年北京高考696分以上考生104人，650分以上考生累计2754人，600分以上考生累计10348人；达到特殊类型招生控制分数线、即527分及以上考生累计25551；达到普通本科录取控制分数线、即448分及以上考生累计42347人。2022年北京市高考考生
				分数分布表则显示，北京高考700以上考生106人，650分以上考生累计2744人，600分以上考生累计8671人；达到特殊类型招生控制分数线、即518分及以上考生累计22413人；达到普通本科录取控制分数线、即425分及以上考生累计37319人。北京市教委公布的数据显示，今年北京市参加6月高考考生共计5.8万余人，2022年为5.4万余人。澎湃新闻从北京市教委了解到，今年北京统考考生的本科志愿填报和单考考生单招志愿填报工作将于6月27日至7月1日进行，专科志愿填报工作于7月21日至22日进行。考生志愿填报实行网上填报方式，考生须在规定时间内登录北京教育考试院网站（www.bjeea.edu.cn；www.bjeea.cn）进行填报。今年，北京继续坚持首善标准，不断优化考生服务，进一步畅通考生咨询和申诉渠道。为了做好广大考生的咨询宣传、答疑解惑工作，北京分别在考试院和各区高招办设立高招录
				取咨询接待点，考试院咨询电话为89193989转1键，安排专人接听和接待考生及家长的咨询和申诉，依法依规维护考生合法权益。同时，北京加强录取组织管理，严把录取公平关，严格执行招生工作“八项基本要求”和“30个不得”招生工作禁令，在录取现场设立监察小组；邀请部分市人大代表、市政协委员检查录取工作，接受社会监督`,
			},
			// wantNum_tokens: 1720,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Print(NumTokensSingleString(tt.args.str))
		})
	}
}

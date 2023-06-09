package config

const (
	// ExchangeRateGatherUrl 欧元汇率采集网站
	ExchangeRateGatherUrl = "https://www.nbb.be/en/about-national-bank/eurosystem/exchange-rates"
	// ExchangeRateHtmlTagSelector 欧元汇率网站HTML标签选择器
	ExchangeRateHtmlTagSelector = "#block-nbb-exchange-rates-exchange-rates-full tbody > tr"

	// ExchangeRateRequestAPIForNl 欧元汇率请求API（荷兰）
	ExchangeRateRequestAPIForNl = "https://www.belastingdienst.nl/data/douane_wisselkoersen/wks.douane.wisselkoersen.dd202304.xml"
)

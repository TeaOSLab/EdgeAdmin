Tea.context(function () {
    this.$delay(function () {
        let that = this

        let countryUnit = this.processMaxUnit(this.countryStats)
        this.reloadChart("country-chart", "地区", this.countryStats, function (v) {
            return v.country.name
        }, function (args) {
            return that.countryStats[args.dataIndex].country.name + ": " + teaweb.formatNumber(that.countryStats[args.dataIndex].rawCount)
        }, countryUnit)

        let provinceUnit = this.processMaxUnit(this.provinceStats)
        this.reloadChart("province-chart", "省市", this.provinceStats, function (v) {
            return v.province.name
        }, function (args) {
            return that.provinceStats[args.dataIndex].country.name + ": " + that.provinceStats[args.dataIndex].province.name + " " + teaweb.formatNumber(that.provinceStats[args.dataIndex].rawCount)
        }, provinceUnit)

        let cityUnit = this.processMaxUnit(this.cityStats)
        this.reloadChart("city-chart", "城市", this.cityStats, function (v) {
            return v.city.name
        }, function (args) {
            return that.cityStats[args.dataIndex].country.name + ": " + that.cityStats[args.dataIndex].province.name + " " + that.cityStats[args.dataIndex].city.name + " " + teaweb.formatNumber(that.cityStats[args.dataIndex].rawCount)
        }, cityUnit)

        window.addEventListener("resize", function () {
            that.resizeChart("country-chart")
            that.resizeChart("province-chart")
            that.resizeChart("city-chart")
        })
    })

    this.reloadChart = function (chartId, name, stats, xFunc, tooltipFunc, unit) {
        let chartBox = document.getElementById(chartId)
        if (chartBox == null) {
            return
        }
        let chart = echarts.init(chartBox)
        let option = {
            xAxis: {
                data: stats.map(xFunc),
                axisLabel: {
                    interval: 0
                }
            },
            yAxis: {
                axisLabel: {
                    formatter: function (value) {
                        return value + unit
                    }
                }
            },
            tooltip: {
                show: true,
                trigger: "item",
                formatter: tooltipFunc
            },
            grid: {
                left: 40,
                top: 10,
                right: 20,
                bottom: 20
            },
            series: [
                {
                    name: name,
                    type: "bar",
                    data: stats.map(function (v) {
                        return v.count;
                    }),
                    itemStyle: {
                        color: "#9DD3E8"
                    },
                    barWidth: "20em"
                }
            ],
            animation: true
        }
        chart.setOption(option)
        chart.resize()
    }

    this.resizeChart = function (chartId) {
        let chartBox = document.getElementById(chartId)
        if (chartBox == null) {
            return
        }
        let chart = echarts.init(chartBox)
        chart.resize()
    }

    this.processMaxUnit = function (stats) {
        let max = stats.$map(function (k, v) {
            return v.count
        }).$max()
        let divider = 0
        let unit = ""
        if (max >= 1000 * 1000 * 1000) {
            unit = "B"
            divider = 1000 * 1000 * 1000
        } else if (max >= 1000 * 1000) {
            unit = "M"
            divider = 1000 * 1000
        } else if (max >= 1000) {
            unit = "K"
            divider = 1000
        }
        stats.forEach(function (v) {
            v.rawCount = v.count
            if (divider > 0) {
                v.count /= divider
            }
        })
        return unit
    }
})

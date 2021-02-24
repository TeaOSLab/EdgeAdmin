Tea.context(function () {
    this.$delay(function () {
        let that = this

        let systemUnit = this.processMaxUnit(this.systemStats)
        this.reloadChart("system-chart", "操作系统", this.systemStats, function (v) {
            return v.system.name
        }, function (args) {
            return that.systemStats[args.dataIndex].system.name + ": " + teaweb.formatNumber(that.systemStats[args.dataIndex].rawCount)
        }, systemUnit)
        window.addEventListener("resize", function () {
            that.resizeChart("system-chart")
        })

        let browserUnit = this.processMaxUnit(this.browserStats)
        this.reloadChart("browser-chart", "浏览器", this.browserStats, function (v) {
            return v.browser.name
        }, function (args) {
            return that.browserStats[args.dataIndex].browser.name + ": " + teaweb.formatNumber(that.browserStats[args.dataIndex].rawCount)
        }, browserUnit)

        window.addEventListener("resize", function () {
            that.resizeChart("system-chart")
            that.resizeChart("browser-chart")
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

Tea.context(function () {
    this.$delay(function () {
        let that = this

        this.totalDailyStats = this.logDailyStats.map(function (v, k) {
            return {
                day: v.day,
                count: that.logDailyStats[k].count + that.blockDailyStats[k].count + that.captchaDailyStats[k].count
            }
        })
        let dailyUnit = this.processMaxUnit(this.totalDailyStats)
        this.reloadLineChart("daily-chart", "规则分组", this.totalDailyStats, function (v) {
            return v.day.substring(4, 6) + "-" + v.day.substring(6)
        }, function (args) {

            return that.logDailyStats[args.dataIndex].day.substring(4, 6) + "-" + that.logDailyStats[args.dataIndex].day.substring(6) + ": 拦截: "
                + teaweb.formatNumber(that.blockDailyStats[args.dataIndex].count) + ", 验证码: " + teaweb.formatNumber(that.captchaDailyStats[args.dataIndex].count) + ", 记录: " + teaweb.formatNumber(that.logDailyStats[args.dataIndex].count)
        }, dailyUnit)

        let groupUnit = this.processMaxUnit(this.groupStats)
        let total = this.groupStats.$sum(function (k, v) {
            return v.rawCount
        })
        this.reloadBarChart("group-chart", "规则分组", this.groupStats, function (v) {
            return v.group.name
        }, function (args) {
            let percent = ""
            if (total > 0) {
                percent = ", 占比: " + (Math.round(that.groupStats[args.dataIndex].rawCount * 100 * 100 / total) / 100) + "%"
            }
            return that.groupStats[args.dataIndex].group.name + ": " + teaweb.formatNumber(that.groupStats[args.dataIndex].rawCount) + percent
        }, groupUnit)

        window.addEventListener("resize", function () {
            that.resizeChart("daily-chart")
            that.resizeChart("group-chart")
        })
    })

    this.reloadLineChart = function (chartId, name, stats, xFunc, tooltipFunc, unit) {
        let chartBox = document.getElementById(chartId)
        if (chartBox == null) {
            return
        }
        let that = this
        let chart = echarts.init(chartBox)
        let option = {
            xAxis: {
                data: stats.map(xFunc)
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
                    type: "line",
                    data: this.totalDailyStats.map(function (v, index) {
                        return that.totalDailyStats[index].count;
                    }),
                    areaStyle: {},
                    itemStyle: {
                        color: "#9DD3E8"
                    }
                },
                {
                    name: name,
                    type: "line",
                    data: this.logDailyStats.map(function (v) {
                        return v.count;
                    }),
                    itemStyle: {
                        color: "#879BD7"
                    }
                },
                {
                    name: name,
                    type: "line",
                    data: this.blockDailyStats.map(function (v) {
                        return v.count;
                    }),
                    itemStyle: {
                        color: "#F39494"
                    }
                },
                {
                    name: name,
                    type: "line",
                    data: this.captchaDailyStats.map(function (v) {
                        return v.count;
                    }),
                    itemStyle: {
                        color: "#FBD88A"
                    }
                }
            ],
            animation: true
        }
        chart.setOption(option)
        chart.resize()
    }

    this.reloadBarChart = function (chartId, name, stats, xFunc, tooltipFunc, unit) {
        let chartBox = document.getElementById(chartId)
        if (chartBox == null) {
            return
        }
        let chart = echarts.init(chartBox)
        let option = {
            xAxis: {
                data: stats.map(xFunc)
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
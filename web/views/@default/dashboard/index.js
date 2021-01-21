Tea.context(function () {
    this.trafficTab = "hourly"


    this.$delay(function () {
        this.reloadHourlyTrafficChart()
    })

    this.selectTrafficTab = function (tab) {
        this.trafficTab = tab
        if (tab == "hourly") {

        } else if (tab == "daily") {
            this.$delay(function () {
                this.reloadDailyTrafficChart()
            })
        }
    }

    this.reloadHourlyTrafficChart = function () {
        let chartBox = document.getElementById("hourly-traffic-chart-box")
        let chart = echarts.init(chartBox)
        let option = {
            xAxis: {
                data: this.hourlyTrafficStats.map(function (v) {
                    return v.hour;
                })
            },
            yAxis: {},
            tooltip: {
                show: true,
                trigger: "item",
                formatter: "{c} GB"
            },
            grid: {
                left: 40,
                top: 10,
                right: 20
            },
            series: [
                {
                    name: "流量",
                    type: "line",
                    data: this.hourlyTrafficStats.map(function (v) {
                        return v.count;
                    }),
                    itemStyle: {
                        color: "#9DD3E8"
                    },
                    lineStyle: {
                        color: "#9DD3E8"
                    },
                    areaStyle: {
                        color: "#9DD3E8"
                    }
                }
            ],
            animation: false
        }
        chart.setOption(option)
    }

    this.reloadDailyTrafficChart = function () {
        let chartBox = document.getElementById("daily-traffic-chart-box")
        let chart = echarts.init(chartBox)
        let option = {
            xAxis: {
                data: this.dailyTrafficStats.map(function (v) {
                    return v.day;
                })
            },
            yAxis: {},
            tooltip: {
                show: true,
                trigger: "item",
                formatter: "{c} GB"
            },
            grid: {
                left: 40,
                top: 10,
                right: 20
            },
            series: [
                {
                    name: "流量",
                    type: "line",
                    data: this.dailyTrafficStats.map(function (v) {
                        return v.count;
                    }),
                    itemStyle: {
                        color: "#9DD3E8"
                    },
                    lineStyle: {
                        color: "#9DD3E8"
                    },
                    areaStyle: {
                        color: "#9DD3E8"
                    }
                }
            ],
            animation: false
        }
        chart.setOption(option)
    }
})

// 给Table增加排序功能
function sortTable(callback) {
	// 引入js
	let jsFile = document.createElement("script")
	jsFile.setAttribute("src", "/js/sortable.min.js")
	jsFile.addEventListener("load", function () {
		// 初始化
		let box = document.querySelector("#sortable-table")
		if (box == null) {
			return
		}
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {
			},
			onUpdate: function (event) {
				let rows = box.querySelectorAll("tbody")
				let rowIds = []
				rows.forEach(function (row) {
					rowIds.push(parseInt(row.getAttribute("v-id")))
				})
				callback(rowIds)
			}
		})
	})
	document.head.appendChild(jsFile)
}

function sortLoad(callback) {
	let jsFile = document.createElement("script")
	jsFile.setAttribute("src", "/js/sortable.min.js")
	jsFile.addEventListener("load", function () {
		if (typeof (callback) == "function") {
			callback()
		}
	})
	document.head.appendChild(jsFile)
}

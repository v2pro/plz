package witch

import "net/http"

func Start(addr string) {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/more-events", moreEvents)
	http.ListenAndServe(addr, nil)
}

func homepage(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <!-- 引入样式 -->
  <link rel="stylesheet" href="http://cdn.jsdeliver.net/npm/element-ui/lib/theme-chalk/index.css">
  <style>
body {
	font-family: Helvetica
}

.icon {
	width: 1em;
	height: 1em;
	vertical-align: -0.15em;
	fill: currentColor;
	overflow: hidden;
}

.main {
	margin-left: 120px;
}
  </style>
</head>
<body>
	<div id="app">
    	<ide></ide>
	</div>
</body>
  <!-- 先引入 Vue -->
  <script src="http://cdn.jsdeliver.net/npm/vue/dist/vue.js"></script>
  <!-- 引入组件库 -->
  <script src="http://cdn.jsdeliver.net/npm/element-ui/lib/index.js"></script>
  <script src="http://cdn.jsdeliver.net/npm/axios/dist/axios.min.js"></script>
	` + compIDE() + compLogViewer() + `
  <script>
	var $vue = new Vue({
      el: '#app',
      data: function() {
        return { visible: false }
      }
    });
	axios.interceptors.response.use(function (response) {
        return response;
    }, function (error) {
        $vue.$notify.error({
            title: error.message,
            message: error.response
        });
        return error;
    });

  </script>
</html>
	`))
}

func compIDE() string {
	return `
<script type="text/x-template" id="ide-template">
    <div>
        <el-menu default-active="1" :collapse="isCollapse"
                 style="position:fixed; _position:absolute; top:8px; z-index:999">
            <el-menu-item index="1" @click="switchView('log')">
                <i class="el-icon-document"></i>
                <span slot="title">Log</span>
            </el-menu-item>
        </el-menu>
        <div class="main">
            <div id="log-view" v-if="currentView == 'log'">
                <log-viewer></log-viewer>
            </div>
        </div>
    </div>
</script>
<script>
    Vue.component('ide', {
        template: '#ide-template',
        data: function () {
            return {
				isCollapse: false,
				currentView: 'log',
            }
        },
	    methods: {
            switchView: function(view) {
                this.currentView = view;
            },
		}
    })
</script>

	`
}

func compLogViewer() string {
	return `
<script type="text/x-template" id="log-viewer-template">
    <div>
		<el-table
			:data="tableData"
			row-key="timestamp"
			stripe
			style="width: 100%">
		<el-table-column
			prop="timestamp"
			label="Time"
			width="180">
		</el-table-column>
		<el-table-column
			:formatter="formatSummary"
			label="Summary">
			</el-table-column>
		</el-table>
	</div>
</script>
<script>
    Vue.component('log-viewer', {
        template: '#log-viewer-template',
        data: function () {
            return {
				events: [],
            }
        },
		methods: {
			formatSummary: function(row, column, cellValue) {
				var desc = row.event;
				for (var key in row) {
					if (key == 'event' || key == 'lineNumber' || key == 'level' || key == 'timestamp') {
						continue;
					}
					desc = desc + ' ' + key + ':' + row[key];
				}
				return desc;
			}
		},
		computed: {
			tableData: function() {
				if (this.events.length < 1000) {
					return this.events;
				}
				return this.events.slice(0, 1000);
			}
		},
		created: function() {
			var me = this;
			(function(){
				axios.get('/more-events?ts=' + Date.now())
					.then(function (resp) {
						me.events = resp.data.reverse().concat(me.events);
					});
				setTimeout(arguments.callee, 1000);
			})();
		}
    })
</script>
	`
}

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
.el-tag + .el-tag {
margin-left: 10px;
}
.el-tag {
line-height: 30px;
}
.button-new-tag {
margin-left: 10px;
height: 32px;
line-height: 30px;
padding-top: 0;
padding-bottom: 0;
}
.input-new-tag {
width: 90px;
line-height: 30px;
margin-left: 10px;
vertical-align: bottom;
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
	<el-row>
		<el-col :span="20">
			<el-table
				:data="tableData"
				row-key="timestamp"
				stripe
				border
				style="width: 100%;">
			<el-table-column label="timestamp" :formatter="formatTimestamp" />
			<el-table-column v-for="item in userDefinedColumns" :prop="item" :label="item" resizable />
			<el-table-column label="event" prop="event" />
			<el-table-column label="details" min-width="400">
				  <template slot-scope="scope">
						<table>
						<tr v-for="(propValue, propKey) in scope.row" v-if="shouldShowProp(propKey)">
							<td>{{ propKey }}</td><td>{{ propValue }}</td>
						</tr>
						</table>
				  </template>
				</el-table-column>
			</el-table>
		</el-col>
		<el-col :span="4" style="pdding-right: 1em;">
			<ul style="position:fixed; _position:absolute; top:8px; z-index:999; list-style-type: none;">
				<li @click="scrollToTop" style="cursor: pointer;">
                <i class="el-icon-upload2"></i> Scroll to Top
				</li>
				<li><h3>Excluded Properties</h3>
				<div>
					<el-tag
					  :key="tag"
					  v-for="tag in excludedProperties"
					  closable
					  :disable-transitions="false"
					  @close="handleClose(tag)">
					  {{tag}}
					</el-tag>
					<el-input
					  class="input-new-tag"
					  v-if="inputVisible"
					  v-model="inputValue"
					  ref="saveTagInput"
					  size="small"
					  @keyup.enter.native="handleInputConfirm"
					  @blur="handleInputConfirm">
					</el-input>
					<el-button v-else class="button-new-tag" size="small" @click="showInput">+ New Tag</el-button>
				</div>
				</li>
			</ul>
		</el-col>
	</el-row>
</script>
<script>
    Vue.component('log-viewer', {
        template: '#log-viewer-template',
        data: function () {
            return {
				events: [],
				userDefinedColumns: [],
				excludedProperties: ['response']
            }
        },
		methods: {
			formatTimestamp: function(row, column, cellValue) {
				var d = new Date(row.timestamp / 1000000);
				return d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds() + '.' + d.getMilliseconds();
			},
			scrollToTop: function() {
				window.scrollTo(0, 0);
			},
			handleClose(tag) {
				this.excludedProperties.splice(this.excludedProperties.indexOf(tag), 1);
			},

			showInput() {
				this.inputVisible = true;
				this.$nextTick(_ => {
				  this.$refs.saveTagInput.$refs.input.focus();
				});
			},
			handleInputConfirm() {
				let inputValue = this.inputValue;
				if (inputValue) {
				  this.excludedProperties.push(inputValue);
				}
				this.inputVisible = false;
				this.inputValue = '';
			},
			shouldShowProp(propKey) {
				if (this.excludedProperties.indexOf(propKey) !== -1) {
					return false;
				}
				return propKey != 'lineNumber' && propKey != 'level' && propKey != 'event' && propKey != 'timestamp'
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

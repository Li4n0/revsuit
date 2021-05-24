(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-49fbda44"],{a299:function(t,e,a){},a2a5:function(t,e,a){"use strict";a("a299")},f832:function(t,e,a){"use strict";a.r(e);var n=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("div",[a("a-button",{attrs:{id:"add-rule",type:"primary"},on:{click:t.addRule}},[a("a-icon",{attrs:{type:"plus"}}),t._v(" New Rule ")],1),a("a-drawer",{attrs:{title:t.formAction+" FTP rule",width:490,visible:t.formVisible,"body-style":{paddingBottom:"80px"}},on:{close:t.closeDrawer}},[a("a-form-model",{ref:"form",attrs:{model:t.form,layout:"vertical"},on:{submit:t.handleSubmit}},[a("BasicRule",{attrs:{form:t.form,readOnly:t.formReadOnly,flagFormat:"user or password"}}),a("a-row",{attrs:{gutter:24}},[a("a-col",{attrs:{span:24}},[a("a-form-model-item",[a("span",{attrs:{slot:"label"},slot:"label"},[t._v(" Pasv Address "),a("a-tooltip",{attrs:{title:"1. Support template such as ${user}/${password}/${varname}.\n2.For setting the rebind mode, please use ',' to separate addresses."}},[a("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),a("a-input",{staticStyle:{width:"100%"},attrs:{placeholder:"Use `external_ip:pasv_port` by default",readOnly:t.formReadOnly},model:{value:t.form.pasv_address,callback:function(e){t.$set(t.form,"pasv_address",e)},expression:"form.pasv_address"}})],1)],1)],1),a("a-row",{attrs:{gutter:24}},[a("a-col",{attrs:{span:24}},[a("a-form-model-item",[a("span",{attrs:{slot:"label"},slot:"label"},[t._v(" Data "),a("a-tooltip",{attrs:{title:"1. The data to be returned when the client executes the download request.\n                      2. Please use base64 encoding.\n                      3. Support template such as ${user}/${password}/${varname}"}},[a("a-icon",{attrs:{type:"question-circle-o"}})],1)],1),a("a-textarea",{attrs:{placeholder:"The data to be returned when the client executes the download request. Please use base64 encoding.",readOnly:t.formReadOnly,"auto-size":{minRows:10,maxRows:30}},model:{value:t.form.data,callback:function(e){t.$set(t.form,"data",e)},expression:"form.data"}})],1)],1)],1)],1),a("div",{style:{position:"absolute",right:0,bottom:0,width:"100%",borderTop:"1px solid #e9e9e9",padding:"10px 16px",background:"#fff",textAlign:"right",zIndex:1}},[a("a-button",{style:{marginRight:"8px"},on:{click:t.handleCancel}},[t._v(" Cancel ")]),a("a-button",{attrs:{type:"primary",disabled:t.formReadOnly},on:{click:t.handleSubmit}},[t._v(" Submit ")])],1)],1),a("a-table",{staticStyle:{"overflow-x":"auto"},attrs:{columns:t.columns,"data-source":t.data,loading:t.loading,pagination:t.pagination},on:{change:t.handleTableChange},scopedSlots:t._u([{key:"filterDropdown",fn:function(e){var n=e.setSelectedKeys,o=e.selectedKeys,r=(e.clearFilters,e.column);return a("div",{staticStyle:{padding:"8px"}},[a("a-input",{staticStyle:{width:"188px","margin-bottom":"8px",display:"block"},attrs:{placeholder:"Search "+r.dataIndex,value:o[0]},on:{change:function(t){return n(t.target.value?[t.target.value]:[])},pressEnter:function(){t.filters[r.dataIndex]=o[0],t.fetch()}}}),a("a-button",{staticStyle:{width:"90px","margin-right":"8px"},attrs:{type:"primary",icon:"search",size:"small"},on:{click:function(){t.filters[r.dataIndex]=o[0],t.fetch()}}},[t._v(" Search ")])],1)}},{key:"filterIcon",fn:function(t){return a("a-icon",{style:{color:t?"#108ee9":void 0},attrs:{type:"search"}})}},{key:"rank",fn:function(e){return a("span",{},[a("a-tag",{attrs:{color:"#"+(2996213+80*e).toString(16)}},[t._v(" "+t._s(e)+" ")])],1)}},{key:"switchRender",fn:function(e,n,o,r){return a("span",{},[a("a-switch",{attrs:{checked:e},on:{click:function(e){return t.clickSwitch(n,r.dataIndex)}}})],1)}},{key:"valueRender",fn:function(e){return a("span",{},t._l(e.split(","),(function(e){return a("span",{key:e},[t._v(t._s(e)),a("br")])})),0)}},{key:"action",fn:function(e,n,o){return a("span",{},[a("a-button",{staticStyle:{color:"#67C23A","background-color":"transparent","border-color":"#67C23A","text-shadow":"none",margin:"0 10px 3px 0"},attrs:{size:"small",ghost:""},on:{click:function(e){return t.viewRule(n)}}},[t._v("View")]),a("a-button",{staticStyle:{color:"#909399","background-color":"transparent","border-color":"#909399","text-shadow":"none",margin:"0 10px 3px 0"},attrs:{size:"small",ghost:""},on:{click:function(e){return t.editRule(n,o)}}},[t._v("Edit")]),a("a-popconfirm",{attrs:{title:"Are you sure delete this task?","ok-text":"Yes","cancel-text":"No"},on:{confirm:function(e){return t.deleteRule(n,o)}}},[a("a-button",{attrs:{type:"danger",size:"small",ghost:""}},[t._v("Delete")])],1)],1)}}])})],1)},o=[],r=a("5530"),i=(a("a434"),a("34c6")),s=a("56d7"),c=a("2084"),l="View",d="Edit",u="Create",p=[{title:"ID",dataIndex:"id",key:"id",sorter:!0,sortDirections:["descend","ascend"]},{title:"NAME",dataIndex:"name",key:"name",scopedSlots:{filterDropdown:"filterDropdown",filterIcon:"filterIcon"}},{title:"FLAG FORMAT",dataIndex:"flag_format",key:"flag_format",ellipsis:!0},{title:"RANK",dataIndex:"rank",key:"rank",scopedSlots:{customRender:"rank"}},{title:"PASV ADDRESS",dataIndex:"pasv_address",key:"pasv_address",ellipsis:!0},{title:"PUSH TO CLIENT",dataIndex:"push_to_client",key:"push_to_client",scopedSlots:{customRender:"switchRender"}},{title:"NOTICE",dataIndex:"notice",key:"notice",scopedSlots:{customRender:"switchRender"}},{title:"Action",key:"action",scopedSlots:{customRender:"action"}}],f={name:"FtpRules",data:function(){return{store:s["store"],data:[],formVisible:!1,pagination:{current:1,showSizeChanger:!0,pageSize:s["store"].pageSize,onShowSizeChange:function(t,e){s["store"].pageSize=e}},filters:{},loading:!1,columns:p,form:{},formReadOnly:!1,formAction:""}},methods:{handleTableChange:function(t,e,a){var n=Object(r["a"])({},this.pagination);n.current=t.current,this.pagination=n,this.order="ascend"===a.order?"asc":"desc",this.fetch()},fetch:function(){var t=this,e=Object(r["a"])(Object(r["a"])({},this.filters),{},{page:this.pagination.current,pageSize:this.pagination.pageSize,order:this.order});this.loading=!0,Object(i["g"])(e).then((function(e){var a=e.data.result;t.data=a.data;var n=Object(r["a"])({},t.pagination);n.total=a.count,t.pagination=n,t.loading=!1})).catch((function(e){403!==e.response.status&&t.$message.error("Unknown error with status code: "+e.response.status)}))},clickSwitch:function(t,e){var a=this;t[e]=!t[e],Object(i["l"])(t).then().catch((function(t){a.$notification.error({message:"Edit failed",description:t.response.data.error,style:{width:"600px",marginLeft:"".concat(-265,"px")},duration:4})}))},addRule:function(){this.form={},this.showForm(u)},viewRule:function(t){this.form=t,this.showForm(l)},editRule:function(t){this.form=JSON.parse(JSON.stringify(t)),this.showForm(d)},deleteRule:function(t,e){var a=this;Object(i["b"])(t).then((function(){a.data.splice(e,1)})).catch((function(t){a.$notification.error({message:"Error",description:t.response.data.error,style:{width:"600px",marginLeft:"".concat(-265,"px")},duration:4})}))},showForm:function(t){this.formAction=t,this.formReadOnly=t===l,this.formVisible=!0},closeDrawer:function(){this.formVisible=!1},handleSubmit:function(){var t=this;this.$refs.form.validate((function(e){e&&Object(i["l"])(t.form).then((function(){t.closeDrawer(),t.fetch({page:t.pagination.current}),t.$notification.info({message:"Success",style:{width:"600px",marginLeft:"".concat(-265,"px")},duration:2.5})})).catch((function(e){t.$notification.error({message:t.formAction+" failed",description:e.response.data.error,style:{width:"600px",marginLeft:"".concat(-265,"px")},duration:4})}))}))},handleCancel:function(){this.form={},this.closeDrawer()}},mounted:function(){this.fetch({page:"1"})},components:{BasicRule:c["a"]}},h=f,m=(a("a2a5"),a("2877")),g=Object(m["a"])(h,n,o,!1,null,"6f29426c",null);e["default"]=g.exports}}]);
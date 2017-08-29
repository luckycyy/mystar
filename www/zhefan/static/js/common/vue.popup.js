Vue.component('my-popup', {
	/*
	页面引用
	 <my-popup v-show="showPopUp" v-on:close="doClosePopUp()" v-on:done="do()"  pop-up-title="标题"  multi-select="2"  pop-up-okbtn="Done"></my-popup>
	* */
	template:'<div class="popUp">\
					<div class="popUpPanel">\
					<div class="popUpBody">{{popUpContent}}</div>\
					<div class="popUpFooter">\
						<div class="popUpBtn" @click="close">{{popUpCancelbtn}}</div>\
						<div class="cutLine" v-show="multiSelect==2"></div>\
						<div class="popUpBtn popUpCancelBtn" v-show="multiSelect==2" @click="done">{{popUpOkbtn}}</div>\
					</div>\
				</div>\
			</div>',
	props:{
		popUpTitle:{
			type:String,
			default:'标题'
		},
		popUpContent:{
			type:String,
			default:'确定删除吗？'
		},
		popUpOkbtn:{
			type:String,
			default:'删除'
		},
		popUpCancelbtn:{
			type:String,
			default:'取消'
		},
		multiSelect:{
			type:String,
			default:"2"
		}
	},
	methods:{
		close:function(){
			this.$emit('close');
		},
		done:function(){
			this.$emit("done");
		}
	}	
})
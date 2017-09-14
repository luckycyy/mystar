(function (Vue) {
    const template = '<transition name="alert">\
                          <div v-if="showprompt" class="prompt">\
                              <span class="content" v-text="promptext"></span>\
                          </div>\
                       </transition>';
    var element = document.createElement('div');
    element.id = 'V-prompt'
    element.innerHTML = template
    document.body.appendChild(element)
    var $platformprompt = new Vue({
        el: '#V-prompt',
        data: {
            showprompt: false,
            promptext: ''
        },
        methods: {
            prompt:function (promptext) {
                var $this = this;
                if($this.showprompt === false){
                    $this.promptext = promptext;
                    $this.showprompt = true;
                    setTimeout(function () {
                        $this.showprompt = false;
                    },2000)
                }
            }
        }
    })
    window.$platformprompt = $platformprompt;
})(Vue)
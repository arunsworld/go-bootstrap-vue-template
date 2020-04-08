var app = new Vue({
    el: '#login',
    data: {
      username: '',
      password: '',
      processing: false,
      badCredentials: false,
    },
    methods: {
      login: function (event) {
        this.badCredentials = false;
        this.processing = true;
        self = this
        $.post("/login", {username: this.username, password: this.password})
          .done(function(data){
            window.location.href = data.Next;
          })
          .fail(function(response){
            self.processing = false;
            if (response.status==400) {
              alert("looks like we have a bug, server said: "+response.responseText)
            } else if (response.status == 403) {
              self.badCredentials = true;
            }
          });
      }
    },
    computed:{
      isNotValid: function(){
        if (this.processing) {
          return true
        }
        if (this.username=='') {
          return true
        }
        if (this.password=='') {
          return true
        }
      }
    },
    directives:{
      focus: {
        // When the bound element is inserted into the DOM...
        inserted: function (el) {
          // Focus the element
          el.focus()
        }
      }
    }
})


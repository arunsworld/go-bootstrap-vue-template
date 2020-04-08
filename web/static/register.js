var app = new Vue({
    el: '#registration',
    data: {
        username:'',
        password:'',
        name:'',
        processing: false,
        userAlreadyExists: false,
    },
    methods: {
        register: function (event) {
          this.badCredentials = false;
          this.processing = true;
          this.userAlreadyExists=false;
          self = this
          $.post("/register", {username: this.username, password: this.password, name: this.name})
            .done(function(){
              window.location.href = "/";
            })
            .fail(function(response){
              self.processing = false;
              if (response.responseText == "user already exists\n") {
                self.userAlreadyExists=true;
              } else {
                alert("looks like we have a bug, server said: "+response.responseText);
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
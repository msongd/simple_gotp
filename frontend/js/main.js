var KC_AUTHENTICATED = false;
var KC = "";

function sleep (time) {
    return new Promise((resolve) => setTimeout(resolve, time));
}

function initKeycloak(app) {
    var keycloak = new Keycloak(keycloakConfig);
    keycloak.init({
        enableLogging:true, 
        onLoad: 'login-required'
    }).then(function(authenticated) {
        //alert(authenticated ? 'authenticated' : 'not authenticated');
        //console.log(authenticated ? 'authenticated' : 'not authenticated');
        //app.authenticated = authenticated ;
        //app.tokenParsed = keycloak.tokenParsed ;
        KC_AUTHENTICATED = authenticated ;
        KC = keycloak ;
        createVueApp(keycloak);
        Vue.prototype.$keycloak = keycloak;
    }).catch(function() {
        console.log('keycloak failed to initialize');
    });
}

function createVueApp(kc) {
    var app = new Vue({
        el: '#app',
        data: {
          selectedUsername:"",
          authenticated: false,
          tokenParsed: '',
        },
        mounted: function() {
          //initKeycloak(this);
          this.authenticated = kc.authenticated;
          this.tokenParsed = kc.tokenParsed;
        },
        computed: {
          authenticatedUsername() {
            console.log(this.tokenParsed.preferred_username);
            return this.tokenParsed.preferred_username ;
          }
        },
        methods: {
          isAuthenticated() {
            return authenticated;
          },
          onSelectUser: function (user) {
            var self = this;
            self.selectedUsername = user;
            self.$refs.tokens.fetchTokens(user);
          },
          onDeleteToken: function(data) {
            var self = this;
            userToDelete = data.username;
            tokenToDelete = data.token;
            console.log("In MAIN():ondeletetoken:username:",userToDelete,":token:",tokenToDelete);
            self.$refs.users.removeTokenFromUser(userToDelete,tokenToDelete);
          },
          onChangeUser: function(data) {
            var self = this;
            console.log("In MAIN():onchangeuser:username:",data);
            self.$refs.users.refreshUser(data);
          }
        }
      });
}

function MakeHeader(obj) {
    var header = {
        "Content-Type": "application/json"
    }
    if (('$keycloak' in obj) && ('authenticated' in obj.$keycloak) && ('token' in obj.$keycloak)) {
        header["Authorization"] = obj.$keycloak.authenticated?"Bearer "+obj.$keycloak.token:"none" ;
    } else {
        header["Authorization"] = "none" ;
    }
    
    return header;    
}

//createVueApp();
window.onload = initKeycloak();
var KC_AUTHENTICATED = false;
var KC = "";

function sleep (time) {
    return new Promise((resolve) => setTimeout(resolve, time));
}

function initKeycloak(app) {
    if ((typeof keycloakConfig === 'undefined') || (!keycloakConfig.hasOwnProperty('url'))) {
        // no config -> skip keycloak
        createVueApp({});
        return 
    }
    var keycloak = new Keycloak(keycloakConfig);

    keycloak.onAuthError = function() { alert("xxx");};
    keycloak.init({
        enableLogging: true,
        onLoad: 'login-required'
    }).then(function(authenticated) {
      if (!authenticated) {
        window.location.reload();
      } else {
        //Vue.$log.info("Authenticated");
        KC_AUTHENTICATED = authenticated ;
        KC = keycloak ;
        createVueApp(keycloak);
        Vue.prototype.$keycloak = keycloak;
        //Token Refresh
        setInterval(() => {
          keycloak.updateToken(70).then((refreshed) => {
            if (refreshed) {
              console.log('Token refreshed' + refreshed);
            } else {
              console.log('Token not refreshed, valid for '
                + Math.round(keycloak.tokenParsed.exp + keycloak.timeSkew - new Date().getTime() / 1000) + ' seconds');
            }
          }).catch(() => {
            console.log('Failed to refresh token');
          });
        }, 6000)
      }
    }).catch(function(error) {
        console.log('keycloak failed to initialize:', error);
    });
}

function createVueApp(kc) {
    var app = new Vue({
        el: '#app',
        data: {
          selectedUsername:"",
          authenticated: false,
          tokenParsed: '',
          isAdmin: false
        },
        mounted: function() {
          //initKeycloak(this);
          if ('authenticated' in kc) {
            this.authenticated = kc.authenticated;
          }
          if ('tokenParsed' in kc) {
            this.tokenParsed = kc.tokenParsed;
          }
          if (ADMIN_ROLE != "" && kc.tokenParsed.realm_access.roles.includes(ADMIN_ROLE)) {
            this.isAdmin = true;
          } else {
            this.isAdmin = false;
          }
        },
        computed: {
          authenticatedUsername() {
            //console.log(this.tokenParsed.preferred_username);
            if (this.tokenParsed == '') {
                return "Anonymous" ;
            }
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
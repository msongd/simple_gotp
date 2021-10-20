Vue.component('list-users', {
  data: function() {
    return {
        users: [],
        selectedUser: "",
        showNew: false,
        newUsername: "",
        timer: ""
    }
  },
  mounted: function() {
    this.fetchUsers();
    this.timer = setInterval(this.fetchUserOTPs, 30000);
  },
  methods: {
    cancelAutoUpdate () {
      clearInterval(this.timer);
    },
    beforeDestroy () {
      this.cancelAutoUpdate();
    },
    fetchUsers: function(event) {
      var self = this ;
      //console.log('inside fetchUsers()');
      fetch('/auth/user', {
        method: 'GET',
        headers: {
          "Content-Type": "application/json",
          "Authorization": "none"
        }
      }).then(
        function(response) {
          //console.log('inside fetchUsers()-> response');
          if (response.status !== 200) {
              console.log('Looks like there was a problem. Status Code: ' + response.status);
              return;
            }
            // Examine the text in the response
            response.json().then(function(data) {
              //console.log(data);
              self.users = data;
            });
          }
        )
        .catch(function(err) {
          console.log('Fetch Error :-S', err);
        }
      );
    },
    fetchUserOTPs() {
      var self = this ;
      //console.log('inside fetchUsers()');
      fetch('/auth/otp', {
        method: 'GET',
        headers: {
          "Content-Type": "application/json",
          "Authorization": "none"
        }
      }).then(
        function(response) {
          //console.log('inside fetchUsers()-> response');
          if (response.status !== 200) {
              console.log('Looks like there was a problem. Status Code: ' + response.status);
              return;
            }
            // Examine the text in the response
            response.json().then(function(data) {
              //console.log(data);
              for (i=0;i<self.users.length;i++) {
                for (j=0;j<data.length;j++) {
                  if (self.users[i].username == data[j].username) {
                    self.users[i].current_code = data[j].current_code;
                  }
                }
              }
            });
          }
        )
        .catch(function(err) {
          console.log('Fetch Error :-S', err);
        }
      );
    },
    fetchTokens: function(event) {
      var self = this ;
      //console.log('inside fetchUsers()');
      fetch('/auth/user', {
        method: 'GET',
        headers: {
          "Content-Type": "application/json",
          "Authorization": "none"
        }
      }).then(
        function(response) {
          //console.log('inside fetchUsers()-> response');
          if (response.status !== 200) {
              console.log('Looks like there was a problem. Status Code: ' + response.status);
              return;
            }
            // Examine the text in the response
            response.json().then(function(data) {
              console.log(data);
              self.users = data;
            });
          }
        )
        .catch(function(err) {
          console.log('Fetch Error :-S', err);
        }
      );
    },
    isSelected: function(id) {
      var self = this ;
      console.log("isSelected:",id,"while:",self.selectedUser);
      return id == self.selectedUser ;
    },
    selectUser(id) {
      var self = this ;
      console.log("inside selectUser():", id);
      self.selectedUser = id;
      self.$emit('select-user', id);
    },
    addUser() {
      var self = this ;
      console.log("Inside addUser():newUsername:",self.newUsername);
      if (self.newUsername == undefined || self.newUsername=="") {
        console.log("Empty new user -> ignore");
      } else {
        u = { username: self.newUsername};
        console.log("Prepare to POST:", u)
        fetch('/auth/user', {
          method: 'POST',
          headers: {
            "Content-Type": "application/json",
            "Authorization": "none"
          },
          body: JSON.stringify(u)
        }).then(
          function(response) {
            if (response.status !== 201) {
              console.log('Looks like there was a problem. Status Code: ' + response.status);
              return;
            } else {
              self.users.push(u);
              self.newUsername = "";
            }
        });
      }
    },
    selectActiveToken(username,tokenId) {
      var self = this ;
      console.log("choosing active token", tokenId, "for user", self.selectedUser);
      if (username == undefined || username=="") {
        console.log("Empty username -> ignore");
        return;
      }
      let t={ active_token: tokenId}
      console.log("Prepare to POST:", t)
      fetch('/auth/user/'+username, {
        method: 'POST',
        headers: {
          "Content-Type": "application/json",
          "Authorization": "none"
        },
        body: JSON.stringify(t)
      }).then(
        function(response) {
          if (response.status !== 200) {
            console.log('Looks like there was a problem. Status Code: ' + response.status);
            return;
          }
          response.json().then(function(data) {
            console.log(data);
            //self.tokens = data;
            self.users[username].active_token = tokenId;
          });
      });
    },
    removeTokenFromUser(userToDelete,tokenToDelete) {
      var self = this ;
      self.users = self.users.filter(function(anUser){
        if (anUser.username != userToDelete) {
          return anUser;
        }
        if (anUser.active_token == tokenToDelete)
          anUser.active_token = ""
        anUser.tokens = anUser.tokens.filter(function(aToken){
          if (aToken.id != tokenToDelete)
            return aToken;
        });
        return anUser;        
      });
    },
    refreshUser(username) {
      var self = this ;
      var updatedUser ;
      console.log("in component tokens:refreshUser:username:", self.selectedUser);

      
      fetch('/auth/user/'+username.username, {
        method: 'GET',
        headers: {
          "Content-Type": "application/json",
          "Authorization": "none"
        }
      }).then(function(response) {
        if (response.status !== 200) {
          console.log('Looks like there was a problem. Status Code: ' + response.status);
          return nil;
        }
        response.json().then(function(data) {
          updatedUser = data;
          console.log('xxxxxx');
          //console.log(self.users);
          for (let i = 0; i < self.users.length; i++) {
            if (self.users[i].username == username.username) {
              console.log("should update token list for user:", username.username);
              console.log(data);
              self.$set(self.users,i,data);
            }
          }
        });
      });
      //self.users[0].username = "testtest";
      //self.users[0].tokens = [{issuer:'issue1', id:"1"},{issuer:'issue2', id:"2"}];
    }
},
  template: `
    <div class="mb-4 table-responsive">
      <div class="panel panel-default">
        <div class="panel-heading">Users</div>
        <div class="panel-body">
          <div class="col-md-4">
            <div class="input-group">
              <span class="input-group-btn">
                <button class="btn btn-primary" type="button" @click="addUser">Add user</button>
              </span>
              <input type="text" class="form-control" placeholder="New username..." v-model="newUsername">
            </div><!-- /input-group -->
          </div> <!-- class="col-md-4" -->
          <table class="table table-striped table-sm table-condensed ">
            <thead>
              <tr>
                <th>Username</th>
                <th>Active token</th>
                <th>Total</th>
                <th>Code</th>
                <th>Command</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="u in users" :class="{'info' : isSelected(u.username)}">
                <td>{{u.username}}</td>
                <td>
                  <select v-model="u.active_token" v-on:change="selectActiveToken(u.username, u.active_token)">
                    <option v-for="option in u.tokens" v-bind:value="option.id">
                      {{ option.issuer }}
                    </option>
                  </select>
                </td>
                <td>
                  {{u.tokens.length}}
                </td>
                <td>
                  {{u.current_code}}
                </td>
                <td>
                <!-- <button type="button" class="btn btn-info btn-sm" @click="">Tokens</button> -->
                <button type="button" class="btn btn-info btn-sm" @click="deleteUser(u.username)">Delete</button>
                <button type="button" class="btn btn-info btn-sm" @click="selectUser(u.username)">Detail</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    `
})


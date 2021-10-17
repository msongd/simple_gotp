Vue.component('users', {
  data: function() {
    return {
        users: [],
        selectedUser: "",
        showNew: false,
        newUsername: ""
    }
  },
  mounted: function() {
    this.fetchUsers();
  },
  methods: {
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
            }
        });
      }
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
                <th>Command</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="u in users" @click="selectUser(u.username)" :class="{'info' : isSelected(u.username)}">
                <td>{{u.username}}</td>
                <td>{{u.active_token}}</td>
                <td>
                <button type="button" class="btn btn-info btn-sm" @click="fetchToken(u.username)">Tokens</button>
                <!-- <button type="button" class="btn btn-info btn-sm" @click="">Tokens</button> -->
                <button type="button" class="btn btn-info btn-sm" @click="deleteUser(u.username)">Delete</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    `
})



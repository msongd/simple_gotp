Vue.component('users', {
  data: function() {
    return {
        users: [],
    }
  },
  prop: {
    selectedUser: String
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
      return id == self.selectedUser ;
    },
    selectUser(id) {
      var self = this ;
      console.log("inside selectUser():", id);
      self.selectedUser = id;
      //$emit('select-user', id);
    }
},
  template: `
          <div class="mb-4 table-responsive">
            <h4 class="mb-3">Users</h4>
            <table class="table table-striped table-sm table-condensed ">
              <thead>
                <tr>
                  <th>Username</th>
                  <th>Active token</th>
                  <th>Command</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="u in users" @click="selectUser(u.username)" v-on:click="$emit('select-user', u.username)" :class="{'info' : isSelected(u.username)}">
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
          `
})



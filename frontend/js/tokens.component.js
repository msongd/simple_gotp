Vue.component('tokens', {
  data: function() {
    return {
        tokens: [],
        selectedToken: ""
    }
  },
  mounted: function() {
    var self = this;
    this.fetchTokens(self.selectedUsername);
  },
  props: {selectedUsername: String},
  methods: {
    fetchTokens: function(username) {
      var self = this ;
      console.log('inside fetchToken(',username,')');
      if (username === undefined || username === "") return;
      var fetchUrl = '/auth/token/'+username ;
      fetch(fetchUrl, {
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
              self.tokens = data;
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
      return id == self.selectedToken ;
    },
    selectToken(id) {
      var self = this ;
      self.selectedToken = id;
    }
},
  template: `
          <div class="mb-4 table-responsive">
            <h4 class="mb-3">{{selectedUsername}}</h4>
            <table class="table table-striped table-sm table-condensed ">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Issuer</th>
                  <th>URL</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="t in tokens" @click="selectToken(u.id)" :class="{'info' : isSelected(u.id)}">
                  <td>{{u.id}}</td>
                  <td>{{u.issuer}}</td>
                  <td>{{u.url}}</td>
                  <td>
                  <button type="button" class="btn btn-info btn-sm" >Tokens</button>
                  <!-- <button type="button" class="btn btn-info btn-sm" >Tokens</button> -->
                  <button type="button" class="btn btn-info btn-sm" @click="deleteToken(username,u.id)">Delete</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          `
})



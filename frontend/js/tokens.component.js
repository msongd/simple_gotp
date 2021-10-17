Vue.component('tokens', {
  data: function() {
    return {
        tokens: [],
        selectedToken: "",
        newTokenIssuer: "",
        newTokenUrl: ""
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
    },
    addToken() {
      var self = this ;
      console.log("Inside addToken():username,newTokenIssuer:",self.selectedUsername,self.newTokenIssuer);
      if (self.newTokenIssuer == undefined || self.newTokenIssuer=="") {
        console.log("Empty new token issuer -> ignore");
        return;
      }
      if (self.selectedUsername == undefined || self.selectedUsername=="") {
        console.log("Empty selectedUsername -> ignore");
        return;
      }
      t = { issuer: self.newTokenIssuer};
      console.log("Prepare to POST:", t)
      fetch('/auth/token/'+self.selectedUsername, {
        method: 'POST',
        headers: {
          "Content-Type": "application/json",
          "Authorization": "none"
        },
        body: JSON.stringify(t)
      }).then(
        function(response) {
          if (response.status !== 201) {
            console.log('Looks like there was a problem. Status Code: ' + response.status);
            return;
          } else {
            self.fetchTokens(self.selectedUsername);
            self.newTokenIssuer = "" ;
          }
      });
    },
    importToken() {
      var self = this ;
      console.log("Inside importToken():username,newTokenUrl:",self.selectedUsername,self.newTokenUrl);
      if (self.newTokenUrl == undefined || self.newTokenUrl=="") {
        console.log("Empty new token url -> ignore");
        return;
      }
      if (self.selectedUsername == undefined || self.selectedUsername=="") {
        console.log("Empty selectedUsername -> ignore");
        return;
      }
      var newU = new URL(self.newTokenUrl);
      t = { url: self.newTokenUrl};
      console.log("Prepare to POST:", t)
      fetch('/auth/token/'+self.selectedUsername+'/import', {
        method: 'POST',
        headers: {
          "Content-Type": "application/json",
          "Authorization": "none"
        },
        body: JSON.stringify(t)
      }).then(
        function(response) {
          if (response.status !== 201) {
            console.log('Looks like there was a problem. Status Code: ' + response.status);
            return;
          } else {
            //self.tokens.push(importedToken);
            self.fetchTokens(self.selectedUsername);
            self.newTokenUrl = "";
          }
      });
    }
},
  template: `
        <div class="mb-4 table-responsive">
          <div class="panel panel-default">
            <div class="panel-heading">Token for {{selectedUsername}}</div>
            <div class="panel-body">
              <div class="col-md-4">
                <div class="input-group">
                  <span class="input-group-btn">
                    <button class="btn btn-primary" type="button" @click="addToken">Generate new token for issuer:</button>
                  </span>
                  <input type="text" class="form-control" placeholder="Issuer..." v-model="newTokenIssuer">
                </div><!-- /input-group -->
              </div> <!-- class="col-md-4" -->
              <div class="col-md-4">
                <div class="input-group">
                  <span class="input-group-btn">
                    <button class="btn btn-primary" type="button" @click="importToken">Import token from:</button>
                  </span>
                  <input type="text" class="form-control" placeholder="Token..." v-model="newTokenUrl">
                </div><!-- /input-group -->
              </div> <!-- class="col-md-4" -->
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
                  <tr v-for="t in tokens" @click="selectToken(t.id)" :class="{'info' : isSelected(t.id)}">
                    <td>{{t.id}}</td>
                    <td>{{t.issuer}}</td>
                    <td>{{t.url}}</td>
                    <td>
                    <button type="button" class="btn btn-info btn-sm" >Delete</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
          `
})



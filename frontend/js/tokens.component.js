Vue.component('tokens', {
  data: function() {
    return {
        tokens: [],
        selectedToken: "",
        newTokenIssuer: "",
        newTokenUrl: "",
        messageErr:"",
        visible: false
    }
  },
  mounted: function() {
    var self = this;
    console.log("when mounted:selectedUsername:", self.selectedUsername);
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
        headers: MakeHeader(self)
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
              self.visible = true;
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
        self.messageErr = "Empty issuer";
        return;
      }
      if (self.selectedUsername == undefined || self.selectedUsername=="") {
        console.log("Empty selectedUsername -> ignore");
        self.messageErr = "No user selected";
        return;
      }
      t = { issuer: self.newTokenIssuer};
      console.log("Prepare to POST:", t)
      fetch('/auth/token/'+self.selectedUsername, {
        method: 'POST',
        headers: MakeHeader(self),
        body: JSON.stringify(t)
      }).then(
        function(response) {
          if (response.status !== 201) {
            console.log('Looks like there was a problem. Status Code: ' + response.status);
            self.messageErr = "Request to server error";
            return;
          } else {
            self.fetchTokens(self.selectedUsername);
            self.newTokenIssuer = "" ;
            self.$emit('change-user',{username:self.selectedUsername});
          }
      });
    },
    importToken() {
      var self = this ;
      console.log("Inside importToken():username,newTokenUrl:",self.selectedUsername,self.newTokenUrl);
      if (self.newTokenUrl == undefined || self.newTokenUrl=="") {
        console.log("Empty new token url -> ignore");
        self.messageErr = "No token url input";
        return;
      }
      if (self.selectedUsername == undefined || self.selectedUsername=="") {
        console.log("Empty selectedUsername -> ignore");
        self.messageErr = "No user selected";
        return;
      }
      var newU = new URL(self.newTokenUrl);
      t = { url: self.newTokenUrl};
      console.log("Prepare to POST:", t)
      fetch('/auth/token/'+self.selectedUsername+'/import', {
        method: 'POST',
        headers: MakeHeader(self),
        body: JSON.stringify(t)
      }).then(
        function(response) {
          if (response.status !== 201) {
            console.log('Looks like there was a problem. Status Code: ' + response.status);
            self.messageErr = "Request to server error";
            return;
          } else {
            //self.tokens.push(importedToken);
            self.fetchTokens(self.selectedUsername);
            self.newTokenUrl = "";
            self.$emit('change-user',{username:self.selectedUsername});
          }
      });
    },
    getModalId(id) {
      return 'qrmodal-'+id;
    },
    showQRModal(id) {
      //$('#qrmodal-'+id).modal('show');
      var self = this ;
      console.log("Inside showQRModal():username,id:",self.selectedUsername,id);
      if (id == undefined || id == "") {
        console.log("Empty token url id -> ignore");
        self.messageErr = "Emtpy token id";
        return;
      }
      if (self.selectedUsername == undefined || self.selectedUsername=="") {
        console.log("Empty selectedUsername -> ignore");
        self.messageErr = "No user selected";
        return;
      }
      let t={}
      console.log("Prepare to POST:", t)
      fetch('/auth/qr/'+self.selectedUsername+'/'+id, {
        method: 'POST',
        headers: MakeHeader(self),
        body: JSON.stringify(t)
      }).then(
        function(response) {
          if (response.status !== 200) {
            console.log('Looks like there was a problem. Status Code: ' + response.status);
            self.messageErr = "Request to server error";
            return;
          }
          response.json().then(function(data) {
            //console.log(data);
            //self.tokens = data;
            let imgData = data.Img;
            console.log(imgData);
            $("#qrmodal-"+id+" img").attr("src",imgData);
            ///////
            $('#qrmodal-'+id).modal('show');
          });
      });
    },
    deleteToken(tokenId) {
      var self = this ;
      console.log("Request to delete token:username:",self.selectedUsername,":token:",tokenId);
      //self.$emit('delete-token',{username:self.selectedUsername,token:tokenId});
      
      if (tokenId == undefined || tokenId == "") {
        console.log("Empty token url id -> ignore");
        self.messageErr = "Emtpy token id";
        return;
      }
      if (self.selectedUsername == undefined || self.selectedUsername=="") {
        console.log("Empty selectedUsername -> ignore");
        self.messageErr = "No user selected";
        return;
      }
      let t={}
      console.log("Prepare to POST:", t);
      let isConfirm = confirm("Are you sure?");
      if (!isConfirm) {
        return;
      }
      fetch('/auth/token/'+self.selectedUsername+'/'+tokenId, {
        method: 'DELETE',
        headers: MakeHeader(self),
        body: JSON.stringify(t)
      }).then(
        function(response) {
          if (response.status !== 200) {
            console.log('Looks like there was a problem. Status Code: ' + response.status);
            self.messageErr = "Request to server error";
            return;
          }
          response.json().then(function(data) {
            //console.log(data);
            self.tokens = self.tokens.filter(function(item) {
              if (item.id != tokenId)
                return item;
            });
            self.$emit('delete-token',{username:self.selectedUsername,token:tokenId});
          });
      });
    }
},
  template: `
        <div class="mb-4 table-responsive" :class="[visible?'show':'hide']" >
          <div class="panel panel-default">
            <div class="panel-heading">Token for {{selectedUsername}}</div>
            <div class="panel-body">
              <div class="row">
                <div class="form-inline col-md-12">
                  <!-- todo: protect token url by PIN 
                  <div class="form-group">
                    <label>PIN</label>
                    <input type="text" class="form-control" placeholder="PIN">
                  </div>
                  -->
                  <div class="form-group">
                    <label>Issuer</label>
                    <input type="text" class="form-control" placeholder="Issuer..." v-model="newTokenIssuer">
                  </div>
                  <button class="btn btn-primary" type="button" @click="addToken">Generate</button>
                  <!-- todo: protect token url by PIN 
                  <div class="form-group">
                    <label>PIN</label>
                    <input type="text" class="form-control" placeholder="PIN">
                  </div>
                  -->
                  <div class="form-group">
                    <label>Token Url</label>
                    <input type="text" class="form-control" placeholder="Token..." v-model="newTokenUrl">
                  </div>
                  <button class="btn btn-primary" type="button" @click="importToken">Import</button>
                </div>                
              </div> <!-- div row -->
              <div class="alert alert-danger alert-dismissible" role="alert" v-if="messageErr!=''">
                <button type="button" class="close" data-dismiss="alert" aria-label="Close" @click="messageErr=''"><span aria-hidden="true">&times;</span></button>
                {{messageErr}}
              </div>
              <table class="table table-striped table-sm table-condensed " v-if="tokens && tokens.length>0">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Issuer</th>
                    <th>URL</th>
                    <th></th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="t in tokens" >
                    <td>{{t.id}}</td>
                    <td>{{t.issuer}}</td>
                    <td>
                      {{t.url}}
                      <div class="modal fade" :id="getModalId(t.id)" tabindex="-1" role="dialog">
                        <div class="modal-dialog modal-sm" role="document">
                          <div class="modal-content">
                            <div class="modal-header">
                              <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">×</span></button>
                              <h4 class="modal-title">{{selectedUsername}} - {{t.issuer}}</h4>
                            </div>
                            <div class="modal-body" style="text-align:center">
                            <img src=""/>
                            </div>
                          </div>
                        </div>
                      </div>
                    </td>
                    <td>
                    <button type="button" class="btn btn-info btn-sm" @click="showQRModal(t.id)">QR</button>
                    <button type="button" class="btn btn-info btn-sm btn-danger" @click="deleteToken(t.id)"><span class="glyphicon glyphicon-trash" aria-hidden="true"></span></button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
          `
})



Vue.component('confirm', {
    data: function() {
      return {
          dlg: 'confirm',
          message: '',
          title:'',
          isVisible: false,
          clicked:''
      }
    },
    props: ["name"],
    methods: {
        open() {
            this.isVisible = true;
            $("#"+ this.dlg).modal('show');
        },
        onOK() {
            this.isVisible = false;
            $("#"+ this.dlg).modal('hide');
            this.clicked="ok";
        },
        onCancel() {
            this.isVisible = false;
            $("#"+ this.dlg).modal('hide');
            this.clicked="cancel";
        },
        close() {
            this.isVisible = false;
            $("#"+ this.dlg).modal('hide');
            this.clicked="cancel";
        },
    },
    template: `
        <div class="modal fade" :id="dlg" tabindex="-1" role="dialog">
        <div class="modal-dialog modal-sm" role="document">
        <div class="modal-content">
            <div class="modal-header">
            <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">×</span></button>
            <h4 class="modal-title">{{title}}</h4>
            </div>
            <div class="modal-body">
            {{message}}
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" data-dismiss="modal" :click="onCancel">Cancel</button>
                <button type="button" class="btn btn-primary" :click="onOK">OK</button>
            </div>
        </div>
        </div>
        </div>
    `
})
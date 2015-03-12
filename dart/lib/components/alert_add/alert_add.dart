library alert_add_comp;

import 'package:blckur/models/alert_type.dart' as alrt_typ;
import 'package:blckur/models/account.dart' as acct;
import 'package:blckur/collections/alert_types.dart' as alrt_typs;
import 'package:blckur/loading.dart' as lodin;
import 'package:blckur/enter_aware.dart' as ent_awr;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alrt;

import 'package:angular/angular.dart' show Component, NgTwoWay;

@Component(
  selector: 'x-alert-add',
  templateUrl: 'packages/blckur/components/alert_add/alert_add.html',
  cssUrl: 'packages/blckur/components/alert_add/alert_add.css'
)
class AlertAddComp extends ent_awr.EnterAware with lodin.Loading {
  bool selected;
  dynamic typeValue;
  alrt_typ.AlertType typeModel;
  alrt_typs.AlertTypes alertTypes;

  @NgTwoWay('active')
  bool active;

  @NgTwoWay('model')
  acct.Account model;

  AlertAddComp(this.alertTypes);

  void clear() {
    this.active = false;
    this.selected = false;
    this.typeValue = null;
  }

  void onAdd() {
    if (!this.setLoading()) {
      return;
    }

    this.active = true;
    this.alertTypes.fetch().catchError((err) {
      logger.severe('Failed to load alert types', err);
      new alrt.Alert('Failed to load alert types');
    }).whenComplete(() {
      this.clearLoading();
    });
  }

  void onCancel() {
    this.clear();
    this.clearLoading();
  }

  void onClick(alrt_typ.AlertType model) {
    this.typeModel = model;

    if (model.valueType == "" || model.valueType == null) {
      this.onSave();
    }

    this.selected = true;
  }

  void onEnter() {
    if (this.selected) {
      this.onSave();
    }
  }

  void onSave() {
    if (!this.setLoading()) {
      return;
    }

    if (this.model.alerts == null) {
      this.model.alerts = [];
    }

    this.model.alerts.add({
      'type': typeModel.type,
      'value': this.typeValue,
    });

    this.model.save(['alerts']).catchError((err) {
      logger.severe('Failed to add alert', err);
      new alrt.Alert('Failed to add alert');
    }).whenComplete(() {
      this.clear();
      this.clearLoading();
    });
  }
}

library alert_add_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/loading.dart' as loading;
import 'package:blckur/enter_aware.dart' as enter_aware;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgTwoWay;

@Component(
  selector: 'x-account-alerts',
  templateUrl: 'packages/blckur/components/account_alerts/account_alerts.html',
  cssUrl: 'packages/blckur/components/account_alerts/account_alerts.css'
)
class AccountAlertsComp extends enter_aware.EnterAware with loading.Loading {
  bool selected;
  bool active;
  dynamic typeValue;
  models.AlertType typeModel;
  collections.AlertTypes alertTypes;

  @NgTwoWay('model')
  models.Account model;

  AccountAlertsComp(this.alertTypes);

  String get activeCls {
    if (this.active == true) {
      return 'adding-alert';
    }
    return '';
  }

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
      new alert.Alert('Failed to load alert types');
    }).whenComplete(() {
      this.clearLoading();
    });
  }

  void onCancel() {
    this.clear();
    this.clearLoading();
  }

  void onClick(models.AlertType model) {
    this.typeModel = model;

    if (model.valueType == "" || model.valueType == null) {
      this.onSave();
      return;
    }

    this.selected = true;
  }

  void onEnter() {
    if (this.selected) {
      this.onSave();
    }
  }

  void onDel(Map<String, String> alert) {
    if (!this.setLoading()) {
      return;
    }

    this.model.alerts.remove(alert);
    this.model.save(['alerts']).catchError((err) {
      logger.severe('Failed to remove alert', err);
      new alert.Alert('Failed to remove alert');
    }).whenComplete(() {
      this.clearLoading();
    });
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
      new alert.Alert('Failed to add alert');
    }).whenComplete(() {
      this.clear();
      this.clearLoading();
    });
  }
}

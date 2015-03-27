library alert_add_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/loading.dart' as loading;
import 'package:blckur/enter_aware.dart' as enter_aware;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgTwoWay;
import 'dart:async' as async;

const NONE = 0;
const ALERTS = 1;
const ADD_ALERT = 2;

@Component(
  selector: 'x-account-alerts',
  templateUrl: 'packages/blckur/components/account_alerts/account_alerts.html',
  cssUrl: 'packages/blckur/components/account_alerts/account_alerts.css'
)
class AccountAlertsComp extends enter_aware.EnterAware with loading.Loading {
  int _mode;
  async.Future _modeWait;
  bool showAlerts;
  bool showAddAlert;
  bool selected;
  dynamic typeValue;
  models.AlertType typeModel;
  collections.AlertTypes alertTypes;

  @NgTwoWay('model')
  models.Account model;

  AccountAlertsComp(this.alertTypes) {
    this.mode = ALERTS;
  }

  void set mode (int mode) {
    if (this._modeWait != null) {
      this._modeWait.then((_) {
        this._modeWait = null;
        this.mode = mode;
      });
      return;
    }

    if (this._mode == mode) {
      return;
    }
    this._mode = mode;

    switch (mode) {
      case ALERTS:
        this.showAddAlert = false;
        break;
      case ADD_ALERT:
        this.showAlerts = false;
        break;
      default:
        this.showAlerts = false;
        this.showAddAlert = false;

        this._modeWait = new async.Future.delayed(
          const Duration(milliseconds: 400), () {
            this._modeWait = null;
          });

        return;
    }

    this._modeWait = new async.Future.delayed(
      const Duration(milliseconds: 800), () {});

    new async.Timer(const Duration(milliseconds: 400), () {
      switch (mode) {
        case ALERTS:
          this.showAlerts = true;
          break;
        case ADD_ALERT:
          this.showAddAlert = true;
          break;
      }
      this.mode = mode;
      this._modeWait = null;
    });
  }

  int get mode {
    return this._mode;
  }

  void clear() {
    this.selected = false;
    this.typeValue = null;
  }

  void onAdd() {
    if (!this.setLoading()) {
      return;
    }
    this.mode = NONE;

    this.alertTypes.acctType = this.model.type;
    this.alertTypes.fetch().catchError((err) {
      logger.severe('Failed to load alert types', err);
      new alert.Alert('Failed to load alert types');
    }).whenComplete(() {
      this.mode = ADD_ALERT;
      this.clearLoading();
    });
  }

  void onCancel() {
    this.mode = ALERTS;
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

  void onDel(Map<String, String> alrt) {
    if (!this.setLoading()) {
      return;
    }

    this.model.alerts.remove(alrt);
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
    this.mode = NONE;

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
      this.mode = ALERTS;
      this.clear();
      this.clearLoading();
    });
  }
}

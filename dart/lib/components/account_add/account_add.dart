library account_add_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/loading.dart' as loading;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgTwoWay;
import 'dart:html' as dom;

@Component(
  selector: 'x-account-add',
  templateUrl: 'packages/blckur/components/account_add/account_add.html',
  cssUrl: 'packages/blckur/components/account_add/account_add.css'
)
class AccountAddComp extends loading.Loading {
  models.AccountAdd model;
  collections.AccountTypes accountTypes;

  @NgTwoWay('add-account')
  bool active;

  AccountAddComp(this.model, this.accountTypes);

  void onAdd() {
    if (!this.setLoading()) {
      return;
    }

    this.active = true;
    this.accountTypes.fetch().catchError((err) {
      logger.severe('Failed to load account types', err);
      new alert.Alert('Failed to load account types');
    }).whenComplete(() {
      this.clearLoading();
    });
  }

  void onCancel() {
    this.active = false;
    this.clearLoading();
  }

  void onClick(String type) {
    if (!this.setLoading()) {
      return;
    }

    this.model.type = type;
    this.model.create().then((_) {
      if (this.active != true) {
        return;
      }
      dom.window.location.replace(this.model.redirect);
    }).catchError((err) {
      if (this.active != true) {
        return;
      }
      logger.severe('Failed to add account', err);
      new alert.Alert('Unable to add account, try again later.', () {
        this.onClick(type);
      });
    }).whenComplete(() {
      this.clearLoading();
    });
  }
}

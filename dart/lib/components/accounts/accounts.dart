library accounts_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/injectables/injectables.dart' as injectables;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'dart:html' as dom;

const NONE = 0;
const ACCOUNTS = 1;
const ADD_ACCOUNT = 2;

@Component(
  selector: 'x-accounts',
  templateUrl: 'packages/blckur/components/accounts/accounts.html',
  cssUrl: 'packages/blckur/components/accounts/accounts.css'
)
class AccountsComp {
  collections.Accounts accounts;
  models.AccountAdd model;
  collections.AccountTypes accountTypes;
  injectables.Loading loading;
  injectables.ModeSwitch modeswitch;

  AccountsComp(this.accounts, this.accountTypes, this.loading,
      this.modeswitch) {
    this.update(true);
  }

  void update([bool initMode]) {
    if (!this.loading.set()) {
      return;
    }

    this.accounts.fetch().catchError((err) {
      logger.severe('Failed to load accounts', err);
      new alert.Alert('Error loading accounts', () {
        this.update();
      });
    }).whenComplete(() {
      if (initMode == true) {
        this.modeswitch.mode = ACCOUNTS;
      }
      this.loading.clear();
    });
  }

  void onAdd() {
    if (!this.loading.set()) {
      return;
    }
    this.modeswitch.mode = NONE;

    this.accountTypes.fetch().catchError((err) {
      logger.severe('Failed to load account types', err);
      new alert.Alert('Failed to load account types');
    }).whenComplete(() {
      this.loading.clear();
      this.modeswitch.mode = ADD_ACCOUNT;
    });
  }

  void onCancel() {
    this.modeswitch.mode = ACCOUNTS;
    this.loading.clear();
  }

  void onClick(String type) {
    if (!this.loading.set()) {
      return;
    }

    this.model.type = type;
    this.model.create().then((_) {
      if (this.modeswitch.mode != ADD_ACCOUNT) {
        return;
      }
      dom.window.location.replace(this.model.redirect);
    }).catchError((err) {
      if (this.modeswitch.mode != ADD_ACCOUNT) {
        return;
      }
      logger.severe('Failed to add account', err);
      new alert.Alert('Unable to add account, try again later.', () {
        this.onClick(type);
      });
    }).whenComplete(() {
      this.loading.clear();
    });
  }
}

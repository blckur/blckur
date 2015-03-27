library accounts_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/loading.dart' as loading;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'dart:html' as dom;
import 'dart:async' as async;

const NONE = 0;
const ACCOUNTS = 1;
const ADD_ACCOUNT = 2;

@Component(
  selector: 'x-accounts',
  templateUrl: 'packages/blckur/components/accounts/accounts.html',
  cssUrl: 'packages/blckur/components/accounts/accounts.css'
)
class AccountsComp extends loading.Loading {
  int _mode;
  async.Future _modeWait;
  bool showAccounts;
  bool showAddAcount;
  collections.Accounts accounts;
  models.AccountAdd model;
  collections.AccountTypes accountTypes;

  AccountsComp(this.accounts, this.accountTypes) {
    this.update(true);
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
      case ACCOUNTS:
        this.showAddAcount = false;
        break;
      case ADD_ACCOUNT:
        this.showAccounts = false;
        break;
      default:
        this.showAccounts = false;
        this.showAddAcount = false;

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
        case ACCOUNTS:
          this.showAccounts = true;
          break;
        case ADD_ACCOUNT:
          this.showAddAcount = true;
          break;
      }
      this.mode = mode;
      this._modeWait = null;
    });
  }
  int get mode {
    return this._mode;
  }

  void update([bool initMode]) {
    if (!this.setLoading()) {
      return;
    }

    this.accounts.fetch().catchError((err) {
      logger.severe('Failed to load accounts', err);
      new alert.Alert('Error loading accounts', () {
        this.update();
      });
    }).whenComplete(() {
      if (initMode == true) {
        this.mode = ACCOUNTS;
      }
      this.clearLoading();
    });
  }

  void onAdd() {
    if (!this.setLoading()) {
      return;
    }
    this.mode = NONE;

    this.accountTypes.fetch().catchError((err) {
      logger.severe('Failed to load account types', err);
      new alert.Alert('Failed to load account types');
    }).whenComplete(() {
      this.clearLoading();
      this.mode = ADD_ACCOUNT;
    });
  }

  void onCancel() {
    this.mode = ACCOUNTS;
    this.clearLoading();
  }

  void onClick(String type) {
    if (!this.setLoading()) {
      return;
    }

    this.model.type = type;
    this.model.create().then((_) {
      if (this.mode != ADD_ACCOUNT) {
        return;
      }
      dom.window.location.replace(this.model.redirect);
    }).catchError((err) {
      if (this.mode != ADD_ACCOUNT) {
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

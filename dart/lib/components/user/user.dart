library user_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/alert.dart' as alert;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Component;
import 'package:angular/angular.dart'as ng;
import 'package:crypto/crypto.dart' as crypto;
import 'dart:convert' as convert;
import 'dart:html' as dom;

@Component(
  selector: 'x-user',
  templateUrl: 'packages/blckur/components/user/user.html',
  cssUrl: 'packages/blckur/components/user/user.css'
)
class UserComp {
  bool settings;
  models.User model;
  models.User settingsModel;
  models.Auth authModel;
  ng.Router router;
  var _curEmail;
  var _avatar;

  UserComp(this.router) {
    this.model = new models.User();
    this.authModel = new models.Auth();
    this.update();
  }

  String get avatar {
    if (this.model != this._curEmail && this.model.email != null) {
      var enc = new crypto.MD5();
      enc.add(convert.UTF8.encoder.convert(this.model.email));
      var hash = crypto.CryptoUtils.bytesToHex(enc.close());

      this._avatar = '//www.gravatar.com/avatar/${hash}?s=170&d=mm';
    }
    return this._avatar;
  }

  void update() {
    this.model.fetch().catchError((err) {
      logger.severe('Failed to load user', err);
      new alert.Alert('Failed to load user');
    });
  }

  void onRefreshApi() {
    models.User model = this.model.clone();

    model.refreshApi().then((_) {
      this.model.apikey = model.apikey;
      this.settingsModel.apikey = model.apikey;
      print('refresh');
    }).catchError((err) {
      logger.severe('Failed to refresh apikey', err);
      new alert.Alert('Failed to refresh apikey');
    });
  }

  void onSettings() {
    this.settingsModel = this.model.clone();
    this.settings = this.settings != true;
  }

  void onSave() {
    this.settingsModel.save().then((_) {
      this.settingsModel.password = null;
      this.settings = false;
      this.model = this.settingsModel;
    }).catchError((err) {
      logger.severe('Failed to save user', err);
      new alert.Alert('Failed to save user');
    });
  }

  void onLogout() {
    this.authModel.logout().then((_) {
      dom.window.location.replace('#/login');
    }).catchError((err) {
      logger.severe('Failed to logout', err);
      new alert.Alert('Failed to logout');
    });
  }
}

library user_comp;

import 'package:blckur/models/auth.dart' as ath;
import 'package:blckur/models/user.dart' as usr;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Component;
import 'package:angular/angular.dart'as ng;

@Component(
  selector: 'x-user',
  templateUrl: 'packages/blckur/components/user/user.html',
  cssUrl: 'packages/blckur/components/user/user.css'
)
class UserComp {
  bool settings;
  usr.User model;
  usr.User settingsModel;
  ath.Auth authModel;
  ng.Router router;

  UserComp(this.model, this.authModel, this.router) {
    this.update();
  }

  void update() {
    this.model.fetch().catchError((err) {
      logger.severe('Failed to load user', err);
    });
  }

  void onSettings() {
    this.settingsModel = this.model.clone();
    settings = this.settings != true;
  }

  void onSave() {
    this.settingsModel.save().catchError((err) {
      logger.severe('Failed to save user', err);
    });

    settings = false;
    this.model = this.settingsModel;
  }

  void onLogout() {
    this.authModel.logout().then((_) {
      this.router.gotoUrl("/login");
    }).catchError((err) {
      logger.severe('Failed to logout', err);
    });
  }
}

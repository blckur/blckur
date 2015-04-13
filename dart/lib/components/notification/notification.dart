library notification_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/utils/utils.dart' as utils;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;

@Component(
  selector: 'x-notification',
  templateUrl: 'packages/blckur/components/notification/notification.html',
  cssUrl: 'packages/blckur/components/notification/notification.css'
)
class NotificationComp {
  utils.Loading loading;

  NotificationComp() {
    this.loading = new utils.Loading();
  }

  @NgOneWay('model')
  models.Notification model;

  void toggleRead() {
    if (!this.loading.set()) {
      return;
    }

    this.model.markRead().catchError((err) {
      logger.severe('Failed to mark notification as read', err);
      new alert.Alert('Failed to mark notification as read');
    }).whenComplete(() {
      this.loading.clear();
    });
  }
}

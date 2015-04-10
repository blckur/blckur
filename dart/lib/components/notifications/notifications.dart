library notifications_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/utils/utils.dart' as utils;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;

@Component(
  selector: 'x-notifications',
  templateUrl: 'packages/blckur/components/notifications/notifications.html',
  cssUrl: 'packages/blckur/components/notifications/notifications.css'
)
class NotificationsComp {
  models.Notification model;
  collections.Notifications notifications;
  utils.Loading loading;

  NotificationsComp() {
    this.model = new models.Notification();
    this.notifications = new collections.Notifications();
    this.loading = new utils.Loading();
    this.update();
  }

  void update() {
    if (!this.loading.set()) {
      return;
    }

    this.notifications.fetch().catchError((err) {
      logger.severe('Failed to load notifications', err);
      new alert.Alert('Error loading notifications', () {
        this.update();
      });
    }).whenComplete(() {
      this.loading.clear();
    });
  }
}

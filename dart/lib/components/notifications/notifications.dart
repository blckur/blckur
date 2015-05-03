library notifications_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/collections/collections.dart' as collections;
import 'package:blckur/utils/utils.dart' as utils;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;
import 'package:blckur/notification.dart' as notification;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;

@Component(
  selector: 'x-notifications',
  templateUrl: 'packages/blckur/components/notifications/notifications.html',
  cssUrl: 'packages/blckur/components/notifications/notifications.css'
)
class NotificationsComp implements ng.ScopeAware, ng.ShadowRootAware {
  models.Notification model;
  collections.Notifications notifications;
  utils.Loading loading;
  dom.Element list;

  NotificationsComp() {
    this.model = new models.Notification();
    this.notifications = new collections.Notifications();
    this.loading = new utils.Loading();
    this.update();
  }

  void set scope(ng.Scope scope) {
    scope.on('update_all').listen((evt) {
      this.update();
    });

    scope.on('notf_new').listen((evt) {
      var model = this.notifications.prepend(evt.data.data);

      if (model != null) {
        notification.notify(model.subject, model.type, model.body);
      }

      this.list.click();
    });

    scope.on('notf_update').listen((evt) {
      this.notifications.update(evt.data.data);
      this.list.click();
    });

    scope.on('notf_rem').listen((evt) {
      this.notifications.remove(evt.data.data);
      this.list.click();
    });
  }

  void onShadowRoot(dom.ShadowRoot root) {
    this.list = root.querySelector('.list');
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
      this.notifications = this.notifications;
      this.loading.clear();
    });
  }
}

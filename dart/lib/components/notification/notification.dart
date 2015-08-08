library notification_comp;

import 'package:blckur/models/models.dart' as models;
import 'package:blckur/utils/utils.dart' as utils;
import 'package:blckur/logger.dart' as logger;
import 'package:blckur/alert.dart' as alert;

import 'package:angular/angular.dart' show Component, NgOneWay;
import 'dart:math' as math;

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

  String get date {
    var sIndex = this.model.timestamp.lastIndexOf('.');
    var eIndex = this.model.timestamp.lastIndexOf('-');
    var timestamp = this.model.timestamp.substring(0, sIndex) +
      this.model.timestamp.substring(eIndex);
    var notf = DateTime.parse(timestamp);
    var now = new DateTime.now();
    var diff = now.difference(notf);

    if (diff < const Duration(hours: 1)) {
      return '${math.max(1, diff.inMinutes)} mins ago';
    } else if (diff < const Duration(hours: 24)) {
      return '${diff.inHours} hours ago';
    } else if (diff < const Duration(days: 2)) {
      return 'Yesterday';
    } else if (diff < const Duration(days: 7)) {
      switch (notf.weekday) {
        case DateTime.MONDAY:
          return 'Monday';
          break;
        case DateTime.TUESDAY:
          return 'Tuesday';
          break;
        case DateTime.WEDNESDAY:
          return 'Wednesday';
          break;
        case DateTime.THURSDAY:
          return 'Thursday';
          break;
        case DateTime.FRIDAY:
          return 'Friday';
          break;
        case DateTime.SATURDAY:
          return 'Saturday';
          break;
        case DateTime.SUNDAY:
          return 'Sunday';
          break;
      }
    } else {
      var month;

      switch (notf.month) {
        case DateTime.JANUARY:
          month = 'Jan';
          break;
        case DateTime.FEBRUARY:
          month = 'Feb';
          break;
        case DateTime.MARCH:
          month = 'Mar';
          break;
        case DateTime.APRIL:
          month = 'Apr';
          break;
        case DateTime.MAY:
          month = 'May';
          break;
        case DateTime.JUNE:
          month = 'Jun';
          break;
        case DateTime.JULY:
          month = 'Jul';
          break;
        case DateTime.AUGUST:
          month = 'Aug';
          break;
        case DateTime.SEPTEMBER:
          month = 'Sep';
          break;
        case DateTime.OCTOBER:
          month = 'Oct';
          break;
        case DateTime.NOVEMBER:
          month = 'Nov';
          break;
        case DateTime.DECEMBER:
          month = 'Dec';
          break;
      }
      return '$month ${notf.day}';
    }
  }

  void toggleRead() {
    if (!this.loading.set()) {
      return;
    }

    this.model.read = this.model.read != true;

    this.model.save(['read']).catchError((err) {
      logger.severe('Failed to mark notification as read', err);
      new alert.Alert('Failed to mark notification as read');
      this.model.read = this.model.read != true;
    }).whenComplete(() {
      this.loading.clear();
    });
  }

  void onDelete() {
    if (!this.loading.set()) {
      return;
    }

    this.model.destroy().catchError((err) {
      logger.severe('Failed to delete notification', err);
      new alert.Alert('Failed to delete notification');
      this.loading.clear();
    });
  }
}

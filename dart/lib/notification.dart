library notification;

import 'package:blckur/constants.dart';

import 'dart:html' as dom;

void notify(String title, String icon, String body) {
  dom.Notification.requestPermission().then((String permission) {
    if (permission == 'denied') {
      return;
    }

    var iconData = icons[icon];
    if (iconData != null) {
      icon = '${dom.window.location.protocol}//${dom.window.location.host}' +
        iconData['url'];
      icon = iconData['url'];
    }

    print(icon);

    try {
      new dom.Notification(title, body:body, icon:icon);
      return;
    } catch(_) {}

    try {
      new dom.Notification(title, body:body, iconUrl:icon);
      return;
    } catch(_) {}
  });
}

void init() {
  dom.Notification.requestPermission();
}

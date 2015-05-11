library events_col;

import 'package:blckur/model.dart' as model;
import 'package:blckur/collection.dart' as collection;
import 'package:blckur/models/models.dart' as models;

import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;
import 'dart:html' as dom;
import 'dart:convert' as convert;

class Events extends collection.Collection {
  String cursor;
  ng.RootScope rootScope;
  dom.WebSocket socket;
  ng.Router router;
  bool state;

  Events(this.rootScope, this.router) : super();

  String url = '/events';

  model.Model newModel() {
    return new models.Event();
  }

  collection.Collection newCollection() {
    return new Events(this.rootScope, this.router);
  }

  void onEvent(dom.MessageEvent evt) {
    var event = new models.Event();
    event.init();

    event.import(convert.JSON.decode(evt.data));

    this.rootScope.broadcast(event.type, event);
  }

  void onClose(_) {
    this.socket = null;

    if (this.state == true) {
      this.openDelay();
    }
  }

  void openSocket() {
    if (this.router.activePath[0].name != 'feed') {
      this.openDelay();
      return;
    }

    var proto;
    if (dom.window.location.protocol == 'http:') {
      proto = 'ws';
    } else {
      proto = 'wss';
    }

    this.stop();

    this.socket = new dom.WebSocket(
      '${proto}://${dom.window.location.host}${this.url}');
    this.socket.onMessage.listen(this.onEvent);
    this.socket.onClose.listen(this.onClose);
  }

  void openDelay() {
    new async.Timer(const Duration(seconds: 1), () {
      this.openSocket();
    });
  }

  void start() {
    this.state = true;
    this.openSocket();
  }

  void stop() {
    this.state = false;

    if (this.socket != null) {
      this.socket.close();
    }
  }
}

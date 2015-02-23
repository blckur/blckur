library events_col;

import 'package:blckur/collection.dart' as collec;
import 'package:blckur/models/event.dart' as evnt;
import 'package:blckur/logger.dart' as logger;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;
import 'dart:async' as async;
import 'dart:html' as dom;
import 'dart:convert' as convert;

@Injectable()
class Events extends collec.Collection {
  Type model = evnt.Event;
  String cursor;
  ng.RootScope rootScope;
  dom.WebSocket socket;

  Events(ng.Http http, this.rootScope) : super(http);

  String url = 'ws://localhost:3000/event';

  void onEvent(dom.MessageEvent evt) {
    var event = new evnt.Event(this.http);
    event.init();

    event.import(convert.JSON.decode(evt.data));

    this.rootScope.broadcast(event.type, event);
  }

  void onClose(_) {
    new async.Timer(const Duration(seconds: 1), () {
      this.openSocket();
    });
  }

  void openSocket() {
    this.socket = new dom.WebSocket(this.url);
    this.socket.onMessage.listen(this.onEvent);
    this.socket.onClose.listen(this.onClose);
  }

  void start() {
    this.openSocket();
  }
}

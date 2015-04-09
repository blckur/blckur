library notification_mod;

import 'package:blckur/model.dart' as model;

class Notification extends model.Model {
  String id;
  String timestamp;
  String type;
  String origin;
  String link;
  String subject;
  String body;

  model.Model newModel() {
    return new Notification();
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'timestamp': () => this.timestamp,
      'type': () => this.type,
      'origin': () => this.origin,
      'link': () => this.link,
      'subject': () => this.subject,
      'body': () => this.body,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'timestamp': (x) => this.timestamp = x,
      'type': (x) => this.type = x,
      'origin': (x) => this.origin = x,
      'link': (x) => this.link = x,
      'subject': (x) => this.subject = x,
      'body': (x) => this.body = x,
    };
  }

  String get url {
    return '/notifications/' + this.id;
  }
}

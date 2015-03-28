library event_mod;

import 'package:blckur/model.dart' as model;

class Event extends model.Model {
  String id;
  String type;
  String resourceId;
  double timestamp;

  model.Model newModel() {
    return new Event();
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'type': () => this.type,
      'resource_id': () => this.resourceId,
      'timestamp': () => this.timestamp,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'type': (x) => this.type = x,
      'resource_id': (x) => this.resourceId = x,
      'timestamp': (x) => this.timestamp = x,
    };
  }
}

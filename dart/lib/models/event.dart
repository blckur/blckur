library event_mod;

import 'package:blckur/model.dart' as model;

class Event extends model.Model {
  String id;
  String type;
  dynamic data;

  model.Model newModel() {
    return new Event();
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'type': () => this.type,
      'data': () => this.data,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'type': (x) => this.type = x,
      'data': (x) => this.data = x,
    };
  }
}

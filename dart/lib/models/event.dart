library event_mod;

import 'package:blckur/model.dart' as model;

import 'package:angular/angular.dart' show Injectable;
import 'package:angular/angular.dart' as ng;

@Injectable()
class Event extends model.Model {
  String id;
  String type;
  String resourceId;
  double timestamp;

  Event(ng.Http http) : super(http);

  model.Model newModel() {
    return new Event(this.http);
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

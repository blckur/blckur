library filter_type_mod;

import 'package:blckur/model.dart' as model;

class FilterType extends model.Model {
  int id;
  String label;
  String type;
  String valueType;
  String valueLabel;
  String valueHolder;

  model.Model newModel() {
    return new FilterType();
  }

  Map<String, Function> get getters {
    return {
      'id': () => this.id,
      'label': () => this.label,
      'type': () => this.type,
      'value_type': () => this.valueType,
      'value_label': () => this.valueLabel,
      'value_holder': () => this.valueHolder,
    };
  }

  Map<String, Function> get setters {
    return {
      'id': (x) => this.id = x,
      'label': (x) => this.label = x,
      'type': (x) => this.type = x,
      'value_type': (x) => this.valueType = x,
      'value_label': (x) => this.valueLabel = x,
      'value_holder': (x) => this.valueHolder = x,
    };
  }
}

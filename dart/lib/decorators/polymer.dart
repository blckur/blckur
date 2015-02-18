library polymer_dec;

import 'dart:html';
import 'package:angular/angular.dart' show Decorator;
import 'package:angular/angular.dart' as ng;
import 'package:observe/observe.dart' as obsrv;
import 'package:template_binding/template_binding.dart' as tmpl_bnd;

typedef dynamic _Nullary();

@Decorator(
  selector: '[py-*]'
)
class PolymerDec {
  PolymerDec(Node node, ng.Parser parser, ng.Scope scope) {
    Element element = node;

    element.attributes.forEach((key, value) {
      if (!key.startsWith('py-')) {
        return;
      }
      var bindVal = new ValueBindable();
      var expression = parser(value);

      tmpl_bnd.nodeBind(node).bind(key.substring(3), bindVal);

      if (expression.isAssignable) {
        bindVal.onChange = (value) {
          expression.assign(scope.context, value);
        };
      }

      scope.watch(value, bindVal.update);
    });
  }
}

class ValueBindable implements obsrv.Bindable {
  var _value;
  Function callback;
  Function onChange;

  void close() {
    this.callback = null;
  }

  dynamic open(Function callback) {
    this.callback = callback;
    return this._value;
  }

  void update(dynamic newValue, dynamic oldValue) {
    this._value = newValue;
    if (this.callback != null) {
      if (this.callback is _Nullary) {
        this.callback();
      }
      else {
        this.callback(this._value);
      }
    }
  }

  void set value(dynamic value) {
    this._value = value;
    if (this.onChange != null) {
      this.onChange(this._value);
    }
  }
  dynamic get value {
    return this._value;
  }

  void deliver() {}
}

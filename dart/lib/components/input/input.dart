library input_comp;

import 'package:angular/angular.dart' show Component, NgAttr, NgTwoWay;
import 'package:angular/angular.dart' as ng;
import 'dart:html' as dom;

@Component(
  selector: 'x-input',
  templateUrl: 'packages/blckur/components/input/input.html',
  cssUrl: 'packages/blckur/components/input/input.css'
)
class InputComp implements ng.ShadowRootAware {
  String label;
  String labelClass;
  bool labelFloat;

  @NgAttr('type')
  String type;

  var _placeholder;
  @NgAttr('placeholder')
  void set placeholder(String val) {
    this._placeholder = val;

    if (this.error == null) {
      this.label = val;
    }
  }
  String get placeholder {
    return this._placeholder;
  }

  @NgTwoWay('model')
  String model;

  var _error;
  @NgTwoWay('error')
  void set error(String val) {
    if (this._error != val) {
      this._error = val;
      if (val == null) {
        this.label = this.placeholder;
        this.labelClass = null;
        this.labelFloat = false;
      }
      else {
        this.label = val;
        this.labelClass = 'invalid';
        this.labelFloat = true;
      }
    }
  }
  String get error {
    return this._error;
  }

  InputComp() : labelFloat = false;

  void onShadowRoot(dom.ShadowRoot root) {
    root.querySelector('input').onInput.listen((_) {
      this.error = null;
    });
  }
}
